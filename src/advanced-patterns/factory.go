package main

import "fmt"

// Product interface defines the contract for all products
type Product interface {
	GetName() string
	GetPrice() float64
	Display()
}

// ConcreteProductA implements Product
type ConcreteProductA struct {
	name  string
	price float64
}

func (p *ConcreteProductA) GetName() string {
	return p.name
}

func (p *ConcreteProductA) GetPrice() float64 {
	return p.price
}

func (p *ConcreteProductA) Display() {
	fmt.Printf("Product A: %s - $%.2f\n", p.name, p.price)
}

// ConcreteProductB implements Product
type ConcreteProductB struct {
	name  string
	price float64
}

func (p *ConcreteProductB) GetName() string {
	return p.name
}

func (p *ConcreteProductB) GetPrice() float64 {
	return p.price
}

func (p *ConcreteProductB) Display() {
	fmt.Printf("Product B: %s - $%.2f\n", p.name, p.price)
}

// ProductType represents different types of products
type ProductType string

const (
	TypeA ProductType = "A"
	TypeB ProductType = "B"
)

// ProductFactory is responsible for creating products
type ProductFactory struct{}

// NewProductFactory creates a new product factory
func NewProductFactory() *ProductFactory {
	return &ProductFactory{}
}

// CreateProduct creates a product based on the type
func (f *ProductFactory) CreateProduct(productType ProductType, name string, price float64) (Product, error) {
	switch productType {
	case TypeA:
		return &ConcreteProductA{
			name:  name,
			price: price,
		}, nil
	case TypeB:
		return &ConcreteProductB{
			name:  name,
			price: price,
		}, nil
	default:
		return nil, fmt.Errorf("unknown product type: %s", productType)
	}
}

// Abstract Factory Pattern Example

// GUIFactory is an abstract factory interface
type GUIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
}

// Button interface
type Button interface {
	Render()
}

// Checkbox interface
type Checkbox interface {
	Render()
}

// Windows implementations
type WindowsButton struct{}

func (b *WindowsButton) Render() {
	fmt.Println("Rendering Windows button")
}

type WindowsCheckbox struct{}

func (c *WindowsCheckbox) Render() {
	fmt.Println("Rendering Windows checkbox")
}

type WindowsFactory struct{}

func (f *WindowsFactory) CreateButton() Button {
	return &WindowsButton{}
}

func (f *WindowsFactory) CreateCheckbox() Checkbox {
	return &WindowsCheckbox{}
}

// Mac implementations
type MacButton struct{}

func (b *MacButton) Render() {
	fmt.Println("Rendering Mac button")
}

type MacCheckbox struct{}

func (c *MacCheckbox) Render() {
	fmt.Println("Rendering Mac checkbox")
}

type MacFactory struct{}

func (f *MacFactory) CreateButton() Button {
	return &MacButton{}
}

func (f *MacFactory) CreateCheckbox() Checkbox {
	return &MacCheckbox{}
}

// GetGUIFactory returns the appropriate GUI factory based on OS
func GetGUIFactory(os string) GUIFactory {
	switch os {
	case "windows":
		return &WindowsFactory{}
	case "mac":
		return &MacFactory{}
	default:
		return &WindowsFactory{}
	}
}

func main() {
	fmt.Println("=== Simple Factory Pattern ===")
	factory := NewProductFactory()

	productA, _ := factory.CreateProduct(TypeA, "Laptop", 999.99)
	productA.Display()

	productB, _ := factory.CreateProduct(TypeB, "Mouse", 29.99)
	productB.Display()

	fmt.Println("\n=== Abstract Factory Pattern ===")
	
	// Create Windows UI
	fmt.Println("\nCreating Windows UI:")
	windowsFactory := GetGUIFactory("windows")
	winButton := windowsFactory.CreateButton()
	winCheckbox := windowsFactory.CreateCheckbox()
	winButton.Render()
	winCheckbox.Render()

	// Create Mac UI
	fmt.Println("\nCreating Mac UI:")
	macFactory := GetGUIFactory("mac")
	macButton := macFactory.CreateButton()
	macCheckbox := macFactory.CreateCheckbox()
	macButton.Render()
	macCheckbox.Render()

	fmt.Println("\n=== Factory Method Pattern ===")
	
	// Using factory method pattern
	creators := []Creator{
		&ConcreteCreatorA{},
		&ConcreteCreatorB{},
	}

	for _, creator := range creators {
		product := creator.FactoryMethod()
		product.Display()
	}
}

// Factory Method Pattern

type Creator interface {
	FactoryMethod() Product
}

type ConcreteCreatorA struct{}

func (c *ConcreteCreatorA) FactoryMethod() Product {
	return &ConcreteProductA{
		name:  "Factory Method Product A",
		price: 599.99,
	}
}

type ConcreteCreatorB struct{}

func (c *ConcreteCreatorB) FactoryMethod() Product {
	return &ConcreteProductB{
		name:  "Factory Method Product B",
		price: 399.99,
	}
}
