package main

import (
    "log"
    "golang.org/x/net/websocket"
)

// T represents a data structure for JSON messages.
type T struct {
    Msg   string `json:"msg"`
    Count int    `json:"count"`
}

// receiveJSON receives a JSON message of type T from the WebSocket connection.
func receiveJSON(ws *websocket.Conn) (T, error) {
    var data T
    err := websocket.JSON.Receive(ws, &data)
    if err != nil {
        return T{}, err
    }
    return data, nil
}

// sendJSON sends a JSON message of type T through the WebSocket connection.
func sendJSON(ws *websocket.Conn, data T) error {
    return websocket.JSON.Send(ws, data)
}

// receiveText receives a text frame from the WebSocket connection.
func receiveText(ws *websocket.Conn) (string, error) {
    var message string
    err := websocket.Message.Receive(ws, &message)
    if err != nil {
        return "", err
    }
    return message, nil
}

// sendText sends a text frame through the WebSocket connection.
func sendText(ws *websocket.Conn, message string) error {
    return websocket.Message.Send(ws, message)
}

// receiveBinary receives a binary frame from the WebSocket connection.
func receiveBinary(ws *websocket.Conn) ([]byte, error) {
    var data []byte
    err := websocket.Message.Receive(ws, &data)
    if err != nil {
        return nil, err
    }
    return data, nil
}

// sendBinary sends a binary frame through the WebSocket connection.
func sendBinary(ws *websocket.Conn, data []byte) error {
    return websocket.Message.Send(ws, data)
}

func main() {
    // Example usage
    ws, err := websocket.Dial("ws://example.com/socket", "", "http://example.com/")
    if err != nil {
        log.Fatal(err)
    }
    defer ws.Close()

    // Receiving and sending JSON
    data, err := receiveJSON(ws)
    if err != nil {
        log.Println("Error receiving JSON:", err)
    }
    err = sendJSON(ws, data)
    if err != nil {
        log.Println("Error sending JSON:", err)
    }

    // Receiving and sending text
    text, err := receiveText(ws)
    if err != nil {
        log.Println("Error receiving text:", err)
    }
    err = sendText(ws, text)
    if err != nil {
        log.Println("Error sending text:", err)
    }

    // Receiving and sending binary
    binaryData, err := receiveBinary(ws)
    if err != nil {
        log.Println("Error receiving binary:", err)
    }
    err = sendBinary(ws, binaryData)
    if err != nil {
        log.Println("Error sending binary:", err)
    }
}
