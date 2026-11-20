# Task Manager learning roadmap

## Learning method

Implement one product flow at a time. For each topic, explain the mechanism, implement it behind a suitable port/adapter, trigger a failure or boundary case, and capture evidence. The mentor reviews decisions and asks for explanations; implementation assistance is requested explicitly per task.

Do not postpone testing, security, or basic observability to the last stage. Add them to each slice and broaden their coverage as the system grows. The [backlog](BACKLOG.md) defines the sequence; this matrix ensures that no requested learning topic is lost.

## Coverage and evidence

| Topic | Product application | First slice | Evidence of understanding |
| --- | --- | --- | --- |
| PRD | User flows, release scope, acceptance criteria | S00 | Trace a requirement to a test and explain an excluded feature |
| ADR | Service, storage, consistency, security decisions | S00 onward | Compare alternatives and identify a condition that would reverse a decision |
| Microservices and hexagonal architecture | Six owners; ports and adapters | S01 | Run one service independently; replace an outbound adapter without changing domain rules |
| Cloud-native lifecycle | External config, probes, graceful stop, stateless processes | S01 | Restart during a request and show bounded draining/recovery |
| Middleware and timing | HTTP request pipeline | S01 | Trace success, rejection, and panic; verify final status and monotonic duration |
| Validators and error handling | Transport and domain validation | S01–S03 | Show why valid JSON can still violate a business rule; map errors consistently |
| Logging | Redacted structured logs and audit events | S01 | Correlate a failure without exposing a token or password |
| Authentication, JWT, RS256 | Credentials, sessions, signatures, JWKS | S01 | Rotate signing keys; reject wrong audience/algorithm; explain logout's access-token window |
| Authorization, RBAC, ACL | Workspace roles and private-project membership | S02–S03 | Prove that membership/ACL revocation applies across adapters |
| CORS and CSRF | Explicit bearer mode and cookie comparison fixture | S02 | Demonstrate why CORS is not authorization and when CSRF protection is required |
| Security and sanitization | Safe queries, text ingestion, bounded input | S01–S08 | Run malicious input and prompt fixtures with no privilege expansion |
| Rate limiting and secret management | Login/model limits, mounted credentials | S01–S02 | Demonstrate multi-replica limits, limiter outage policy, and secret rotation |
| SQL | Ownership, constraints, transactions, joins, indexes | S01–S06 | Explain a query plan and a concurrency invariant enforced by the database |
| NoSQL | MongoDB conversations/documents/memory and Redis keys | S06–S08 | Design indexes, compare embedding vs referencing, resume a schema backfill |
| Vector databases | Qdrant chunks and embedding generations | S08 | Compare recall/latency and show tenant filters plus canonical permission checks |
| Configuration management | Typed startup configuration and precedence | S01 | Fail on missing/conflicting secrets with a redacted error |
| Migrations | Owned SQL migrations, document backfills, vector generation changes | S01 onward | Upgrade prior schema, resume a backfill, restore before a destructive change |
| Connection pooling | Bounded pools across processes | S01–S06 | Saturate the pool and distinguish pool wait from slow SQL |
| Unit tests | Domain and application invariants | S01 onward | Test role rules and time-based behavior with fake ports/clock |
| Integration tests | Real storage, broker, middleware, adapters | S01 onward | Catch a transaction/redelivery error a mock cannot reproduce |
| E2E tests | Registration through assignment and notification, later AI | S03–S08 | Run the acceptance story over actual service boundaries |
| Message-driven and event-driven design | Work commands and committed facts | S02–S04 | Explain consumer groups vs independent subscribers and replay effects |
| Idempotency and retrying | Mutations, outbox/inbox, remote effects | S02–S04 | Kill a process at ambiguity points and recover without another business effect |
| Workers and background tasks | Notification consumers and AI ingestion | S04, S08 | Bound concurrency, acknowledgment, cancellation, and retry behavior |
| Cron jobs and schedulers | Reminders, reconciliation, maintenance | S04, S08, S10 | Survive downtime and overlapping executions; distinguish a clock from an executor |
| Celery | Python ingestion workers and Beat schedules | S08 | Recover a lost dispatch from the durable job ledger |
| Redis caching | Measured hot reads, TTL/invalidation/stampede | S06 | Compare cached/uncached runs and reproduce stale/failed-cache behavior |
| Gzip | Eligible JSON response compression | S06 | Measure bytes, CPU, and latency; explain excluded auth/stream routes |
| Bottleneck diagnosis and optimization | Fixed task workload, traces, profiles, query plans | S06 | Identify the limiting resource and compare before/after with equal workload |
| HTTP | Versioned product API | S01 onward | Handle errors, deadlines, pagination, retries, and conflicts consistently |
| gRPC | Membership check | S02 | Explain schema compatibility and deadline propagation |
| WebSocket | Live inbox IDs and AI stream | S05, S08 | Recover after disconnect and bound slow-client memory |
| MQTT | Scoped notification test subscriber | S05 | Enforce broker ACL and handle QoS duplicates |
| MCP | Agent-facing Task/Search tools | S07 | Pass tool/schema/auth tests and preserve the same business policy as HTTP |
| AI tools and structured outputs | Typed task actions | S07 | Reject malformed arguments and invented IDs without a side effect |
| LangGraph | Checkpointed workflow and approval interrupt | S07 | Resume after restart without silently reapproving or repeating a write |
| Short- and long-term memory | Conversation context and opt-in preferences | S08 | Delete memory and prevent revoked content returning through a summary |
| RAG | Text ingestion, retrieval, citations | S08 | Evaluate recall and citation correctness against a versioned corpus |
| Agent security and evaluation | Prompt injection, scope, budgets, quality | S07–S09 | Prove no unapproved mutation and report quality/cost/latency separately |
| OpenTelemetry | Trace propagation and async span links | S01–S08 | Follow a request through outbox, worker, and tool call |
| Prometheus and Grafana | Service/worker/AI metrics and dashboards | S04–S09 | Detect lag/saturation with bounded-cardinality metrics |
| Docker and Compose | Reproducible, selective local runtime | S01 onward | Bootstrap clean volumes and fit the selected profile in its budget |
| Kubernetes | Deployments, Jobs, CronJobs, stateful storage | S10 | Roll out and terminate pods without corrupting work |
| CI/CD | Tested immutable builds and gated rollout | S11 | Diagnose a failed pipeline and roll back a schema-compatible application |

## Advanced extensions

After required gates pass, compare single-agent and multi-agent designs using the same evaluation corpus. Other optional exercises include social OAuth, an email sender, hybrid retrieval/reranking, a dedicated policy service, automated certificate/secrets management, and an external MCP integration. Add an extension only with a concrete learning question and measured completion criterion.
