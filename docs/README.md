# Task Manager engineering documentation

This specification connects the product requirements to architecture, development work, and learning evidence. The baseline was established on 2026-11-20 for an educational project operated locally, without a public launch.

## Reading order

| Document | Purpose |
| --- | --- |
| [Product requirements](PRD.md) | Product behavior, permissions, vocabulary, and acceptance criteria |
| [Architecture](ARCHITECTURE.md) | Service boundaries, dependencies, data ownership, and failure flows |
| [Data and migrations](DATA.md) | Logical models, consistency, storage lifecycle, and migrations |
| [Contracts](CONTRACTS.md) | API conventions, protocol allocation, events, and compatibility |
| [Security](SECURITY.md) | Threat model, authentication, authorization, and trust boundaries |
| [AI design](AI.md) | Agent execution, approval, memory, retrieval, and evaluation |
| [Operations](OPERATIONS.md) | Configuration, containers, telemetry, and recovery |
| [Test strategy](TESTING.md) | Acceptance scenarios, failure injection, and performance experiments |
| [Learning roadmap](LEARNING.md) | Coverage of every requested engineering topic |
| [Development backlog](BACKLOG.md) | Ordered slices, dependencies, deliverables, and completion gates |
| [Architecture decisions](adr/README.md) | Decisions, alternatives, and consequences |

## How to interpret the specification

- **Agreed constraint:** educational project; microservices; hexagonal architecture; messaging and events; Go/Echo except the Python AI Engine; personal and team work; project ACL; English documentation; model API; Compose before Kubernetes and CI/CD; no frontend decision yet.
- **Accepted design:** the six services, selected storage families, protocol allocation, and AI confirmation policy recorded in the ADRs.
- **Implementation default:** numeric limits, timeouts, and initial role capabilities in these documents. These make the first implementation unambiguous and can change through a documented decision and updated tests.
- **Future deliverable:** schemas, binaries, dashboards, deployment files, or test results named in the backlog. A description here does not mean that deliverable already exists.

Implementers create executable OpenAPI, Protobuf, JSON Schema, and MCP definitions before implementing the corresponding adapter. Their compatibility rules come from [Contracts](CONTRACTS.md). Changes to product behavior update the PRD and acceptance cases in the same change; architectural tradeoffs receive an ADR.

## Collaboration and completion

The default collaboration mode is mentoring: explain tradeoffs and review work. Implementation assistance is requested per task. A completed slice includes its executable behavior, meaningful tests, operational evidence, and updated documentation. Installing a dependency alone does not demonstrate learning.

The original [diagram](design.drawio) and [UI reference](ui.jpeg) remain available as inputs. Neither defines additional requirements or overrides the written baseline.
