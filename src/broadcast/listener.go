type BroadcastServer interface {
  Subscribe() <-chan int
  CancelSubscription(<-chan int)
}

type broadcastServer struct {
  source <-chan int
  listeners []chan int
  addListener chan chan int
  removeListener chan (<-chan int)
}

func (s *broadcastServer) Subscribe() <-chan int {
  newListener := make(chan int)
  s.addListener <- newListener
  return newListener
}

func (s *broadcastServer) CancelSubscription(channel <-chan int) {
  s.removeListener <- channel
}

func NewBroadcastServer(ctx context.Context, source <-chan int) BroadcastServer {
  service := &broadcastServer{
    source: source,
    listeners: make([]chan int, 0),
    addListener: make(chan chan int),
    removeListener: make(chan (<-chan int)),
  }
  go service.serve(ctx)
  return service
}
