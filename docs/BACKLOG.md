# Task Manager development backlog

## Delivery rules

S00 is the documentation baseline. S01–S11 are future implementation work, in dependency order. Each slice is a small demonstrable capability rather than a demand to finish a whole service first. All code changes include relevant contracts, tests, migrations, and basic telemetry. A slice is complete only when its exit evidence exists.

The project has no calendar deadline. Record actual effort and revise later slice size from experience. Do not create all infrastructure before delivering the first task flow. Service extraction follows the agreed six-service ownership from the beginning; selective startup keeps the local stack manageable.

## S00 Product and architecture baseline

**Deliverables:** English PRD, architecture, ADRs, data/contracts/security/AI/operations/test designs, topic coverage, and this backlog.

**Exit:** documents link correctly; product requirements have acceptance scenarios; accepted decisions are distinct from future executable artifacts. The existing Go scaffold and original diagrams remain intact.

## S01 Register and sign in

**Depends on:** S00. **Product value:** a user can register, log in, inspect their account, rotate credentials, and log out through HTTP.

- Introduce the Identity composition root and hexagonal packages, reusing the useful scaffold without importing transport types into domain. Add executable Identity OpenAPI and separate account response/password-hash models.
- Add PostgreSQL migrations, bounded pools, typed configuration, mounted development secrets, and RS256/JWKS/session rotation.
- Add a minimal Compose setup for Identity, PostgreSQL, and limiter Redis; no AI, Search, or full telemetry stack required.
- Add validation, structured errors/logging, timing/tracing, body limits, recovery, and rate limiting with tested middleware order.

**Exit:** T01; real database migration/credential integration tests; JWT negative tests; redacted logs; clean shutdown; commands to reproduce registration and login. Basic probes/metrics exist even without collectors.

## S02 Provision workspaces and membership

**Depends on:** S01. **Product value:** users receive personal spaces and can create teams and invite members.

- Add Workspace HTTP contracts, schema, role rules, invitation lifecycle, ownership transfer, and local invitation-link delivery.
- Add Identity outbox, NATS JetStream, Workspace inbox/provisioning and the pending/ready onboarding flow.
- Add Workspace membership Protobuf contract and gRPC server with mTLS/user context, deadlines, and error mapping.
- Add key-rotation exercise and bearer/CORS/CSRF comparison fixture. Keep social login outside the slice.

**Exit:** T02–T04 and Workspace role boundaries from T05; broker-down provisioning recovery; duplicate invitation acceptance; gRPC deadline and identity tests; one personal workspace after replay.

## S03 Create projects and manage tasks

**Depends on:** S02. **Product value:** personal and team users can manage real work with private projects.

- Add Task service contracts, owned migrations, project visibility/ACL, task CRUD/filtering/cursor pagination, assignee validation, version conflicts, audit, and outbox.
- Enforce current Workspace membership through gRPC on every protected use case. Add internal authorized batch hydration and reminder-eligibility contracts.
- Extend Compose with Task and provide executable API acceptance fixtures for multiple users/workspaces.

**Exit:** T05–T06; cross-tenant and unauthorized list tests; archive behavior; concurrent edit test; Workspace outage fails closed. The first complete personal-task flow runs without Search, AI, or frontend.

## S04 Deliver durable notifications and reminders

**Depends on:** S03. **Product value:** invitation/assignment notifications and due reminders appear in the user's inbox. This completes Release 1.

- Add Notifications API, event consumers, durable inbox, uniqueness rules, and read state. Route invitations using the identity projection and reconcile out-of-order account/invitation arrivals.
- Add Go scheduler/worker leases, reminder evaluation commands, canonical eligibility checks, retries, quarantine, and replay tooling with a dry run.
- Add OTel propagation through NATS, lag/delay metrics, and the first optional Prometheus/Grafana dashboard.

**Exit:** T07–T08; kill-before/after-commit experiments; two schedulers produce one eligible reminder; broker outage drains safely; demonstrate correlation from task write to inbox.

## S05 Add live and MQTT delivery

**Depends on:** S04. **Product value:** protocol clients receive notification signals without polling continuously.

- Add authenticated notification WebSocket stream, reconnect synchronization, heartbeats, queue bounds, and token-expiry handling.
- Add optional Mosquitto profile, per-user credentials/ACL provisioning and revocation, Notifications bridge, and CLI test subscriber. Send inbox IDs only.
- Add the one-shot maintenance command for the later cron exercise.

**Exit:** real WS/MQTT tests, denied foreign/wildcard subscription, duplicate QoS deliveries, slow-consumer bounds, and successful inbox recovery after disconnect. No frontend framework is selected.

## S06 Search and measured optimization

**Depends on:** S03–S05. **Product value:** users search accessible work efficiently. This completes Release 2.

- Add Search's versioned PostgreSQL projection, indexing consumer, canonical hydration, tombstones, and rebuild/catch-up procedure.
- Establish the fixed workload and diagnose one actual bottleneck using traces, pprof, pool metrics, and SQL plans.
- Implement cache-aside for a measured hot read, tested invalidation/TTL/stampede behavior, pool tuning as warranted, and eligible response gzip.

**Exit:** T09; revoked access despite deliberately stale projection/cache; reproducible before/after performance report and memory measurements; search rebuild after consumer downtime. Record improvements and remaining bottlenecks without assuming every intervention helps.

## S07 Read and mutate tasks through an approved AI tool flow

**Depends on:** S06. **Product value:** a deterministic test agent, then a configured live model, can look up a task and propose an update that the user confirms.

- Add Task/Search MCP executable contracts and adapters, resource-audience validation, OAuth discovery/authorization-code PKCE flow, and scoped internal delegation.
- Add FastAPI/LangGraph with a fake model first, owned PostgreSQL checkpoints/run/approval ledger, MongoDB conversation content, and provider ports.
- Add structured tool arguments, interrupts, approval decisions, idempotent remote execution, ownership checks, and ambiguous-result recovery.

**Exit:** T11–T13; restart while waiting and after a remote write; reject/tamper/expire tests; MCP conformance and domain equivalence; no unapproved writes. Live model selection is explicit and optional for deterministic gates.

## S08 Add documents, memory, RAG, and ingestion workers

**Depends on:** S07. **Product value:** the assistant answers from authorized project documents and remembers explicit user preferences. This completes Release 3.

- Add text/Markdown storage port and local volume, MongoDB document/memory models, Qdrant generation/scoping, Celery broker/worker/Beat, and durable ingestion reconciliation.
- Add permission-aware chunk retrieval, source/version citations, deletion tombstones, explicit memory inspect/delete, and source dependency tracking for conversation access.
- Add authenticated AI streams, budget controls, fake embeddings, adversarial fixtures, and a live evaluation harness recording model/prompt/fixture versions and price basis.

**Exit:** T10 and T14 plus full T11–T13 regression; kill ingestion at every durable stage, delete during indexing, lose dispatch, revoke access during a run; deterministic security gates pass and live evaluation report states measured quality/latency/cost.

## S09 Harden and recover

**Depends on:** S08. **Product value:** failures are diagnosable and authoritative data is recoverable.

- Complete full telemetry profile, dashboards, diagnostic thresholds, log redaction, external-call budgets, and bounded concurrency/load shedding.
- Exercise expand-contract migrations, resumable Mongo backfill, embedding generation replacement, backup/restore, retention/purge, key/secret rotation, and schema-compatible application rollback.
- Repeat the acceptance story with dependency outages and complete the learning coverage evidence. Measure selected Compose profiles against resource budgets.

**Exit:** documented restore into isolated volumes, failure-injection report, no unresolved access/idempotency failures, reproducible performance results, and an updated list of limitations.

## S10 Move the working system to Kubernetes

**Depends on:** S09. **Product value:** the same contracts and workflows run in a local cluster.

- Add Deployments/Services, probes, resources, mounted config/secrets, storage, migration Jobs, maintenance CronJobs, and NetworkPolicies.
- Define stateful backup/restore and certificate/secret access; migrate local fixture data through documented export/restore rather than shared host paths.
- Exercise rolling updates, scale changes, pod termination during processing, and Job/CronJob overlaps.

**Exit:** acceptance E2E and recovery tests pass on the cluster; resource and pool budgets account for replicas; no domain code depends on Compose networking.

## S11 Add CI/CD

**Depends on:** S10. **Product value:** changes produce reproducible tested images and a controlled deployment.

- Choose the repository's CI provider and registry, then implement the platform-neutral stages in [operations](OPERATIONS.md).
- Add deterministic test/schema/compatibility gates, immutable builds, security scans, image publication, and a deployment gate. Keep live-model evaluation opt-in and budgeted.
- Record rollback constraints, migration order, and build provenance; demonstrate one failed gate and one successful rollout.

**Exit:** a clean checkout reproduces the pipeline; unsafe changes fail before deployment; the final acceptance story passes against the deployed environment. This completes Release 4.
