package services_test

import (
	"testing"
	"time"

	. "github.com/smaruf/go-lang-study/src/ai-flight-tracker/services"
)

func TestPricingCalculate(t *testing.T) {
	svc := NewPricingService(300)
	route := Route{
		Origin:        "Dhaka (DAC)",
		Destination:   "Warsaw (WAW)",
		BasePrice:     780,
		Currency:      "USD",
		DurationHours: 11.5,
		Priority:      1,
		Airlines:      []string{"Qatar Airways"},
	}
	departure := time.Now().AddDate(0, 0, 30)
	result := svc.Calculate("DAC-WAW", route, departure)
	if result.Price <= 0 {
		t.Error("expected positive price")
	}
	if result.Currency != "USD" {
		t.Errorf("expected USD, got %s", result.Currency)
	}
	if result.IsLive {
		t.Error("expected simulated price, not live")
	}
}

func TestPricingLastMinuteSurcharge(t *testing.T) {
	svc := NewPricingService(300)
	route := Route{BasePrice: 780, Currency: "USD"}
	lastMinute := time.Now().Add(3 * 24 * time.Hour) // 3 days away
	result := svc.Calculate("DAC-WAW", route, lastMinute)
	if result.Price < 780 {
		t.Errorf("expected last-minute price >= base price, got %.2f", result.Price)
	}
}

func TestPricingEarlyBirdDiscount(t *testing.T) {
	svc := NewPricingService(300)
	route := Route{BasePrice: 780, Currency: "USD"}
	earlyBird := time.Now().AddDate(0, 0, 90)
	result := svc.Calculate("DAC-WAW", route, earlyBird)
	// Early bird should be at most base price (0.85 * base * 1.05 = ~0.89 * base)
	if result.Price > route.BasePrice {
		t.Errorf("expected early bird price <= base price, got %.2f (base=%.2f)", result.Price, route.BasePrice)
	}
}

func TestPricingCaching(t *testing.T) {
	svc := NewPricingService(300)
	route := Route{BasePrice: 780, Currency: "USD"}
	departure := time.Now().AddDate(0, 0, 30)
	result1 := svc.Calculate("DAC-WAW", route, departure)
	result2 := svc.Calculate("DAC-WAW", route, departure)
	if result1.Price != result2.Price {
		t.Error("expected same price on second call (caching)")
	}
}
