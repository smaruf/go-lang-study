package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Server represents our HTTP server with dependencies
type Server struct {
	router *mux.Router
	logger *logrus.Logger
	config *Config
}

// Config holds application configuration
type Config struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	LogLevel     string
	Environment  string
}

// User represents a user in our system
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// In-memory storage for demo purposes
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
}
var nextUserID = 3

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	// Create server instance
	server := NewServer()

	// Setup routes
	server.setupRoutes()

	// Start server with graceful shutdown
	server.start()
}

// NewServer creates a new server instance
func NewServer() *Server {
	config := &Config{
		Port:         getEnv("PORT", "8080"),
		Host:         getEnv("HOST", "localhost"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		Environment:  getEnv("ENV", "development"),
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	return &Server{
		router: mux.NewRouter(),
		logger: logger,
		config: config,
	}
}

// setupRoutes configures all the routes for the server
func (s *Server) setupRoutes() {
	// Apply middleware to all routes
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.corsMiddleware)
	s.router.Use(s.contentTypeMiddleware)

	// API routes
	api := s.router.PathPrefix("/api/v1").Subrouter()
	
	// Health check
	api.HandleFunc("/health", s.healthHandler).Methods("GET")
	
	// User endpoints
	api.HandleFunc("/users", s.getUsersHandler).Methods("GET")
	api.HandleFunc("/users", s.createUserHandler).Methods("POST")
	api.HandleFunc("/users/{id:[0-9]+}", s.getUserHandler).Methods("GET")
	api.HandleFunc("/users/{id:[0-9]+}", s.updateUserHandler).Methods("PUT")
	api.HandleFunc("/users/{id:[0-9]+}", s.deleteUserHandler).Methods("DELETE")

	// Static routes
	s.router.HandleFunc("/", s.indexHandler).Methods("GET")
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
}

// Middleware functions

// loggingMiddleware logs all incoming requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a response writer wrapper to capture status code
		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapper, r)
		
		duration := time.Since(start)
		s.logger.WithFields(logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"status_code": wrapper.statusCode,
			"duration":    duration,
			"user_agent":  r.UserAgent(),
		}).Info("HTTP request")
	})
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// contentTypeMiddleware sets default content type
func (s *Server) contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// Handler functions

// healthHandler provides health check endpoint
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
		},
	}
	s.sendJSONResponse(w, http.StatusOK, response)
}

// indexHandler serves the main page
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Go Web Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { background: #f4f4f4; padding: 10px; margin: 10px 0; border-radius: 5px; }
        .method { font-weight: bold; color: #2196F3; }
    </style>
</head>
<body>
    <h1>Go Web Server with Modern Patterns</h1>
    <p>This server demonstrates modern Go web development patterns including middleware, structured logging, and RESTful APIs.</p>
    
    <h2>Available Endpoints:</h2>
    <div class="endpoint">
        <span class="method">GET</span> /api/v1/health - Health check
    </div>
    <div class="endpoint">
        <span class="method">GET</span> /api/v1/users - Get all users
    </div>
    <div class="endpoint">
        <span class="method">POST</span> /api/v1/users - Create a new user
    </div>
    <div class="endpoint">
        <span class="method">GET</span> /api/v1/users/{id} - Get user by ID
    </div>
    <div class="endpoint">
        <span class="method">PUT</span> /api/v1/users/{id} - Update user by ID
    </div>
    <div class="endpoint">
        <span class="method">DELETE</span> /api/v1/users/{id} - Delete user by ID
    </div>
    
    <h3>Try it out:</h3>
    <p><a href="/api/v1/health">Health Check</a></p>
    <p><a href="/api/v1/users">Get Users</a></p>
</body>
</html>`
	fmt.Fprint(w, html)
}

// getUsersHandler returns all users
func (s *Server) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Success: true,
		Data:    users,
		Message: fmt.Sprintf("Found %d users", len(users)),
	}
	s.sendJSONResponse(w, http.StatusOK, response)
}

// getUserHandler returns a specific user
func (s *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.sendErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	for _, user := range users {
		if user.ID == id {
			response := APIResponse{
				Success: true,
				Data:    user,
			}
			s.sendJSONResponse(w, http.StatusOK, response)
			return
		}
	}

	s.sendErrorResponse(w, http.StatusNotFound, "User not found")
}

// createUserHandler creates a new user
func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Validate input
	if user.Name == "" || user.Email == "" {
		s.sendErrorResponse(w, http.StatusBadRequest, "Name and email are required")
		return
	}

	// Assign ID and add to users
	user.ID = nextUserID
	nextUserID++
	users = append(users, user)

	s.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
	}).Info("User created")

	response := APIResponse{
		Success: true,
		Data:    user,
		Message: "User created successfully",
	}
	s.sendJSONResponse(w, http.StatusCreated, response)
}

// updateUserHandler updates an existing user
func (s *Server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.sendErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		s.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	for i, user := range users {
		if user.ID == id {
			updatedUser.ID = id
			users[i] = updatedUser
			
			response := APIResponse{
				Success: true,
				Data:    updatedUser,
				Message: "User updated successfully",
			}
			s.sendJSONResponse(w, http.StatusOK, response)
			return
		}
	}

	s.sendErrorResponse(w, http.StatusNotFound, "User not found")
}

// deleteUserHandler deletes a user
func (s *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.sendErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			
			s.logger.WithFields(logrus.Fields{
				"user_id": id,
			}).Info("User deleted")
			
			response := APIResponse{
				Success: true,
				Message: "User deleted successfully",
			}
			s.sendJSONResponse(w, http.StatusOK, response)
			return
		}
	}

	s.sendErrorResponse(w, http.StatusNotFound, "User not found")
}

// Helper functions

// sendJSONResponse sends a JSON response
func (s *Server) sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.WithError(err).Error("Failed to encode JSON response")
	}
}

// sendErrorResponse sends an error response
func (s *Server) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := APIResponse{
		Success: false,
		Error:   message,
	}
	s.sendJSONResponse(w, statusCode, response)
}

// start starts the HTTP server with graceful shutdown
func (s *Server) start() {
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.config.Host, s.config.Port),
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	// Channel to listen for interrupt signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		s.logger.WithFields(logrus.Fields{
			"host": s.config.Host,
			"port": s.config.Port,
			"env":  s.config.Environment,
		}).Info("Starting HTTP server")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	s.logger.Info("Server is ready to handle requests")

	// Block until signal is received
	<-done
	s.logger.Info("Server is shutting down...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := server.Shutdown(ctx); err != nil {
		s.logger.WithError(err).Error("Server forced to shutdown")
	} else {
		s.logger.Info("Server exited gracefully")
	}
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}