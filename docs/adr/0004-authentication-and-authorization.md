# RS256 sessions and resource authorization

Status: Accepted. Date: 2026-11-20.

## Context

The user requires authentication, JWT/RS256, RBAC, and ACL. Private projects and delayed search/AI projections mean an edge-only authorization check or embedded role snapshot would permit stale access.

## Decision

Identity issues short-lived RS256 access tokens and rotating opaque refresh tokens. Services validate issuer/audience/type/time/algorithm and use configured JWKS. Workspace owns active membership and role; Task checks current membership through gRPC and combines it with local project ACL. Initial permission checks do not cache positive decisions.

Service calls use mTLS and validated user context. MCP tokens have resource-specific audiences and require scoped delegation when crossing to another internal resource. Logout revokes refresh capability; already issued access tokens have a bounded expiry window. [Security](../SECURITY.md) defines revocation and stream limits.

## Alternatives and consequences

A central authorization service would introduce a separate policy API and state synchronization before it is needed. Gateway-only checks fail to cover workers/internal calls. Encoding roles only in JWTs delays revocation. Canonical checks are easier to reason about but make Workspace latency/availability part of protected Task requests; failures must deny access safely.

## Verification and revisit

Run the role matrix across HTTP, internal calls, search, MCP, and AI. Remove membership with stale projections and valid tokens. Test key rotation and refresh reuse. Revisit a policy service or carefully versioned permission cache only after profiling and an explicit revocation-consistency design.
