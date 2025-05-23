# Spacecraft Telemetry System

## Getting Started

### Prerequisites

- Docker and Docker Compose installed on your system
- Git (to clone this repository)

### Quick Start

1. Clone this repository:

```bash
git clone https://github.com/hsu-isaac/turion-takehome.git
cd turion-takehome
```

2. Set up environment variables:

```bash
# Copy the frontend environment example file
cp telemetry-frontend/.env.example telemetry-frontend/.env
# Edit .env file with your specific configuration if needed
```

3. Start all services using Docker Compose:

```bash
docker-compose up
```

4. Access the services:

- Frontend Dashboard: http://localhost:3001
- API Service: http://localhost:8080
- Database: localhost:5432 (PostgreSQL)

---

About this starter code:
This is a sample telemetry generator that sends spacecraft data over UDP. While you're
welcome to use it as-is, you can also implement a simpler solution. The key
requirement is the ability to:
Serialize data into bytes
Send those bytes over UDP
Deserialize the received bytes back into structured data
That's all you need to know to get started!

## Project Requirements

### Part 1: Telemetry Ingestion Service (Required)

#### Requirements

1. Create a service that:
   - [x] Listens for UDP packets containing spacecraft telemetry
   - [x] Decodes CCSDS-formatted packets according to provided structure
   - [x] Validates telemetry values against defined ranges:
     - [x] Temperature: 20.0°C to 30.0°C (normal), >35.0°C (anomaly)
     - [x] Battery: 70-100% (normal), <40% (anomaly)
     - [x] Altitude: 500-550km (normal), <400km (anomaly)
     - [x] Signal Strength: -60 to -40dB (normal), <-80dB (anomaly)
   - [x] Persists data to a database (Timescale or PostgreSQL preferred but not required)
   - [x] Implements an alerting mechanism for out-of-range values (Anomalies)

### Part 2: Telemetry API Service (Required)

#### Requirements

1. Create a REST API using:
   - [x] Fiber/Echo (Go)
   - [ ] ~~FastAPI (Python)~~
   - [ ] ~~Express/Fastify (TypeScript)~~

#### API Endpoints (Minimum Required)

- [x] `GET /api/v1/telemetry`
  - [x] Query Parameters:
    - [x] `start_time` (ISO8601)
    - [x] `end_time` (ISO8601)
- [x] `GET /api/v1/telemetry/current`
  - [x] Returns latest telemetry values
- [x] `GET /api/v1/telemetry/anomalies`
  - [x] Query Parameters:
    - [x] `start_time` (ISO8601)
    - [x] `end_time` (ISO8601)

### Part 3: Front End Implementation

#### Requirements

Create a telemetry dashboard that:

- [x] Real-time updates: Display the most recent telemetry values in real time
- [x] Historical graphs or tables: Show historical telemetry data
- [x] Anomaly notifications: Provide real-time anomaly notifications

#### Technical Requirements

- [x] Use React (You can use another front end tool if you do not understand React)

### Optional Requirements

#### Frontend-Focused Optional Requirements

- [x] Error handling: Implement basic error handling and loading states
- [x] Responsive design: Ensure the dashboard works on desktop and mobile
- [x] User experience: Add features like:
  - [x] Search/filter for telemetry data
  - [x] Dark mode
  - [x] Theming
- [ ] Telemetry visualization: Include charts for telemetry metrics (embedded Grafana is acceptable)

#### Backend-Focused Optional Requirements

- [x] Database migrations: Implement migrations for storing telemetry data and managing schema evolution
  - [x] Setting up the system and having one migration is acceptable
- [x] Observability: Use OpenTelemetry to instrument backend APIs and pipelines
  - [ ] Optional visualization using Grafana Tempo, Loki, Prometheus/Mimir
- [ ] Integration test: Write integration tests to ensure the API correctly:
  - [ ] Serves telemetry data
  - [ ] Handles edge cases (e.g., real-time updates, data gaps)
- [ ] Performance testing: Include performance benchmarks for:
  - [ ] Real-time update pipelines
  - [ ] Historical queries

### Bonus Points

- [x] Docker Compose: Provide a working Docker Compose file for local development with all dependencies:
  - [x] Frontend
  - [x] Backend
  - [x] Database
  - [ ] Observability tools
- [ ] Comprehensive tests:
  - [ ] Unit tests
  - [ ] Integration tests
  - [ ] End-to-end tests
- [ ] Performance testing results: Provide evidence of load testing or benchmarking
  - [ ] Using tools like JMeter, k6, or Locust
