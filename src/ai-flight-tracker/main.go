package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/smaruf/go-lang-study/src/ai-flight-tracker/config"
	"github.com/smaruf/go-lang-study/src/ai-flight-tracker/services"
)

var (
	cfg       *config.Config
	opensky   *services.OpenSkyClient
	pricing   *services.PricingService
	aiService *services.AIService
	icao24Re  = regexp.MustCompile(`^[a-fA-F0-9]{6}$`)
)

func main() {
	cfg = config.Load()
	opensky = services.NewOpenSkyClient(cfg.OpenSkyURL, cfg.OpenSkyUsername, cfg.OpenSkyPassword, cfg.OpenSkyTimeout)
	pricing = services.NewPricingService(cfg.CacheTTL)
	aiService = services.NewAIService(cfg.OllamaURL, cfg.OllamaModel, cfg.OllamaTimeout, opensky, pricing)

	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/routes", handleGetRoutes)
	mux.HandleFunc("/api/route/", handleGetRoute)
	mux.HandleFunc("/api/search", handleSearchRoutes)
	mux.HandleFunc("/api/flights", handleGetFlights)
	mux.HandleFunc("/api/flights/", handleGetFlight)
	mux.HandleFunc("/api/aviation/ask", handleAsk)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting Flight Tracker AI on %s", addr)

	srv := &http.Server{
		Addr:         addr,
		Handler:      withLogging(withCORS(mux, cfg.CORSOrigins)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	done := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
		close(done)
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
	<-done
	log.Println("Server stopped")
}

// ---- Middleware ----

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("method=%s path=%s status=%d latency_ms=%d",
			r.Method, r.URL.Path, rw.status, time.Since(start).Milliseconds())
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func withCORS(next http.Handler, origins string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origins)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ---- Helpers ----

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func parseFloatParam(r *http.Request, key string) *float64 {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil
	}
	return &f
}

func parseIntParam(r *http.Request, key string) *int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return nil
	}
	return &n
}

// ---- Handlers ----

func handleHealth(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}

func handleGetRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	allRoutes := services.GetAll()
	type routeItem struct {
		RouteCode     string   `json:"route_code"`
		Origin        string   `json:"origin"`
		Destination   string   `json:"destination"`
		DurationHours float64  `json:"duration_hours"`
		Stoppages     []string `json:"stoppages"`
		NumStops      int      `json:"num_stops"`
		Airlines      []string `json:"airlines"`
		BasePrice     float64  `json:"base_price"`
		Currency      string   `json:"currency"`
		IsPriority    bool     `json:"is_priority"`
		IsLive        bool     `json:"is_live"`
	}

	// Sort by priority then code
	entries := make([]services.RouteEntry, 0, len(allRoutes))
	for code, r := range allRoutes {
		entries = append(entries, services.RouteEntry{Code: code, Route: r})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Route.Priority != entries[j].Route.Priority {
			return entries[i].Route.Priority < entries[j].Route.Priority
		}
		return entries[i].Code < entries[j].Code
	})

	departure := time.Now().AddDate(0, 0, 30)
	var data []routeItem
	for _, e := range entries {
		priceInfo := pricing.Calculate(e.Code, e.Route, departure)
		data = append(data, routeItem{
			RouteCode:     e.Code,
			Origin:        e.Route.Origin,
			Destination:   e.Route.Destination,
			DurationHours: e.Route.DurationHours,
			Stoppages:     e.Route.Stoppages,
			NumStops:      len(e.Route.Stoppages),
			Airlines:      e.Route.Airlines,
			BasePrice:     priceInfo.Price,
			Currency:      e.Route.Currency,
			IsPriority:    e.Route.Priority == 1,
			IsLive:        priceInfo.IsLive,
		})
	}
	if data == nil {
		data = []routeItem{}
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"routes": data})
}

func handleGetRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	// Extract route code from path: /api/route/{code}
	code := strings.TrimPrefix(r.URL.Path, "/api/route/")
	code = strings.TrimSuffix(code, "/")
	if code == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Route code required"})
		return
	}
	route, ok := services.GetByCode(code)
	if !ok {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": "Route not found"})
		return
	}
	now := time.Now()
	type priceWindow struct {
		Price    float64 `json:"price"`
		Currency string  `json:"currency"`
		IsLive   bool    `json:"is_live"`
		Source   string  `json:"source"`
	}
	makePW := func(d time.Duration) priceWindow {
		pr := pricing.Calculate(code, route, now.Add(d))
		return priceWindow{Price: pr.Price, Currency: pr.Currency, IsLive: pr.IsLive, Source: pr.Source}
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"route_code":     strings.ToUpper(code),
		"origin":         route.Origin,
		"destination":    route.Destination,
		"duration_hours": route.DurationHours,
		"stoppages":      route.Stoppages,
		"num_stops":      len(route.Stoppages),
		"airlines":       route.Airlines,
		"is_priority":    route.Priority == 1,
		"pricing": map[string]priceWindow{
			"last_minute": makePW(5 * 24 * time.Hour),
			"two_weeks":   makePW(14 * 24 * time.Hour),
			"one_month":   makePW(30 * 24 * time.Hour),
			"early_bird":  makePW(90 * 24 * time.Hour),
		},
	})
}

func handleSearchRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")
	maxStops := parseIntParam(r, "max_stops")
	maxPrice := parseFloatParam(r, "max_price")

	results := services.SearchRoutes(origin, destination, maxStops, maxPrice)
	departure := time.Now().AddDate(0, 0, 30)

	type resultItem struct {
		RouteCode      string   `json:"route_code"`
		Origin         string   `json:"origin"`
		Destination    string   `json:"destination"`
		DurationHours  float64  `json:"duration_hours"`
		Stoppages      []string `json:"stoppages"`
		NumStops       int      `json:"num_stops"`
		Airlines       []string `json:"airlines"`
		EstimatedPrice float64  `json:"estimated_price"`
		Currency       string   `json:"currency"`
		IsPriority     bool     `json:"is_priority"`
		IsLive         bool     `json:"is_live"`
	}

	var data []resultItem
	for _, e := range results {
		priceInfo := pricing.Calculate(e.Code, e.Route, departure)
		data = append(data, resultItem{
			RouteCode:      e.Code,
			Origin:         e.Route.Origin,
			Destination:    e.Route.Destination,
			DurationHours:  e.Route.DurationHours,
			Stoppages:      e.Route.Stoppages,
			NumStops:       len(e.Route.Stoppages),
			Airlines:       e.Route.Airlines,
			EstimatedPrice: priceInfo.Price,
			Currency:       e.Route.Currency,
			IsPriority:     e.Route.Priority == 1,
			IsLive:         priceInfo.IsLive,
		})
	}
	if data == nil {
		data = []resultItem{}
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"results": data, "count": len(data)})
}

func handleGetFlights(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	lamin := parseFloatParam(r, "lamin")
	lamax := parseFloatParam(r, "lamax")
	lomin := parseFloatParam(r, "lomin")
	lomax := parseFloatParam(r, "lomax")

	bboxParams := []*float64{lamin, lamax, lomin, lomax}
	anySet := false
	for _, v := range bboxParams {
		if v != nil {
			anySet = true
			break
		}
	}
	allSet := lamin != nil && lamax != nil && lomin != nil && lomax != nil
	if anySet && !allSet {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "All bbox params required: lamin, lamax, lomin, lomax"})
		return
	}
	if allSet {
		if *lamin < -90 || *lamax > 90 || *lomin < -180 || *lomax > 180 {
			jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid bbox values"})
			return
		}
		if *lamin >= *lamax || *lomin >= *lomax {
			jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "lamin must be < lamax and lomin must be < lomax"})
			return
		}
	}

	flights, err := opensky.FetchFlights(lamin, lamax, lomin, lomax)
	if err != nil {
		log.Printf("Error fetching flights: %v", err)
		flights = []services.FlightState{}
	}
	if flights == nil {
		flights = []services.FlightState{}
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"flights": flights, "count": len(flights)})
}

func handleGetFlight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	icao24 := strings.TrimPrefix(r.URL.Path, "/api/flights/")
	icao24 = strings.TrimSuffix(icao24, "/")
	if icao24 == "" {
		// Redirect to list handler
		handleGetFlights(w, r)
		return
	}
	if !icao24Re.MatchString(icao24) {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid ICAO24 format (must be 6 hex chars)"})
		return
	}
	flight, err := opensky.FetchFlightDetail(strings.ToLower(icao24))
	if err != nil || flight == nil {
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": "Flight not found"})
		return
	}
	jsonResponse(w, http.StatusOK, flight)
}

func handleAsk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "JSON body required"})
		return
	}
	var body struct {
		Question string `json:"question"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON body"})
		return
	}
	question := strings.TrimSpace(body.Question)
	if question == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "question field is required"})
		return
	}
	if len(question) > 500 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "question too long (max 500 chars)"})
		return
	}
	answer := aiService.Ask(question, "")
	jsonResponse(w, http.StatusOK, map[string]string{"answer": answer})
}
