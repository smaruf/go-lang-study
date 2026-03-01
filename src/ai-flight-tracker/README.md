# AI Flight Tracker (Go)

A Go replica of the Python AI flight tracker — a REST API serving live flight data, route search, pricing windows, and an AI assistant powered by Ollama.

## Features

- **GET /health** — health check
- **GET /** — serve static `index.html` (interactive map + route search + AI chat)
- **GET /api/routes** — all routes, priority-sorted with simulated pricing
- **GET /api/route/{code}** — route details + pricing windows (last-minute, 2 weeks, 1 month, early-bird)
- **GET /api/search?origin=&destination=&max_stops=&max_price=** — filter routes
- **GET /api/flights?lamin=&lamax=&lomin=&lomax=** — live flights from OpenSky Network
- **GET /api/flights/{icao24}** — single flight detail by ICAO24 code
- **POST /api/aviation/ask** — AI assistant (JSON body `{"question": "..."}`)

## Tech Stack

- **Language**: Go 1.21 (standard library only — no external dependencies)
- **AI backend**: [Ollama](https://ollama.ai) (llama2 by default)
- **Live flight data**: [OpenSky Network](https://opensky-network.org/)
- **Frontend**: Leaflet.js map + vanilla JS

## Project Structure

```
ai-flight-tracker/
├── main.go               # HTTP server, middleware, handlers
├── go.mod
├── routes.json           # Static route database
├── config/
│   └── config.go         # Environment-based configuration
├── services/
│   ├── routes_repo.go    # Route loading, search, filtering
│   ├── opensky.go        # OpenSky Network API client
│   ├── pricing.go        # Simulated pricing with TTL cache
│   └── ai.go             # AI service with intent routing
├── static/
│   ├── index.html
│   ├── app.js
│   └── style.css
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

## Quick Start

### Run locally

```bash
go build -o flight-tracker .
./flight-tracker
# Server starts on :8080
```

### Run with Docker Compose (includes Ollama)

```bash
docker compose up --build
# Pull a model: docker compose exec ollama ollama pull llama2
```

## Configuration

All settings are via environment variables (see `.env.example`):

| Variable              | Default                                  | Description                    |
|-----------------------|------------------------------------------|--------------------------------|
| `PORT`                | `8080`                                   | HTTP port                      |
| `OPENSKY_URL`         | `https://opensky-network.org/api`        | OpenSky base URL               |
| `OPENSKY_USERNAME`    | _(empty)_                                | Optional OpenSky credentials   |
| `OPENSKY_PASSWORD`    | _(empty)_                                | Optional OpenSky credentials   |
| `OLLAMA_URL`          | `http://localhost:11434/api/generate`    | Ollama API endpoint            |
| `OLLAMA_MODEL`        | `llama2`                                 | Model name                     |
| `CACHE_TTL`           | `300`                                    | Pricing cache TTL (seconds)    |

## Running Tests

```bash
go test ./...
```

## GUI Application

A cross-platform native desktop GUI (built with [Fyne v2](https://fyne.io)) is available at `cmd/gui/`. It provides:

- **Server tab** — start/stop the HTTP server, set port, open the web UI in a browser
- **Configuration tab** — edit OpenSky, Ollama, cache TTL, and CORS settings (applied on next start)
- **Deploy tab** — build/run/stop Docker containers and docker-compose services
- **Logs tab** — real-time server log output

### Prerequisites (Linux/WSL)

```bash
sudo apt-get install libgl1-mesa-dev xorg-dev
```

### Run the GUI

```bash
cd src/ai-flight-tracker
go run ./cmd/gui
```

### Build a native binary

```bash
go build -o flight-tracker-gui ./cmd/gui
./flight-tracker-gui
```

Supported platforms: Windows, macOS, Linux (including WSL with a display), Android, iOS.
