////init: go get -d -u gobot.io/x/gobot/...
package main

import (
        "os"
        "time"

        "gobot.io/x/gobot"
        "gobot.io/x/gobot/platforms/ble"
        "gobot.io/x/gobot/platforms/sphero/ollie"
)

func main() {
        bleAdaptor := ble.NewClientAdaptor(os.Args[1])
        ollie := ollie.NewDriver(bleAdaptor)

        work := func() {
                gobot.Every(1*time.Second, func() {
                        r := uint8(gobot.Rand(255))
                        g := uint8(gobot.Rand(255))
                        b := uint8(gobot.Rand(255))
                        ollie.SetRGB(r, g, b)
                })
        }

        robot := gobot.NewRobot("ollieBot",
                []gobot.Connection{bleAdaptor},
                []gobot.Device{ollie},
                work,
        )

        robot.Start()
}

//// deploy: GODEBUG=cgocheck=0 go run examples/sphero_ollie.go 2B-1247
//// on ubuntu: go build examples/sphero_ollie.go
/// --: sudo ./minidrone 2B-1247
