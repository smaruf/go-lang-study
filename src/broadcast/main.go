func main() {
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()
// Generates a channel sending integers
// From 0 to 9
  range10 := rangeChannel(ctx, 10)

  broadcaster := NewBroadcastServer(ctx, range10)
  listener1 := broadcaster.Subscribe()
  listener2 := broadcaster.Subscribe()
  listener3 := broadcaster.Subscribe()

  var wg sync.WaitGroup
  wg.Add(3)
  go func() {
    defer wg.Done()
    for i := range listener1 {
        fmt.Printf("Listener 1: %v/10 \n", i+1)
    }
  }()
  go func() {
    defer wg.Done()
    for i := range listener2 {
        fmt.Printf("Listener 2: %v/10 \n", i+1)
    }
  }()
  go func() {
    defer wg.Done()
    for i := range listener3 {
        fmt.Printf("Listener 3: %v/10 \n", i+1)
    }
  }()
  wg.Wait()
}
