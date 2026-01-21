package main

import (
	"fmt"
	"sync"
)

// Observer interface defines the contract for observers
type Observer interface {
	Update(subject Subject, data interface{})
	GetID() string
}

// Subject interface defines the contract for subjects (observables)
type Subject interface {
	Attach(observer Observer)
	Detach(observer Observer)
	Notify(data interface{})
}

// ConcreteSubject implements Subject
type ConcreteSubject struct {
	observers []Observer
	mu        sync.RWMutex
}

// NewConcreteSubject creates a new concrete subject
func NewConcreteSubject() *ConcreteSubject {
	return &ConcreteSubject{
		observers: make([]Observer, 0),
	}
}

// Attach adds an observer
func (s *ConcreteSubject) Attach(observer Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
	fmt.Printf("Observer '%s' attached\n", observer.GetID())
}

// Detach removes an observer
func (s *ConcreteSubject) Detach(observer Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for i, obs := range s.observers {
		if obs.GetID() == observer.GetID() {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			fmt.Printf("Observer '%s' detached\n", observer.GetID())
			return
		}
	}
}

// Notify notifies all observers
func (s *ConcreteSubject) Notify(data interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	fmt.Printf("Notifying %d observers...\n", len(s.observers))
	for _, observer := range s.observers {
		observer.Update(s, data)
	}
}

// ConcreteObserver implements Observer
type ConcreteObserver struct {
	id   string
	data interface{}
}

// NewConcreteObserver creates a new concrete observer
func NewConcreteObserver(id string) *ConcreteObserver {
	return &ConcreteObserver{
		id: id,
	}
}

// Update receives updates from the subject
func (o *ConcreteObserver) Update(subject Subject, data interface{}) {
	o.data = data
	fmt.Printf("Observer '%s' received update: %v\n", o.id, data)
}

// GetID returns the observer's ID
func (o *ConcreteObserver) GetID() string {
	return o.id
}

// GetData returns the observer's current data
func (o *ConcreteObserver) GetData() interface{} {
	return o.data
}

// Real-world example: Stock Market

// Stock represents a stock subject
type Stock struct {
	symbol    string
	price     float64
	observers []Observer
	mu        sync.RWMutex
}

// NewStock creates a new stock
func NewStock(symbol string, initialPrice float64) *Stock {
	return &Stock{
		symbol:    symbol,
		price:     initialPrice,
		observers: make([]Observer, 0),
	}
}

// Attach adds an observer
func (s *Stock) Attach(observer Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
}

// Detach removes an observer
func (s *Stock) Detach(observer Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for i, obs := range s.observers {
		if obs.GetID() == observer.GetID() {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			return
		}
	}
}

// Notify notifies all observers
func (s *Stock) Notify(data interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for _, observer := range s.observers {
		observer.Update(s, data)
	}
}

// SetPrice updates the stock price and notifies observers
func (s *Stock) SetPrice(newPrice float64) {
	s.mu.Lock()
	oldPrice := s.price
	s.price = newPrice
	s.mu.Unlock()

	priceChange := map[string]interface{}{
		"symbol":    s.symbol,
		"oldPrice":  oldPrice,
		"newPrice":  newPrice,
		"change":    newPrice - oldPrice,
		"changePct": ((newPrice - oldPrice) / oldPrice) * 100,
	}

	fmt.Printf("\n[%s] Price changed: $%.2f -> $%.2f\n", s.symbol, oldPrice, newPrice)
	s.Notify(priceChange)
}

// GetPrice returns the current stock price
func (s *Stock) GetPrice() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.price
}

// Investor represents a stock investor observer
type Investor struct {
	name      string
	portfolio map[string]float64
}

// NewInvestor creates a new investor
func NewInvestor(name string) *Investor {
	return &Investor{
		name:      name,
		portfolio: make(map[string]float64),
	}
}

// Update receives stock price updates
func (i *Investor) Update(subject Subject, data interface{}) {
	if priceChange, ok := data.(map[string]interface{}); ok {
		symbol := priceChange["symbol"].(string)
		newPrice := priceChange["newPrice"].(float64)
		changePct := priceChange["changePct"].(float64)

		i.portfolio[symbol] = newPrice

		if changePct > 5 {
			fmt.Printf("  [%s] 📈 BUY signal! %s up %.2f%%\n", i.name, symbol, changePct)
		} else if changePct < -5 {
			fmt.Printf("  [%s] 📉 SELL signal! %s down %.2f%%\n", i.name, symbol, changePct)
		} else {
			fmt.Printf("  [%s] ➡️  %s changed %.2f%% - HOLD\n", i.name, symbol, changePct)
		}
	}
}

// GetID returns the investor's name
func (i *Investor) GetID() string {
	return i.name
}

func main() {
	fmt.Println("=== Basic Observer Pattern ===")
	
	subject := NewConcreteSubject()
	
	observer1 := NewConcreteObserver("Observer-1")
	observer2 := NewConcreteObserver("Observer-2")
	observer3 := NewConcreteObserver("Observer-3")

	subject.Attach(observer1)
	subject.Attach(observer2)
	subject.Attach(observer3)

	fmt.Println("\nSending notification 1:")
	subject.Notify("Hello, observers!")

	fmt.Println("\nDetaching Observer-2:")
	subject.Detach(observer2)

	fmt.Println("\nSending notification 2:")
	subject.Notify("Observer-2 should not receive this")

	fmt.Println("\n=== Stock Market Observer Pattern ===")
	
	// Create stocks
	appleStock := NewStock("AAPL", 150.00)
	googleStock := NewStock("GOOGL", 2800.00)

	// Create investors
	investor1 := NewInvestor("Warren")
	investor2 := NewInvestor("Charlie")
	investor3 := NewInvestor("Peter")

	// Subscribe investors to stocks
	appleStock.Attach(investor1)
	appleStock.Attach(investor2)
	appleStock.Attach(investor3)

	googleStock.Attach(investor1)
	googleStock.Attach(investor3)

	// Simulate price changes
	fmt.Println("\n--- Market Opens ---")
	appleStock.SetPrice(160.00) // +6.67% increase
	
	googleStock.SetPrice(2700.00) // -3.57% decrease
	
	appleStock.SetPrice(145.00) // -9.38% decrease from 160
	
	googleStock.SetPrice(2900.00) // +7.41% increase from 2700

	fmt.Println("\n--- Investor 2 stops watching AAPL ---")
	appleStock.Detach(investor2)
	
	appleStock.SetPrice(155.00) // Only investor1 and investor3 notified
}
