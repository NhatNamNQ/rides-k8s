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

## Stage 6: Prometheus Locally

Date: 2026-05-16

### Summary

This stage adds Prometheus as a local monitoring service. The Go API already exposed a `GET /metrics` endpoint in Stage 3, but until now nothing was collecting those metrics over time.

Prometheus is the first tool in the project that regularly scrapes the API and stores metric values so we can query them.

### What Changed

- Added `observability/prometheus/prometheus.yml`.
- Added a `prometheus` service to `deploy/docker-compose/docker-compose.yml`.
- Configured Prometheus to scrape:
  - itself at `prometheus:9090`
  - the Go API at `api:8080`

### Concepts Learned

- **Prometheus**: a monitoring system that pulls metrics from HTTP endpoints.
- **Scrape target**: a service Prometheus tries to read metrics from.
- **`up` metric**: built-in Prometheus metric that shows whether a scrape succeeded.
- **`prometheus.yml`**: the configuration file that tells Prometheus what to scrape.
- **Compose service DNS**: inside Docker Compose, `api` and `prometheus` work like hostnames.

### How I Verified It

- Verified the Docker Compose file parses successfully with the Prometheus service added.
- Checked that the Prometheus config points to `api:8080`, which matches the Compose service name.

The next manual runtime check is:

- `docker compose -f deploy/docker-compose/docker-compose.yml up --build`
- open `http://localhost:9090`
- query `up`
- query `rides_created_total`

### Problems and Fixes

- No new code issue appeared in the Go API because Stage 6 only adds local observability infrastructure.
- The main beginner pitfall here is using `localhost` as the API target inside Prometheus config. That would be wrong because Prometheus runs in its own container. The correct target is `api:8080`.

### Beginner Notes

Think of Stage 6 like this:

```text
Before Stage 6:
The API exposes metrics, but nobody collects them.

After Stage 6:
The API exposes metrics.
Prometheus visits /metrics on a schedule and stores what it sees.
```

That is why Prometheus comes before Grafana. First we need to collect data. Later Grafana will read from Prometheus and turn the numbers into dashboards.

## Stage 7: Grafana Locally

Date: 2026-05-16

### Summary

This stage adds Grafana on top of Prometheus so the collected metrics can be viewed as a dashboard instead of only through Prometheus queries. The main learning goal is to understand the difference between collecting metrics and visualizing them.

### What Changed

- Added a Grafana service to `deploy/docker-compose/docker-compose.yml`.
- Added Grafana provisioning files for:
  - a Prometheus datasource
  - a dashboard provider
- Added a starter dashboard at `observability/grafana/dashboards/rides-overview.json`.

### Concepts Learned

- **Grafana datasource**: the system Grafana reads metrics from.
- **Provisioning**: automatic setup of datasources and dashboards when Grafana starts.
- **Dashboard panels**: visual blocks that show one metric or one operational question.
- **PromQL in Grafana**: the same Prometheus query language, but used inside dashboard panels.

### How I Verified It

- Verified the Compose file now includes Grafana on port `3000`.
- Verified Grafana is configured to read from Prometheus at `http://prometheus:9090`.
- Added a dashboard that tracks:
  - request rate
  - p95 latency
  - active rides
  - rides created
  - error rate

### Problems and Fixes

- No backend code change was needed.
- The main beginner trap here is expecting Grafana to collect metrics by itself. It does not. Grafana only visualizes data that Prometheus already collected.

### Beginner Notes

Think of the stack like this:

```text
Go API -> exposes /metrics
Prometheus -> collects and stores the metrics
Grafana -> turns the metrics into charts
```

This is the point where monitoring starts to feel useful. Prometheus shows that the app is being scraped, but Grafana makes the system easier to read at a glance.

## Stage 8: Simulator

Date: 2026-05-18

### Summary

This stage adds a simple simulator service that continuously creates rides through the Go API. The goal is to generate realistic traffic so Prometheus and Grafana show changing behavior instead of mostly idle metrics.

### What Changed

- Added a new Go service at `services/simulator`.
- Added a background loop that sends `POST /api/rides` requests to the API.
- Added simulator metrics and a `GET /metrics` endpoint.
- Added the simulator to Docker Compose.
- Added the simulator as a Prometheus scrape target.
- Added a Grafana panel for simulator create-ride events.

### Concepts Learned

- **Background loop**: a long-running process that performs work on a schedule.
- **Service-to-service HTTP**: one service calling another by service name inside Docker Compose.
- **Load generation**: producing traffic automatically instead of clicking endpoints manually.
- **Simulator metrics**: metrics for the service that creates traffic, not only the service that receives traffic.

### How I Verified It

- Added a simulator unit test for successful ride creation requests.
- Verified the Grafana dashboard JSON still parses.
- Updated Prometheus to scrape both the API and the simulator.
- Updated Docker Compose to run the simulator with `API_BASE_URL=http://api:8080`.

The next runtime check for this stage is:

- `docker compose -f deploy/docker-compose/docker-compose.yml up --build`
- open `http://localhost:3000`
- confirm rides are created without manual requests
- confirm the simulator panel moves in Grafana

### Problems and Fixes

- The simulator needed its own metrics instead of reusing API metrics, otherwise Grafana would only show the effect on the API and not the source of the load.
- The simulator config needed a separate `LogLevelRaw` field so the parsed log-level helper would not collide with a field name.

### Beginner Notes

Think of the simulator like a robot user:

```text
simulator -> sends ride requests -> api -> writes data -> Prometheus scrapes both services
```

This stage matters because a monitoring stack is much easier to understand when something in the system is actually changing on its own.

## Stage 9: Kubernetes Deployment with Minikube

Date: 2026-05-18

### Summary

This stage prepares the first Kubernetes deployment for the API. Instead of moving PostgreSQL into the cluster, this version keeps the existing Supabase database and focuses the Kubernetes work on the API itself.

That keeps the first Minikube step smaller. The learning goal is to understand the basic Kubernetes objects needed to run the API: a Deployment, a Service, a ConfigMap, and a Secret.

### What Changed

- Added `deploy/k8s/configmap.yaml` for non-secret API settings.
- Added `deploy/k8s/secret.yaml` for `DATABASE_URL`.
- Added `deploy/k8s/api-deployment.yaml` for the `rides-api` Deployment.
- Added `deploy/k8s/api-service.yaml` for the in-cluster API Service.
- Updated `docs/architecture.md` to describe the Stage 9 Minikube setup.

### Concepts Learned

- **Deployment**: tells Kubernetes how many API Pods to run and which container image to use.
- **Service**: gives the API a stable in-cluster name even if Pods are replaced.
- **ConfigMap**: stores non-secret environment variables such as `PORT` and `LOG_LEVEL`.
- **Secret**: stores sensitive values such as the database connection string.
- **External database from Kubernetes**: the app can run in Minikube while still talking to a managed database outside the cluster.

### How I Verified It

- Reviewed the API config code to confirm it reads `PORT`, `LOG_LEVEL`, and `DATABASE_URL`.
- Matched the Kubernetes image tag to the Stage Guide build flow: `rides-api:dev`.
- Checked that the Service selector matches the Deployment Pod labels.
- Confirmed the manifest set stays scoped to the API only, leaving probes for Stage 10.

Manual runtime verification is still required for this stage:

- `minikube start`
- `eval $(minikube docker-env)`
- `docker build -t rides-api:dev services/api`
- `kubectl apply -f deploy/k8s`
- `kubectl get pods`
- `kubectl get svc`
- `kubectl port-forward svc/rides-api 8080:8080`
- `curl localhost:8080/healthz`

### Problems and Fixes

- Problem: the original Stage 9 guide assumes PostgreSQL also runs inside Kubernetes.
- Fix: for this repository, Stage 9 uses the existing Supabase database so the first Minikube deployment can stay focused on core Kubernetes objects.

### Beginner Notes

Think of this stage like this:

```text
Docker image already exists
Kubernetes decides how to run it
ConfigMap gives normal settings
Secret gives the database URL
Service gives the Pod a stable network name
```

This is the first point where the API stops being only a local process or a Docker Compose container and starts being a Kubernetes workload.

## Stage 10: Kubernetes Health Checks

Date: 2026-05-19

### Summary

This stage teaches Kubernetes how to judge whether the API process is alive and whether it is actually ready to receive traffic. The key idea is that those are not the same thing.

### What Changed

- Updated `deploy/k8s/api-deployment.yaml` to add a liveness probe.
- Updated `deploy/k8s/api-deployment.yaml` to add a readiness probe.
- Named the API container port `http` so the probes can reference it clearly.
- Updated `docs/architecture.md` to describe the Stage 10 probe behavior.
- Updated `STAGE_GUIDE.md` with a repository-specific readiness failure exercise.

### Concepts Learned

- **Liveness probe**: asks whether the process is still healthy enough to keep running.
- **Readiness probe**: asks whether the app should receive traffic right now.
- **Probe path choice**: `/healthz` should stay lightweight, while `/readyz` can check important dependencies like PostgreSQL.
- **Restart loop risk**: if a readiness-style dependency check is put into liveness, Kubernetes may keep restarting a healthy process for the wrong reason.

### How I Verified It

- Reviewed the Go API handlers to confirm `/healthz` returns process health and `/readyz` checks `store.Ping(...)`.
- Validated that the Deployment now points liveness to `/healthz` and readiness to `/readyz`.
- Kept the probe setup scoped to the existing Supabase-based Kubernetes layout rather than introducing an in-cluster database.

Manual runtime verification for this stage:

- `kubectl apply -f deploy/k8s`
- `kubectl describe pod <pod-name>`
- `kubectl port-forward svc/rides-api 8080:8080`
- `curl localhost:8080/healthz`
- `curl localhost:8080/readyz`
- temporarily break `DATABASE_URL` in `deploy/k8s/secret.yaml`
- re-apply and observe readiness fail without the container being killed for liveness

### Problems and Fixes

- Problem: Kubernetes needs different answers for “is the process alive?” and “can this Pod receive traffic?”
- Fix: point liveness to `/healthz` and readiness to `/readyz` so database issues affect traffic routing, not unnecessary restarts.

### Beginner Notes

Think of it like this:

```text
/healthz -> "Is the app process still alive?"
/readyz  -> "Is the app ready to serve real requests right now?"
```

If Supabase has a temporary problem, the API process may still be running. In that case, Kubernetes should usually stop sending it traffic, not immediately restart it.

## Stage 11: Prometheus and Grafana in Minikube

Date: 2026-05-20

### Summary

This stage moves monitoring into Minikube so Prometheus and Grafana run inside the cluster instead of only in Docker Compose. The first pass stays deliberately small: monitor the API from inside Kubernetes and leave the simulator out of the cluster for now.

### What Changed

- Added `deploy/k8s/prometheus-configmap.yaml` with a simple in-cluster Prometheus config.
- Added `deploy/k8s/prometheus-deployment.yaml` and `deploy/k8s/prometheus-service.yaml`.
- Added `deploy/k8s/grafana-datasource-configmap.yaml` for datasource and dashboard-provider provisioning.
- Added `deploy/k8s/grafana-dashboard-configmap.yaml` with an API-only starter dashboard.
- Added `deploy/k8s/grafana-deployment.yaml` and `deploy/k8s/grafana-service.yaml`.
- Updated `docs/architecture.md` and `STAGE_GUIDE.md` to reflect the API-first in-cluster monitoring path.

### Concepts Learned

- **Monitoring inside Kubernetes**: Prometheus can scrape Services from inside the cluster, not only containers in Docker Compose.
- **ClusterIP Service**: Prometheus and Grafana do not need public access inside Minikube; they can be reached by service name.
- **Grafana provisioning from ConfigMaps**: dashboards and datasources can be created automatically at pod startup.
- **API-only scope**: keeping the first cluster-monitoring step narrow makes it easier to debug the stack.

### How I Verified It

- Reused the Prometheus and Grafana local configuration patterns and adapted them to Kubernetes service names.
- Kept the in-cluster Prometheus scrape target scoped to `rides-api:8080/metrics`.
- Built the Grafana dashboard to query only `job="rides-api"` metrics so it matches the Stage 11 scope.
- Parsed the dashboard JSON successfully before embedding it into the ConfigMap.

Manual runtime verification for this stage:

- `kubectl apply -f deploy/k8s`
- `kubectl get pods`
- `kubectl get svc`
- `kubectl logs deploy/prometheus`
- `kubectl logs deploy/grafana`
- `kubectl port-forward svc/grafana 3000:3000`
- open `http://localhost:3000`
- generate API traffic and confirm the dashboard moves

### Problems and Fixes

- Problem: the local Grafana dashboard included simulator panels, but the simulator is not in Minikube yet.
- Fix: the in-cluster dashboard is API-only and filters queries with `job="rides-api"`.

### Beginner Notes

Think of this stage like this:

```text
rides-api runs in Minikube
Prometheus runs in Minikube and scrapes rides-api
Grafana runs in Minikube and reads from Prometheus
your browser reaches Grafana through kubectl port-forward
```

This stage is where monitoring becomes part of the cluster itself, not just an external local tool.

## Stage 11.5: Helm for Monitoring

Date: 2026-05-21

### Summary

This mini-stage repeats the monitoring setup with Helm so the same Prometheus and Grafana ideas can be managed as chart releases instead of only raw Kubernetes YAML.

The goal was not to replace the earlier learning work. The goal was to compare two deployment styles:

- raw manifests in `deploy/k8s`
- Helm values in `deploy/helm/monitoring`

### What Changed

- Added `deploy/helm/monitoring/prometheus-values.yaml`.
- Added `deploy/helm/monitoring/grafana-values.yaml`.
- Configured the Prometheus Helm chart to scrape `rides-api.default.svc.cluster.local:8080/metrics`.
- Configured the Grafana Helm chart to provision Prometheus as the default datasource.
- Added dashboard provisioning in Helm values so Grafana loads a `Rides Overview` dashboard automatically.
- Aligned the Grafana datasource UID with the dashboard JSON so the dashboard panels know which datasource to use.

### Concepts Learned

- **Helm repo**: an online source of Helm charts. It is similar to an app store for Kubernetes packages.
- **Helm chart**: a reusable Kubernetes application template that creates objects like Deployments, Services, Secrets, and ConfigMaps.
- **Helm release**: one installed instance of a chart in a namespace.
- **Values file**: a YAML file that changes chart behavior without rewriting the chart itself.
- **Provisioned datasource**: Grafana can create the Prometheus connection automatically when the pod starts.
- **Provisioned dashboard**: Grafana can load a dashboard from configuration instead of requiring manual UI setup every time.
- **Service DNS in Kubernetes**: Prometheus can scrape another service by cluster DNS name such as `rides-api.default.svc.cluster.local`.

### How I Verified It

- Confirmed the Helm CLI worked with `helm version`.
- Added and updated the Prometheus and Grafana chart repositories.
- Installed Prometheus and Grafana into the `monitoring` namespace with Helm.
- Validated `deploy/helm/monitoring/grafana-values.yaml` as YAML.
- Rendered the Grafana chart with `helm template` to confirm the values file was accepted.
- Verified Grafana could be opened with `kubectl port-forward -n monitoring svc/grafana 3000:80`.
- Verified the Helm values now define:
  - a Prometheus datasource
  - a dashboard provider
  - a `Rides Overview` dashboard for API metrics

Manual runtime verification for this stage:

- `helm list -n monitoring`
- `helm status grafana -n monitoring`
- `helm upgrade prometheus prometheus-community/prometheus -n monitoring -f deploy/helm/monitoring/prometheus-values.yaml`
- `helm upgrade grafana grafana/grafana -n monitoring -f deploy/helm/monitoring/grafana-values.yaml`
- `kubectl port-forward -n monitoring svc/prometheus-server 9090:80`
- `kubectl port-forward -n monitoring svc/grafana 3000:80`
- query `up{job="rides-api"}` in Prometheus
- query `http_requests_total{job="rides-api"}` in Prometheus
- open Grafana and confirm the dashboard appears

### Problems and Fixes

- Problem: Homebrew failed while trying to build Helm because the machine ran out of local build space.
- Fix: installed Helm manually from the official binary release instead of depending on the Homebrew build.
- Problem: Grafana initially showed an empty dashboard page.
- Fix: the Helm values only provisioned a datasource at first, so dashboard provisioning was added explicitly.
- Problem: the dashboard JSON expected datasource UID `prometheus`, but the Grafana datasource definition did not set a UID.
- Fix: added `uid: prometheus` to the datasource config so the dashboard panels can bind to the correct datasource.
- Problem: Minikube became unhealthy during Helm setup because the local Docker-backed cluster stopped responding.
- Fix: restarted the Minikube container and continued after the cluster became reachable again.

### Beginner Notes

Think of this stage like this:

```text
raw YAML = I write each Kubernetes object myself
Helm = the chart writes the Kubernetes objects for me
values.yaml = I control the chart without changing the chart source
```

Helm does not replace Kubernetes knowledge. It sits on top of Kubernetes. If a Grafana chart creates a Deployment, Service, Secret, and ConfigMap, those are still normal Kubernetes objects underneath.

This is why it helped to do Stage 11 manually first. After the manual version made sense, the Helm version became much easier to understand.
