# SR-71 Simulator - Architecture Documentation

## System Architecture Overview

```
┌────────────────────────────────────────────────────────────────┐
│                        Client Layer                            │
├────────────┬────────────┬────────────┬────────────┬───────────┤
│  Web App   │ Mobile App │  Console   │  REST API  │  Hardware │
│  Browser   │ iOS/Android│    CLI     │   Clients  │  RPI/Ardu │
└─────┬──────┴─────┬──────┴─────┬──────┴─────┬──────┴─────┬─────┘
      │            │            │            │            │
      └────────────┴────────────┴────────────┴────────────┘
                              │
      ┌───────────────────────▼──────────────────────────┐
      │            API Gateway / Router                  │
      │  • HTTP/HTTPS Endpoints                          │
      │  • WebSocket Support                             │
      │  • Authentication/Authorization                  │
      │  • Rate Limiting                                 │
      └───────────────────────┬──────────────────────────┘
                              │
      ┌───────────────────────▼──────────────────────────┐
      │         Application Layer (Go)                   │
      ├──────────────────────────────────────────────────┤
      │  Simulation Controller                           │
      │  • Mission Management                            │
      │  • State Management                              │
      │  • Event Handling                                │
      └───────────────────────┬──────────────────────────┘
                              │
      ┌───────────────────────▼──────────────────────────┐
      │         Domain Layer (Core Modules)              │
      ├──────────┬──────────┬──────────┬────────┬───────┤
      │ Avionics │  Engine  │  Flying  │ Fueling│Landing│
      └──────────┴──────────┴──────────┴────────┴───────┘
                              │
      ┌───────────────────────▼──────────────────────────┐
      │         Data Layer                               │
      ├──────────┬──────────┬──────────┬────────────────┤
      │PostgreSQL│  Redis   │ InfluxDB │  File Storage  │
      │(Metadata)│ (Cache)  │(Telemetry)│ (Logs/Exports)│
      └──────────┴──────────┴──────────┴────────────────┘
```

## Module Structure

### 1. Core Modules (Domain Layer)

#### Avionics Module
**Purpose:** Manage aircraft systems and instruments

**Components:**
- Navigation System (GPS, INS)
- Communication System
- Autopilot
- Environmental Control (Cabin Pressure)
- G-Force Recovery
- Fuel Leaching Detection

**Data Structure:**
```go
type AvionicsState struct {
    Altitude      float64
    Speed         float64
    Heading       float64
    GPSActive     bool
    INSActive     bool
    AutopilotOn   bool
    CabinPressure float64
    GForce        float64
    FuelLeaching  bool
}
```

#### Engine Module
**Purpose:** Simulate Pratt & Whitney J58 engine

**Components:**
- Velocity Control
- Air Intake Management
- Combustion Chamber Simulation
- Exhaust Pattern Modeling
- Mode Switching (Turbojet/Ramjet/Scramjet)

**Engine Modes:**
- **Turbojet:** Mach 0 - 2.0 (subsonic to supersonic)
- **Ramjet:** Mach 2.0 - 5.0 (supersonic to hypersonic)
- **Scramjet:** Mach 5.0+ (hypersonic, theoretical)

**Data Structure:**
```go
type EngineState struct {
    Mode            string  // "Turbojet", "Ramjet", "Scramjet"
    Velocity        float64 // mph
    Mach            float64
    AirIntake       float64 // percentage
    CombustionTemp  float64 // Fahrenheit
    Thrust          float64 // pounds
    FuelFlow        float64 // gallons/hour
}
```

#### Flying Module
**Purpose:** Flight control and dynamics

**Components:**
- Altitude Management
- Velocity Control
- Mission Types (Reconnaissance, High-Speed, Stealth)
- Flight Maneuvers

**Data Structure:**
```go
type SR71 struct {
    CurrentAltitude float64
    CurrentSpeed    float64
    TargetAltitude  float64
    TargetSpeed     float64
    Pitch           float64
    Roll            float64
    Yaw             float64
}

type Mission struct {
    Type        string // "reconnaissance", "high_speed", "stealth"
    Duration    time.Duration
    Waypoints   []Waypoint
    Objectives  []Objective
}
```

#### Fueling Module
**Purpose:** Fuel management and refueling

**Components:**
- Fuel Level Monitoring
- Consumption Rate Calculation
- Air Refueling Logic
- Fuel System Leaching Detection

**Characteristics:**
- JP-7 fuel usage
- Fuel as hydraulic coolant
- Pre-flight fuel leaching (expected behavior)
- Consumption: 5,600+ gal/hour at high speeds

#### Landing/Takeoff Modules
**Purpose:** Ground operations

**Components:**
- Runway approach
- Landing gear management
- Taxi operations
- Emergency procedures

## Data Flow Architecture

### Simulation Loop
```
┌─────────────────────────────────────────────────────────┐
│  1. Initialize Simulation                               │
│     • Load configuration                                │
│     • Initialize all modules                            │
│     • Set initial state                                 │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│  2. Main Loop (60 Hz default)                           │
│     ┌─────────────────────────────────────────┐         │
│     │ a. Read Inputs (controls, commands)     │         │
│     └──────────────┬──────────────────────────┘         │
│     ┌──────────────▼──────────────────────────┐         │
│     │ b. Update State                         │         │
│     │    • Engine calculations                │         │
│     │    • Avionics updates                   │         │
│     │    • Flight dynamics                    │         │
│     │    • Fuel consumption                   │         │
│     └──────────────┬──────────────────────────┘         │
│     ┌──────────────▼──────────────────────────┐         │
│     │ c. Check Conditions                     │         │
│     │    • Engine mode transitions            │         │
│     │    • Alerts and warnings                │         │
│     │    • Mission objectives                 │         │
│     └──────────────┬──────────────────────────┘         │
│     ┌──────────────▼──────────────────────────┐         │
│     │ d. Output Data                          │         │
│     │    • Telemetry to database              │         │
│     │    • WebSocket to clients               │         │
│     │    • Hardware outputs (GPIO)            │         │
│     └──────────────┬──────────────────────────┘         │
│                    │                                     │
│     Loop until stop signal                              │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│  3. Cleanup                                             │
│     • Save final state                                  │
│     • Export data                                       │
│     • Close connections                                 │
└─────────────────────────────────────────────────────────┘
```

## API Design

### REST Endpoints

#### Simulation Management
```
POST   /api/v1/simulation
GET    /api/v1/simulation/:id
DELETE /api/v1/simulation/:id
POST   /api/v1/simulation/:id/start
POST   /api/v1/simulation/:id/stop
POST   /api/v1/simulation/:id/pause
POST   /api/v1/simulation/:id/resume
```

#### Control Endpoints
```
POST   /api/v1/control/throttle
POST   /api/v1/control/altitude
POST   /api/v1/control/heading
POST   /api/v1/control/autopilot
```

#### Data Endpoints
```
GET    /api/v1/telemetry/current
GET    /api/v1/telemetry/history?from=&to=
GET    /api/v1/engine/state
GET    /api/v1/avionics/state
GET    /api/v1/flight/state
```

#### Mission Endpoints
```
GET    /api/v1/missions
GET    /api/v1/missions/:id
POST   /api/v1/missions
PUT    /api/v1/missions/:id
DELETE /api/v1/missions/:id
```

### WebSocket Protocol

**Connection:** `ws://host:port/ws/telemetry`

**Message Format:**
```json
{
  "type": "telemetry|control|event",
  "timestamp": "2026-01-22T05:36:31Z",
  "data": {
    // Type-specific data
  }
}
```

**Telemetry Message:**
```json
{
  "type": "telemetry",
  "timestamp": "2026-01-22T05:36:31Z",
  "data": {
    "altitude": 78450,
    "speed": 2156,
    "mach": 3.21,
    "engine_mode": "ramjet",
    "fuel_level": 85.3,
    "cabin_pressure": 14.7
  }
}
```

**Control Message:**
```json
{
  "type": "control",
  "timestamp": "2026-01-22T05:36:31Z",
  "data": {
    "command": "set_throttle",
    "value": 95
  }
}
```

**Event Message:**
```json
{
  "type": "event",
  "timestamp": "2026-01-22T05:36:31Z",
  "data": {
    "event": "engine_mode_switch",
    "from": "turbojet",
    "to": "ramjet",
    "mach": 2.05
  }
}
```

## Hardware Integration Architecture

### Raspberry Pi Setup

```
┌──────────────────────────────────────────────────────────┐
│                   Raspberry Pi 4                         │
├──────────────────────────────────────────────────────────┤
│  SR-71 Simulator Application (Go Binary)                 │
│                                                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │
│  │ GPIO Driver │  │ I2C Driver  │  │ SPI Driver  │      │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘      │
└─────────┼─────────────────┼─────────────────┼────────────┘
          │                 │                 │
    ┌─────▼─────┐     ┌────▼─────┐     ┌────▼─────┐
    │   GPIO    │     │   I2C    │     │   SPI    │
    │   Pins    │     │  Devices │     │ Devices  │
    └─────┬─────┘     └────┬─────┘     └────┬─────┘
          │                │                 │
    ┌─────▼──────────┬─────▼─────────┬──────▼──────┐
    │ LEDs/Buttons   │ LCD Display   │ Arduino/    │
    │ Buzzer         │ Sensors       │ Other MCUs  │
    └────────────────┴───────────────┴─────────────┘
```

### Arduino Communication

```
┌──────────────────┐         Serial/SPI         ┌──────────────────┐
│   Main System    │◄──────────────────────────►│   Arduino        │
│   (Go App)       │    Custom Protocol         │   (C++ Firmware) │
└──────────────────┘                            └──────────────────┘
                                                          │
                                                ┌─────────▼─────────┐
                                                │ Hardware I/O      │
                                                ├───────────────────┤
                                                │ • LCD Display     │
                                                │ • LED Indicators  │
                                                │ • Servo Gauges    │
                                                │ • Input Controls  │
                                                │ • Analog Sensors  │
                                                └───────────────────┘
```

## Deployment Architectures

### Cloud Deployment (AWS/GCP/Azure)

```
                        ┌─────────────────┐
                        │  Load Balancer  │
                        │   (ELB/GCLB)    │
                        └────────┬────────┘
                                 │
                 ┌───────────────┼───────────────┐
                 │               │               │
        ┌────────▼────────┐ ┌───▼──────────┐ ┌─▼───────────────┐
        │   App Server 1  │ │ App Server 2 │ │ App Server N    │
        │  (Container)    │ │ (Container)  │ │  (Container)    │
        └────────┬────────┘ └───┬──────────┘ └─┬───────────────┘
                 │               │               │
                 └───────────────┼───────────────┘
                                 │
                 ┌───────────────┼───────────────┐
                 │               │               │
        ┌────────▼────────┐ ┌───▼──────────┐ ┌─▼───────────────┐
        │   PostgreSQL    │ │    Redis     │ │    InfluxDB     │
        │    (RDS/SQL)    │ │  (Cache)     │ │  (Time Series)  │
        └─────────────────┘ └──────────────┘ └─────────────────┘
```

### Edge Deployment (RPI Standalone)

```
┌─────────────────────────────────────────────────────────────┐
│                    Raspberry Pi 4                           │
├─────────────────────────────────────────────────────────────┤
│  ┌───────────────────────────────────────────────────────┐  │
│  │         SR-71 Simulator (Go Binary)                   │  │
│  ├───────────────────────────────────────────────────────┤  │
│  │  • Embedded SQLite Database                           │  │
│  │  • Local Web Server (Optional)                        │  │
│  │  • Console UI                                         │  │
│  │  • GPIO Control                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                              │
│  Display: HDMI ────────────► Monitor/TV                     │
│  Input: USB ───────────────► Keyboard/Joystick              │
│  GPIO ─────────────────────► Arduino/Sensors/LEDs           │
└──────────────────────────────────────────────────────────────┘
```

## Security Architecture

### Authentication Flow

```
┌──────────┐                                    ┌──────────┐
│  Client  │                                    │  Server  │
└─────┬────┘                                    └─────┬────┘
      │                                               │
      │  1. POST /api/v1/auth/login                  │
      │  {username, password}                        │
      ├──────────────────────────────────────────────►│
      │                                               │
      │                      2. Validate credentials  │
      │                         Generate JWT token    │
      │                                               │
      │  3. Response {token, expires_at}             │
      │◄──────────────────────────────────────────────┤
      │                                               │
      │  4. GET /api/v1/simulation/status            │
      │  Header: Authorization: Bearer <token>       │
      ├──────────────────────────────────────────────►│
      │                                               │
      │                      5. Validate token        │
      │                         Check permissions     │
      │                                               │
      │  6. Response {data}                          │
      │◄──────────────────────────────────────────────┤
      │                                               │
```

### Authorization Levels

| Role        | Permissions                                      |
|-------------|--------------------------------------------------|
| Admin       | Full access, user management, configuration      |
| Instructor  | Start/stop sims, view all data, create missions  |
| Trainee     | Participate in assigned missions, view own data  |
| Viewer      | Read-only access to public simulations           |
| API Client  | Limited programmatic access with rate limits     |

## Performance Considerations

### Optimization Strategies

1. **Simulation Loop:**
   - Fixed time step (60 Hz recommended)
   - Separate rendering from simulation logic
   - Use goroutines for concurrent module updates

2. **Database:**
   - Write telemetry in batches (buffer 1000 records)
   - Use time-series database for high-frequency data
   - Index frequently queried fields
   - Archive old simulations

3. **API:**
   - Rate limiting (100 req/min per client)
   - Response caching (Redis)
   - Pagination for large datasets
   - Compression (gzip)

4. **WebSocket:**
   - Throttle updates to 10-30 Hz for UI
   - Binary encoding for efficiency
   - Client-side interpolation
   - Connection pooling

5. **Hardware:**
   - Optimize GPIO operations
   - Reduce serial communication overhead
   - Buffer sensor readings
   - Use DMA where available

### Scalability

**Horizontal Scaling:**
- Stateless API servers (can run multiple instances)
- Session affinity for WebSocket connections
- Shared database and cache layer

**Vertical Scaling:**
- Increase CPU for more concurrent simulations
- More RAM for caching and buffering
- Faster storage for database performance

## Testing Strategy

### Test Pyramid

```
                    ┌─────────┐
                    │   E2E   │  (10%)
                    │  Tests  │
                ┌───┴─────────┴───┐
                │   Integration   │  (30%)
                │      Tests      │
            ┌───┴─────────────────┴───┐
            │      Unit Tests         │  (60%)
            │                         │
            └─────────────────────────┘
```

**Unit Tests:** Each module independently
**Integration Tests:** Module interactions, API endpoints
**E2E Tests:** Complete simulation scenarios

## Monitoring & Observability

### Metrics to Track

**Application Metrics:**
- Active simulations count
- API response times
- WebSocket connections
- Error rates
- Request throughput

**Simulation Metrics:**
- Average mission duration
- Engine mode distribution
- Altitude distribution
- Speed distribution
- Fuel consumption patterns

**Infrastructure Metrics:**
- CPU usage
- Memory usage
- Disk I/O
- Network bandwidth
- Database query performance

### Logging Levels

- **DEBUG:** Detailed simulation state changes
- **INFO:** Simulation start/stop, mode changes
- **WARN:** Non-critical issues (fuel low, altitude warning)
- **ERROR:** Critical failures (engine failure, system errors)
- **FATAL:** Unrecoverable errors

## Future Architecture Considerations

1. **Microservices:** Split into separate services (simulation, telemetry, missions)
2. **Event Sourcing:** Store all state changes as events
3. **CQRS:** Separate read and write models
4. **Message Queue:** RabbitMQ/Kafka for async communication
5. **Service Mesh:** Istio for advanced networking
6. **Kubernetes:** Container orchestration for cloud deployment

## Conclusion

This architecture provides:
- **Modularity:** Easy to extend and maintain
- **Scalability:** Can handle multiple concurrent simulations
- **Flexibility:** Supports multiple platforms and interfaces
- **Robustness:** Proper error handling and recovery
- **Observability:** Comprehensive monitoring and logging

The design allows for gradual implementation, starting with core functionality and progressively adding advanced features.
