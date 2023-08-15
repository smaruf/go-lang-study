import "websocket"

type T struct {
	Msg string
	Count int
}

// receive JSON type T
var data T
websocket.JSON.Receive(ws, &data)

// send JSON type T
websocket.JSON.Send(ws, data)

var Message = Codec{marshal, unmarshal}

// receive text frame
var message string
websocket.Message.Receive(ws, &message)

// send text frame
message = "hello"
websocket.Message.Send(ws, message)

// receive binary frame
var data []byte
websocket.Message.Receive(ws, &data)

// send binary frame
data = []byte{0, 1, 2}
websocket.Message.Send(ws, data)
