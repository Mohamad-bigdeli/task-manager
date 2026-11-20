# Task Manager

An educational, cloud-native task manager for personal and team work. The target architecture uses six hexagonal microservices, Go and Echo for product services, and Python, FastAPI, and LangGraph for the AI Engine.

## Current state

The repository contains early Go scaffolding and an English development specification. The services, containers, and API contracts are not implemented yet. The existing drawings in `docs/` are design inputs; the written architecture below is the current baseline.

## Start here

Read the [documentation index](docs/README.md), then the [product requirements](docs/PRD.md), [architecture](docs/ARCHITECTURE.md), and [development backlog](docs/BACKLOG.md).

Development proceeds through runnable vertical slices. Docker Compose comes first; Kubernetes and CI/CD are final stages. Frontend implementation is outside the current scope.
