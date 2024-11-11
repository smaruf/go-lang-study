package main

import (
    "errors"
    "fmt"
)

func mayPanic() {
    panic("a severe problem occurred")
}

func handlePanic() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from:", r)
        }
    }()
    mayPanic()
    fmt.Println("handlePanic exited normally")
}

func someOperation() error {
    return errors.New("an expected error occurred")
}

func main() {
    handlePanic()
    
    if err := someOperation(); err != nil {
        fmt.Println("Error in someOperation:", err)
    }
    
    fmt.Println("main completed")
}
