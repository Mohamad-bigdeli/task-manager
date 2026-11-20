# Independent storage ownership

Status: Accepted. Date: 2026-11-20.

## Context

The curriculum requires SQL, document-oriented NoSQL, vector databases, and Redis caching. The user selected independent storage families rather than treating key-value caching as the entire NoSQL exercise.

## Decision

Use PostgreSQL for authoritative relational state, MongoDB for AI content/documents/memory, Qdrant for derived embeddings, and Redis for cache/limits. Use separate Redis processes for evictable cache and durable Celery transport. AI also owns PostgreSQL ledgers/checkpoints so approval and execution integrity do not depend on graph replay or vector state.

Share a PostgreSQL process locally, but isolate service databases and credentials. Prohibit cross-service SQL joins and cross-owner migrations. Use explicit lifecycle records and reconciliation for writes spanning AI stores. [Data design](../DATA.md) defines which store is authoritative for each concern.

## Alternatives and consequences

PostgreSQL JSON/vector extensions would reduce runtime components but provide less independent database experience. One shared schema would simplify joins but destroy service ownership. Polyglot storage costs more memory, migrations, backup work, and failure handling; AI dependencies therefore run in a separate profile and derived indexes remain rebuildable.

## Verification and revisit

Test credential isolation, migrations, partial AI ingestion, deletion races, and authoritative-data restore followed by index rebuild. Revisit if measured memory or operational cost prevents meaningful local exercises; a consolidation would require an explicit replacement learning exercise and superseding ADR.
