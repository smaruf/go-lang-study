func (s *broadcastServer) serve(ctx context.Context) {
  defer func () {
    for _, listener := range s.listeners {
      if listener != nil {
          close(listener)
      }
    }
  } ()

  for {
    select {
      case <-ctx.Done():
        return
      case newListener := <- s.addListener:
        s.listeners = append(s.listeners, newListener)
      case listenerToRemove := <- s.removeListener:
        for i, ch := range s.listeners {
          if ch == listenerToRemove {
              s.listeners[i] = s.listeners[len(s.listeners)-1]
              s.listeners = s.listeners[:len(s.listeners)-1]
              close(ch)
              break
          }
        }
      case val, ok := <-s.source:
        if !ok {
          return
        }
        for _, listener := range s.listeners {
          if listener != nil {
            select {
             case listener <- val:
             case <-ctx.Done():
              return
            }
            
          }
        }
    }
  }
}
