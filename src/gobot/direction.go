package main

  import (
      "time"

      "gobot.io/x/gobot"
      "gobot.io/x/gobot/platforms/dji/tello"
  )

  func main() {
      drone := tello.NewDriver("8888")

      work := func() {
          drone.TakeOff()

          gobot.After(3*time.Second, func() {
              drone.Forward(10)
          })

          gobot.After(6*time.Second, func() {
              drone.Backward(10)
          })

          gobot.After(9*time.Second, func() {
              drone.Land()
          })
      }

      robot := gobot.NewRobot("tello",
          []gobot.Connection{},
          []gobot.Device{drone},
          work,
      )

      robot.Start()
  }
