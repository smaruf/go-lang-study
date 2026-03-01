// Package server wraps the Flight Tracker AI HTTP server with programmatic start/stop.
package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/smaruf/go-lang-study/src/ai-flight-tracker/config"
	"github.com/smaruf/go-lang-study/src/ai-flight-tracker/services"
)

var icao24Re = regexp.MustCompile(`^[a-fA-F0-9]{6}$`)

// LogBuffer is a thread-safe in-memory log buffer with a notification channel.
type LogBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
	ch  chan string
}

// NewLogBuffer creates a LogBuffer.
func NewLogBuffer() *LogBuffer {
	return &LogBuffer{ch: make(chan string, 256)}
}

// Write implements io.Writer so it can be used as a log output.
func (lb *LogBuffer) Write(p []byte) (n int, err error) {
	lb.mu.Lock()
	lb.buf.Write(p)
	lb.mu.Unlock()
	line := string(p)
	select {
	case lb.ch <- line:
	default:
	}
	return len(p), nil
}

// Lines returns all accumulated log text.
func (lb *LogBuffer) Lines() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return lb.buf.String()
}

// Ch returns the channel that receives new log lines.
func (lb *LogBuffer) Ch() <-chan string {
	return lb.ch
}

// Server wraps the HTTP server with programmatic start/stop lifecycle.
type Server struct {
	cfg       *config.Config
	opensky   *services.OpenSkyClient
	pricing   *services.PricingService
	aiService *services.AIService
	httpSrv   *http.Server
	LogBuf    *LogBuffer
	logger    *log.Logger
}

// New creates a Server from configuration.
func New(cfg *config.Config) *Server {
	lb := NewLogBuffer()
	logger := log.New(lb, "", log.LstdFlags)

	opensky := services.NewOpenSkyClient(cfg.OpenSkyURL, cfg.OpenSkyUsername, cfg.OpenSkyPassword, cfg.OpenSkyTimeout)
	pricing := services.NewPricingService(cfg.CacheTTL)
	aiSvc := services.NewAIService(cfg.OllamaURL, cfg.OllamaModel, cfg.OllamaTimeout, opensky, pricing)

	s := &Server{
		cfg:       cfg,
		opensky:   opensky,
		pricing:   pricing,
		aiService: aiSvc,
		LogBuf:    lb,
		logger:    logger,
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/api/routes", s.handleGetRoutes)
	mux.HandleFunc("/api/route/", s.handleGetRoute)
	mux.HandleFunc("/api/search", s.handleSearchRoutes)
	mux.HandleFunc("/api/flights", s.handleGetFlights)
	mux.HandleFunc("/api/flights/", s.handleGetFlight)
	mux.HandleFunc("/api/aviation/ask", s.handleAsk)

	s.httpSrv = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      withLogging(withCORS(mux, cfg.CORSOrigins), logger),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		ErrorLog:     logger,
	}
	return s
}

// Start begins serving (non-blocking).
func (s *Server) Start() error {
	s.logger.Printf("Starting Flight Tracker AI on :%d", s.cfg.Port)
	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Printf("Server error: %v", err)
		}
	}()
	return nil
}

// Stop gracefully shuts down the HTTP server.
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.logger.Printf("Shutting down server...")
	return s.httpSrv.Shutdown(ctx)
}

// WebURL returns the base URL the server is listening on.
func (s *Server) WebURL() string {
	return fmt.Sprintf("http://localhost:%d", s.cfg.Port)
}

// ---- Middleware ----

func withLogging(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		logger.Printf("method=%s path=%s status=%d latency_ms=%d",
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

func jsonResp(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data) //nolint:errcheck
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

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	jsonResp(w, http.StatusOK, map[string]string{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}

func (s *Server) handleGetRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResp(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
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

	entries := make([]services.RouteEntry, 0, len(allRoutes))
	for code, route := range allRoutes {
		entries = append(entries, services.RouteEntry{Code: code, Route: route})
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
		priceInfo := s.pricing.Calculate(e.Code, e.Route, departure)
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
	jsonResp(w, http.StatusOK, map[string]interface{}{"routes": data})
}

func (s *Server) handleGetRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResp(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	code := strings.TrimPrefix(r.URL.Path, "/api/route/")
	code = strings.TrimSuffix(code, "/")
	if code == "" {
		jsonResp(w, http.StatusBadRequest, map[string]string{"error": "Route code required"})
		return
	}
	route, ok := services.GetByCode(code)
	if !ok {
		jsonResp(w, http.StatusNotFound, map[string]string{"error": "Route not found"})
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
		pr := s.pricing.Calculate(code, route, now.Add(d))
		return priceWindow{Price: pr.Price, Currency: pr.Currency, IsLive: pr.IsLive, Source: pr.Source}
	}
	jsonResp(w, http.StatusOK, map[string]interface{}{
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

func (s *Server) handleSearchRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResp(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
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
		priceInfo := s.pricing.Calculate(e.Code, e.Route, departure)
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
	jsonResp(w, http.StatusOK, map[string]interface{}{"results": data, "count": len(data)})
}

func (s *Server) handleGetFlights(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResp(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
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
		jsonResp(w, http.StatusBadRequest, map[string]string{"error": "All bbox params required: lamin, lamax, lomin, lomax"})
		return
	}
	if allSet {
		if *lamin < -90 || *lamax > 90 || *lomin < -180 || *lomax > 180 {
			jsonResp(w, http.StatusBadRequest, map[string]string{"error": "Invalid bbox values"})
			return
		}
		if *lamin >= *lamax || *lomin >= *lomax {
			jsonResp(w, http.StatusBadRequest, map[string]string{"error": "lamin must be < lamax and lomin must be < lomax"})
			return
		}
	}

	flights, err := s.opensky.FetchFlights(lamin, lamax, lomin, lomax)
	if err != nil {
		s.logger.Printf("Error fetching flights: %v", err)
		flights = []services.FlightState{}
	}
	if flights == nil {
		flights = []services.FlightState{}
	}
	jsonResp(w, http.StatusOK, map[string]interface{}{"flights": flights, "count": len(flights)})
}

func (s *Server) handleGetFlight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResp(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	icao24 := strings.TrimPrefix(r.URL.Path, "/api/flights/")
	icao24 = strings.TrimSuffix(icao24, "/")
	if icao24 == "" {
		s.handleGetFlights(w, r)
		return
	}
	if !icao24Re.MatchString(icao24) {
		jsonResp(w, http.StatusBadRequest, map[string]string{"error": "Invalid ICAO24 format (must be 6 hex chars)"})
		return
	}
	flight, err := s.opensky.FetchFlightDetail(strings.ToLower(icao24))
	if err != nil || flight == nil {
		jsonResp(w, http.StatusNotFound, map[string]string{"error": "Flight not found"})
		return
	}
	jsonResp(w, http.StatusOK, flight)
}

func (s *Server) handleAsk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResp(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		jsonResp(w, http.StatusBadRequest, map[string]string{"error": "JSON body required"})
		return
	}
	var body struct {
		Question string `json:"question"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonResp(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON body"})
		return
	}
	question := strings.TrimSpace(body.Question)
	if question == "" {
		jsonResp(w, http.StatusBadRequest, map[string]string{"error": "question field is required"})
		return
	}
	if len(question) > 500 {
		jsonResp(w, http.StatusBadRequest, map[string]string{"error": "question too long (max 500 chars)"})
		return
	}
	answer := s.aiService.Ask(question, "")
	jsonResp(w, http.StatusOK, map[string]string{"answer": answer})
}
