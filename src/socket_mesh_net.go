// Package main provides a basic demonstration of a simple peer-to-peer
// network setup in Go that can serve as a starting point for building
// more complex mesh network applications. Each instance of the program
// can function as both a client and a server.
//
// Building and Running the Program:
//
// To build the program:
//   go build main.go
//
// This will compile the code into an executable named 'main' (or 'main.exe' on Windows).
//
// To run the executable on Linux or similar OS:
//   ./main
//
// Upon execution, each instance functions as a server listening on port 8080,
// and prompts the user to connect to a peer by entering an IP address. You will need
// to enter the address in the form 'ip:port' (e.g., '192.168.1.5:8080') to connect
// nodes in a peer-to-peer manner.
//
// The program can be terminated using standard interrupt signals (e.g., Ctrl+C).
package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func main() {
    // Start the server in a new goroutine to manage incoming connections concurrently.
    go startServer()
    // Allow user to connect to peers by specifying an IP address.
    connectToPeer()

    // Prevent the main function from terminating prematurely.
    select {}
}

// startServer initializes a TCP server that listens on port 8080 and handles incoming
// connections. Each connection is managed in a separate goroutine to handle multiple
// clients simultaneously.
func startServer() {
    // Listen on all interfaces at port 8080.
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    defer listener.Close()
    fmt.Println("Server started on :8080")

    // Continuously accept and handle incoming connections.
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        // Handle the connection in a separate goroutine.
        go handleConnection(conn)
    }
}

// connectToPeer allows the user to establish a connection to a peer by entering
// the peer's address. It manages sending user-entered text to the connected peer.
func connectToPeer() {
    fmt.Print("Enter peer address (or 'none' to skip): ")
    reader := bufio.NewReader(os.Stdin)
    address, _ := reader.ReadString('\n')
    address = strings.TrimSpace(address)

    if address != "none" {
        // Establish a connection to the peer.
        conn, err := net.Dial("tcp", address)
        if err != nil {
            fmt.Println("Error connecting to peer:", err)
            return
        }
        defer conn.Close()
        fmt.Println("Connected to peer", address)

        // Send text entered by the user to the connected peer.
        for {
            text, _ := reader.ReadString('\n')
            fmt.Fprint(conn, text)
        }
    }
}

// handleConnection reads data from the connected peer and prints it to the standard output.
// It is run in a separate goroutine for each connection to handle concurrent clients.
func handleConnection(conn net.Conn) {
    defer conn.Close()
    fmt.Println("Connected from", conn.RemoteAddr().String())

    // Read and display messages from the connected peer.
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        message := scanner.Text()
        fmt.Println("Received:", message)
    }
}
