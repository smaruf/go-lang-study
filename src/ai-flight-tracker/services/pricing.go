package services

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// PriceResult represents the pricing result for a route.
type PriceResult struct {
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Route       string  `json:"route"`
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	IsLive      bool    `json:"is_live"`
	Source      string  `json:"source"`
}

type cachedEntry struct {
	data PriceResult
	ts   time.Time
}

// PricingService calculates prices with a TTL cache.
type PricingService struct {
	mu       sync.RWMutex
	cache    map[string]cachedEntry
	cacheTTL time.Duration
}

// NewPricingService creates a new pricing service.
func NewPricingService(cacheTTLSecs int) *PricingService {
	return &PricingService{
		cache:    make(map[string]cachedEntry),
		cacheTTL: time.Duration(cacheTTLSecs) * time.Second,
	}
}

func (p *PricingService) getCached(key string) (PriceResult, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	entry, ok := p.cache[key]
	if !ok {
		return PriceResult{}, false
	}
	if time.Since(entry.ts) > p.cacheTTL {
		return PriceResult{}, false
	}
	return entry.data, true
}

func (p *PricingService) setCached(key string, result PriceResult) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cache[key] = cachedEntry{data: result, ts: time.Now()}
}

// Calculate computes the price for a route given a departure date.
func (p *PricingService) Calculate(routeCode string, route Route, departureDate time.Time) PriceResult {
	if departureDate.IsZero() {
		departureDate = time.Now().AddDate(0, 0, 30)
	}
	cacheKey := fmt.Sprintf("%s_%s", routeCode, departureDate.Format("2006-01-02"))
	if cached, ok := p.getCached(cacheKey); ok {
		return cached
	}

	base := route.BasePrice
	days := int(time.Until(departureDate).Hours() / 24)
	switch {
	case days < 7:
		base *= 1.4
	case days < 14:
		base *= 1.2
	case days > 60:
		base *= 0.85
	}
	// Apply ±5% random variation
	price := base * (0.95 + rand.Float64()*0.10)
	price = float64(int(price*100)) / 100 // round to 2 decimal places

	result := PriceResult{
		Price:       price,
		Currency:    route.Currency,
		Route:       routeCode,
		Origin:      route.Origin,
		Destination: route.Destination,
		IsLive:      false,
		Source:      "simulated",
	}
	p.setCached(cacheKey, result)
	return result
}
