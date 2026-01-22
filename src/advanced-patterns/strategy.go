package main

import "fmt"

// Strategy interface defines the contract for all strategies
type Strategy interface {
	Execute(a, b int) int
	GetName() string
}

// AddStrategy implements addition
type AddStrategy struct{}

func (s *AddStrategy) Execute(a, b int) int {
	return a + b
}

func (s *AddStrategy) GetName() string {
	return "Addition"
}

// SubtractStrategy implements subtraction
type SubtractStrategy struct{}

func (s *SubtractStrategy) Execute(a, b int) int {
	return a - b
}

func (s *SubtractStrategy) GetName() string {
	return "Subtraction"
}

// MultiplyStrategy implements multiplication
type MultiplyStrategy struct{}

func (s *MultiplyStrategy) Execute(a, b int) int {
	return a * b
}

func (s *MultiplyStrategy) GetName() string {
	return "Multiplication"
}

// DivideStrategy implements division
type DivideStrategy struct{}

func (s *DivideStrategy) Execute(a, b int) int {
	if b == 0 {
		return 0 // Simple error handling
	}
	return a / b
}

func (s *DivideStrategy) GetName() string {
	return "Division"
}

// Context uses a strategy
type Context struct {
	strategy Strategy
}

// NewContext creates a new context with a strategy
func NewContext(strategy Strategy) *Context {
	return &Context{
		strategy: strategy,
	}
}

// SetStrategy changes the strategy at runtime
func (c *Context) SetStrategy(strategy Strategy) {
	c.strategy = strategy
}

// ExecuteStrategy executes the current strategy
func (c *Context) ExecuteStrategy(a, b int) int {
	return c.strategy.Execute(a, b)
}

// Real-world example: Payment processing

// PaymentStrategy defines payment method interface
type PaymentStrategy interface {
	Pay(amount float64) string
	GetType() string
}

// CreditCardPayment implements credit card payment
type CreditCardPayment struct {
	cardNumber string
	cvv        string
}

func NewCreditCardPayment(cardNumber, cvv string) *CreditCardPayment {
	return &CreditCardPayment{
		cardNumber: cardNumber,
		cvv:        cvv,
	}
}

func (p *CreditCardPayment) Pay(amount float64) string {
	return fmt.Sprintf("Paid $%.2f using Credit Card ending in %s", 
		amount, p.cardNumber[len(p.cardNumber)-4:])
}

func (p *CreditCardPayment) GetType() string {
	return "Credit Card"
}

// PayPalPayment implements PayPal payment
type PayPalPayment struct {
	email string
}

func NewPayPalPayment(email string) *PayPalPayment {
	return &PayPalPayment{
		email: email,
	}
}

func (p *PayPalPayment) Pay(amount float64) string {
	return fmt.Sprintf("Paid $%.2f using PayPal account %s", amount, p.email)
}

func (p *PayPalPayment) GetType() string {
	return "PayPal"
}

// CryptoPayment implements cryptocurrency payment
type CryptoPayment struct {
	walletAddress string
	currency      string
}

func NewCryptoPayment(walletAddress, currency string) *CryptoPayment {
	return &CryptoPayment{
		walletAddress: walletAddress,
		currency:      currency,
	}
}

func (p *CryptoPayment) Pay(amount float64) string {
	return fmt.Sprintf("Paid $%.2f using %s to wallet %s", 
		amount, p.currency, p.walletAddress[:10]+"...")
}

func (p *CryptoPayment) GetType() string {
	return "Cryptocurrency"
}

// ShoppingCart uses payment strategy
type ShoppingCart struct {
	items          []string
	total          float64
	paymentMethod  PaymentStrategy
}

func NewShoppingCart() *ShoppingCart {
	return &ShoppingCart{
		items: make([]string, 0),
	}
}

func (c *ShoppingCart) AddItem(item string, price float64) {
	c.items = append(c.items, item)
	c.total += price
}

func (c *ShoppingCart) SetPaymentMethod(method PaymentStrategy) {
	c.paymentMethod = method
}

func (c *ShoppingCart) Checkout() string {
	if c.paymentMethod == nil {
		return "Error: No payment method selected"
	}
	
	result := fmt.Sprintf("\nCheckout Summary:\n")
	result += fmt.Sprintf("Items: %v\n", c.items)
	result += fmt.Sprintf("Total: $%.2f\n", c.total)
	result += fmt.Sprintf("Payment Method: %s\n", c.paymentMethod.GetType())
	result += c.paymentMethod.Pay(c.total)
	
	return result
}

// Compression strategy example

// CompressionStrategy defines compression algorithm interface
type CompressionStrategy interface {
	Compress(data string) string
	GetAlgorithm() string
}

// ZipCompression implements ZIP compression
type ZipCompression struct{}

func (z *ZipCompression) Compress(data string) string {
	return fmt.Sprintf("[ZIP compressed: %s...]", data[:min(len(data), 10)])
}

func (z *ZipCompression) GetAlgorithm() string {
	return "ZIP"
}

// RARCompression implements RAR compression
type RARCompression struct{}

func (r *RARCompression) Compress(data string) string {
	return fmt.Sprintf("[RAR compressed: %s...]", data[:min(len(data), 10)])
}

func (r *RARCompression) GetAlgorithm() string {
	return "RAR"
}

// GzipCompression implements GZIP compression
type GzipCompression struct{}

func (g *GzipCompression) Compress(data string) string {
	return fmt.Sprintf("[GZIP compressed: %s...]", data[:min(len(data), 10)])
}

func (g *GzipCompression) GetAlgorithm() string {
	return "GZIP"
}

// FileCompressor uses compression strategy
type FileCompressor struct {
	strategy CompressionStrategy
}

func NewFileCompressor(strategy CompressionStrategy) *FileCompressor {
	return &FileCompressor{
		strategy: strategy,
	}
}

func (fc *FileCompressor) SetCompressionStrategy(strategy CompressionStrategy) {
	fc.strategy = strategy
}

func (fc *FileCompressor) CompressFile(filename, data string) string {
	compressed := fc.strategy.Compress(data)
	return fmt.Sprintf("File '%s' compressed using %s: %s", 
		filename, fc.strategy.GetAlgorithm(), compressed)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println("=== Basic Strategy Pattern ===")
	
	// Create context with different strategies
	context := NewContext(&AddStrategy{})
	fmt.Printf("%s: 10 + 5 = %d\n", context.strategy.GetName(), context.ExecuteStrategy(10, 5))

	context.SetStrategy(&SubtractStrategy{})
	fmt.Printf("%s: 10 - 5 = %d\n", context.strategy.GetName(), context.ExecuteStrategy(10, 5))

	context.SetStrategy(&MultiplyStrategy{})
	fmt.Printf("%s: 10 * 5 = %d\n", context.strategy.GetName(), context.ExecuteStrategy(10, 5))

	context.SetStrategy(&DivideStrategy{})
	fmt.Printf("%s: 10 / 5 = %d\n", context.strategy.GetName(), context.ExecuteStrategy(10, 5))

	fmt.Println("\n=== Payment Strategy Pattern ===")
	
	cart := NewShoppingCart()
	cart.AddItem("Laptop", 999.99)
	cart.AddItem("Mouse", 29.99)
	cart.AddItem("Keyboard", 79.99)

	// Pay with credit card
	cart.SetPaymentMethod(NewCreditCardPayment("1234567890123456", "123"))
	fmt.Println(cart.Checkout())

	// Create new cart and pay with PayPal
	cart2 := NewShoppingCart()
	cart2.AddItem("Headphones", 149.99)
	cart2.SetPaymentMethod(NewPayPalPayment("user@example.com"))
	fmt.Println(cart2.Checkout())

	// Create new cart and pay with crypto
	cart3 := NewShoppingCart()
	cart3.AddItem("Monitor", 399.99)
	cart3.SetPaymentMethod(NewCryptoPayment("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb", "Bitcoin"))
	fmt.Println(cart3.Checkout())

	fmt.Println("\n=== Compression Strategy Pattern ===")
	
	data := "This is a large file with lots of data that needs to be compressed"
	
	// Compress with ZIP
	compressor := NewFileCompressor(&ZipCompression{})
	fmt.Println(compressor.CompressFile("document.txt", data))

	// Compress with RAR
	compressor.SetCompressionStrategy(&RARCompression{})
	fmt.Println(compressor.CompressFile("archive.txt", data))

	// Compress with GZIP
	compressor.SetCompressionStrategy(&GzipCompression{})
	fmt.Println(compressor.CompressFile("backup.txt", data))
}
