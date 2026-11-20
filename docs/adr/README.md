# Task Manager architecture decisions

ADRs capture why a decision was made and when to revisit it. The baseline decisions below are accepted for this educational project; acceptance does not imply the implementation exists. Supersede an ADR with a new record instead of silently rewriting its historical rationale. Small tunable limits can change in their owning specification with evidence.

| ADR | Decision |
| --- | --- |
| [0001](0001-service-boundaries.md) | Six hexagonal services in a monorepo |
| [0002](0002-storage-ownership.md) | Independent data ownership and selected database families |
| [0003](0003-messaging-and-jobs.md) | JetStream business messaging and Celery AI jobs |
| [0004](0004-authentication-and-authorization.md) | RS256 sessions with canonical RBAC/project ACL |
| [0005](0005-protocol-allocation.md) | One concrete use case per required protocol |
| [0006](0006-ai-execution.md) | One resumable AI service with confirmed writes |
| [0007](0007-deployment-progression.md) | Selective Compose before Kubernetes and CI/CD |

## Template for a new decision

Include: descriptive title, status, date, context, decision, alternatives, consequences, verification, and revisit conditions. Link affected requirements, contracts, and prior ADRs. A proposed ADR names the unresolved choice; an accepted ADR states the chosen behavior.
