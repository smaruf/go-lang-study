package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// AIService handles AI questions using intent routing.
type AIService struct {
	ollamaURL   string
	ollamaModel string
	timeout     time.Duration
	opensky     *OpenSkyClient
	pricing     *PricingService
}

// NewAIService creates a new AI service.
func NewAIService(ollamaURL, ollamaModel string, timeoutSecs int, opensky *OpenSkyClient, pricing *PricingService) *AIService {
	return &AIService{
		ollamaURL:   ollamaURL,
		ollamaModel: ollamaModel,
		timeout:     time.Duration(timeoutSecs) * time.Second,
		opensky:     opensky,
		pricing:     pricing,
	}
}

var (
	icaoRe     = regexp.MustCompile(`\b([a-f0-9]{6})\b`)
	callsignRe = regexp.MustCompile(`\b([a-z]{2,3}\d{1,4}[a-z]?)\b`)
	bboxRe     = regexp.MustCompile(`(-?\d+\.?\d*)[,\s]+(-?\d+\.?\d*)[,\s]+(-?\d+\.?\d*)[,\s]+(-?\d+\.?\d*)`)
)

type intentResult struct {
	tool   string
	result interface{}
}

func (a *AIService) routeIntent(question string) intentResult {
	q := strings.ToLower(question)

	// Flight lookup
	icaoMatch := icaoRe.FindString(q)
	callsignMatch := callsignRe.FindString(q)
	if (strings.Contains(q, "flight ") || strings.Contains(q, "aircraft ") || strings.Contains(q, "plane ")) &&
		(icaoMatch != "" || callsignMatch != "") {
		ident := icaoMatch
		if ident == "" {
			ident = callsignMatch
		}
		result, _ := a.opensky.FetchFlightDetail(ident)
		if result == nil {
			return intentResult{tool: "lookup_flight", result: map[string]string{"error": fmt.Sprintf("Flight %s not found", ident)}}
		}
		return intentResult{tool: "lookup_flight", result: result}
	}

	// BBox
	bboxMatch := bboxRe.FindStringSubmatch(q)
	if len(bboxMatch) == 5 {
		bboxKW := []string{"area", "region", "bbox", "over", "around"}
		for _, kw := range bboxKW {
			if strings.Contains(q, kw) {
				var lamin, lamax, lomin, lomax float64
				fmt.Sscanf(bboxMatch[1], "%f", &lamin)
				fmt.Sscanf(bboxMatch[2], "%f", &lamax)
				fmt.Sscanf(bboxMatch[3], "%f", &lomin)
				fmt.Sscanf(bboxMatch[4], "%f", &lomax)
				flights, _ := a.opensky.FetchFlights(&lamin, &lamax, &lomin, &lomax)
				if len(flights) > 20 {
					flights = flights[:20]
				}
				return intentResult{tool: "flights_in_bbox", result: flights}
			}
		}
	}

	// Route search
	routeKW := []string{"route", "flight from", "fly from", "ticket", "price", "cheapest", "cost"}
	for _, kw := range routeKW {
		if strings.Contains(q, kw) {
			known := []string{"dhaka", "dac", "gdansk", "gdn", "warsaw", "waw", "london", "lhr", "dubai", "dxb"}
			var found []string
			for _, o := range known {
				if strings.Contains(q, o) {
					found = append(found, strings.ToUpper(o))
				}
			}
			origin, dest := "", ""
			if len(found) > 0 {
				origin = found[0]
			}
			if len(found) > 1 {
				dest = found[1]
			}
			results := SearchRoutes(origin, dest, nil, nil)
			return intentResult{tool: "search_routes", result: results}
		}
	}

	// Default: list all routes
	return intentResult{tool: "list_routes", result: GetAll()}
}

func buildContext(toolName string, toolResult interface{}) string {
	switch toolName {
	case "list_routes":
		routes, ok := toolResult.(map[string]Route)
		if !ok {
			return ""
		}
		lines := []string{"Available routes:"}
		for code, r := range routes {
			lines = append(lines, fmt.Sprintf("  %s: %s -> %s, %.1fh, %d stop(s), ~$%.0f %s",
				code, r.Origin, r.Destination, r.DurationHours, len(r.Stoppages), r.BasePrice, r.Currency))
		}
		return strings.Join(lines, "\n")

	case "search_routes":
		entries, ok := toolResult.([]RouteEntry)
		if !ok || len(entries) == 0 {
			return "No matching routes found."
		}
		lines := []string{"Matching routes:"}
		for _, e := range entries {
			lines = append(lines, fmt.Sprintf("  %s: %s -> %s, %.1fh, %d stop(s), ~$%.0f, Airlines: %s",
				e.Code, e.Route.Origin, e.Route.Destination, e.Route.DurationHours,
				len(e.Route.Stoppages), e.Route.BasePrice, strings.Join(e.Route.Airlines, ", ")))
		}
		return strings.Join(lines, "\n")

	case "flights_in_bbox":
		flights, ok := toolResult.([]FlightState)
		if !ok || len(flights) == 0 {
			return "No flights found in that area."
		}
		lines := []string{fmt.Sprintf("Flights in area (%d shown):", len(flights))}
		limit := 10
		if len(flights) < limit {
			limit = len(flights)
		}
		for _, f := range flights[:limit] {
			lat, lon, alt := 0.0, 0.0, 0.0
			if f.Lat != nil {
				lat = *f.Lat
			}
			if f.Lon != nil {
				lon = *f.Lon
			}
			if f.Altitude != nil {
				alt = *f.Altitude
			}
			lines = append(lines, fmt.Sprintf("  %s %s from %s at lat=%.4f, lon=%.4f, alt=%.0fm",
				f.ICAO24, f.Callsign, f.OriginCountry, lat, lon, alt))
		}
		return strings.Join(lines, "\n")

	case "lookup_flight":
		if m, ok := toolResult.(map[string]string); ok {
			if errMsg, has := m["error"]; has {
				return errMsg
			}
		}
		f, ok := toolResult.(*FlightState)
		if !ok || f == nil {
			return "Flight not found."
		}
		lat, lon, alt, vel, hdg := 0.0, 0.0, 0.0, 0.0, 0.0
		if f.Lat != nil {
			lat = *f.Lat
		}
		if f.Lon != nil {
			lon = *f.Lon
		}
		if f.Altitude != nil {
			alt = *f.Altitude
		}
		if f.Velocity != nil {
			vel = *f.Velocity
		}
		if f.Heading != nil {
			hdg = *f.Heading
		}
		return fmt.Sprintf("Flight %s (%s): from %s, lat=%.4f, lon=%.4f, altitude=%.0fm, speed=%.0fm/s, heading=%.0f°",
			f.ICAO24, f.Callsign, f.OriginCountry, lat, lon, alt, vel, hdg)
	}
	return fmt.Sprintf("%v", toolResult)
}

// Ask processes an AI question with intent routing.
func (a *AIService) Ask(question, sessionContext string) string {
	intent := a.routeIntent(question)
	context := buildContext(intent.tool, intent.result)
	if sessionContext != "" {
		context = sessionContext + "\n\n" + context
	}

	prompt := fmt.Sprintf(
		"You are a helpful flight assistant. Use the following data to answer the question.\n\n%s\n\nQuestion: %s\nAnswer:",
		context, question,
	)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":  a.ollamaModel,
		"prompt": prompt,
		"stream": false,
	})

	client := &http.Client{Timeout: a.timeout}
	resp, err := client.Post(a.ollamaURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		log.Printf("Ollama request failed: %v", err)
		return fmt.Sprintf("AI unavailable. Based on available data: %s", context)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var ollamaResp struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &ollamaResp); err != nil || ollamaResp.Response == "" {
		return fmt.Sprintf("AI unavailable. Based on available data: %s", context)
	}
	return ollamaResp.Response
}
