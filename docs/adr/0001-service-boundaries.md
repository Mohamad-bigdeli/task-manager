# Six hexagonal services

Status: Accepted. Date: 2026-11-20.

## Context

The project must teach microservices and hexagonal architecture. The initial diagram contains generic Core, independent authentication/authorization boxes, and two AI boxes without defined data ownership. The local runtime has a limited budget.

## Decision

Use Identity, Workspace, Task Management, Notifications, Search, and AI Engine as the six service owners in one repository. Five use Go/Echo; AI uses Python/FastAPI/LangGraph. Each has independent processes, ports/adapters, migrations, and data ownership. API/worker separation inside an owner does not create another domain service.

Domain and application code depend on ports, never framework or remote-service internals. Shared technical libraries must not become a shared domain model. Replace the generic Core responsibility with explicit Workspace and Task boundaries.

## Alternatives and consequences

A modular monolith would simplify local development but would not satisfy the agreed microservices constraint. Splitting authorization or AI into additional services increases network coordination before their independent responsibilities are demonstrated. Keeping both logical and deployment boundaries explicit provides distributed-system learning while limiting unnecessary services.

## Verification and revisit

Review imports, run service-specific tests, and prove that a service cannot read another's database. Revisit when a component has independently owned policy, scaling, security, or release needs. [Architecture](../ARCHITECTURE.md) defines the current ownership map.
