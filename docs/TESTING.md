# Task Manager test strategy

## Test layers

Unit tests cover domain invariants and application decisions with fake ports and a controllable clock. Integration tests run real owned databases, NATS, Redis, and applicable adapters in isolated test containers. Contract tests validate executable HTTP/event schemas, Protobuf compatibility, and MCP behavior. System E2E tests run multiple service processes and real transports; they do not require a frontend.

Use Go's testing tools and race detector, Python pytest, Testcontainers or isolated Compose test projects, and k6 for workload experiments. Never run destructive fixtures against a developer's regular volumes. Give each run unique database/project names and clean up only resources it created. Pin tool versions when their slice starts.

Avoid arbitrary sleeps: poll observable outcomes with bounded deadlines and capture traces/logs on failure. Unit tests use fake clocks for expiry/reminders. Tests of retry behavior assert both bounded attempts and the eventual effect. Use fake model and embedding clients for reproducible gates; paid live evaluations remain opt-in.

## Acceptance matrix

| Case | Requirements | Scenario and expected outcome |
| --- | --- | --- |
| T01 | FR-01 | Unique registration, wrong password, rotation, logout, refresh-token reuse; no credential appears in a response/log |
| T02 | FR-02 | Deliver account-created twice and restart Workspace between deliveries; one personal workspace, pending then ready onboarding |
| T03 | FR-03, FR-05 | Owner transfer is atomic; sole owner cannot leave; personal workspace rejects invitations and role changes |
| T04 | FR-04 | Correct/wrong email acceptance, expiry, revoked link, concurrent acceptance, member duplication; one membership and consumed token |
| T05 | FR-05, FR-06 | Exercise every role against visible/private projects; test creator removal, Viewer ACL, admin boundaries, and unknown/inaccessible IDs |
| T06 | FR-07, FR-08 | Task defaults/validation, assignee eligibility, archive behavior, filters/pagination, two concurrent versioned updates; one winner |
| T07 | FR-09 | Duplicate assignment event and crash before acknowledgment; exactly one durable notification for that assignment |
| T08 | FR-09 | Fake clock and two schedulers; due changes, reassignment, completion, deletion, downtime catch-up; no obsolete reminder; title-only edits do not create new assignment/reminder generations |
| T09 | FR-10 | Stale/out-of-order events and index rebuild; versions/tombstones prevent resurrection; canonical access controls snippets/counts |
| T10 | FR-11, FR-13 | Ingest/delete/retry documents, isolate vector and memory namespaces, cite current sources, handle missing evidence |
| T11 | FR-12 | Read-only answer requires no mutation approval; model-proposed write cannot execute before confirmation |
| T12 | FR-12 | Reject/expire/tamper proposal, duplicate approval, concurrent resume, revoke access or change resource before execute; no unauthorized write |
| T13 | FR-12 | Kill worker after remote write before checkpoint; resume same key and recover same result without another mutation |
| T14 | FR-11, FR-13 | Remove project access while conversation is stored/streaming; block grounded history and future prompts using revoked sources |

## Transport and security cases

Test real HTTP middleware composition, including panic recovery, status/timing logs, malformed JSON, oversized bodies, content types, gzip negotiation, CORS preflight, and allowed/disallowed origins. Compare explicit bearer and cookie CSRF behavior in the educational fixture; never infer one from the other.

JWT tests cover wrong issuer/audience/signature, expired/not-yet-valid tokens, unsupported algorithm, untrusted key URL, key rotation, and cache refresh bounds. gRPC tests cover mTLS identity, delegated user context, unauthorized target checks, deadline cancellation, and fail-closed Workspace outage. Rate-limit tests cover multiple API replicas and Redis failure policy.

WebSocket tests cover authenticated connect, foreign run ID, heartbeat timeout, disconnect/reconnect, missed notifications, slow consumer bounds, token expiry, and midstream permission removal. MQTT tests use the real broker: own topic succeeds, another user's or wildcard subscription fails, QoS duplication changes no inbox state, and credential revocation closes access.

MCP tests include initialize/tool discovery, JSON schema validation, authenticated/unauthenticated calls, resource-specific audience, PKCE/discovery flow, delegated calls, tool errors, and domain-rule equivalence with HTTP. Test that a raw foreign MCP token cannot be passed to another resource and that tool annotations are not authorization controls.

## Failure injection

Stop NATS after a business commit; restore it and drain the outbox. Kill a consumer before/after database commit and before acknowledgment. Deliver older aggregate versions after a deletion. Exhaust delivery attempts and verify durable quarantine plus an operator-visible alert. Simulate a failing external call after local deduplication to verify the effect ledger rather than only inbox behavior.

Stop cache Redis, Celery Redis, MongoDB, Qdrant, and model responses separately. Verify isolation: core task writes still work during AI failure; expensive routes fail safely during limiter loss; ingestion reconciliation recovers a lost dispatch. Simulate connection-pool exhaustion and deadline expiry without goroutine/task leaks. Force termination after shutdown grace and confirm safe restart.

Restore authoritative data into new volumes and rebuild Search/Qdrant. Test migrations from empty and prior schema, partial/resumed Mongo backfills, old/new application compatibility during expand-contract, and an intentionally failed migration. No destructive rollback is considered safe without a restoration test.

## Performance and bottleneck experiments

Use the PRD workload and record a reproducible seed, environment, image versions, CPU/RAM limits, warmup, concurrency, sample count, and cache state. Save p50/p95/p99, successful throughput, unexpected error rate, memory, CPU, pool wait, query count, event lag, and model timing where applicable. A result from another host is not a local measurement.

Begin uncached. Diagnose one bottleneck from traces and profiles, then compare one intervention at a time: missing index, N+1 elimination, pool tuning, payload compression, cache-aside, or worker concurrency. Use Go pprof and the corresponding Python profiler. Inspect [PostgreSQL query plans](https://www.postgresql.org/docs/current/sql-explain.html); run mutating `EXPLAIN ANALYZE` only in a disposable database or an intentionally rolled-back experiment.

Security/correctness tests are hard gates. Performance targets are initial objectives and must produce an evidence-backed report, including regressions and the next limiting resource. A cache that accelerates a request while leaking data fails the slice.

## Completion evidence

Each slice records commands actually run, tool/image versions, passing/failing cases, representative redacted logs/traces, migration status, and known limits in its implementation PR or work log. Do not mark dashboards, benchmarks, restore procedures, or live-model quality as verified merely because their files exist.
