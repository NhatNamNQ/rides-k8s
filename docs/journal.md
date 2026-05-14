# Learning Journal

Use this file after every stage.

## Stage 0: Plan the System

Date:

### What I Built

- Created the project learning plan.
- Created the detailed stage guide.
- Created agent rules.
- Created review agent playbook.
- Created initial architecture documentation.

### Commands I Used

```bash
mkdir -p docs services/api services/simulator services/web deploy/k8s deploy/docker-compose observability/prometheus observability/grafana infra/terraform
```

### What Broke

- Nothing yet.

### How I Fixed It

- Not applicable.

### What I Learned

- The project should start with a small Go API before adding Kubernetes, Jenkins, AWS, or Terraform.
- Minikube is the local Kubernetes target.
- Prometheus and Grafana should be added after the API exposes useful metrics.
- Jenkins should automate a manual process that already works.
- Terraform should define AWS resources after those resources are understood manually.
- Security/auth ideas from `Go_Secure_Auth_Pro` are useful, but should be added later.

### Questions

- Should the simulator start in Go for simplicity, or Node.js to match the original Rides project?
- Should authentication be added after the core API works, or kept out of scope until the Kubernetes/AWS path is complete?

## Stage 1: First Go HTTP API

Date: 2026-05-14

### What I Built

- Created the first Go API in `services/api`.
- Refactored the Go API into a small `cmd/` and `internal/` structure inspired by `golang-standards/project-layout`.
- Added an in-memory ride store.
- Added health and readiness endpoints.
- Added endpoints to create and list rides.
- Added focused unit tests for health, readiness, ride creation, validation, and ride listing.

### Commands I Used

```bash
cd services/api
go mod init rides-api
gofmt -w cmd/rides-api/main.go internal/httpapi/server.go internal/httpapi/server_test.go internal/ride/store.go
go test ./...
go run ./cmd/rides-api
curl -s -i http://127.0.0.1:8080/healthz
curl -s -i http://127.0.0.1:8080/readyz
curl -s -i http://127.0.0.1:8080/api/rides
curl -s -i -X POST http://127.0.0.1:8080/api/rides \
  -H 'Content-Type: application/json' \
  -d '{"rider_id":"rider-1","pickup":"zone-a","dropoff":"zone-b"}'
```

### What Broke

- `go test ./...` initially failed because the Go build cache is outside the workspace sandbox.
- `curl` from the default sandbox could not reach the API process started with elevated execution.

### How I Fixed It

- Reran `go test ./...` with permission to use the Go build cache.
- Ran endpoint verification in the same execution context as the local API process.

### What I Learned

- The first API can stay simple with Go's standard `net/http` package.
- `GET /healthz` should confirm the process is alive.
- `GET /readyz` currently returns ready, but later it should check PostgreSQL.
- The ride store is intentionally in-memory for Stage 1.
- Unit tests give a stable base before adding config, metrics, and PostgreSQL.

### Questions

- Should Stage 2 keep using only the standard library for logging, or introduce structured logging with `slog`?
- Should ride IDs stay simple strings until PostgreSQL, or move to UUIDs before persistence?

## Stage 2: Configuration and Logging

Date: 2026-05-14

### What I Built

- Added `internal/config` for environment-based configuration.
- Added support for:
  - `PORT`
  - `LOG_LEVEL`
  - `DATABASE_URL`
- Switched the API startup logs to structured JSON logs with `log/slog`.
- Added request logging for method, path, status, and duration.
- Added ride creation logs.
- Added tests for config defaults and environment overrides.

### Commands I Used

```bash
cd services/api
gofmt -w cmd/rides-api/main.go internal/config/config.go internal/config/config_test.go internal/httpapi/server.go internal/httpapi/server_test.go
go test ./...
PORT=9090 LOG_LEVEL=debug DATABASE_URL=postgres://example go run ./cmd/rides-api
curl -s -i http://127.0.0.1:9090/healthz
curl -s -i -X POST http://127.0.0.1:9090/api/rides \
  -H 'Content-Type: application/json' \
  -d '{"rider_id":"rider-2","pickup":"zone-c","dropoff":"zone-d"}'
```

### What Broke

- Nothing in the application code.
- Local verification still needed elevated execution because Go writes to its build cache outside the workspace and the test API process was started in that same execution context.

### How I Fixed It

- Ran tests and manual endpoint checks with the required execution permission.
- Stopped the temporary API process after verification.

### What I Learned

- `PORT` lets the service run on a configurable address instead of hardcoding `:8080`.
- `LOG_LEVEL` controls how verbose structured logs should be.
- `DATABASE_URL` is loaded now but will not be used until the PostgreSQL stage.
- `slog` gives structured logs that are easier to search in Docker and Kubernetes.
- Request logs are useful later for debugging failures in Minikube.

### Questions

- Should invalid `LOG_LEVEL` values fail startup, or default to `info`?
- Should request logs include remote address and user agent in a later stage?
