# Review Agents Playbook

This file defines the review agents for the Rides K8s learning project.

Use these agents as checklists, review roles, or prompts for focused review sessions. They do not need to run at the same time. For early stages, use only the reviewers that match the work.

## Review Philosophy

The goal is to catch real engineering problems, not to create ceremony.

Every review should answer:

- Does it work?
- Is it understandable?
- Is it safe enough for the current stage?
- Is it observable?
- Is it documented?
- What would fail first?

## Core Review Agents

## 1. Code Review Agent

### Purpose

Review application code quality, correctness, tests, and maintainability.

### Applies To

- Go API
- simulator
- React frontend
- scripts
- Jenkinsfile logic when code-like

### Review Checklist

- Code is simple and readable.
- Function names explain intent.
- Error handling is explicit.
- HTTP status codes are correct.
- Input validation exists where needed.
- Tests cover core behavior.
- No large unrelated refactors.
- No hardcoded secrets or environment-specific values.
- Logs contain useful debugging context.
- Public endpoints are documented.

### Common Findings To Look For

- ignored errors
- shared mutable state without protection
- missing request timeouts
- unclear handler responsibilities
- no tests for state transitions
- duplicated config parsing
- logs that hide the useful error

## 2. Security Review Agent

### Purpose

Review secrets, permissions, attack surface, dependency risk, and unsafe defaults.

### Applies To

- Go API
- Kubernetes manifests
- Jenkins
- Terraform
- AWS IAM
- Dockerfiles
- frontend config

### Review Checklist

- No secrets committed to Git.
- Secrets are not printed in logs.
- Kubernetes Secrets are used for sensitive values.
- ConfigMaps are not used for passwords or tokens.
- IAM policies use least privilege.
- Jenkins credentials are stored in Jenkins credentials manager.
- Docker images do not run as root when avoidable.
- Containers expose only required ports.
- External endpoints validate input.
- CORS is not overly permissive unless intentionally local-only.
- Terraform state is not committed.

### Common Findings To Look For

- `AWS_ACCESS_KEY_ID` or `AWS_SECRET_ACCESS_KEY` in files
- database passwords in YAML
- wildcard IAM permissions such as `s3:*` or `*:*`
- `latest` image tags in Kubernetes manifests
- public S3 buckets
- missing `.gitignore` entries for state and secrets

## 3. Kubernetes Review Agent

### Purpose

Review Kubernetes correctness, deployability, reliability, and debuggability.

### Applies To

- Minikube manifests
- EKS manifests
- services
- deployments
- config
- probes
- resource settings

### Review Checklist

- Labels and selectors match.
- Container ports and Service ports are correct.
- Readiness probe uses `/readyz`.
- Liveness probe uses `/healthz`.
- Resource requests and limits are defined after the first working deployment.
- Config is provided through ConfigMaps or Secrets.
- Pods can start without manual intervention.
- Rollout status can be checked.
- Logs are accessible with `kubectl logs`.
- Manifests are stage-appropriate and not overcomplicated.

### Common Findings To Look For

- Service selector does not match pod labels
- wrong container port
- readiness and liveness using the same endpoint incorrectly
- missing env vars
- image pull failures
- unbounded CPU/memory
- depending on local Docker images after switching to ECR

## 4. Observability Review Agent

### Purpose

Review metrics, dashboards, logs, and debugging usefulness.

### Applies To

- Go API metrics
- simulator metrics
- Prometheus config
- Grafana dashboards
- Kubernetes monitoring

### Review Checklist

- `/metrics` works.
- Prometheus target is `UP`.
- Metrics have clear names.
- Counters, gauges, and histograms are used correctly.
- Dashboards answer operational questions.
- Error rate and latency are visible.
- Ride-specific metrics are visible.
- Dashboard panels are not decorative-only.
- Logs and metrics can be correlated during failures.

### Common Findings To Look For

- metrics that never change
- missing labels
- high-cardinality labels such as raw user IDs
- dashboards with no failure signal
- Prometheus scraping the wrong host or port
- histograms created but not used in Grafana

## 5. Terraform AWS Review Agent

### Purpose

Review AWS infrastructure, Terraform structure, IAM, cost, and cleanup safety.

### Applies To

- Terraform files
- AWS service usage
- ECR
- S3
- SSM Parameter Store
- Secrets Manager
- IAM
- EKS
- RDS

### Review Checklist

- `terraform fmt` passes.
- `terraform validate` passes.
- Terraform state is not committed.
- Variables and outputs are clear.
- Resource names identify the project and environment.
- IAM permissions are scoped narrowly.
- S3 buckets are private by default.
- ECR lifecycle policy exists eventually.
- Expensive resources are documented.
- Cleanup steps are documented.

### Common Findings To Look For

- committed `.tfstate`
- broad IAM permissions
- public S3 access
- no region variable
- no environment name
- EKS or RDS created before cost controls are understood
- missing destroy instructions

## 6. CI/CD Review Agent

### Purpose

Review Jenkins pipeline correctness and deployment safety.

### Applies To

- `Jenkinsfile`
- build scripts
- deployment scripts
- image tagging
- Kubernetes rollout verification

### Review Checklist

- Pipeline runs tests before build.
- Pipeline fails on test failure.
- Docker image tag is deterministic.
- Deployment uses the intended image.
- Rollout status is checked.
- Credentials are not stored in the repository.
- Build logs are useful.
- Failed deployment fails the pipeline.
- Manual deployment still works.

### Common Findings To Look For

- pipeline deploys even after tests fail
- image tag mismatch between build and manifest
- credentials in shell commands
- no rollout verification
- deployment depends on local machine state without documentation

## 7. Documentation Review Agent

### Purpose

Review whether another person can understand, run, and evaluate the project.

### Applies To

- README
- `docs/architecture.md`
- `docs/journal.md`
- stage notes
- diagrams
- portfolio write-up

### Review Checklist

- Setup steps are accurate.
- Commands are copyable.
- Architecture is explained.
- Each stage has notes.
- Known limitations are documented.
- Debugging notes include symptoms, cause, and fix.
- Screenshots are included when useful.
- The project explains what was learned.

### Common Findings To Look For

- stale commands
- missing prerequisites
- no cleanup instructions
- docs say Kubernetes but only Docker works
- no explanation of observability or CI/CD decisions

## Review Agents by Stage

Use this table to decide which agents are required.

| Stage | Work | Required Review Agents |
|---|---|---|
| 0 | planning and architecture | Documentation, Code Review lightly |
| 1 | first Go API | Code Review |
| 2 | config and logging | Code Review, Security |
| 3 | Prometheus metrics | Code Review, Observability |
| 4 | PostgreSQL | Code Review, Security |
| 5 | Docker Compose | Code Review, Security |
| 6 | Prometheus local | Observability |
| 7 | Grafana local | Observability, Documentation |
| 8 | simulator | Code Review, Observability |
| 9 | Minikube deployment | Kubernetes, Security |
| 10 | health checks | Kubernetes, Observability |
| 11 | monitoring in Minikube | Kubernetes, Observability |
| 12 | scaling and rollouts | Kubernetes, Documentation |
| 13 | React frontend | Code Review, Security, Documentation |
| 14 | matching logic | Code Review, Observability |
| 15 | failure scenarios | Kubernetes, Observability, Documentation |
| 16 | production-style improvements | Depends on feature |
| 17 | Jenkins CI/CD | CI/CD, Security, Kubernetes |
| 18 | AWS from Minikube | Terraform AWS, Security, Documentation |
| 19 | Terraform | Terraform AWS, Security |
| 20 | EKS | Kubernetes, Terraform AWS, Security, Observability |

## Per-Stage Review Prompts

Use these prompts when asking for review.

### Code Review Prompt

```text
Review this stage as the Code Review Agent.
Focus on correctness, simplicity, test coverage, error handling, and maintainability.
List findings by severity with file and line references.
```

### Security Review Prompt

```text
Review this stage as the Security Review Agent.
Focus on secrets, unsafe defaults, IAM, container security, exposed ports, and dependency risk.
List concrete risks and how to fix them.
```

### Kubernetes Review Prompt

```text
Review this stage as the Kubernetes Review Agent.
Focus on labels/selectors, probes, services, env vars, resources, rollout behavior, and Minikube compatibility.
List anything that can break deployment or debugging.
```

### Observability Review Prompt

```text
Review this stage as the Observability Review Agent.
Focus on metrics correctness, Prometheus scrape config, Grafana usefulness, logs, and failure visibility.
Identify missing signals.
```

### Terraform AWS Review Prompt

```text
Review this stage as the Terraform AWS Review Agent.
Focus on Terraform structure, state safety, IAM least privilege, AWS cost, naming, and cleanup.
List risks before apply.
```

### CI/CD Review Prompt

```text
Review this stage as the CI/CD Review Agent.
Focus on Jenkins pipeline order, failure behavior, image tagging, credentials, deployment verification, and rollback.
List pipeline risks.
```

### Documentation Review Prompt

```text
Review this stage as the Documentation Review Agent.
Focus on whether a new person can reproduce the work, understand the design, and learn from the debugging notes.
List missing or stale documentation.
```

## Minimum Review Policy

Use this policy if you want to keep reviews lightweight:

- Every Go stage: Code Review
- Every Kubernetes stage: Kubernetes Review
- Every AWS or Terraform stage: Security Review and Terraform AWS Review
- Every monitoring stage: Observability Review
- Every completed milestone: Documentation Review

## Strong Portfolio Review Policy

Use this policy before presenting the project to recruiters:

1. Full Code Review
2. Full Security Review
3. Full Kubernetes Review
4. Full Observability Review
5. Full Terraform AWS Review
6. Full CI/CD Review
7. Full Documentation Review

The final portfolio should show not only that the system works, but that you can reason about reliability, safety, deployment, and operations.

