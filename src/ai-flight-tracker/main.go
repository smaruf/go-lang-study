package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/smaruf/go-lang-study/src/ai-flight-tracker/config"
	"github.com/smaruf/go-lang-study/src/ai-flight-tracker/server"
)

func main() {
	cfg := config.Load()
	srv := server.New(cfg)

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := srv.Stop(); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
