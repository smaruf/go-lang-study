package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"

    "gobot.io/x/gobot"
    "gobot.io/x/gobot/drivers/gpio"
    "gobot.io/x/gobot/platforms/raspi"
)

func main() {
    // Initialize Gobot with a Raspberry Pi Adaptor
    adaptor := raspi.NewAdaptor()
    led := gpio.NewLedDriver(adaptor, "7") // Assuming an LED is connected to pin 7

    work := func() {
        gobot.Every(gobot.Second, func() {
            led.Toggle() // Toggle LED every second
        })
    }

    robot := gobot.NewRobot("bot",
        []gobot.Connection{adaptor},
        []gobot.Device{led},
        work,
    )

    go robot.Start() // Start Gobot routine in a goroutine
    go startServer() // Start TCP server
    connectToPeer()  // Act as a TCP client

    select {} // Prevent main from exiting
}

func startServer() {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    defer listener.Close()
    fmt.Println("Server started on :8080")

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(conn)
    }
}

func connectToPeer() {
    fmt.Print("Enter peer address (or 'none' to skip): ")
    reader := bufio.NewReader(os.Stdin)
    address, _ := reader.ReadString('\n')
    address = strings.TrimSpace(address)

    if address != "none" {
        conn, err := net.Dial("tcp", address)
        if err != nil {
            fmt.Println("Error connecting to peer:", err)
            return
        }
        defer conn.Close()
        fmt.Println("Connected to peer", address)

        for {
            text, _ := reader.ReadString('\n')
            fmt.Fprint(conn, text)
        }
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    fmt.Println("Connected from", conn.RemoteAddr().String())
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        message := scanner.Text()
        fmt.Println("Received:", message)
    }
}
