# Task Manager operations

## Configuration and process lifecycle

Load typed configuration at the composition root, validate it once, and fail startup with a redacted explanation when required values are missing. Initial precedence is defaults, then environment, with explicit `*_FILE` secret paths; setting both a secret and its file form is an error. Configuration changes require restart unless a component explicitly implements safe reload.

| Group | Examples |
| --- | --- |
| Process | `SERVICE_NAME`, `APP_ENV`, `HTTP_ADDR`, `LOG_LEVEL`, shutdown grace |
| Storage | Owned database DSN file, pool limits, connection/acquisition timeouts |
| Messaging | NATS URL/credential file, stream/consumer names, replay horizon |
| Identity | Issuer, JWKS URL, audience allowlist, signing-key file |
| Security | CORS allowlist, limiter Redis credential, trusted proxy list |
| AI | Model/embedding provider and model names, API-key file, run budgets, Mongo/Qdrant and Celery credentials |
| Telemetry | OTLP endpoint, sample rate, metrics listener, resource attributes |

Future `.env.example` files contain nonsecret placeholders, never working credentials. Pin dependencies, container images, and protocol versions when each component is introduced; record upgrades and run compatibility checks. Avoid floating `latest` image tags in reproducible development.

Implement `/livez` for process health, `/readyz` for ability to serve core routes, and an internal `/metrics`. Readiness checks required owned storage and schema readiness, without contacting a paid model endpoint. Optional telemetry/search/cache failures do not trigger liveness restarts. Report degraded asynchronous dependencies separately.

On SIGTERM, stop readiness and work claims, drain requests within 20 seconds, finish or abandon leased work safely, flush bounded telemetry buffers, then close pools. Compose/Kubernetes termination grace starts at 30 seconds. A forced stop must still permit redelivery and reconciliation.

## Middleware behavior

Compose HTTP behavior deliberately rather than relying on accidental registration order: request context and validated request ID, panic recovery, tracing/timing/access log, body limits, CORS/preflight, rate limits, authentication, input validation, application authorization, and error mapping. Ensure recovery and final status/duration logging still work on all failures. Middleware integration tests verify the actual Echo execution order chosen by the pinned version.

Gzip applies only to sufficiently large eligible JSON responses (initial threshold 1 KiB), sets `Vary: Accept-Encoding`, and is excluded from auth/token endpoints, WebSocket upgrades, and AI streams. Test encoded and unencoded clients and verify actual byte/latency tradeoffs. Timing uses monotonic duration measurements. Logs contain service, route template, status, duration, request/trace IDs, and redacted error code.

## Compose profiles and resource budgets

The following profiles are implementation targets; this documentation does not provide runnable Compose files yet.

| Profile | Components | Initial aggregate memory budget |
| --- | --- | --- |
| `core` | Identity, Workspace, Task, Notifications, Go relays/workers/scheduler, PostgreSQL, cache Redis, NATS | 4 GiB |
| `search` | Search API/indexer, using its owned database in local PostgreSQL | Additional 512 MiB |
| `ai` | FastAPI, one Celery worker, one Beat, MongoDB, Qdrant, dedicated broker Redis, document volume | Additional 4.5 GiB |
| `observability` | OTel Collector, Prometheus, Grafana, Tempo, Loki | Additional 2 GiB |
| `mqtt` | Mosquitto and protocol test client | Additional 256 MiB |

The full stack exceeds the approximately 8 GiB target; run selected profiles and reduce nonessential processes during AI or telemetry exercises. Budget actual container limits, steady/peak memory, startup behavior, and OOM recovery before claiming the target is met. Keep basic JSON logs and metrics instrumentation enabled even when collectors are off.

Bind exposed APIs to loopback. Databases, NATS, Redis, MongoDB, Qdrant, metrics, and admin ports stay on private networks unless an explicit loopback troubleshooting override is used. Use named volumes for durable state. Check health before starting dependent traffic; container start order is not service readiness.

## Workers and scheduling

Go workers consume NATS events with bounded concurrency and acknowledge after durable effects. Outbox relays are explicit processes with backlog metrics. Notifications uses a Go scheduler scanning due records every 30 seconds; workers claim records with leases/row locks and publish durable evaluation commands. Before creating the notification, the worker asks Task whether the task/due-version/assignee/due-time combination remains eligible. A uniqueness constraint protects against overlapping schedulers and delivery retries.

Do not create one long-lived goroutine or sleeping queue message per task deadline. Store deadlines in PostgreSQL so restarts and downtime can catch up. Start with one scheduler, then intentionally run two to verify coordination. Completed, edited, deleted, or reassigned tasks invalidate old reminder generations. During recovery, stale reminders are suppressed through canonical eligibility checks.

Celery workers and Beat belong only to AI. Beat handles periodic reconciliation/retention, while the job ledger holds business progress. A one-shot Go maintenance command provides the cron learning exercise, invoked manually during Compose development and as a Kubernetes CronJob later; it never replaces the due-reminder workflow.

## Caching and connection pools

First establish uncached measurements. Add cache-aside only to a measured hot read, with service-scoped keys, a 30-second initial TTL plus jitter, version metadata, and invalidation after committed updates. Prevent stampedes with bounded coalescing. Do not negative-cache authorization decisions or place private payloads in global keys. Current authorization runs before a cache hit is returned.

Configure eviction for data cache Redis. The Celery broker is a separate persistent, non-evicting instance; a logical Redis database number alone is not memory/failure isolation. Lost queue dispatch is reconciled against PostgreSQL. Exercise stale reads, invalidation loss, Redis restarts, and failover to database reads within concurrency limits.

Record pool usage, wait time, acquisition errors, active/idle connections, and total budget across APIs/workers. Pool saturation should produce bounded timeouts rather than unbounded goroutines or more connections. Use request cancellation and deadlines on SQL, gRPC, HTTP, model, and broker calls.

## Observability

Instrument using OpenTelemetry and propagate W3C trace context through HTTP, gRPC, NATS headers, Celery job metadata, and tool calls. Use span links for asynchronous work that resumes long after its originating request. Signals and their roles follow the [OpenTelemetry concepts](https://opentelemetry.io/docs/concepts/signals/).

Prometheus measures request rate/errors/duration, database and pool waits, outbox age, consumer lag/redeliveries/quarantine, scheduler delay, job depth/age, cache hit rate, model tokens/cost estimates, and tool/approval outcomes. Use bounded labels such as service, route template, operation, and error class. User, task, workspace, request, and run IDs belong in redacted logs/traces, not metric labels.

Grafana dashboards cover request health, dependencies/pools, event processing, and AI quality/latency/cost. Tempo receives traces and Loki receives redacted logs. Start with local warning thresholds: oldest unpublished outbox over 60 seconds, quarantine count above zero, due-reminder delay over two minutes, unexpected error ratio above 1% over five minutes, and sustained pool saturation. Tune from workload evidence; these are diagnostic thresholds rather than availability promises.

## Recovery runbooks

| Symptom | Inspect | Recovery and proof |
| --- | --- | --- |
| Notifications missing | Task transaction/outbox, relay, consumer lag, inbox, socket | Restore dependency; verify one inbox record and refetch after reconnect |
| Quarantined message | Schema, consumer error, aggregate version, effect ledger | Fix cause; dry-run scoped replay; preserve original event ID; verify no duplicate effect |
| Search stale | Consumer lag and source versions | Rebuild into a new projection generation; compare samples and switch only after catch-up |
| AI job stuck | PostgreSQL job stage, worker lease, broker, document generation | Reconcile and resume same job; exclude tombstoned sources |
| Tool result unknown | Approval/execution key, remote idempotency result, audit | Recover existing result; never create a fresh execution key to bypass uncertainty |
| High latency | Trace breakdown, pool waits, query plans, CPU/memory | Reproduce fixed workload, change one cause, compare identical runs |
| Database loss | Backup age, schema version, event horizon | Restore into isolated volumes; verify ownership/row counts and run acceptance story |

Back up authoritative PostgreSQL, MongoDB, and document originals. Derived Search/Qdrant indexes may be rebuilt, but record their model/index generation. Test a restore before destructive migrations. Keep local rolling backups for seven days as the initial retention default; private text does not enter general logs or evaluation artifacts.

## Kubernetes and CI/CD exit stages

S10 introduces namespace isolation, Deployments, Services, probes, resource requests/limits, ConfigMaps, mounted Secrets, migration Jobs, maintenance CronJobs, persistent volumes, and NetworkPolicies. Start with one local cluster. Test rollout, pod termination during work, and restoring persistent data. Kubernetes Secrets still require access controls and an explicit at-rest encryption/secret-management decision; base64 encoding is not a vault.

S11 adds a platform-neutral pipeline: formatting/static checks, unit and race tests, contract validation, container integration/E2E tests with fake AI, image build and vulnerability/secret scans, immutable image publication, and a gated local/staging rollout. Live AI evaluation is an opt-in job with a budget. The Git host/CI provider and registry are chosen at S11, before writing provider-specific files. Automated rollback must check schema compatibility; application rollback cannot undo an incompatible migration.
