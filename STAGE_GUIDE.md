# Rides K8s Stage Guide

This guide turns `LEARNING_PLAN.md` into practical execution steps.

Use this rule throughout the project:

```text
Build locally -> containerize -> deploy to Minikube -> observe -> automate -> move selected parts to AWS.
```

Do not skip the local step. If something fails in Kubernetes, you need to know whether the bug is in the application, the container, the manifest, the cluster, or the cloud configuration.

## Stage 0: Plan the System

### Objective

Define the first version of the system before writing code.

### Tasks

1. Create project folders:

```bash
mkdir -p docs services/api services/simulator services/web deploy/k8s deploy/docker-compose observability/prometheus observability/grafana infra/terraform
```

2. Create `docs/architecture.md`.

Include:

```text
Services:
- api: Go HTTP API
- simulator: fake ride generator
- postgres: data store
- prometheus: metrics scraper
- grafana: dashboard

Main flow:
simulator -> api -> postgres
prometheus -> api /metrics
grafana -> prometheus
```

3. Create `docs/journal.md`.

Use this template after every stage:

```text
Date:
Stage:
Summary:
What Changed:
Concepts Learned:
How I Verified It:
Problems and Fixes:
Beginner Notes:
```

### Checkpoint

You can explain the system in 2 minutes without reading notes.

## Stage 1: First Go HTTP API

### Objective

Build the smallest useful Go backend.

### Tasks

1. Initialize Go module:

```bash
cd services/api
go mod init rides-api
```

2. Create the Go API using this small project-layout style:

```text
services/api/
  cmd/
    rides-api/
      main.go
  internal/
    httpapi/
      server.go
      server_test.go
    ride/
      store.go
```

Implement:

- `GET /healthz`
- `GET /readyz`
- `GET /api/rides`
- `POST /api/rides`

Start with in-memory storage.

3. Run locally:

```bash
go run ./cmd/rides-api
```

4. Test manually:

```bash
curl localhost:8080/healthz
curl localhost:8080/readyz
curl localhost:8080/api/rides
curl -X POST localhost:8080/api/rides \
  -H 'Content-Type: application/json' \
  -d '{"rider_id":"rider-1","pickup":"zone-a","dropoff":"zone-b"}'
```

### Learn

Focus on:

- `net/http`
- JSON encode/decode
- HTTP status codes
- handler functions
- simple application state

### Checkpoint

You can create a ride and list it back through HTTP.

## Stage 2: Configuration and Logging

### Objective

Make the service configurable like a real backend.

### Tasks

1. Add environment variables:

```text
PORT
LOG_LEVEL
DATABASE_URL
```

2. Run with a custom port:

```bash
PORT=9090 go run ./cmd/rides-api
```

3. Add structured logs for:

- server startup
- request method/path/status
- ride creation
- errors

### Learn

Focus on:

- why config should not be hardcoded
- logs versus metrics
- how logs help during Kubernetes debugging

### Checkpoint

Changing `PORT` changes where the service runs.

## Stage 3: Prometheus Metrics in Go

### Objective

Expose application metrics.

### Tasks

1. Add Prometheus client library:

```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
```

2. Add `/metrics`.

3. Track:

```text
http_requests_total
http_request_duration_seconds
rides_created_total
rides_active
```

4. Test:

```bash
curl localhost:8080/metrics
```

5. Create a ride and check if `rides_created_total` changes.

### Learn

Focus on:

- Counter: only increases
- Gauge: increases and decreases
- Histogram: measures duration distribution
- Prometheus pulls metrics from your service

### Checkpoint

Your API exposes custom ride metrics.

## Stage 4: PostgreSQL

### Objective

Persist data instead of keeping it in memory.

### Tasks

1. Choose where PostgreSQL runs.

For this project, use Supabase free tier if Docker image downloads are too slow.

Recommended for Supabase:

- Create a Supabase project.
- Open the project dashboard.
- Click **Connect**.
- Copy the **Session pooler** connection string first.
- Add `sslmode=require` if it is not already present.
- Store it locally as `DATABASE_URL`.

Why Session pooler:

- this Go API is a long-running backend service
- Supabase direct connections often require IPv6
- Session pooler supports IPv4 and IPv6

Avoid Transaction pooler for now because it can require extra driver settings around prepared statements.

If your network cannot reach port `5432`, use the Supabase pooler on port `6543` instead:

```text
postgresql://postgres.PROJECT_REF:YOUR_PASSWORD@aws-1-ap-southeast-2.pooler.supabase.com:6543/postgres?sslmode=require&default_query_exec_mode=simple_protocol
```

Use `default_query_exec_mode=simple_protocol` with the Go `pgx` driver when using the transaction pooler.

For `psql`, remove `default_query_exec_mode=simple_protocol` because `psql` does not understand that parameter.

2. Optional local Docker PostgreSQL:

```bash
docker run --name rides-postgres \
  -e POSTGRES_USER=rides \
  -e POSTGRES_PASSWORD=rides \
  -e POSTGRES_DB=rides \
  -p 5432:5432 \
  -d postgres:16
```

3. Add schema file:

```text
services/api/migrations/001_init.sql
```

Start with:

```sql
CREATE SEQUENCE IF NOT EXISTS rides_id_seq;

CREATE TABLE IF NOT EXISTS rides (
  id TEXT PRIMARY KEY DEFAULT ('ride-' || nextval('rides_id_seq')::TEXT),
  rider_id TEXT NOT NULL,
  pickup TEXT NOT NULL,
  dropoff TEXT NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

4. Apply the migration.

For Supabase, use the dashboard SQL Editor:

- open SQL Editor
- paste `services/api/migrations/001_init.sql`
- run the SQL

5. Connect from Go using `DATABASE_URL`.

Example:

```bash
cd services/api
DATABASE_URL='postgres://postgres.PROJECT_REF:YOUR_PASSWORD@aws-0-REGION.pooler.supabase.com:6543/postgres?sslmode=require&default_query_exec_mode=simple_protocol' go run ./cmd/rides-api
```

6. Make `/readyz` check database connectivity.

### Learn

Focus on:

- database connection strings
- `database/sql`
- migrations
- readiness checks
- persistent state

### Checkpoint

- The API starts with `DATABASE_URL`.
- `/readyz` returns ready when Supabase is reachable.
- Created rides are stored in Supabase.
- Rides survive an API restart.

## Stage 5: Docker Compose

### Objective

Run the app and database in containers.

### Tasks

1. Add `services/api/Dockerfile`.

Recommended shape:

```text
multi-stage build
  builder stage -> compile Go binary
  runtime stage -> small final image
```

2. Add `deploy/docker-compose/docker-compose.yml` with:

- api
- postgres

Recommended local `DATABASE_URL` in Compose:

```text
postgres://rides:rides@postgres:5432/rides?sslmode=disable
```

Use the Docker Compose service name `postgres` as the hostname.

Let Postgres apply the migration automatically by mounting:

```text
services/api/migrations/001_init.sql
```

into:

```text
/docker-entrypoint-initdb.d/
```

3. Run:

```bash
docker compose -f deploy/docker-compose/docker-compose.yml up --build
```

4. Test:

```bash
curl localhost:8080/healthz
curl localhost:8080/api/rides
```

### Learn

Focus on:

- image build
- container runtime config
- service names as DNS names
- container logs
- volumes
- multi-stage Docker builds
- why `postgres` works as a hostname inside Docker Compose

### Checkpoint

The API container connects to the Postgres container.
Postgres applies `001_init.sql` automatically on first start.

## Stage 6: Prometheus Locally

### Objective

Run Prometheus and scrape the Go API.

### Tasks

1. Add `observability/prometheus/prometheus.yml`.

Example targets:

```yaml
scrape_configs:
  - job_name: rides-api
    static_configs:
      - targets: ["api:8080"]
```

2. Add Prometheus to Docker Compose.

3. Open Prometheus:

```text
http://localhost:9090
```

4. Query:

```text
up
rides_created_total
http_requests_total
```

### Learn

Focus on:

- scrape targets
- `up`
- PromQL basics
- why labels matter

### Checkpoint

Prometheus target for the API is `UP`.

## Stage 7: Grafana Locally

### Objective

Create a dashboard from Prometheus metrics.

### Tasks

1. Add Grafana to Docker Compose.

2. Open Grafana:

```text
http://localhost:3000
```

3. Add Prometheus data source:

```text
http://prometheus:9090
```

4. Create panels:

- request rate
- p95 latency
- active rides
- rides created
- error rate

### Useful PromQL

```promql
rate(http_requests_total[5m])
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
rides_active
increase(rides_created_total[10m])
```

### Checkpoint

Creating rides changes your Grafana panels.

## Stage 8: Simulator

### Objective

Generate realistic traffic for the API.

### Tasks

1. Create `services/simulator`.

2. Start simple. The simulator can be Go first.

It should:

- call `POST /api/rides`
- sleep between requests
- generate random pickup/dropoff zones
- expose `/metrics`

3. Add simulator to Docker Compose.

4. Add simulator target to Prometheus.

### Learn

Focus on:

- background loops
- service-to-service HTTP
- backoff on errors
- load generation
- simulator metrics

### Checkpoint

Rides are created automatically and visible in Grafana.

## Stage 9: Kubernetes Deployment with Minikube

### Objective

Deploy the local system to Kubernetes.

### Tasks

1. Start Minikube:

```bash
minikube start
```

2. Build image for Minikube:

```bash
eval $(minikube docker-env)
docker build -t rides-api:dev services/api
```

3. Create Kubernetes manifests:

```text
deploy/k8s/api-deployment.yaml
deploy/k8s/api-service.yaml
deploy/k8s/postgres-deployment.yaml
deploy/k8s/postgres-service.yaml
deploy/k8s/configmap.yaml
deploy/k8s/secret.yaml
```

4. Apply:

```bash
kubectl apply -f deploy/k8s
```

5. Inspect:

```bash
kubectl get pods
kubectl get svc
kubectl logs deploy/rides-api
```

6. Access API:

```bash
kubectl port-forward svc/rides-api 8080:8080
curl localhost:8080/healthz
```

### Learn

Focus on:

- Pod
- Deployment
- Service
- ConfigMap
- Secret
- labels and selectors
- port-forwarding

### Checkpoint

The API runs in Minikube and connects to Postgres in the cluster.

## Stage 10: Kubernetes Health Checks

### Objective

Teach Kubernetes when the app is alive and ready.

### Tasks

1. Add liveness probe:

```text
/healthz
```

2. Add readiness probe:

```text
/readyz
```

3. Apply manifests:

```bash
kubectl apply -f deploy/k8s
kubectl describe pod <pod-name>
```

4. Break database config temporarily and observe readiness failure.

### Learn

Focus on:

- liveness means process health
- readiness means can receive traffic
- bad probes cause restart loops

### Checkpoint

You can explain why `/readyz` checks Postgres and `/healthz` does not.

## Stage 11: Prometheus and Grafana in Minikube

### Objective

Run monitoring inside Kubernetes.

### Tasks

1. Deploy Prometheus with a simple config first.

2. Deploy Grafana.

3. Port-forward Grafana:

```bash
kubectl port-forward svc/grafana 3000:3000
```

4. Add dashboards for API and simulator.

### Later Option

After manual manifests work, install the real-world stack:

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install monitoring prometheus-community/kube-prometheus-stack
```

### Checkpoint

Grafana running inside Minikube shows API metrics.

## Stage 12: Scaling and Rollouts

### Objective

Practice Kubernetes operations.

### Tasks

1. Scale API:

```bash
kubectl scale deployment rides-api --replicas=2
kubectl get pods
```

2. Scale simulator:

```bash
kubectl scale deployment rides-simulator --replicas=3
```

3. Delete a pod:

```bash
kubectl delete pod <pod-name>
kubectl get pods -w
```

4. Check rollout:

```bash
kubectl rollout status deployment/rides-api
kubectl rollout history deployment/rides-api
```

### Learn

Focus on:

- desired state
- replicas
- rollout status
- self-healing
- resource requests and limits

### Checkpoint

You can show Kubernetes replacing a deleted pod.

## Stage 13: React Frontend

### Objective

Create a simple UI for the system.

### Tasks

1. Create React + TypeScript app in `services/web`.

2. Add views:

- active rides
- recent rides
- API status
- basic metrics summary

3. Containerize the frontend.

4. Deploy it to Minikube.

### Learn

Focus on:

- API calls from frontend
- environment-based API URL
- frontend Dockerfile
- frontend Kubernetes service

### Checkpoint

The frontend runs in Minikube and displays API data.

## Stage 14: Matching Logic

### Objective

Make the domain logic more realistic.

### Tasks

Add:

- drivers
- driver locations
- ride states
- nearest-driver assignment
- ride event history

Suggested ride states:

```text
requested
assigned
picked_up
completed
cancelled
```

### Learn

Focus on:

- state machines
- database transactions
- race conditions
- business metrics

### Checkpoint

A ride can move through multiple states and every state change is recorded.

## Stage 15: Failure and Debugging Scenarios

### Objective

Practice diagnosing real operational failures.

### Tasks

Create and document these failures:

1. Wrong database password
2. API image does not exist
3. Readiness probe fails
4. Simulator overloads API
5. Memory limit too low
6. Slow database query

For each failure, run:

```bash
kubectl get pods
kubectl describe pod <pod-name>
kubectl logs <pod-name>
kubectl get events --sort-by=.lastTimestamp
```

### Learn

Focus on:

- symptoms
- root cause
- fix
- prevention

### Checkpoint

You have one written debugging note for each failure.

## Stage 16: Production-Style Improvements

### Objective

Add real-world practices after the core system works.

### Options

Pick only one or two at first:

- Helm chart
- database migrations
- OpenTelemetry tracing
- Loki logs
- Horizontal Pod Autoscaler
- k6 load testing
- Ingress

### Checkpoint

You can explain why the improvement was needed and what complexity it added.

## Stage 17: Jenkins CI/CD

### Objective

Automate test, build, and deploy.

### Tasks

1. Create `Jenkinsfile`.

Start with stages:

```text
Checkout
Test
Build Image
Deploy to Minikube
Verify
```

2. Pipeline commands:

```bash
go test ./...
docker build -t rides-api:$BUILD_NUMBER services/api
kubectl apply -f deploy/k8s
kubectl rollout status deployment/rides-api
```

3. Add verification:

```bash
curl localhost:8080/healthz
curl localhost:8080/readyz
curl localhost:8080/metrics
```

### Learn

Focus on:

- pipeline syntax
- build failure versus deploy failure
- credentials
- rollout verification
- reproducible builds

### Checkpoint

A failing Go test stops deployment.

## Stage 18: AWS Integration from Minikube

### Objective

Use AWS services while keeping Kubernetes local.

### Tasks

1. Create ECR repo manually first.

2. Push image:

```bash
aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <account-id>.dkr.ecr.<region>.amazonaws.com
docker tag rides-api:dev <account-id>.dkr.ecr.<region>.amazonaws.com/rides-api:dev
docker push <account-id>.dkr.ecr.<region>.amazonaws.com/rides-api:dev
```

3. Add S3 ride report export endpoint.

4. Add one config value from SSM Parameter Store.

5. Add one secret from Secrets Manager.

### Learn

Focus on:

- IAM permissions
- ECR login
- AWS SDK for Go
- S3 object writes
- external config

### Checkpoint

Minikube deploys an image from ECR and the API writes one object to S3.

## Stage 19: Terraform

### Objective

Create AWS resources as code.

### Tasks

1. Create Terraform folder:

```text
infra/terraform/environments/dev
```

2. Define:

- ECR repository
- S3 bucket
- SSM parameter
- Secrets Manager secret
- IAM policy

3. Run:

```bash
terraform init
terraform fmt
terraform validate
terraform plan
terraform apply
```

4. Destroy test resources when done:

```bash
terraform destroy
```

### Learn

Focus on:

- provider
- resource
- variable
- output
- state
- plan versus apply
- destroy

### Checkpoint

Terraform can recreate the AWS resources you previously made manually.

## Stage 20: EKS

### Objective

Move from local Kubernetes to AWS-managed Kubernetes.

### Tasks

1. Create EKS cluster.

Begin with `eksctl` if you want a simpler first experience, then move to Terraform later.

2. Push app images to ECR.

3. Deploy manifests to EKS.

4. Use RDS PostgreSQL or temporary in-cluster Postgres.

5. Install monitoring.

6. Delete resources when finished.

### Learn

Focus on:

- EKS
- node groups
- ECR image pulls
- AWS Load Balancer Controller
- RDS networking
- cost cleanup

### Checkpoint

The app runs on EKS and you know how to delete everything safely.

## How to Work Through Each Stage

For every stage:

1. Read the objective.
2. Create or edit only the files required for that stage.
3. Run the smallest local test first.
4. Commit working code.
5. Break one thing intentionally.
6. Diagnose it.
7. Write the result in `docs/journal.md`.

Suggested commit style:

```text
stage-01: add initial Go API
stage-02: add config and logging
stage-03: expose prometheus metrics
```

## When to Ask for Help

Ask for help with one concrete failure at a time. Include:

- command you ran
- exact error message
- files you changed
- expected result
- actual result

That makes debugging much faster.
