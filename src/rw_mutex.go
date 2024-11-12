package main

import (
    "fmt"
    "sync"
    "time"
)

var (
    // A map protected by an RWMutex
    sharedMap = make(map[string]string)
    rwMutex   sync.RWMutex
)

// Writes to the map
func write(key, value string) {
    rwMutex.Lock()  // Lock for writing
    sharedMap[key] = value
    rwMutex.Unlock()  // Unlock after writing
}

// Reads from the map
func read(key string) string {
    rwMutex.RLock()  // Lock for reading
    value := sharedMap[key]
    rwMutex.RUnlock()  // Unlock after reading
    return value
}

func main() {
    // Simulate concurrent reading and writing
    go write("foo", "bar")
    time.Sleep(1 * time.Second)  // Ensure the write happens before the read

    result := read("foo")
    fmt.Println("Read value:", result)

    go write("foo", "baz")
    time.Sleep(1 * time.Second)  // Allow write to complete before exiting
}
