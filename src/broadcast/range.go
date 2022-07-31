func rangeChannel(
  ctx context.Context,
  n int,
) <-chan int {
  valueStream := make(chan int)
  go func() {
      defer close(valueStream)
      for i := 0; i<n; i++ {
          select {
          case <-ctx.Done():
              return
          case valueStream <- i:
          }
      }
  }()
  return valueStream
}
