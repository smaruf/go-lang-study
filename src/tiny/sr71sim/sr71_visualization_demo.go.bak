package main

import (
	"fmt"
	"math"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func simulateFlight() ([]float64, []float64, []float64, []string) {
	times := linspace(0, 100, 1000)
	altitudes := make([]float64, len(times))
	speeds := make([]float64, len(times))
	engineModes := make([]string, len(times))

	for i, t := range times {
		altitudes[i] = 25000 + 10000*math.Sin(0.04*t)
		speeds[i] = 1800 + 800*math.Cos(0.05*t)
		if speeds[i] < 2200 {
			engineModes[i] = "Turbojet"
		} else {
			engineModes[i] = "Ramjet"
		}
	}

	return times, altitudes, speeds, engineModes
}

func linspace(start, end float64, num int) []float64 {
	step := (end - start) / float64(num-1)
	result := make([]float64, num)
	for i := range result {
		result[i] = start + step*float64(i)
	}
	return result
}

func main() {
	times, altitudes, speeds, engineModes := simulateFlight()

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "SR-71 Blackbird Flight Simulation",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Time (s)",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Altitude (ft)",
		}),
	)

	line.SetXAxis(times).
		AddSeries("Altitude", generateLineItems(times, altitudes)).
		AddSeries("Speed", generateLineItems(times, speeds))

	line.Render()
}

func generateLineItems(times, values []float64) []opts.LineData {
	items := make([]opts.LineData, len(times))
	for i := range items {
		items[i] = opts.LineData{Value: values[i]}
	}
	return items
}
