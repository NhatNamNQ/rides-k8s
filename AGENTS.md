# AGENTS.md

## Project Goal

Build a ridesharing simulation platform to learn:

- Go backend development
- Kubernetes with Minikube
- PostgreSQL
- Prometheus and Grafana
- Jenkins CI/CD
- AWS services
- Terraform infrastructure as code

The goal is not to build a full commercial ridesharing product. The goal is to build a production-style learning system with clear services, observable behavior, repeatable deployment, and useful debugging exercises.

## Main References

- `LEARNING_PLAN.md`: roadmap and learning stages
- `STAGE_GUIDE.md`: detailed execution steps
- `REVIEW_AGENTS.md`: review agent roles and per-stage review checklist
- `docs/architecture.md`: system design notes
- `docs/journal.md`: learning journal and debugging notes

## Learning Order

Follow this sequence unless the user explicitly changes the plan:

1. Build the Go API locally.
2. Add configuration and logging.
3. Add Prometheus metrics.
4. Add PostgreSQL.
5. Add Docker and Docker Compose.
6. Add Prometheus and Grafana locally.
7. Add the simulator.
8. Deploy to Minikube.
9. Add Kubernetes health checks.
10. Run Prometheus and Grafana in Minikube.
11. Practice scaling, rollouts, and debugging.
12. Add the React frontend.
13. Add Jenkins CI/CD.
14. Integrate basic AWS services from Minikube.
15. Create AWS resources with Terraform.
16. Deploy to EKS only after local Kubernetes is stable.

Do not introduce cloud complexity before the local and Minikube versions work.

## Collaboration Rules

- Prefer making small, working changes over large rewrites.
- Keep the user oriented: state what changed, what was verified, and what remains.
- If a stage is unclear, follow `STAGE_GUIDE.md`.
- Do not skip stages unless the user asks.
- Do not add advanced tooling before the simpler version works.
- Update documentation when a change affects architecture, setup, or learning flow.

## Git and Commit Rules

- Commit after each completed stage or coherent subtask.
- Keep commits small and reviewable.
- Use imperative commit messages.
- Prefer this format:

```text
stage-02: add config and structured logging
```

- Good examples:
  - `stage-00: add architecture and learning journal`
  - `stage-01: add initial Go API`
  - `stage-02: add config and structured logging`
  - `docs: add review agent playbook`
  - `infra: add initial terraform layout`
- Do not commit secrets, `.env` files, kubeconfigs, Terraform state, AWS credentials, or generated tokens.
- Before committing code, run the relevant tests or validation commands.
- Before committing Kubernetes or Terraform changes, run the relevant dry-run, validation, or plan command when available.
- If a commit contains known limitations, document them in `docs/journal.md`.

## Engineering Rules

- Prefer simple, explicit code.
- Avoid abstractions until there is real duplication or complexity.
- Keep service boundaries clear.
- Every important behavior should be testable locally before Kubernetes is involved.
- Every service should have clear configuration through environment variables.
- Do not hardcode secrets.
- Do not commit generated secrets, credentials, kubeconfigs, Terraform state, or cloud tokens.
- Keep commits and changes scoped to the current stage.

## Go Backend Rules

- Use the Go standard library first where reasonable.
- Start with `net/http`; add routers or frameworks only when needed.
- Keep handlers small.
- Separate HTTP handling, business logic, and storage once the code grows.
- Use context-aware database calls.
- Return clear HTTP status codes.
- Add tests for core ride logic.
- Expose:
  - `GET /healthz`
  - `GET /readyz`
  - `GET /metrics`
- `/healthz` should report process health.
- `/readyz` should check dependencies such as PostgreSQL.

## PostgreSQL Rules

- Keep the initial schema small.
- Use migrations once the schema is no longer trivial.
- Prefer explicit SQL for learning.
- Use database constraints where they protect data correctness.
- Do not hide database failures from readiness checks.

## Observability Rules

- Add metrics before dashboards.
- Prometheus metrics should answer operational questions.
- Grafana dashboards should be useful, not decorative.
- Track at minimum:
  - HTTP request count
  - HTTP request duration
  - error count
  - rides created
  - active rides
  - simulator events
- Useful dashboard questions:
  - Is the API healthy?
  - Is latency increasing?
  - Are requests failing?
  - Is the simulator overwhelming the API?
  - Are pods restarting?

## Docker Rules

- Use Docker Compose before Kubernetes.
- Keep Dockerfiles small and readable.
- Prefer multi-stage builds for compiled services.
- Runtime configuration must come from environment variables.
- Container images should not contain secrets.

## Kubernetes Rules

- Use Minikube as the local Kubernetes environment.
- Prefer plain Kubernetes YAML before Helm.
- Every workload should define labels clearly.
- Services must use selectors that match pod labels.
- Add resource requests and limits once the basic deployment works.
- Add liveness and readiness probes.
- Verify deployments with:
  - `kubectl get pods`
  - `kubectl get svc`
  - `kubectl describe pod`
  - `kubectl logs`
  - `kubectl rollout status`
- Do not move to EKS until Minikube deployment and debugging are understood.

## Jenkins Rules

- Jenkins should automate a manual process that already works.
- Start with a simple `Jenkinsfile`.
- Pipeline stages should be:
  - checkout
  - test
  - build
  - image
  - deploy
  - verify
- A failing test must block deployment.
- A failed Kubernetes rollout must fail the pipeline.
- Credentials must be stored in Jenkins credentials, not in source files.

## AWS Rules

- Start with AWS services from Minikube before using EKS.
- Recommended first AWS services:
  - ECR for images
  - S3 for ride report exports
  - SSM Parameter Store for non-secret config
  - Secrets Manager for secrets
- Use least-privilege IAM.
- Document every AWS service used and why it exists.
- Track cleanup steps for every AWS resource.
- Avoid EKS, RDS, Route 53, and production TLS until the local version is stable.

## Terraform Rules

- Use Terraform after understanding the AWS resources manually.
- Start with one `dev` environment.
- Keep the first Terraform version simple before creating reusable modules.
- Always run:
  - `terraform fmt`
  - `terraform validate`
  - `terraform plan`
- Do not apply Terraform changes that are not understood.
- Do not commit Terraform state files.
- Know how to run `terraform destroy` for all created learning resources.

## Frontend Rules

- Use React and TypeScript for the dashboard.
- Keep the first frontend simple.
- Show operational state before adding visual polish.
- Useful first views:
  - active rides
  - recent ride events
  - API status
  - basic metrics summary
- The frontend should use environment-based API configuration.

## Documentation Rules

After each stage, update `docs/journal.md` with:

```text
Date:
Stage:
What I built:
Commands I used:
What broke:
How I fixed it:
What I learned:
Questions:
```

Update `docs/architecture.md` when:

- a new service is added
- data flow changes
- deployment topology changes
- AWS services are introduced
- CI/CD behavior changes

## Suggested Agent Roles

Use these roles conceptually when splitting work:

- Project Architect: roadmap, architecture, service boundaries
- Go Backend Agent: API, storage, metrics, tests
- Kubernetes Agent: manifests, Minikube, probes, services, debugging
- Observability Agent: Prometheus, Grafana, dashboards, metrics
- CI/CD Agent: Jenkins pipeline and deployment verification
- Terraform AWS Agent: AWS resources, IAM, Terraform, EKS
- Frontend Agent: React and TypeScript dashboard
- Reviewer Agent: reviews correctness, security, reliability, and documentation

For early stages, keep it simple:

- Stage 0: Project Architect
- Stages 1-4: Go Backend Agent
- Stages 5-12: Kubernetes Agent and Observability Agent
- Stage 13: Frontend Agent
- Stage 17: CI/CD Agent
- Stages 18-20: Terraform AWS Agent

Use `REVIEW_AGENTS.md` after each stage to decide which review agents should check the work. At minimum:

- Go stages need Code Review.
- Kubernetes stages need Kubernetes Review.
- AWS and Terraform stages need Security Review and Terraform AWS Review.
- Monitoring stages need Observability Review.
- Completed milestones need Documentation Review.

## Definition of Done

A stage is done only when:

- the code or configuration works locally or in the intended environment
- the relevant commands were verified
- docs or journal notes were updated
- the user can explain what was built and why
- known failures or limitations are recorded
