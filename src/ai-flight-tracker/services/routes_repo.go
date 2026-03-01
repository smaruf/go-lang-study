package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// Route represents a flight route.
type Route struct {
	Origin        string   `json:"origin"`
	Destination   string   `json:"destination"`
	BasePrice     float64  `json:"base_price"`
	Currency      string   `json:"currency"`
	DurationHours float64  `json:"duration_hours"`
	Stoppages     []string `json:"stoppages"`
	Priority      int      `json:"priority"`
	Airlines      []string `json:"airlines"`
}

// RouteEntry combines route code + data.
type RouteEntry struct {
	Code  string
	Route Route
}

var (
	routesMu    sync.RWMutex
	routesCache map[string]Route
)

func loadRoutes() map[string]Route {
	routesMu.RLock()
	if routesCache != nil {
		defer routesMu.RUnlock()
		return routesCache
	}
	routesMu.RUnlock()

	routesMu.Lock()
	defer routesMu.Unlock()
	if routesCache != nil {
		return routesCache
	}

	// Find routes.json relative to this file (or via working directory)
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	candidates := []string{
		filepath.Join(dir, "..", "routes.json"),
		"routes.json",
		filepath.Join("src", "ai-flight-tracker", "routes.json"),
	}
	var data []byte
	var err error
	for _, p := range candidates {
		data, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		// Return empty map if file not found
		routesCache = map[string]Route{}
		return routesCache
	}

	var routes map[string]Route
	if err := json.Unmarshal(data, &routes); err != nil {
		routesCache = map[string]Route{}
		return routesCache
	}
	routesCache = routes
	return routesCache
}

// GetAll returns all routes.
func GetAll() map[string]Route {
	return loadRoutes()
}

// GetByCode returns a single route by code.
func GetByCode(code string) (Route, bool) {
	routes := loadRoutes()
	r, ok := routes[strings.ToUpper(code)]
	return r, ok
}

// SearchRoutes filters routes by origin, destination, max stops, and max price.
func SearchRoutes(origin, destination string, maxStops *int, maxPrice *float64) []RouteEntry {
	routes := loadRoutes()
	var results []RouteEntry
	for code, r := range routes {
		if origin != "" && !strings.Contains(strings.ToUpper(r.Origin), strings.ToUpper(origin)) {
			continue
		}
		if destination != "" && !strings.Contains(strings.ToUpper(r.Destination), strings.ToUpper(destination)) {
			continue
		}
		if maxStops != nil && len(r.Stoppages) > *maxStops {
			continue
		}
		if maxPrice != nil && r.BasePrice > *maxPrice {
			continue
		}
		results = append(results, RouteEntry{Code: code, Route: r})
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Route.Priority != results[j].Route.Priority {
			return results[i].Route.Priority < results[j].Route.Priority
		}
		return results[i].Code < results[j].Code
	})
	return results
}
