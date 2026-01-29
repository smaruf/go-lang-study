# SR-71 Simulator Examples

This directory contains example programs demonstrating various features of the SR-71 simulator.

## Examples

### 1. basic_flight.go
Basic flight simulation showing altitude and speed control.

### 2. world_tour_flight.go
World tour flight simulation with GPS navigation and environmental variations.
- Route: Florida → Atlantic → Moscow → Siberia → Pacific → California → Virginia
- Features GPS/INS navigation system switching
- Environmental variations (maritime, extreme cold, etc.)
- Aerial refueling simulation
- Real-time position tracking with coordinates

## Running Examples

```bash
# Run basic flight simulation
cd /home/runner/work/go-lang-study/go-lang-study/src/embedded-os/tiny/sr71sim
go run examples/basic_flight.go

# Run world tour flight simulation
go run examples/world_tour_flight.go
```
