# 🤿 DiveLog

A full-stack personal scuba dive log built in Go. Features a REST API, CLI, and web UI with a nautical-themed dashboard.

## Table of Contents

- [🤿 DiveLog](#-divelog)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Architecture](#architecture)
  - [Features](#features)
  - [Tech Stack](#tech-stack)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Run locally](#run-locally)
  - [CLI Usage](#cli-usage)
  - [REST API](#rest-api)
    - [Dives](#dives)
    - [NDL Calculator](#ndl-calculator)
    - [Export](#export)
    - [System](#system)
    - [Example: Log a dive](#example-log-a-dive)
    - [Example: Calculate NDL](#example-calculate-ndl)
  - [Configuration](#configuration)
  - [Testing](#testing)
    - [Coverage summary](#coverage-summary)
    - [Go patterns showcased in tests](#go-patterns-showcased-in-tests)
  - [Docker](#docker)
    - [Build](#build)
    - [Run](#run)
  - [Kubernetes \& Helm](#kubernetes--helm)
    - [Prerequisites](#prerequisites-1)
    - [Create a local cluster](#create-a-local-cluster)
    - [Load the image into kind](#load-the-image-into-kind)
    - [Deploy](#deploy)
    - [Verify](#verify)
    - [Upgrade after changes](#upgrade-after-changes)
    - [Backup database](#backup-database)
    - [Restore database](#restore-database)
    - [Uninstall](#uninstall)
  - [Metrics](#metrics)
  - [Project Structure](#project-structure)

---

## Overview

DiveLog is a personal dive logging application that allows scuba divers to record and track their dives, calculate no-decompression limits (NDL), and export their dive history. It was built as a portfolio project to showcase Go development patterns, containerization, and Kubernetes deployment.

---

## Architecture

The application follows a clean layered architecture with clear separation of concerns:
HTTP Request
│
▼
Handler        ← decodes requests, encodes responses, maps errors to HTTP status codes
│
▼
Service        ← business logic, validation, orchestration
│
▼
Repository      ← data persistence interface
│
▼
SQLite        ← storage implementation

Each layer only communicates with the layer directly below it. The repository layer is interface-driven — the service layer never depends on a concrete implementation, making the storage backend swappable with zero changes to business logic.

---

## Features

- **Dive Logging** — record dives with depth, duration, location, gas mix, water temp, visibility, tank pressure, and notes
- **NDL Calculator** — no-decompression limit calculator based on PADI recreational dive tables
- **Dive List** — paginated list of all logged dives sorted by date
- **CSV Export** — download your full dive log as a CSV file
- **Health Endpoint** — `/api/health` for liveness and readiness probes
- **Prometheus Metrics** — HTTP request count and duration exposed at `/api/metrics`
- **Web UI** — nautical-themed single-page dashboard
- **CLI** — full command line interface for all operations

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| Router | chi v5 |
| Database | SQLite (via go-sqlite3) |
| CLI | cobra |
| Metrics | Prometheus client_golang |
| Container | Docker (multi-stage build) |
| Orchestration | Kubernetes via kind |
| Packaging | Helm v3 |

---

## Getting Started

### Prerequisites

- Go 1.25+
- GCC (required for CGo/SQLite)
  - macOS: `xcode-select --install`
  - Ubuntu: `apt-get install gcc`
- Task (taskfile.dev)

### Run locally

```bash
# Install dependencies
go mod tidy

# Run the server
task app

# Open the web UI
open http://localhost:8080
```

---

## CLI Usage

```bash
# Log a dive
./divelog log \
  --site "Blue Corner" \
  --location "Palau" \
  --depth 32 \
  --duration 58 \
  --temp 28 \
  --vis 25 \
  --type drift \
  --rating 5 \
  --notes "Schooling barracuda!"

# List recent dives
./divelog list

# Calculate NDL
./divelog ndl --depth 30
./divelog ndl --depth 25 --o2 32
```

---

## REST API

### Dives
POST   /api/dives          Log a new dive
GET    /api/dives          List dives (?limit=20&offset=0)
GET    /api/dives/{id}     Get a single dive by ID
DELETE /api/dives/{id}     Delete a dive

### NDL Calculator
POST /api/ndl

Request body:
```json
{
  "depth": 30,
  "o2_percent": 21
}
```

### Export
GET /api/export/csv    Download dive log as CSV

### System
GET /api/health     Liveness/readiness check
GET /api/metrics    Prometheus metrics

### Example: Log a dive

```bash
curl -X POST http://localhost:8080/api/dives \
  -H "Content-Type: application/json" \
  -d '{
    "site_name": "Blue Corner",
    "location": "Palau",
    "date": "2024-07-15",
    "max_depth": 32.0,
    "duration": 58,
    "water_temp": 28.0,
    "visibility": 25.0,
    "tank_start": 200,
    "tank_end": 60,
    "o2_percent": 21.0,
    "dive_type": "drift",
    "water_type": "salt",
    "rating": 5
  }'
```

### Example: Calculate NDL

```bash
curl -X POST http://localhost:8080/api/ndl \
  -H "Content-Type: application/json" \
  -d '{"depth": 30, "o2_percent": 21}'
```

---

## Configuration

Configuration is loaded from environment variables with sensible defaults:

| Variable | Default | Description |
|---|---|---|
| `ADDR` | `:8080` | HTTP listen address |
| `DB_PATH` | `data/divelog.db` | SQLite database file path |

---

## Testing

```bash
# Run all tests with coverage
task test

# Run tests for a specific package
go test ./internal/calculator/... -v
go test ./internal/service/... -v
```

### Coverage summary

| Package | Coverage |
|---|---|
| calculator | 100% |
| service | 87.9% |

### Go patterns showcased in tests

- **Table-driven tests** — idiomatic Go test structure used throughout
- **Mock repository** — service layer tested without a real database by implementing the `DiveRepository` interface with an in-memory mock
- **Sentinel error assertions** — `errors.Is()` used to verify correct error types are returned

---

## Docker

### Build

```bash
docker build -t divelog .
```

### Run

```bash
docker run -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -e DIVELOG_DB=/app/data/divelog.db \
  divelog
```

The multi-stage Dockerfile keeps the final image lean by only including the compiled binary and web assets — the Go toolchain is discarded after the build stage. `CGO_ENABLED=1` is required because `go-sqlite3` uses C bindings.

---

## Kubernetes & Helm

### Prerequisites

- kind
- helm
- kubectl

### Create a local cluster

```bash
kind create cluster --name divelog
```

### Load the image into kind

```bash
kind load docker-image divelog:latest --name divelog
```

### Deploy

```bash
helm install divelog ./helm/divelog
```

### Verify

```bash
kubectl get pods
kubectl port-forward deployment/divelog 8080:8080
curl http://localhost:8080/api/health
```

### Upgrade after changes

```bash
docker build -t divelog:latest .
kind load docker-image divelog:latest --name divelog
helm upgrade divelog ./helm/divelog
```

### Backup database

```bash
task k8s:db:backup
```

### Restore database

```bash
task k8s:db:restore
```

### Uninstall

```bash
task down
```

---

## Metrics

The app exposes Prometheus metrics at `/api/metrics`. Two metrics are tracked:

| Metric | Type | Description |
|---|---|---|
| `divelog_http_requests_total` | Counter | Total HTTP requests by method, path, and status code |
| `divelog_http_request_duration_seconds` | Histogram | Request duration in seconds by method and path |

These metrics are ready to be scraped by a Prometheus instance. When deploying to a full Kubernetes environment with the Prometheus operator, add a `ServiceMonitor` resource to enable automatic scraping.

---

## Project Structure
divelog/
├── cmd/
│   └── divelog/
│       └── main.go           ← entrypoint, CLI commands, HTTP server
├── internal/
│   ├── calculator/
│   │   ├── ndl.go            ← NDL lookup table calculator
│   │   └── ndl_test.go
│   ├── config/
│   │   └── config.go         ← env-based configuration
│   ├── exporter/
│   │   └── exporter.go       ← CSV export
│   ├── handler/
│   │   ├── handler.go        ← HTTP handlers
│   │   └── metrics.go        ← Prometheus middleware
│   ├── model/
│   │   ├── dive.go           ← domain types
│   │   └── errors.go         ← sentinel errors
│   ├── repository/
│   │   ├── repository.go     ← DiveRepository interface
│   │   └── sqlite.go         ← SQLite implementation
│   └── service/
│       ├── dive_service.go   ← business logic
│       └── dive_service_test.go
├── web/
│   └── static/
│       └── index.html        ← nautical-themed web UI
├── helm/
│   └── divelog/              ← Helm chart
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/
│           ├── deployment.yaml
│           ├── service.yaml
│           ├── pvc.yaml
│           └── _helpers.tpl
├── scripts/
│   └── backup.sh             ← database backup script
├── Dockerfile                ← multi-stage build
├── Taskfile.yml
├── go.mod
└── README.md