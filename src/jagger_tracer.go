package main
import (
    "github.com/labstack/echo-contrib/jaegertracing"
    "github.com/labstack/echo/v4"
)
func main() {
    e := echo.New()
    // Enable tracing middleware
    c := jaegertracing.New(e, nil)
    defer c.Close()
    e.GET("/", func(c echo.Context) error {
        // Do something before creating the child span
        time.Sleep(40 * time.Millisecond)
        sp := jaegertracing.CreateChildSpan(c, "Child span for additional processing")
        defer sp.Finish()
        sp.LogEvent("Test log")
        sp.SetBaggageItem("Test baggage", "baggage")
        sp.SetTag("Test tag", "New Tag")
        time.Sleep(100 * time.Millisecond)
        return c.String(http.StatusOK, "Hello, World!")
    })
    e.Logger.Fatal(e.Start(":1323"))
}
