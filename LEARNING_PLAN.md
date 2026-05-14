# Rides K8s Learning Plan

This project is inspired by Juraj Majerik's Rides simulation series:

- Project overview: https://jurajmajerik.com/
- Starting post: https://jurajmajerik.com/blog/start-here/
- Pragmatic Engineer write-up: https://blog.pragmaticengineer.com/an-educational-side-project/

Juraj's original project used Node.js for the simulation engine, Go for the backend, PostgreSQL for storage, React and TypeScript for the frontend, and Prometheus/Grafana for observability. This plan adapts that idea for a beginner who wants to learn **Go**, **Kubernetes**, **Prometheus**, and **Grafana**.

The goal is not to copy the original project exactly. The goal is to build a smaller production-style system that gives you useful practice with backend services, containers, Kubernetes deployments, service networking, health checks, metrics, dashboards, and debugging.

For step-by-step execution details, use `STAGE_GUIDE.md`.

## Project Goal

Build a small ridesharing simulation platform:

- A Go API manages drivers, riders, rides, and ride events.
- A simulator generates fake riders and drivers.
- PostgreSQL stores application data.
- Prometheus collects metrics.
- Grafana displays dashboards.
- Kubernetes runs the system locally.

## Recommended Prerequisites

Install these before starting:

- Go
- Docker
- kubectl
- minikube
- PostgreSQL client, optional but useful
- Node.js, only needed after the Go API foundation is working
- AWS CLI, only needed for the AWS stages
- Terraform, only needed for the infrastructure-as-code stages
- eksctl, only needed if you later deploy to Amazon EKS
- Jenkins, only needed for the CI/CD stage

Recommended local Kubernetes choice:

- Use `minikube` because you already have it installed.
- Keep the first Kubernetes version local before using AWS.
- Add AWS only after you can deploy and debug the app in Minikube.

## Learning Rules

Use these rules throughout the project:

1. Keep every stage small enough to finish.
2. Write notes after each stage in `docs/journal.md`.
3. Do not add Kubernetes before the app works locally.
4. Do not add Grafana before Prometheus has useful metrics.
5. Prefer boring, explicit configuration over clever abstraction.
6. Each stage must have a clear "done" condition.

## Stage 0: Plan the System

Inspired by Juraj's "Building an Uber clone" post.

### Goal

Understand the system before writing code.

### Build

Create a short design document:

- What services exist?
- What data is stored?
- What requests does the API handle?
- What metrics should be visible?
- What should the first dashboard show?

### Suggested Files

```text
docs/
  architecture.md
  journal.md
```

### Learn

- Basic system design
- Service boundaries
- API-first thinking
- What should be observable in a backend system

### Done When

- `docs/architecture.md` explains the first version of the system.
- You can describe the role of the Go API, simulator, database, Prometheus, and Grafana.

## Stage 1: Build the First Go HTTP API

Inspired by Juraj's early server setup and HTTP server posts.

### Goal

Learn basic Go backend development before adding infrastructure complexity.

### Build

Create a Go service with:

- `GET /healthz`
- `GET /readyz`
- `GET /api/rides`
- `POST /api/rides`

Use in-memory storage for now.

### Suggested Files

```text
services/
  api/
    go.mod
    cmd/
      rides-api/
        main.go
    internal/
      httpapi/
      ride/
```

### Learn

- Go modules
- HTTP handlers
- JSON request/response handling
- Basic project structure
- Health and readiness endpoints

### Done When

- You can run the API locally with `go run ./cmd/rides-api`.
- `curl localhost:8080/healthz` returns success.
- You can create and list fake rides.

## Stage 2: Add Configuration and Logging

Inspired by Juraj's environment variable step.

### Goal

Learn how services are configured in real deployments.

### Build

Add configuration using environment variables:

- `PORT`
- `LOG_LEVEL`
- `DATABASE_URL`, even if unused at first

Add structured logs for:

- server startup
- incoming requests
- ride creation
- errors

### Learn

- Environment-based config
- Why hardcoded configuration is fragile
- Structured logging basics

### Done When

- The API port can be changed with `PORT=9090`.
- Logs are readable and include useful request information.

## Stage 3: Add Prometheus Metrics to Go

This stage moves observability earlier than Juraj's original project because it is one of your main learning goals.

### Goal

Learn how application metrics are exposed.

### Build

Add `GET /metrics` using the Prometheus Go client.

Track:

- total HTTP requests
- request duration
- rides created
- active rides

Example metric names:

```text
http_requests_total
http_request_duration_seconds
rides_created_total
rides_active
```

### Learn

- Counters
- Gauges
- Histograms
- Prometheus scrape model
- Why `/metrics` is different from logs

### Done When

- `curl localhost:8080/metrics` shows Prometheus metrics.
- Creating a ride changes at least one custom metric.

## Stage 4: Add PostgreSQL

Inspired by Juraj's SQL setup and Go database connection posts.

### Goal

Move from in-memory data to persistent storage.

### Build

Create tables:

```sql
drivers
riders
rides
ride_events
```

Update the Go API to store rides in PostgreSQL.

### Learn

- Basic SQL schema design
- Go database access
- Connection strings
- Database migrations
- Difference between app readiness and process health

### Done When

- The API starts only when the database connection works.
- Created rides survive API restarts.
- `/readyz` fails when the database is unavailable.

## Stage 5: Run Locally with Docker Compose

Inspired by Juraj's Docker and Docker Compose posts.

### Goal

Package the service and database so the system is reproducible.

### Build

Add:

- Dockerfile for the Go API
- `docker-compose.yml`
- PostgreSQL service
- API service

### Learn

- Docker images
- Container networking
- Environment variables in containers
- Volumes for database data
- Difference between build-time and runtime config

### Done When

- `docker compose up` starts the API and PostgreSQL.
- The API connects to Postgres using the Compose service name.
- You can create and list rides through the containerized API.

## Stage 6: Add Prometheus Locally

Inspired by Juraj's monitoring and logging post, but simplified.

### Goal

Learn Prometheus before adding Grafana.

### Build

Add Prometheus to Docker Compose.

Configure it to scrape:

- Go API `/metrics`
- Prometheus itself

### Learn

- `prometheus.yml`
- scrape targets
- labels
- PromQL basics
- checking target health

### Done When

- Prometheus UI opens locally.
- The API target is `UP`.
- You can query `rides_created_total`.

## Stage 7: Add Grafana

Inspired by Juraj's monitoring dashboard.

### Goal

Visualize the system.

### Build

Add Grafana to Docker Compose.

Create dashboard panels for:

- request rate
- p95 request latency
- active rides
- rides created
- error rate

### Learn

- Grafana data sources
- dashboard panels
- PromQL in Grafana
- operational dashboards versus business dashboards

### Done When

- Grafana connects to Prometheus.
- A dashboard shows live API metrics.
- Creating rides changes the dashboard.

## Stage 8: Create a Simple Simulator

Inspired by Juraj's simulation engine posts.

### Goal

Generate load so Kubernetes and monitoring have something realistic to show.

### Build

Create a small simulator service.

Start simple. It can be written in Go first. Add Node.js later if you specifically want to practice Node.

The simulator should:

- create fake riders
- create fake ride requests
- randomly complete or cancel rides
- run continuously
- expose its own `/metrics`

### Learn

- Background workers
- Service-to-service HTTP calls
- Load generation
- Failure simulation
- Metrics for non-API services

### Done When

- The simulator creates rides automatically.
- Prometheus scrapes simulator metrics.
- Grafana shows both API and simulator behavior.

## Stage 9: Move to Kubernetes

This is where the project becomes a Kubernetes learning project.

### Goal

Deploy the working local system into a local Kubernetes cluster.

### Build

Create Kubernetes manifests for:

- Go API `Deployment`
- API `Service`
- simulator `Deployment`
- simulator `Service`, if it exposes metrics
- PostgreSQL `StatefulSet` or simple `Deployment`
- PostgreSQL `Service`
- `ConfigMap` for non-secret config
- `Secret` for database credentials

### Learn

- Pods
- Deployments
- Services
- ConfigMaps
- Secrets
- container ports versus service ports
- labels and selectors

### Done When

- `kubectl get pods` shows all core app pods running.
- The Go API can connect to PostgreSQL inside the cluster.
- The simulator can call the Go API using a Kubernetes service name.

## Stage 10: Add Kubernetes Health Checks

### Goal

Learn how Kubernetes decides whether a service is healthy.

### Build

Add probes to the Go API:

- liveness probe uses `/healthz`
- readiness probe uses `/readyz`

### Learn

- Liveness probes
- Readiness probes
- Startup behavior
- Why readiness should check dependencies
- How bad probes cause restart loops

### Done When

- Kubernetes removes the API from service endpoints when `/readyz` fails.
- You can explain the difference between liveness and readiness.

## Stage 11: Add Prometheus and Grafana in Kubernetes

### Goal

Run observability inside the cluster.

### Build

Use either:

- raw manifests, for learning the basics
- Helm chart, for learning the real-world approach

Recommended beginner path:

1. First deploy Prometheus manually with a simple config.
2. Later try `kube-prometheus-stack` with Helm.

Scrape:

- Go API
- simulator
- Kubernetes node/pod metrics if available

### Learn

- monitoring inside Kubernetes
- scraping services
- service discovery
- port-forwarding Grafana
- cluster-level versus app-level metrics

### Done When

- Grafana runs in Kubernetes.
- Prometheus scrapes the API and simulator from inside the cluster.
- You can view app metrics after port-forwarding Grafana.

## Stage 12: Add Scaling Exercises

### Goal

Learn why Kubernetes is useful.

### Build

Run exercises:

- scale the simulator from 1 replica to 3 replicas
- scale the Go API from 1 replica to 2 replicas
- delete an API pod and watch it recover
- deploy a broken image and roll it back
- set CPU and memory requests/limits

### Learn

- replicas
- rolling updates
- rollbacks
- scheduling
- resource requests
- resource limits
- failure recovery

### Done When

- You can scale services with `kubectl scale`.
- You can inspect rollout history.
- You can explain what happened after deleting a pod.

## Stage 13: Add Basic Frontend

Inspired by Juraj's UI and React migration posts.

### Goal

Create a simple view of the system.

### Build

Create a React + TypeScript dashboard showing:

- active rides
- recent ride events
- number of drivers
- number of riders
- API status

Keep the UI simple. Do not build a complex map yet.

### Learn

- React app structure
- TypeScript API types
- polling or server-sent events
- frontend containerization
- frontend service in Kubernetes

### Done When

- The frontend runs locally.
- The frontend can call the Go API.
- The frontend is deployed to Kubernetes.

## Stage 14: Add Routing and Matching Logic

Inspired by Juraj's graph, route planner, destination generation, and driver matching posts.

### Goal

Add enough business logic to make the simulation interesting.

### Build

Add:

- zones
- driver locations
- rider pickup locations
- simple nearest-driver matching
- ride states: requested, assigned, picked_up, completed, cancelled

Optional later:

- graph-based map
- route planner
- simulated movement

### Learn

- state machines
- basic algorithms
- transactional updates
- race conditions
- metrics for business logic

### Done When

- New rides can be matched to available drivers.
- Ride state changes are stored as events.
- Grafana shows completed, cancelled, and active rides.

## Stage 15: Add Failure and Debugging Scenarios

### Goal

Practice operating the system.

### Build

Create scenarios:

- database unavailable
- API returns 500 errors
- simulator sends too many requests
- slow database query
- memory limit too low
- bad config value

### Learn

- `kubectl logs`
- `kubectl describe`
- `kubectl exec`
- `kubectl port-forward`
- Prometheus debugging
- Grafana dashboard interpretation

### Done When

- You can diagnose each failure using Kubernetes commands, logs, and dashboards.
- You document the cause and fix for each scenario.

## Stage 16: Optional Production-Style Improvements

Only do this after the basic system works.

### Options

- Add Ingress
- Add TLS locally with cert-manager
- Add Helm charts
- Add GitHub Actions
- Add database migrations
- Add OpenTelemetry tracing
- Add Loki for logs
- Add Horizontal Pod Autoscaler
- Add load testing with k6
- Deploy to a small cloud Kubernetes cluster

### Learn

- Real deployment workflows
- Release automation
- production observability
- cost and complexity tradeoffs

## Stage 17: CI/CD with Jenkins

Do this after you can manually deploy the Go API to Minikube.

### Goal

Learn CI/CD by automating a deployment process you already understand manually.

### Build

Create a Jenkins pipeline that:

- checks out the repository
- runs Go tests
- runs Go formatting or lint checks
- builds the Go API Docker image
- loads or pushes the image for Minikube
- applies Kubernetes manifests
- waits for the rollout to finish
- verifies `/healthz`, `/readyz`, and `/metrics`

### Learn

- `Jenkinsfile`
- pipeline stages
- build agents
- Docker builds from CI
- Kubernetes deploys from CI
- credentials handling
- rollout verification
- failed deployment debugging

### Done When

- Jenkins can deploy the Go API to Minikube.
- A failing test blocks deployment.
- A failed Kubernetes rollout is visible in the Jenkins build.
- You can explain each pipeline stage.

## Stage 18: AWS Integration from Minikube

Do this before moving the whole application to AWS. The goal is to learn AWS services while keeping Kubernetes debugging local.

### Goal

Connect the local Minikube application to selected AWS services.

### Recommended AWS Services

Start with these:

- **Amazon ECR** for Docker image storage
- **Amazon S3** for exported ride reports or simulator snapshots
- **AWS Secrets Manager** or **SSM Parameter Store** for configuration and secrets
- **Amazon CloudWatch** for optional log or metric experiments

Avoid these at first:

- EKS
- RDS
- ALB Ingress
- Route 53
- production TLS

Those are useful later, but they add too much complexity before the local system is stable.

### Build

Add AWS features in this order:

1. Push the Go API Docker image to Amazon ECR.
2. Update the Minikube Kubernetes manifest to pull the image from ECR.
3. Add an API endpoint that exports a small ride summary to S3.
4. Move one non-secret config value to SSM Parameter Store.
5. Move one secret value to AWS Secrets Manager.

### Learn

- AWS IAM users and roles
- ECR authentication
- image repositories and tags
- S3 buckets and object keys
- AWS SDK for Go
- external secrets versus Kubernetes Secrets
- cloud service permissions

### Done When

- The API image is stored in ECR.
- Minikube can deploy the ECR image.
- The Go API can write a test object to S3.
- One config value is read from AWS instead of a local environment variable.
- You can describe which AWS permissions were required.

## Stage 19: Infrastructure as Code with Terraform

Do this after you understand the AWS resources manually. Terraform should codify a process you can already explain.

### Goal

Learn how to define AWS infrastructure as code instead of creating resources manually in the AWS console.

### Build

Create Terraform code for the first AWS resources:

- ECR repository for the Go API image
- S3 bucket for ride reports
- SSM Parameter Store value for non-secret config
- Secrets Manager secret for one application secret
- IAM policy with the minimum permissions the app needs

Recommended structure:

```text
infra/
  terraform/
    environments/
      dev/
        main.tf
        variables.tf
        outputs.tf
    modules/
      ecr/
      s3/
      app_config/
      iam/
```

Start with one `dev` environment only. Add modules after the first version works.

### Learn

- Terraform providers
- resources
- variables
- outputs
- state files
- remote state, later
- IAM policy design
- planning before applying changes
- destroying unused resources to control cost

### Done When

- `terraform plan` shows the expected AWS resources.
- `terraform apply` creates ECR, S3, SSM, Secrets Manager, and IAM resources.
- `terraform destroy` removes the resources cleanly.
- The Jenkins or local build can push an image to the Terraform-created ECR repository.
- You can explain where Terraform state is stored.

## Stage 20: Deploy to AWS Kubernetes

Only do this after Minikube, Jenkins, basic AWS service integration, and Terraform basics work.

### Goal

Move from local Kubernetes to managed Kubernetes on AWS.

### Build

Deploy the app to Amazon EKS:

- create an EKS cluster with Terraform or `eksctl`
- push images to ECR
- deploy the Go API, simulator, and frontend
- use RDS PostgreSQL or keep a temporary in-cluster Postgres for learning
- expose the frontend through an AWS Load Balancer
- run Prometheus and Grafana in the cluster

Recommended path:

1. Use `eksctl` once if you want to understand the moving parts quickly.
2. Rebuild the same environment with Terraform after you understand the resources.

### Learn

- EKS cluster basics
- node groups
- Terraform-managed infrastructure
- AWS Load Balancer Controller
- ECR with Kubernetes
- RDS connectivity from Kubernetes
- cloud networking basics
- cost awareness

### Done When

- The application runs on EKS.
- Images are pulled from ECR.
- The API can connect to PostgreSQL.
- Grafana shows metrics from the AWS-hosted deployment.
- You know how to delete the cluster and avoid ongoing costs.

## Suggested Study Order

Follow this order if you are new:

1. Go API locally
2. PostgreSQL locally
3. Docker Compose
4. Prometheus
5. Grafana
6. simulator
7. Kubernetes core manifests
8. Kubernetes health checks
9. Kubernetes monitoring
10. scaling and failure exercises
11. frontend
12. advanced simulation logic
13. Jenkins CI/CD
14. AWS integration from Minikube
15. Terraform for AWS infrastructure
16. EKS deployment, optional and only after the local system is stable

This order prevents Kubernetes from hiding basic application bugs.

## Weekly Plan

Use this if you want a steady pace.

### Week 1

- Stage 0: plan the system
- Stage 1: Go API
- Stage 2: config and logging

### Week 2

- Stage 3: Prometheus metrics in Go
- Stage 4: PostgreSQL

### Week 3

- Stage 5: Docker Compose
- Stage 6: Prometheus locally
- Stage 7: Grafana locally

### Week 4

- Stage 8: simulator
- improve metrics
- document local system behavior

### Week 5

- Stage 9: Kubernetes deployment
- Stage 10: Kubernetes health checks

### Week 6

- Stage 11: Prometheus and Grafana in Kubernetes
- Stage 12: scaling exercises

### Week 7

- Stage 13: React frontend
- deploy frontend to Kubernetes

### Week 8

- Stage 14: matching logic
- Stage 15: failure and debugging scenarios

### Week 9

- Stage 17: CI/CD with Jenkins
- automate Go test, Docker build, and Minikube deploy

### Week 10

- Stage 18: AWS integration from Minikube
- push images to ECR
- write a small ride report to S3

### Week 11 and Later

- Stage 19: Infrastructure as Code with Terraform
- create ECR, S3, SSM, Secrets Manager, and IAM resources

### Week 12 and Later

- Stage 20: deploy to AWS Kubernetes
- use EKS only after the local version is reliable

## Minimum Version Worth Finishing

If time is limited, finish this smaller version:

- Go API
- PostgreSQL
- Prometheus metrics
- Grafana dashboard
- simulator
- Kubernetes manifests
- readiness/liveness probes
- one scaling exercise
- one failure debugging exercise
- Jenkins pipeline for test, build, and deploy
- ECR image repository, optional
- S3 export feature, optional
- Terraform-managed AWS resources, optional

This smaller version is enough to demonstrate real learning of Go, Kubernetes, Prometheus, and Grafana.

## What to Document After Each Stage

In `docs/journal.md`, write:

```text
Date:
Stage:
What I built:
What I learned:
How I verified it:
Problems and fixes:
Beginner notes:
```

This matters. The original Rides project stands out partly because the work was documented step by step. Documentation will also make the project useful as a portfolio artifact.

## Final Portfolio Checklist

Before calling the project complete, make sure you can show:

- architecture diagram
- local setup instructions
- Kubernetes setup instructions
- screenshots of Grafana dashboards
- explanation of key metrics
- explanation of readiness and liveness probes
- one scaling demo
- one failure recovery demo
- Jenkins pipeline screenshot or build logs
- AWS architecture notes, if AWS stages are included
- Terraform plan/apply output notes, if Terraform is included
- ECR image repository and S3 export demo, if AWS stages are included
- EKS cleanup notes, if EKS is used
- short write-up of what you learned
