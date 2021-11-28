package main

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
)

func main() {
  e := echo.New()
  // add middleware and routes
  // ...
  s := &http2.Server{
    MaxConcurrentStreams: 250,
    MaxReadFrameSize:     1048576,
    IdleTimeout:          10 * time.Second,
  }
  if err := e.StartH2CServer(":8080", s); err != http.ErrServerClosed {
    log.Fatal(err)
  }
}
