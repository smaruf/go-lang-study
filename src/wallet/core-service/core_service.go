package main

import (
	"context"
	. "core-service/config"
	. "core-service/proto"
	"fmt"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	//"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	log.Printf("Starting Core-Service:....")
	e := echo.New()
	initManage(e)
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"status": "success"})
	})
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = HostConfig.Port
		}
		if err := e.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
	initWalletService(e)
	graceFullShutdown(e)
}

func initWalletService(e *echo.Echo) {
	go func() {
		fmt.Println("Server started...", e.Server)
		_ = e.Server.ListenAndServe()
		s := grpc.NewServer()
		fmt.Println("Service Server:  ", s)
		RegisterWalletServiceServer(s, NewWalletService())
		fmt.Println("Registered Service: ", e.Server)
		e.Logger.Fatal(s.Serve(e.Listener))
		fmt.Println("Listening  ", e.Listener)
	}()
}

func graceFullShutdown(e *echo.Echo) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func initManage(e *echo.Echo) {
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
}
