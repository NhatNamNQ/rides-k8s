# Simulator Workflow

This document explains how the simulator works in the local rides-k8s project.
It is written for someone who is learning Go and DevOps at the same time.

The simulator is not the main API. It is a second Go service whose job is to act
like a fake user. It keeps creating ride requests so the API, Prometheus, and
Grafana have real activity to show.

## Big Picture

The simulator does four things:

1. starts a small HTTP server
2. runs a background loop
3. sends ride requests to the API
4. exposes its own `/metrics` endpoint

That gives you this workflow:

```mermaid
flowchart LR
    A[main.go] --> B[config.Load()]
    A --> C[sim.New()]
    A --> D[http.Server]
    A --> E[signal.NotifyContext]

    C --> F[simulator.Handler()]
    C --> G[simulator.Run()]
    C --> H[createRide()]
    C --> I[metrics.Registry]

    G --> H
    H --> J[POST /api/rides]
    J --> K[Go API]
    K --> L[PostgreSQL]

    I --> M[/metrics]
    M --> N[Prometheus]
    N --> O[Grafana]

    D --> M
    D --> P[/healthz]
```

In plain language:

- `main.go` starts everything.
- `simulator.go` does the work.
- `metrics.go` records numbers about the simulator.
- `Prometheus` reads those numbers.
- `Grafana` shows them as charts.

## File-by-File Guide

### `cmd/rides-simulator/main.go`

This is the program entrypoint. When the simulator container starts, Go runs this
file first.

What it does:

- loads configuration
- creates a structured logger
- creates the simulator service
- starts the HTTP server
- starts the background simulator loop in a goroutine
- shuts everything down cleanly when the process receives a signal

Important concepts here:

- **`main()`**: the startup function for the whole program
- **configuration**: values come from environment variables, not hardcoded values
- **logger**: prints structured JSON logs
- **goroutine**: lets the HTTP server and the simulator loop run at the same time
- **context cancellation**: tells long-running work to stop when the process exits

The practical meaning:

- the server can answer `/metrics` and `/healthz`
- the loop can keep creating rides in the background
- the process can stop gracefully instead of being killed abruptly

### `internal/config/config.go`

This file turns environment variables into typed Go values.

Key fields:

- `PORT`: port for the simulator HTTP server
- `LOG_LEVEL`: logging level
- `API_BASE_URL`: where the simulator sends ride requests
- `SIMULATOR_INTERVAL`: how long to wait between ride attempts

Why this matters:

- Docker Compose can change behavior without changing code
- local dev, tests, and future Kubernetes deployment can use different settings
- it keeps deployment concerns outside the business logic

### `internal/sim/simulator.go`

This is the core of the simulator.

It has three main jobs:

1. serve HTTP endpoints
2. run the ride-creation loop
3. send requests to the API

#### `type Simulator`

This struct holds the simulator’s dependencies:

- `cfg`: runtime settings
- `logger`: prints messages
- `client`: makes HTTP requests
- `metrics`: tracks simulator numbers
- `random`: chooses random zones
- `httpServer`: the handler for local endpoints

This is a common Go pattern:

- keep the dependencies in one struct
- pass them in once during construction
- use the struct methods later

#### `New(...)`

`New` creates a new simulator instance.

It:

- stores config and logger
- creates an `http.Client`
- creates the Prometheus metrics registry
- stores the random number source
- builds the HTTP routes

Why this is useful:

- tests can pass a fake random source
- the HTTP client has a timeout
- the metrics setup stays inside one place

#### `Handler()`

`Handler()` returns the HTTP handler used by the server.

It wraps the internal routes so every request can be measured.

What it does for each HTTP request:

- records the start time
- passes the request to the real handler
- records the status code
- records the duration
- sends that data to the simulator metrics

This means the simulator can count how often its own `/metrics` and `/healthz`
endpoints are used.

#### `Run(ctx context.Context)`

This is the background loop.

It repeats forever until the context is canceled.

Inside each loop:

1. remember the start time
2. call `createRide(ctx)`
3. record success or failure as a simulator event
4. record how long the loop took
5. sleep until the next interval

The important Go ideas here:

- **ticker**: waits a fixed amount of time between loops
- **select**: chooses between “context ended” and “ticker fired”
- **context**: allows graceful shutdown
- **background worker**: work that happens without a user directly clicking anything

Practical meaning:

- the simulator keeps creating load on its own
- you do not need to manually press buttons or run curl commands
- Grafana charts keep changing because the system is active

#### `routes()`

This sets up the HTTP endpoints exposed by the simulator.

Current routes:

- `GET /metrics`
- `GET /healthz`

Why these matter:

- `/metrics` is for Prometheus
- `/healthz` is a simple health check for humans and orchestration tools

#### `createRide(ctx context.Context)`

This is the function that sends one ride request to the API.

What it does:

1. creates a fake rider ID
2. picks a random pickup zone
3. picks a random dropoff zone
4. makes sure pickup and dropoff are not the same
5. converts the request to JSON
6. sends `POST /api/rides`
7. checks that the response status is `201 Created`
8. logs success

Important Go concepts:

- **struct literal**: used to build the JSON request object
- **`json.Marshal`**: converts Go data into JSON
- **`http.NewRequestWithContext`**: creates a request that stops if the context is canceled
- **`http.Client.Do`**: sends the request
- **error handling**: every failure is returned instead of being hidden

Practical meaning:

- the simulator is acting like a fake rider client
- the API receives realistic requests
- failures show up clearly in logs and metrics

#### `pickZone(random randomSource)`

This helper chooses one zone from the list of zones.

Why this exists:

- it keeps randomness in one small function
- it makes tests easier
- it gives the simulator varied traffic

#### `indexOf(zone string)`

This helper finds the position of a zone in the list.

It is used when pickup and dropoff accidentally match.

If the two zones are the same, the code chooses the next zone in the list so
the request stays meaningful.

#### `statusRecorder`

This wrapper remembers the HTTP status code written by a handler.

Why this is needed:

- Go’s `http.ResponseWriter` does not automatically expose the final status code
- the metrics code needs the status code to label requests correctly

#### `writeJSON(...)`

This helper writes JSON responses for the simulator’s local endpoints.

It:

- sets the `Content-Type`
- writes the response status
- encodes the JSON body

### `internal/sim/metrics.go`

This file defines the metrics that Prometheus will scrape.

Think of metrics as counters and measurements that answer questions like:

- how many requests did the simulator handle?
- how many ride creation attempts succeeded?
- how long does one simulator loop take?

#### `type Metrics`

This struct keeps the Prometheus registry and the metric objects together.

It contains:

- `registry`
- `httpRequestsTotal`
- `httpRequestDuration`
- `simulatorEventsTotal`
- `simulatorLoopDuration`

#### `NewMetrics()`

This creates all metrics and registers them.

The metrics are:

- `http_requests_total`
- `http_request_duration_seconds`
- `simulator_events_total`
- `simulator_loop_duration_seconds`

Practical meaning:

- Prometheus can scrape these numbers from `/metrics`
- Grafana can graph them over time
- you can see if the simulator is working, failing, or slowing down

#### `Handler()`

This returns the Prometheus HTTP handler for the simulator registry.

That means `GET /metrics` gives Prometheus the data it expects.

#### `ObserveHTTPRequest(...)`

This records an HTTP request handled by the simulator’s own server.

It stores:

- method
- path
- status code
- duration

This is classic observability data:

- volume
- latency
- success/failure

#### `ObserveEvent(...)`

This records a simulator event such as:

- `create_ride`
- `success`
- `error`

This helps you see whether the simulator is healthy and productive.

#### `ObserveLoopDuration(...)`

This records how long one whole simulator loop takes.

That is useful because:

- if the API gets slower, the loop may get slower too
- if the simulator starts failing often, the loop timing may change

### `internal/sim/simulator_test.go`

This test checks the most important behavior:

- the simulator sends a `POST`
- it uses the `/api/rides` path
- it treats `201 Created` as success

Why this test matters:

- it proves the simulator is calling the API correctly
- it gives you confidence that the request format is right
- it catches regressions if the API path changes

The test uses:

- `httptest.NewServer` to fake the API
- a fixed random source so the test is predictable
- a short interval and a discard logger so the test stays simple

## Workflow Summary

This is the simplest way to remember the whole system:

```text
main.go starts the simulator
config.go reads env vars
simulator.go runs a loop and calls the API
metrics.go records simulator numbers
Prometheus scrapes /metrics
Grafana shows the charts
```

## Why This Stage Exists

This stage teaches an important DevOps lesson:

You do not really understand monitoring until something is actively happening.

Without the simulator:

- the API may be idle
- Prometheus graphs may be flat
- Grafana dashboards may not tell you much

With the simulator:

- traffic appears automatically
- the API gets exercised
- errors and latency become visible
- you can practice debugging real behavior

## Helpful Mental Model

Think of the simulator as a robot user:

- it behaves like a rider client
- it does work in the background
- it creates real pressure on the API
- it also reports on its own behavior

That is why the simulator has both:

- request-sending logic
- metrics-exposing logic

It is both a traffic generator and an observable service.
