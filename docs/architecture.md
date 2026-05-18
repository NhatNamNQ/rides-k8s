# Rides K8s Architecture

## Purpose

This project is a learning system for Go, Kubernetes, Prometheus, Grafana, Jenkins, AWS, and Terraform.

The application simulates a small ridesharing platform. It is intentionally smaller than a real production system, but it should be realistic enough to practice backend design, deployment, monitoring, and debugging.

## Inspiration

This project is inspired by:

- Juraj Majerik's Rides project: https://jurajmajerik.com/
- Pragmatic Engineer educational side project article: https://blog.pragmaticengineer.com/an-educational-side-project/
- Go Secure Auth Pro: https://github.com/fdhhhdjd/Go_Secure_Auth_Pro

Useful ideas from `Go_Secure_Auth_Pro`:

- Clear Go project structure
- Separate application entrypoints
- Docker Compose for local development
- PostgreSQL as the main database
- Migrations folder
- Configuration through environment files
- Security and authentication as explicit concerns
- API documentation through Swagger or similar tooling

For this project, auth/security will be added gradually. The first goal is to build a working Go API, then deploy and observe it.

## First-Version System

```text
simulator -> api -> postgres
prometheus -> api /metrics
grafana -> prometheus
```

Current local Compose stack after Stage 6:

```text
api
simulator
prometheus
grafana
```

Later:

```text
jenkins -> docker image -> minikube
terraform -> aws resources
api -> aws services
eks -> production-like kubernetes deployment
```

## Services

## Go API

Path:

```text
services/api
```

Current structure:

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

This follows the useful parts of `golang-standards/project-layout`: keep the executable entrypoint in `cmd/` and private application code in `internal/`.

Responsibilities:

- expose HTTP API
- manage rides
- manage ride state
- store data in PostgreSQL
- expose health checks
- expose Prometheus metrics

Initial endpoints:

```text
GET  /healthz
GET  /readyz
GET  /metrics
GET  /api/rides
POST /api/rides
```

Later endpoints:

```text
GET  /api/drivers
POST /api/drivers
POST /api/drivers/{id}/location
POST /api/rides/{id}/assign
POST /api/rides/{id}/complete
POST /api/rides/{id}/cancel
```

## Simulator

Path:

```text
services/simulator
```

Responsibilities:

- generate fake riders
- generate fake ride requests
- simulate driver movement
- simulate cancellations and completions
- create load against the Go API
- expose simulator metrics

Current local setup:

- runs as its own Go service in Docker Compose
- calls the API at `http://api:8080`
- exposes `/metrics` for Prometheus scraping at `simulator:8081`

The simulator can start in Go for simplicity. Node.js can be introduced later if practicing Node is still a goal.

## PostgreSQL

Responsibilities:

- store rides
- store drivers
- store riders
- store ride events
- provide persistent state across API restarts

Initial tables:

```text
rides
```

Later tables:

```text
drivers
riders
ride_events
zones
```

## Prometheus

Responsibilities:

- scrape `/metrics` from the API
- scrape `/metrics` from the simulator
- provide queryable operational metrics

Current local setup:

- runs in Docker Compose
- scrapes itself at `prometheus:9090`
- scrapes the API at `api:8080`

Initial metrics:

```text
http_requests_total
http_request_duration_seconds
rides_created_total
rides_active
```

Later metrics:

```text
rides_completed_total
rides_cancelled_total
dispatch_duration_seconds
simulator_events_total
db_query_duration_seconds
```

## Grafana

Responsibilities:

- display operational dashboards
- show API health
- show ride activity
- show simulator behavior
- show latency and error trends

First dashboard panels:

- request rate
- p95 latency
- active rides
- rides created
- error rate

Current local setup:

- runs in Docker Compose
- reads from Prometheus at `http://prometheus:9090`
- loads a provisioned dashboard on startup

## React Frontend

Path:

```text
services/web
```

Responsibilities:

- display active rides
- display recent ride events
- display API status
- display simple operational state

The frontend comes later. It should not block learning Go and Kubernetes.

## Jenkins

Responsibilities:

- run tests
- build Docker images
- deploy to Minikube
- verify rollout status
- later push images to ECR

Jenkins is added after manual Minikube deployment works.

## AWS

First AWS services:

- ECR for container images
- S3 for ride report exports
- SSM Parameter Store for non-secret configuration
- Secrets Manager for secrets

Later AWS services:

- EKS for managed Kubernetes
- RDS PostgreSQL for managed database
- CloudWatch for cloud-side logs and metrics experiments

AWS should not be introduced before the local and Minikube versions work.

## Terraform

Responsibilities:

- define AWS infrastructure as code
- create ECR, S3, SSM, Secrets Manager, and IAM resources
- later create EKS/RDS if needed

Terraform comes after the AWS resources are understood manually.

## Repository Structure

```text
.
├── AGENTS.md
├── LEARNING_PLAN.md
├── REVIEW_AGENTS.md
├── STAGE_GUIDE.md
├── deploy/
│   ├── docker-compose/
│   └── k8s/
├── docs/
│   ├── architecture.md
│   └── journal.md
├── infra/
│   └── terraform/
├── observability/
│   ├── grafana/
│   └── prometheus/
└── services/
    ├── api/
    ├── simulator/
    └── web/
```

## First Learning Milestone

The first milestone is complete when:

- the Go API runs locally
- `/healthz` works
- `/readyz` works
- `/metrics` works
- rides can be created and listed
- basic tests exist
- notes are written in `docs/journal.md`

## Design Decisions

## Start Without Authentication

Authentication is important, and `Go_Secure_Auth_Pro` is useful inspiration for that area. However, this project should not start with JWT, Firebase, Redis, or complex middleware.

Reason:

- the first learning target is Go API basics and Kubernetes deployability
- auth adds complexity before the core service exists

Auth can be added later as a separate learning stage.

## Start With Plain Kubernetes YAML

Helm is useful, but the first Kubernetes version should use plain YAML.

Reason:

- plain manifests teach Deployments, Services, ConfigMaps, Secrets, probes, and labels more directly

## Start With Minikube

Use Minikube before EKS.

Reason:

- Minikube is cheaper
- debugging is faster
- failure recovery is easier
- cloud complexity comes later

## Start With Manual Steps Before Automation

Jenkins and Terraform should automate processes that already work manually.

Reason:

- CI/CD and IaC are easier to debug when the manual workflow is already understood
