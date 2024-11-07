package main

import (
    "fmt"
    "time"
)

func pinger(pings chan<- string) {
    for i := 0; ; i++ {
        pings <- "ping"
    }
}

func ponger(pings <-chan string, pongs chan<- string) {
    for {
        msg := <-pings
        pongs <- msg
    }
}

func printer(pongs <-chan string) {
    for {
        msg := <-pongs
        fmt.Println(msg)
        time.Sleep(time.Second * 1)
    }
}

func main() {
    pings := make(chan string)
    pongs := make(chan string)
    
    go pinger(pings)
    go ponger(pings, pongs)
    go printer(pongs)

    var input string
    fmt.Scanln(&input)
