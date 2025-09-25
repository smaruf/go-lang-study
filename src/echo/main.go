package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// User represents a user entity
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Response represents API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Server wraps the Echo instance with dependencies
type Server struct {
	echo   *echo.Echo
	logger *logrus.Logger
}

// In-memory user storage for demo
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
}
var nextUserID = 3

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using defaults")
	}

	// Initialize server
	server := NewServer()
	
	// Setup routes
	server.setupRoutes()
	
	// Start server with graceful shutdown
	server.start()
}

// NewServer creates a new server instance
func NewServer() *Server {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	
	// Set log level from environment
	if level, err := logrus.ParseLevel(getEnv("LOG_LEVEL", "info")); err == nil {
		logger.SetLevel(level)
	}

	e := echo.New()
	e.HideBanner = true
	
	return &Server{
		echo:   e,
		logger: logger,
	}
}

// setupRoutes configures all routes and middleware
func (s *Server) setupRoutes() {
	// Global middleware
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.CORS())
	s.echo.Use(middleware.Secure())
	s.echo.Use(middleware.RequestID())
	
	// Rate limiting
	s.echo.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))

	// Custom middleware for structured logging
	s.echo.Use(s.structuredLoggingMiddleware())

	// Root route
	s.echo.GET("/", s.homeHandler)
	
	// Health check
	s.echo.GET("/health", s.healthHandler)
	
	// API v1 routes
	api := s.echo.Group("/api/v1")
	{
		// User routes
		users := api.Group("/users")
		users.GET("", s.getUsersHandler)
		users.POST("", s.createUserHandler)
		users.GET("/:id", s.getUserHandler)
		users.PUT("/:id", s.updateUserHandler)
		users.DELETE("/:id", s.deleteUserHandler)
		
		// Echo test routes
		echo := api.Group("/echo")
		echo.GET("/:message", s.echoHandler)
		echo.POST("", s.echoPostHandler)
	}

	// Static files
	s.echo.Static("/static", "static")
}

// Middleware functions

// structuredLoggingMiddleware adds structured logging
func (s *Server) structuredLoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			
			err := next(c)
			
			req := c.Request()
			res := c.Response()
			
			s.logger.WithFields(logrus.Fields{
				"method":      req.Method,
				"uri":         req.RequestURI,
				"status":      res.Status,
				"latency":     time.Since(start),
				"remote_ip":   c.RealIP(),
				"user_agent":  req.UserAgent(),
				"request_id":  c.Response().Header().Get(echo.HeaderXRequestID),
			}).Info("HTTP request processed")
			
			return err
		}
	}
}

// Handler functions

// homeHandler serves the home page
func (s *Server) homeHandler(c echo.Context) error {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Echo Framework Example</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .endpoint { background: #f8f9fa; padding: 15px; margin: 10px 0; border-radius: 5px; border-left: 4px solid #007bff; }
        .method { font-weight: bold; color: #007bff; display: inline-block; width: 60px; }
        .url { font-family: monospace; background: #e9ecef; padding: 2px 6px; border-radius: 3px; }
        h1 { color: #333; border-bottom: 3px solid #007bff; padding-bottom: 10px; }
        h2 { color: #495057; margin-top: 30px; }
        .features { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 15px; margin: 20px 0; }
        .feature { background: #e3f2fd; padding: 10px; border-radius: 5px; text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Echo Framework Example</h1>
        <p>A comprehensive Echo web server demonstrating modern Go web development patterns.</p>
        
        <div class="features">
            <div class="feature">🔒 Security Middleware</div>
            <div class="feature">📊 Structured Logging</div>
            <div class="feature">🛡️ Rate Limiting</div>
            <div class="feature">✅ Health Checks</div>
            <div class="feature">📁 Static Files</div>
            <div class="feature">🌐 CORS Support</div>
        </div>
        
        <h2>API Endpoints</h2>
        
        <div class="endpoint">
            <span class="method">GET</span> <span class="url">/health</span> - Health check endpoint
        </div>
        
        <div class="endpoint">
            <span class="method">GET</span> <span class="url">/api/v1/users</span> - Get all users
        </div>
        
        <div class="endpoint">
            <span class="method">POST</span> <span class="url">/api/v1/users</span> - Create new user
        </div>
        
        <div class="endpoint">
            <span class="method">GET</span> <span class="url">/api/v1/users/:id</span> - Get user by ID
        </div>
        
        <div class="endpoint">
            <span class="method">PUT</span> <span class="url">/api/v1/users/:id</span> - Update user
        </div>
        
        <div class="endpoint">
            <span class="method">DELETE</span> <span class="url">/api/v1/users/:id</span> - Delete user
        </div>
        
        <div class="endpoint">
            <span class="method">GET</span> <span class="url">/api/v1/echo/:message</span> - Echo message
        </div>
        
        <div class="endpoint">
            <span class="method">POST</span> <span class="url">/api/v1/echo</span> - Echo POST data
        </div>
        
        <h2>Quick Test</h2>
        <p>Try these links:</p>
        <ul>
            <li><a href="/health">Health Check</a></li>
            <li><a href="/api/v1/users">Get Users</a></li>
            <li><a href="/api/v1/echo/Hello%20World">Echo Hello World</a></li>
        </ul>
        
        <h2>Features Demonstrated</h2>
        <ul>
            <li>Echo framework with middleware stack</li>
            <li>Structured JSON logging with Logrus</li>
            <li>Rate limiting and security headers</li>
            <li>RESTful API design patterns</li>
            <li>Environment-based configuration</li>
            <li>Graceful server shutdown</li>
            <li>Request ID tracking</li>
            <li>CORS and security middleware</li>
        </ul>
    </div>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}

// healthHandler provides health check
func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
			"service":   "echo-example",
		},
		Message: "Service is running normally",
	})
}

// getUsersHandler returns all users
func (s *Server) getUsersHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
		Message: fmt.Sprintf("Retrieved %d users", len(users)),
	})
}

// getUserHandler returns a specific user
func (s *Server) getUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
		})
	}

	for _, user := range users {
		if user.ID == id {
			return c.JSON(http.StatusOK, Response{
				Success: true,
				Data:    user,
				Message: "User found",
			})
		}
	}

	return c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error:   "User not found",
	})
}

// createUserHandler creates a new user
func (s *Server) createUserHandler(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request payload",
		})
	}

	// Validate required fields
	if user.Name == "" || user.Email == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Name and email are required",
		})
	}

	// Assign ID and add user
	user.ID = nextUserID
	nextUserID++
	users = append(users, user)

	s.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
	}).Info("User created")

	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    user,
		Message: "User created successfully",
	})
}

// updateUserHandler updates an existing user
func (s *Server) updateUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
		})
	}

	var updatedUser User
	if err := c.Bind(&updatedUser); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request payload",
		})
	}

	for i, user := range users {
		if user.ID == id {
			updatedUser.ID = id
			users[i] = updatedUser
			
			return c.JSON(http.StatusOK, Response{
				Success: true,
				Data:    updatedUser,
				Message: "User updated successfully",
			})
		}
	}

	return c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error:   "User not found",
	})
}

// deleteUserHandler deletes a user
func (s *Server) deleteUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
		})
	}

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			
			s.logger.WithFields(logrus.Fields{
				"user_id": id,
			}).Info("User deleted")
			
			return c.JSON(http.StatusOK, Response{
				Success: true,
				Message: "User deleted successfully",
			})
		}
	}

	return c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error:   "User not found",
	})
}

// echoHandler echoes back the message from URL parameter
func (s *Server) echoHandler(c echo.Context) error {
	message := c.Param("message")
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"original_message": message,
			"echo_message":     fmt.Sprintf("Echo: %s", message),
			"timestamp":        time.Now().Unix(),
			"request_id":       c.Response().Header().Get(echo.HeaderXRequestID),
		},
		Message: "Message echoed successfully",
	})
}

// echoPostHandler echoes back POST data
func (s *Server) echoPostHandler(c echo.Context) error {
	var payload map[string]interface{}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid JSON payload",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"received_data": payload,
			"echo_response": "Data received and echoed back",
			"timestamp":     time.Now().Unix(),
			"request_id":    c.Response().Header().Get(echo.HeaderXRequestID),
		},
		Message: "POST data echoed successfully",
	})
}

// start starts the server with graceful shutdown
func (s *Server) start() {
	// Get port from environment
	port := getEnv("PORT", "1323")
	
	// Start server in a goroutine
	go func() {
		s.logger.WithFields(logrus.Fields{
			"port": port,
		}).Info("Starting Echo server")
		
		if err := s.echo.Start(":" + port); err != nil && err != http.ErrServerClosed {
			s.logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	
	<-quit
	s.logger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		s.logger.WithError(err).Fatal("Server forced to shutdown")
	}

	s.logger.Info("Server exited")
}

// Helper functions

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}