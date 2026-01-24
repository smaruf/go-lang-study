package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/smaruf/go-lang-study/src/3-weeks-plan/internal/config"
	"github.com/smaruf/go-lang-study/src/3-weeks-plan/pkg/database"
	"github.com/smaruf/go-lang-study/src/3-weeks-plan/pkg/handlers"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize handlers
	h := handlers.New(db)

	// Setup router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/tasks", h.CreateTask).Methods("POST")
	api.HandleFunc("/tasks", h.ListTasks).Methods("GET")
	api.HandleFunc("/tasks/{id}", h.GetTask).Methods("GET")
	api.HandleFunc("/tasks/{id}", h.UpdateTask).Methods("PUT")
	api.HandleFunc("/tasks/{id}", h.DeleteTask).Methods("DELETE")
	api.HandleFunc("/stats", h.GetStats).Methods("GET")

	// Health check
	router.HandleFunc("/health", h.HealthCheck).Methods("GET")

	// Add logging middleware
	router.Use(loggingMiddleware)

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server
	addr := cfg.Address()
	log.Printf("Starting server on http://%s", addr)
	log.Printf("Environment: %s", cfg.Env)
	log.Println("API endpoints:")
	log.Println("  POST   /api/tasks       - Create a task")
	log.Println("  GET    /api/tasks       - List all tasks")
	log.Println("  GET    /api/tasks/{id}  - Get a task")
	log.Println("  PUT    /api/tasks/{id}  - Update a task")
	log.Println("  DELETE /api/tasks/{id}  - Delete a task")
	log.Println("  GET    /api/stats       - Get task statistics")
	log.Println("  GET    /health          - Health check")

	go func() {
		if err := http.ListenAndServe(addr, router); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop
	log.Println("\nShutting down server...")
}

// loggingMiddleware logs incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
