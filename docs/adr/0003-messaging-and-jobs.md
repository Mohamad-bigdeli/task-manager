# JetStream messaging and Celery jobs

Status: Accepted. Date: 2026-11-20.

## Context

The system must teach message-driven/event-driven design, workers, retries, idempotency, and Celery. The supplied diagram already suggests NATS. Business events must survive disconnected consumers and process restarts.

## Decision

Use NATS JetStream for cross-service facts and targeted work commands. Core NATS alone has no durable replay; JetStream supplies persistence and redelivery, as described in its [official documentation](https://docs.nats.io/concepts/jetstream). Design consumers for duplicates rather than claiming globally exactly-once effects.

Use PostgreSQL outbox/inbox and business uniqueness constraints. Implement explicit quarantine and controlled replay; external side effects require a durable execution/delivery ledger. Use Celery with its [Redis transport](https://docs.celeryq.dev/en/stable/getting-started/backends-and-brokers/redis.html) only inside AI, plus a durable ingestion ledger/reconciler. Do not assume Celery can use the selected NATS server as its supported broker.

Go schedules due reminders from persisted deadlines. Celery Beat schedules AI reconciliation; one active Beat owns a schedule set, consistent with [Celery's scheduler guidance](https://docs.celeryq.dev/en/stable/userguide/periodic-tasks.html).

## Alternatives and consequences

RabbitMQ could support Celery and broader messaging through one broker, but would replace the diagram's NATS learning path. Kafka would add an unnecessary operational model for the current flows. Database polling alone misses the required broker exercises. Two messaging mechanisms require separate monitoring and recovery, but their ownership and purpose remain explicit.

## Verification and revisit

Kill relays/consumers at commit boundaries, duplicate/reorder events, exhaust retries, and recover a lost Celery dispatch. Revisit broker consolidation if maintaining both mechanisms becomes the dominant learning cost. [Contracts](../CONTRACTS.md) defines delivery rules.
