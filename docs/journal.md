# Learning Journal

Use this file after every stage.

The journal should explain the project in beginner-friendly language. It should not just list commands. Write down what changed, why it matters, what new concepts appeared, and what was verified.

Recommended format:

```text
## Stage X: Stage Name

Date:

### Summary

Short explanation of what this stage added.

### What Changed

- Important change 1
- Important change 2

### Concepts Learned

- Concept 1: beginner explanation
- Concept 2: beginner explanation

### How I Verified It

- Test or manual check that proved the stage worked.

### Problems and Fixes

- Problem: what went wrong
- Fix: how it was handled

### Beginner Notes

Extra explanation for future me.
```

## Stage 0: Plan the System

Date: 2026-05-14

### Summary

This stage defined the project direction before writing application code. The project will be a small ridesharing simulation platform used to learn Go, Kubernetes, Prometheus, Grafana, Jenkins, AWS, and Terraform.

The main decision was to keep the first version small. The project should start with a simple Go API, then gradually add storage, containers, Kubernetes, monitoring, CI/CD, AWS, and Terraform.

### What Changed

- Created the project learning plan.
- Created the detailed stage guide.
- Created agent rules for future coding sessions.
- Created review agent rules for code, security, Kubernetes, observability, CI/CD, Terraform, and documentation.
- Created the initial architecture document.

### Concepts Learned

- **Architecture first**: before writing code, define the main services and how they communicate.
- **Small stages**: learning is easier when every stage has a clear goal and finish line.
- **Local before cloud**: the app should work locally before Docker, Kubernetes, or AWS are introduced.
- **Manual before automation**: Jenkins and Terraform should automate workflows that already work manually.
- **Security later, but not forgotten**: authentication ideas from `Go_Secure_Auth_Pro` are useful, but auth should not block the first Go API.

### How I Verified It

- Confirmed the main project folders exist.
- Confirmed `docs/architecture.md` explains the first version of the system.
- Confirmed `AGENTS.md`, `REVIEW_AGENTS.md`, `LEARNING_PLAN.md`, and `STAGE_GUIDE.md` exist.

### Problems and Fixes

- No technical problem yet. This was a planning stage.

### Beginner Notes

The project is intentionally not starting with everything at once. Kubernetes, Jenkins, AWS, and Terraform are powerful, but they add many ways for things to fail. Starting with a small Go API makes future failures easier to understand.

## Stage 1: First Go HTTP API

Date: 2026-05-14

### Summary

This stage created the first working backend service. The API can respond to health checks, readiness checks, create rides, and list rides.

The data is stored in memory for now. That means rides disappear when the API process stops. This is acceptable for Stage 1 because the goal is to learn the HTTP flow before adding PostgreSQL.

### What Changed

- Created the Go module in `services/api`.
- Created the API entrypoint at `services/api/cmd/rides-api/main.go`.
- Created HTTP server code in `services/api/internal/httpapi`.
- Created ride storage code in `services/api/internal/ride`.
- Added an in-memory ride store.
- Added these endpoints:
  - `GET /healthz`
  - `GET /readyz`
  - `GET /api/rides`
  - `POST /api/rides`
- Added tests for health, readiness, ride creation, validation, and ride listing.

### Concepts Learned

- **Go module**: a Go project is usually managed with `go.mod`.
- **`cmd/` folder**: contains executable entrypoints. In this project, `cmd/rides-api/main.go` starts the API.
- **`internal/` folder**: contains private application code that should not be imported by other Go modules.
- **HTTP handler**: a function that receives a request and writes a response.
- **`http.ResponseWriter`**: used to write the HTTP response.
- **`*http.Request`**: contains information about the incoming request.
- **In-memory store**: temporary storage inside the running process.
- **Mutex**: protects shared memory when multiple requests happen at the same time.
- **Unit test**: code that verifies behavior automatically.

### How I Verified It

- Ran the Go test suite successfully.
- Started the API locally.
- Verified `GET /healthz` returned `200 OK`.
- Verified `GET /readyz` returned `200 OK`.
- Verified `GET /api/rides` returned a JSON list.
- Verified `POST /api/rides` created a ride and returned `201 Created`.

### Problems and Fixes

- Problem: Go tests initially could not write to the Go build cache because the cache is outside the workspace sandbox.
- Fix: reran tests with permission to use the Go build cache.
- Problem: `curl` from the default sandbox could not reach the API process started in another execution context.
- Fix: ran endpoint verification in the same execution context as the API process.

### Beginner Notes

The `Store` type is acting like a tiny database for now:

```text
POST /api/rides -> validate request -> Store.Create -> save ride in memory
GET /api/rides  -> Store.List -> return all rides
```

This is not the final design. Later, PostgreSQL will store the rides. Starting with memory keeps the first API easy to understand.

## Stage 2: Configuration and Logging

Date: 2026-05-14

### Summary

This stage made the API more realistic by adding environment-based configuration and structured logs.

Before this stage, the API always listened on port `8080`. Now the port can be changed with `PORT`. The API also reads `LOG_LEVEL` and `DATABASE_URL`.

### What Changed

- Added `services/api/internal/config`.
- Added support for:
  - `PORT`
  - `LOG_LEVEL`
  - `DATABASE_URL`
- Switched startup logs to Go's `log/slog`.
- Added JSON-formatted logs.
- Added one log for every HTTP request.
- Added a log when a ride is created.
- Added tests for config defaults and environment variable overrides.

### Concepts Learned

- **Environment variable**: configuration passed to the process from outside the code.
- **Default value**: value used when an environment variable is not set.
- **Structured logging**: logs written as key-value data instead of plain text.
- **JSON logs**: useful later because Docker, Kubernetes, and log tools can parse them.
- **Request log**: records method, path, status code, and duration for each HTTP request.
- **Startup log**: confirms how the service started and which config was detected.

### How I Verified It

- Ran the Go test suite successfully.
- Started the API with a custom port:

```text
PORT=9090
LOG_LEVEL=debug
DATABASE_URL=postgres://example
```

- Verified the API listened on `9090`.
- Verified `GET /healthz` worked on the custom port.
- Verified `POST /api/rides` created a ride.
- Confirmed JSON logs appeared for startup, request handling, and ride creation.

### Problems and Fixes

- No application bug was found.
- The same local sandbox limitation from Stage 1 applied when running Go and checking local endpoints.

### Beginner Notes

Hardcoding configuration is a bad habit for services that will run in Docker or Kubernetes. In Kubernetes, values like port, database URL, and log level usually come from ConfigMaps, Secrets, or environment variables.

`DATABASE_URL` is loaded now but not used yet. It will become important in the PostgreSQL stage.

## Stage 3: Prometheus Metrics in Go

Date: 2026-05-14

### Summary

This stage added Prometheus metrics to the Go API. The API now exposes a `GET /metrics` endpoint that Prometheus can scrape later.

Metrics are different from logs. Logs explain individual events. Metrics show numeric system behavior over time, such as request count, request duration, and number of active rides.

### What Changed

- Added the Prometheus Go client dependency.
- Added `services/api/internal/metrics`.
- Added `GET /metrics`.
- Added HTTP request metrics:
  - `http_requests_total`
  - `http_request_duration_seconds`
- Added ride metrics:
  - `rides_created_total`
  - `rides_active`
- Updated the HTTP server to record request count, status code, route path, and duration.
- Updated ride creation to increment ride metrics.
- Added a test that verifies the metrics endpoint exposes expected values.

### Concepts Learned

- **Prometheus**: a monitoring system that collects metrics by scraping HTTP endpoints.
- **`/metrics` endpoint**: an HTTP endpoint that returns metrics in Prometheus text format.
- **Counter**: a metric that only goes up, such as `rides_created_total`.
- **Gauge**: a metric that can go up or down, such as `rides_active`.
- **Histogram**: a metric that groups values into buckets, useful for request duration.
- **Label**: extra information attached to a metric, such as method, path, and status.
- **Scraping**: Prometheus pulls metrics from the service instead of the service pushing metrics to Prometheus.

### How I Verified It

- Ran the Go test suite successfully.
- Started the API locally on port `9090`.
- Opened `GET /metrics`.
- Created a ride with `POST /api/rides`.
- Confirmed the metrics output included:
  - `rides_created_total 1`
  - `rides_active 1`
  - `http_requests_total{method="POST",path="/api/rides",status="201"} 1`

### Problems and Fixes

- Problem: adding the Prometheus dependency initially failed because the sandbox could not resolve `proxy.golang.org`.
- Fix: reran the dependency download with network permission.
- Problem: Go commands still needed access to the Go build cache outside the workspace.
- Fix: ran tests with the required permission.

### Beginner Notes

The API now has two kinds of observability:

```text
Logs    -> event details, useful for debugging one request
Metrics -> numbers over time, useful for dashboards and alerts
```

Example:

```text
POST /api/rides
  -> creates a ride
  -> logs "ride created"
  -> increments rides_created_total
  -> increments rides_active
  -> records HTTP request duration
```

This prepares the project for the later Prometheus and Grafana stages. Prometheus will scrape `GET /metrics`, and Grafana will visualize those values.

## Stage 4: PostgreSQL with Supabase

Date: 2026-05-15

### Summary

This stage started moving ride storage from memory to PostgreSQL. Instead of using a local Docker Postgres container, the project will use Supabase free tier for the database.

The code now supports two storage modes:

```text
No DATABASE_URL -> use in-memory storage
DATABASE_URL set -> connect to PostgreSQL
```

This keeps local development easy while allowing real persistence when a Supabase database is configured.

### What Changed

- Added PostgreSQL driver support with `pgx`.
- Added a `ride.Repository` interface so the HTTP API does not depend directly on memory storage.
- Renamed the old memory store to `MemoryStore`.
- Added `PostgresStore` for real PostgreSQL persistence.
- Added `services/api/migrations/001_init.sql`.
- Updated `/readyz` so it checks storage availability.
- Added a test for readiness failure when storage is unavailable.
- Added `services/api/.env.example` with a safe Supabase connection string template.

### Concepts Learned

- **Repository interface**: lets the API use different storage implementations without changing handlers.
- **Memory store**: useful for tests and early local development.
- **Postgres store**: stores data in a real database so it survives process restarts.
- **Migration**: SQL file that creates or changes database tables.
- **Readiness check**: should verify dependencies like the database.
- **Supabase Session pooler**: recommended for this long-running Go API, especially when direct IPv6 connections are not available.
- **Supabase Transaction pooler**: useful fallback when the local network cannot reach port `5432`; it uses port `6543`.
- **`default_query_exec_mode=simple_protocol`**: pgx setting needed for compatibility with Supabase transaction pooler.

### How I Verified It

- Ran the Go test suite successfully.
- Confirmed the API still works with the in-memory store when `DATABASE_URL` is not set.
- Confirmed `/readyz` can return `503 Service Unavailable` when the storage dependency fails.
- Tested Supabase network reachability:
  - port `5432` failed from the current network
  - port `6543` succeeded
- Verified `psql` can connect to Supabase through port `6543`.
- Verified the Go API can connect to Supabase using port `6543` plus `default_query_exec_mode=simple_protocol`.
- Verified:
  - `GET /readyz` returned `200 OK`
  - `GET /api/rides` returned `200 OK`

### Problems and Fixes

- Problem: Docker tried to download the `postgres:16` image, but the network was weak.
- Fix: stopped the Docker download and switched the Stage 4 plan to Supabase free tier.
- Problem: Supabase pooler port `5432` timed out from the current network.
- Fix: used Supabase pooler port `6543`.
- Problem: the Go `pgx` driver needs extra compatibility settings for the transaction pooler.
- Fix: added `default_query_exec_mode=simple_protocol` to the Go API connection string.

### Beginner Notes

The most important design change is this:

```text
HTTP handlers depend on ride.Repository
MemoryStore implements ride.Repository
PostgresStore implements ride.Repository
```

That means the API handlers do not care where rides are stored. They call:

```text
Create
List
Ping
```

The actual storage can be memory or PostgreSQL.

For Supabase, prefer the Session pooler if port `5432` works. On this network, port `5432` was unreachable, so the project uses the pooler on port `6543` instead.

Use this shape for the Go API:

```text
postgresql://USER:PASSWORD@HOST:6543/postgres?sslmode=require&default_query_exec_mode=simple_protocol
```

Do not commit the real connection string because it contains the database password.

## Stage 5: Docker Compose

Date: 2026-05-16

### Summary

This stage packages the Go API and PostgreSQL into containers so the system can run in a repeatable local environment.

The main idea is that we no longer depend only on tools installed directly on the laptop. Instead, Docker provides a defined runtime for the API and the database.

### What Changed

- Added `services/api/Dockerfile`.
- Added `services/api/.dockerignore`.
- Added `deploy/docker-compose/docker-compose.yml`.
- Configured a multi-stage Docker build for the Go API.
- Configured Docker Compose services for:
  - `api`
  - `postgres`
- Configured the API container to use:
  - `PORT=8080`
  - `LOG_LEVEL=debug`
  - `DATABASE_URL=postgres://rides:rides@postgres:5432/rides?sslmode=disable`
- Mounted `services/api/migrations/001_init.sql` into Postgres auto-init.
- Added a named Docker volume for Postgres data.

### Concepts Learned

- **Dockerfile**: instructions for building a container image.
- **Multi-stage build**: use one image to compile the Go binary and another smaller image to run it.
- **Docker Compose**: runs multiple related containers together.
- **Service name as hostname**: inside Compose, the API can reach Postgres using `postgres` as the host.
- **Volume**: keeps Postgres data alive across container restarts.
- **Container startup dependency**: `depends_on` plus health checks helps the API wait for Postgres.
- **Auto-init SQL**: Postgres runs SQL files from `/docker-entrypoint-initdb.d/` on first database initialization.

### How I Verified It

- Ran the Go test suite successfully after the Stage 5 changes.
- Verified the Docker Compose file structure is in place.

Full container runtime verification is still pending if Docker image downloads or builds are slow. The intended next manual check is:

- `docker compose -f deploy/docker-compose/docker-compose.yml up --build`
- `GET /healthz`
- `POST /api/rides`
- restart the stack
- confirm rides persist through Postgres

### Problems and Fixes

- No application bug was introduced during the Stage 5 file setup.
- Earlier weak network conditions made large Docker image downloads undesirable, so the stage was prepared with config-first validation.

### Beginner Notes

The important relationship in Compose is:

```text
api container -> connects to postgres container
```

Inside Docker Compose, this works:

```text
postgres://rides:rides@postgres:5432/rides?sslmode=disable
```

because `postgres` is the service name, and Docker Compose provides internal DNS for service-to-service communication.

The Dockerfile uses two stages:

```text
builder stage -> compiles the Go binary
runtime stage -> runs only the final binary
```

That keeps the final image smaller and cleaner than shipping the whole Go toolchain in production-style containers.
