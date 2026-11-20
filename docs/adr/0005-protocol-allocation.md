# Allocate protocols by use case

Status: Accepted. Date: 2026-11-20.

## Context

HTTP, WebSocket, gRPC, MQTT, and MCP must all be used somewhere in the project. The user explicitly does not require every service to expose every protocol.

## Decision

Use HTTP for product APIs, Workspace gRPC for membership checks, WebSocket for inbox signals and AI streams, MQTT for a scoped notification test subscriber, and Task/Search MCP for agent tools. Every inbound adapter calls the same owning application policy; protocols are not independent implementations of business rules.

WebSocket/MQTT inbox signals contain IDs, with HTTP as the reconnect/synchronization source. MQTT requires real broker ACL tests. HTTP MCP follows the [authorization specification](https://modelcontextprotocol.io/specification/2026-07-28/basic/authorization); generic JWT support alone is insufficient for an interoperability claim.

## Alternatives and consequences

HTTP alone is simpler but misses the curriculum. Replicating CRUD across every transport multiplies compatibility/security work without a product reason. Limited concrete uses support comparisons while requiring real transport-specific tests. MQTT is an educational client integration, not a fabricated IoT domain.

## Verification and revisit

Exercise every transport over a real connection, including deadlines, reconnect, audience validation, and denied cross-user MQTT subscriptions. Revisit allocation when a product requirement calls for another transport. Executable schemas follow [contracts](../CONTRACTS.md).
