package services_test

import (
	"strings"
	"testing"

	. "github.com/smaruf/go-lang-study/src/ai-flight-tracker/services"
)

func TestGetAll(t *testing.T) {
	routes := GetAll()
	if len(routes) == 0 {
		t.Skip("routes.json not found in test environment - skipping")
	}
	if _, ok := routes["DAC-GDN"]; !ok {
		t.Error("expected DAC-GDN route to exist")
	}
	if _, ok := routes["WAW-DAC"]; !ok {
		t.Error("expected WAW-DAC route to exist")
	}
}

func TestGetByCode(t *testing.T) {
	route, ok := GetByCode("DAC-GDN")
	if !ok {
		t.Skip("routes.json not found in test environment - skipping")
	}
	if route.Origin == "" {
		t.Error("expected non-empty origin")
	}
	if route.Priority != 1 {
		t.Errorf("expected priority 1, got %d", route.Priority)
	}
}

func TestSearchRoutes(t *testing.T) {
	results := SearchRoutes("Dhaka", "", nil, nil)
	if len(results) == 0 {
		t.Skip("routes.json not found in test environment - skipping")
	}
	for _, e := range results {
		if !strings.Contains(strings.ToUpper(e.Route.Origin), "DHAKA") {
			t.Errorf("expected origin to contain Dhaka, got %s", e.Route.Origin)
		}
	}
}

func TestSearchRoutesMaxStops(t *testing.T) {
	maxStops := 0
	results := SearchRoutes("", "", &maxStops, nil)
	if len(results) == 0 {
		t.Skip("no direct routes or routes.json not found - skipping")
	}
	for _, e := range results {
		if len(e.Route.Stoppages) > maxStops {
			t.Errorf("expected max %d stops, got %d for route %s", maxStops, len(e.Route.Stoppages), e.Code)
		}
	}
}

func TestSearchRoutesMaxPrice(t *testing.T) {
	maxPrice := 500.0
	results := SearchRoutes("", "", nil, &maxPrice)
	if len(results) == 0 {
		t.Skip("no routes below $500 or routes.json not found - skipping")
	}
	for _, e := range results {
		if e.Route.BasePrice > maxPrice {
			t.Errorf("expected price <= %.0f, got %.0f for route %s", maxPrice, e.Route.BasePrice, e.Code)
		}
	}
}
