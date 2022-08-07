package main

import "fmt"

func main() {
  fmt.Println("Start SR-71 Simulaton")
  ///steps
  go StartSimultaion()
  defer CloseSimulation()
  StoreSimulationData()
  PlotSimulation()
  // this may be external controller or base controller
  fmt.Println("Enging SR-71 Simulation")
}  
