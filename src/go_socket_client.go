package main
import {
  "fmt"
  "net"
}

const (
  SERVER_HOST = "localhost"
  SERVER_PORT = "54001"
  SERVER_TYPE = "tcp"
)

func main() {
  conn, err := net.Dial(SERVER_TYPE, SERVER_HOST + ":" + SERVER_PORT)
  if err != nil {
    panic(err)
  }

  buffer := make([]byte, 1024)

  mLen, err2 := conn.Read(buffer)
  if err2 != nil {
    fmt.Println("Error in reading: ", err2.Error())
  } else {
    fmt.Println("Recieved: ", string(buffer[:mLen]))
  }

  defer conn.Close()
}
 
