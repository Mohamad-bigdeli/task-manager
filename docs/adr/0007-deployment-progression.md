# Compose before Kubernetes and CI/CD

Status: Accepted. Date: 2026-11-20.

## Context

The project is educational and local, with an approximate 8 GiB container budget and model APIs rather than local inference. Kubernetes and CI/CD are required at the end; a frontend is not currently in scope.

## Decision

Start with selective Compose profiles and add runtime components only when a vertical slice needs them. Keep configuration external, credentials scoped, process state explicit, telemetry bounded, and shutdown recoverable. Run full AI and observability profiles selectively rather than claiming the entire stack fits the base budget.

After failure/recovery and resource measurements, move working images to a local Kubernetes cluster. Add provider-specific CI/CD after selecting the actual Git host and registry in the final slice. Execute manual tests and maintain reproducible commands before that automation exists.

## Alternatives and consequences

Kubernetes-first adds platform debugging before product flows exist. Running everything in one profile strains the agreed budget. Deferring all tests until CI would hide defects; tests and telemetry therefore arrive per slice even though pipeline automation is last.

## Verification and revisit

Bootstrap clean volumes, record selected-profile memory peaks, stop/restart processes during work, then repeat acceptance flows in Kubernetes. Revisit profile budgets based on measurements or changed hardware, without changing domain ownership. [Operations](../OPERATIONS.md) defines the deployment gates.
