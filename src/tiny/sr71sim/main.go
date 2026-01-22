package main

import (
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    
    "github.com/smaruf/go-lang-study/src/tiny/sr71sim/avionics"
    "github.com/smaruf/go-lang-study/src/tiny/sr71sim/engine"
)

type SimulationData struct {
    StartTime      time.Time                `json:"start_time"`
    EndTime        time.Time                `json:"end_time"`
    UniqueTestID   string                   `json:"unique_test_id"`
    Altitude       float64                  `json:"altitude"`
    Speed          float64                  `json:"speed"`
    AvionicsStates []avionics.AvionicsState `json:"avionics_states"`
    EngineStates   []engine.EngineState     `json:"engine_states"`
}

type LineGraphData struct {
    Timestamps []time.Time `json:"timestamps"`
    Values     []float64   `json:"values"`
}

var (
    startTime    time.Time
    uniqueTestID string
    logger       *log.Logger
    simDuration  time.Duration
)

const storageURL = "http://example.com/storeSimulationData"

func init() {
    logger = log.New(os.Stdout, "SR-71 Simulation: ", log.LstdFlags)
    
    // Parse command-line flags
    duration := flag.Int("duration", 5, "Simulation duration in seconds")
    flag.Parse()
    simDuration = time.Duration(*duration) * time.Second
}

func generateUniqueTestID() string {
    return fmt.Sprintf("test-%d", time.Now().UnixNano())
}

func StartSimulation() {
    startTime = time.Now()
    uniqueTestID = generateUniqueTestID()
    logger.Println("Starting SR-71 Simulation...")
    // Add simulation steps here
}

func CloseSimulation() {
    endTime := time.Now()
    logger.Println("Closing SR-71 Simulation...")

    avionicsData, engineData := fetchData()

    data := SimulationData{
        StartTime:      startTime,
        EndTime:        endTime,
        UniqueTestID:   uniqueTestID,
        Altitude:       10000, // Example data
        Speed:          300,   // Example data
        AvionicsStates: avionicsData,
        EngineStates:   engineData,
    }

    if err := storeSimulationData(data); err != nil {
        logger.Println("Error storing simulation data:", err)
    } else {
        logger.Println("Simulation data stored successfully")
    }
}

func fetchData() ([]avionics.AvionicsState, []engine.EngineState) {
    avionicsDataCh := make(chan []avionics.AvionicsState)
    engineDataCh := make(chan []engine.EngineState)

    go func() {
        avionicsData, err := avionics.Test()
        if err != nil {
            logger.Println("Error getting avionics data:", err)
            close(avionicsDataCh)
            return
        }
        avionicsDataCh <- avionicsData
    }()

    go func() {
        engineData, err := engine.Test()
        if err != nil {
            logger.Println("Error getting engine data:", err)
            close(engineDataCh)
            return
        }
        engineDataCh <- engineData
    }()

    return <-avionicsDataCh, <-engineDataCh
}

func storeSimulationData(data SimulationData) error {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("error marshalling JSON: %w", err)
    }

    resp, err := http.Post(storageURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("error posting JSON data: %w", err)
    }
    defer resp.Body.Close()

    return nil
}

func PlotSimulation() {
    logger.Println("Plotting Simulation Data...")

    graphData := LineGraphData{
        Timestamps: []time.Time{time.Now().Add(-5 * time.Minute), time.Now()},
        Values:     []float64{10000, 15000},
    }

    jsonData, err := json.Marshal(graphData)
    if err != nil {
        logger.Println("Error marshalling JSON:", err)
        return
    }

    // Print or return the JSON data
    logger.Println(string(jsonData))
}

func main() {
    logger.Println("Start SR-71 Simulation")
    logger.Printf("Simulation will run for %v", simDuration)
    StartSimulation()
    
    // Run simulation for configured duration
    time.Sleep(simDuration)
    
    CloseSimulation()
    PlotSimulation()
    logger.Println("Ending SR-71 Simulation")
}
