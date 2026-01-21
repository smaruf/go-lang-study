package main

import (
	"errors"
	"testing"
)

// Example interfaces for mocking

// EmailSender interface for sending emails
type EmailSender interface {
	SendEmail(to, subject, body string) error
}

// Database interface for database operations
type Database interface {
	GetUser(id string) (*User, error)
	SaveUser(user *User) error
}

// PaymentGateway interface for payment processing
type PaymentGateway interface {
	ProcessPayment(amount float64, cardNumber string) (string, error)
}

// Mock implementations

// MockEmailSender is a mock implementation of EmailSender
type MockEmailSender struct {
	SendEmailCalled bool
	LastTo          string
	LastSubject     string
	LastBody        string
	ReturnError     error
}

func (m *MockEmailSender) SendEmail(to, subject, body string) error {
	m.SendEmailCalled = true
	m.LastTo = to
	m.LastSubject = subject
	m.LastBody = body
	return m.ReturnError
}

// MockDatabase is a mock implementation of Database
type MockDatabase struct {
	GetUserCalled bool
	SaveUserCalled bool
	Users         map[string]*User
	ReturnError   error
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		Users: make(map[string]*User),
	}
}

func (m *MockDatabase) GetUser(id string) (*User, error) {
	m.GetUserCalled = true
	if m.ReturnError != nil {
		return nil, m.ReturnError
	}
	user, exists := m.Users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockDatabase) SaveUser(user *User) error {
	m.SaveUserCalled = true
	if m.ReturnError != nil {
		return m.ReturnError
	}
	m.Users[user.ID] = user
	return nil
}

// MockPaymentGateway is a mock implementation of PaymentGateway
type MockPaymentGateway struct {
	ProcessPaymentCalled bool
	LastAmount          float64
	LastCardNumber      string
	ReturnTransactionID string
	ReturnError         error
}

func (m *MockPaymentGateway) ProcessPayment(amount float64, cardNumber string) (string, error) {
	m.ProcessPaymentCalled = true
	m.LastAmount = amount
	m.LastCardNumber = cardNumber
	if m.ReturnError != nil {
		return "", m.ReturnError
	}
	return m.ReturnTransactionID, nil
}

// UserService demonstrates dependency injection for testing
type UserService struct {
	db      Database
	emailer EmailSender
}

func NewUserService(db Database, emailer EmailSender) *UserService {
	return &UserService{
		db:      db,
		emailer: emailer,
	}
}

func (s *UserService) RegisterUser(name, email string) (*User, error) {
	user := &User{
		ID:    "user-123",
		Name:  name,
		Email: email,
		Age:   0,
	}

	if err := s.db.SaveUser(user); err != nil {
		return nil, err
	}

	if err := s.emailer.SendEmail(email, "Welcome!", "Thank you for registering"); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(id string) (*User, error) {
	return s.db.GetUser(id)
}

// Tests using mocks

func TestUserService_RegisterUser(t *testing.T) {
	tests := []struct {
		name          string
		userName      string
		userEmail     string
		dbError       error
		emailError    error
		expectError   bool
	}{
		{
			name:        "successful registration",
			userName:    "John Doe",
			userEmail:   "john@example.com",
			dbError:     nil,
			emailError:  nil,
			expectError: false,
		},
		{
			name:        "database error",
			userName:    "Jane Doe",
			userEmail:   "jane@example.com",
			dbError:     errors.New("database connection failed"),
			emailError:  nil,
			expectError: true,
		},
		{
			name:        "email error",
			userName:    "Bob Smith",
			userEmail:   "bob@example.com",
			dbError:     nil,
			emailError:  errors.New("email service unavailable"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockDB := NewMockDatabase()
			mockDB.ReturnError = tt.dbError

			mockEmailer := &MockEmailSender{
				ReturnError: tt.emailError,
			}

			// Create service with mocks
			service := NewUserService(mockDB, mockEmailer)

			// Execute
			user, err := service.RegisterUser(tt.userName, tt.userEmail)

			// Verify
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user == nil {
					t.Errorf("expected user, got nil")
				}
				if user.Name != tt.userName {
					t.Errorf("expected user name %q, got %q", tt.userName, user.Name)
				}
				if user.Email != tt.userEmail {
					t.Errorf("expected user email %q, got %q", tt.userEmail, user.Email)
				}

				// Verify mock interactions
				if !mockDB.SaveUserCalled {
					t.Errorf("expected SaveUser to be called")
				}
				if !mockEmailer.SendEmailCalled {
					t.Errorf("expected SendEmail to be called")
				}
				if mockEmailer.LastTo != tt.userEmail {
					t.Errorf("expected email to %q, got %q", tt.userEmail, mockEmailer.LastTo)
				}
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		setupUser   *User
		dbError     error
		expectError bool
	}{
		{
			name:   "existing user",
			userID: "user-1",
			setupUser: &User{
				ID:    "user-1",
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
			},
			dbError:     nil,
			expectError: false,
		},
		{
			name:        "non-existing user",
			userID:      "user-999",
			setupUser:   nil,
			dbError:     nil,
			expectError: true,
		},
		{
			name:        "database error",
			userID:      "user-2",
			setupUser:   nil,
			dbError:     errors.New("database connection failed"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock
			mockDB := NewMockDatabase()
			if tt.setupUser != nil {
				mockDB.Users[tt.setupUser.ID] = tt.setupUser
			}
			mockDB.ReturnError = tt.dbError

			// Create service with mock
			service := NewUserService(mockDB, &MockEmailSender{})

			// Execute
			user, err := service.GetUserByID(tt.userID)

			// Verify
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user == nil {
					t.Errorf("expected user, got nil")
				}
				if user.ID != tt.userID {
					t.Errorf("expected user ID %q, got %q", tt.userID, user.ID)
				}
			}

			// Verify mock interaction
			if !mockDB.GetUserCalled {
				t.Errorf("expected GetUser to be called")
			}
		})
	}
}

// OrderService demonstrates more complex mocking scenarios
type OrderService struct {
	db             Database
	paymentGateway PaymentGateway
	emailer        EmailSender
}

func NewOrderService(db Database, pg PaymentGateway, emailer EmailSender) *OrderService {
	return &OrderService{
		db:             db,
		paymentGateway: pg,
		emailer:        emailer,
	}
}

func (s *OrderService) PlaceOrder(userID string, amount float64, cardNumber string) (string, error) {
	// Get user
	user, err := s.db.GetUser(userID)
	if err != nil {
		return "", err
	}

	// Process payment
	transactionID, err := s.paymentGateway.ProcessPayment(amount, cardNumber)
	if err != nil {
		return "", err
	}

	// Send confirmation email
	subject := "Order Confirmation"
	body := "Your order has been placed successfully. Transaction ID: " + transactionID
	if err := s.emailer.SendEmail(user.Email, subject, body); err != nil {
		return "", err
	}

	return transactionID, nil
}

func TestOrderService_PlaceOrder(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		amount           float64
		cardNumber       string
		setupUser        *User
		dbError          error
		paymentError     error
		emailError       error
		transactionID    string
		expectError      bool
	}{
		{
			name:       "successful order",
			userID:     "user-1",
			amount:     99.99,
			cardNumber: "1234-5678-9012-3456",
			setupUser: &User{
				ID:    "user-1",
				Name:  "John Doe",
				Email: "john@example.com",
			},
			dbError:       nil,
			paymentError:  nil,
			emailError:    nil,
			transactionID: "txn-12345",
			expectError:   false,
		},
		{
			name:          "user not found",
			userID:        "user-999",
			amount:        99.99,
			cardNumber:    "1234-5678-9012-3456",
			setupUser:     nil,
			dbError:       nil,
			paymentError:  nil,
			emailError:    nil,
			transactionID: "",
			expectError:   true,
		},
		{
			name:       "payment failed",
			userID:     "user-1",
			amount:     99.99,
			cardNumber: "1234-5678-9012-3456",
			setupUser: &User{
				ID:    "user-1",
				Name:  "John Doe",
				Email: "john@example.com",
			},
			dbError:       nil,
			paymentError:  errors.New("insufficient funds"),
			emailError:    nil,
			transactionID: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockDB := NewMockDatabase()
			if tt.setupUser != nil {
				mockDB.Users[tt.setupUser.ID] = tt.setupUser
			}
			mockDB.ReturnError = tt.dbError

			mockPayment := &MockPaymentGateway{
				ReturnTransactionID: tt.transactionID,
				ReturnError:         tt.paymentError,
			}

			mockEmailer := &MockEmailSender{
				ReturnError: tt.emailError,
			}

			// Create service
			service := NewOrderService(mockDB, mockPayment, mockEmailer)

			// Execute
			txnID, err := service.PlaceOrder(tt.userID, tt.amount, tt.cardNumber)

			// Verify
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if txnID != tt.transactionID {
					t.Errorf("expected transaction ID %q, got %q", tt.transactionID, txnID)
				}

				// Verify all mocks were called
				if !mockDB.GetUserCalled {
					t.Errorf("expected GetUser to be called")
				}
				if !mockPayment.ProcessPaymentCalled {
					t.Errorf("expected ProcessPayment to be called")
				}
				if !mockEmailer.SendEmailCalled {
					t.Errorf("expected SendEmail to be called")
				}

				// Verify mock arguments
				if mockPayment.LastAmount != tt.amount {
					t.Errorf("expected payment amount %.2f, got %.2f", tt.amount, mockPayment.LastAmount)
				}
			}
		})
	}
}
