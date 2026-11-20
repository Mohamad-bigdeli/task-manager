# Task Manager service contracts

## Contract status and ownership

These are implementation requirements and endpoint inventories. Executable schemas are deliverables of the corresponding slices, not existing generated artifacts. Store HTTP contracts under `contracts/http/<service>.yaml`, gRPC under `contracts/proto`, event schemas under `contracts/events`, and MCP definitions under `contracts/mcp`.

The owner changes its contract and consumer compatibility tests together. Use `/v1` for product HTTP routes, `workspace.v1` for Protobuf packages, and versioned event subjects. Additive optional fields are compatible; removing fields, changing meaning, or making optional fields mandatory requires a new version and a migration window. Protobuf field numbers are never reused. Freeze actual protocol/library versions when each adapter is introduced.

## HTTP conventions

- JSON request/response bodies; UUID IDs; RFC 3339 UTC timestamps; unknown fields rejected on mutation requests.
- List responses: `{"items":[],"next_cursor":null}`. Default page size 20, maximum 100. Cursors encode ordering, scope, and filters and are validated; they are not permission tokens.
- POST creation returns `201` and `Location`; asynchronous job creation returns `202` and a status URL; successful deletion returns `204`.
- Create operations and remote tool mutations accept `Idempotency-Key`; require it for AI writes, invitation creation, and job/run creation. Clients retry network-ambiguous mutations only with the same key. Updates carry `expected_version`.
- Default JSON body limit 256 KiB; document uploads have a separate 1 MiB text/Markdown limit. Reject oversized bodies before decoding.
- `X-Request-ID` is generated or validated; W3C trace context is propagated. Error responses never expose stacks, SQL, secrets, or internal hostnames.

Example error:

```json
{
  "error": {
    "code": "version_conflict",
    "message": "The task changed. Fetch it and retry with the current version.",
    "request_id": "example-request-id",
    "details": {"resource_type": "task"}
  }
}
```

| HTTP status | Stable code family | Retry behavior |
| --- | --- | --- |
| 400 / 413 / 422 | `invalid_request`, `body_too_large`, `validation_failed` | Fix input |
| 401 | `unauthenticated` | Refresh or sign in; do not loop indefinitely |
| 403 | `forbidden` | Known workspace-level capability denied |
| 404 | `not_found` | Unknown or inaccessible resource |
| 409 | `version_conflict`, `idempotency_conflict`, `operation_in_progress` | Refresh state or follow retry hint |
| 429 | `rate_limited` | Honor `Retry-After` |
| 500 | `internal_error` | Investigate with request ID |
| 503 / 504 | `dependency_unavailable`, `deadline_exceeded` | Bounded retry when operation is safe |

Domain errors map consistently across adapters. gRPC uses corresponding canonical status codes; WS/MCP send structured errors and do not serialize internal exceptions.

## Endpoint inventory

| Owner | Endpoint family | Principal and behavior |
| --- | --- | --- |
| Identity | `POST /v1/auth/register`, `/login`, `/refresh`, `/logout` | Registration/login public; refresh/logout use session credential |
| Identity | `GET /v1/me`, `GET /v1/me/onboarding` | Current user and personal-workspace provisioning state |
| Workspace | `GET/POST /v1/workspaces` | List memberships/create team |
| Workspace | `GET/PATCH /v1/workspaces/{id}` | Read/rename/archive subject to role |
| Workspace | `POST /v1/workspaces/{id}/ownership-transfer` | Owner transfer, atomic membership change |
| Workspace | `GET /v1/workspaces/{id}/members`; `PATCH/DELETE .../members/{user_id}` | Membership and role management |
| Workspace | `POST /v1/workspaces/{id}/invitations`; `DELETE .../invitations/{invitation_id}`; `POST /v1/invitations/accept` | Create/revoke/accept invitation; token in request body |
| Task | `GET/POST /v1/workspaces/{id}/projects`; `GET/PATCH /v1/projects/{id}` | List/create/read/update/archive project |
| Task | `GET/POST /v1/projects/{id}/members`; `DELETE .../members/{user_id}` | Private-project ACL |
| Task | `GET/POST /v1/projects/{id}/tasks`; `GET/PATCH/DELETE /v1/tasks/{id}` | Task CRUD and filters |
| Notifications | `GET /v1/notifications`; `PATCH /v1/notifications/{id}` | Own inbox; set `read` boolean |
| Search | `GET /v1/search?q=...&workspace_id=...` | Authorized project/task search |
| AI | `POST /v1/conversations`; `GET/DELETE /v1/conversations/{id}`; `POST .../runs` | Own conversation and run lifecycle |
| AI | `GET /v1/runs/{id}`; `POST .../approvals/{approval_id}/decision` | Run state; approve/reject current proposal |
| AI | `POST /v1/projects/{id}/documents`; `GET/DELETE /v1/documents/{id}` | Project documents and lifecycle |
| AI | `GET/POST /v1/me/memories`; `DELETE /v1/me/memories/{id}` | Explicit personal memory management |

Internal Task batch endpoints authorize/hydrate up to 100 project/task IDs and check reminder eligibility for `(task_id, due_version, assignee, due_at)`. They require service authentication and delegated user context where applicable. User-facing clients cannot call service-only routes. No internal endpoint returns unrestricted data based only on a submitted actor ID.

## Event and command contracts

Use an application envelope with these fields:

```json
{
  "event_id": "uuid",
  "event_type": "task.updated",
  "schema_version": 1,
  "producer": "task",
  "occurred_at": "2026-11-20T12:00:00Z",
  "aggregate_id": "uuid",
  "aggregate_version": 3,
  "workspace_id": "uuid",
  "actor_id": "uuid",
  "correlation_id": "uuid",
  "causation_id": "uuid-or-null",
  "data": {"task_id": "uuid", "project_id": "uuid"}
}
```

Global identity events allow null workspace IDs; system activity allows a null actor with an explicit system origin. `traceparent` travels in message headers. Events contain the minimum fields consumers need and never credentials. Task projection events carry a complete authorized-to-service indexing snapshot plus the aggregate version; clients never subscribe directly to raw business subjects.

| Subject | Owner | Consumers |
| --- | --- | --- |
| `tm.evt.identity.account-created.v1` | Identity | Workspace provisioning/identity projection; Notifications routing |
| `tm.evt.workspace.invitation-created.v1` | Workspace | Notifications |
| `tm.evt.workspace.membership-changed.v1` | Workspace | Task cleanup, Search/AI invalidation, Notifications suppression |
| `tm.evt.workspace.archived.v1` | Workspace | Task write suppression checks and derived cleanup |
| `tm.evt.task.project-changed.v1` | Task | Search, AI scope invalidation |
| `tm.evt.task.task-created.v1` | Task | Search, AI, Notifications |
| `tm.evt.task.task-updated.v1` | Task | Search, AI, Notifications |
| `tm.evt.task.task-deleted.v1` | Task | Search, AI, Notifications cancellation |
| `tm.cmd.notifications.evaluate-reminder.v1` | Notifications scheduler | Notifications worker group |

Configure separate durable consumers per use case. Allow at most five processing attempts, including the first immediate attempt; redelivery delays are 1 s, 5 s, 30 s, and 2 min. Distinguish permanent validation failures from transient failures. Application quarantine handling durably stores the original event, error classification, and consumer before acknowledging the original. Observe maximum-delivery advisories and reconcile stranded messages; do not assume a broker automatically moves them to a dead-letter queue.

Consumers ignore stale aggregate versions and retain deletion tombstones. Do not infer global sequence from timestamps. Different aggregates may be processed concurrently. Assignment notification uniqueness is `(task_id, assignment_version, recipient)`; reminder uniqueness is `(task_id, due_version, recipient)`. Replay requires an operator-selected consumer, scope, and dry-run report before side effects are enabled.

## gRPC membership checks

Workspace exposes `MembershipService.CheckMembership` in package `workspace.v1`. Input contains workspace ID and target user ID. Output contains active state, role, workspace state, and membership version. The authenticated caller may check only within the scope allowed by its delegated actor or service purpose; it cannot obtain an arbitrary directory of users.

Initial client deadline is 500 ms within a two-second product-request budget. Retry a read once with jitter only when the remaining deadline permits. Workspace outages map to unavailable, never an implicit Viewer/Member role. Use contract tests for missing membership, inactive workspace, invalid IDs, deadlines, and unauthenticated calls.

## WebSocket delivery

Notifications hosts `/v1/notifications/stream`; AI hosts `/v1/runs/{id}/stream`. Initial non-browser test clients authenticate with bearer headers. Browser authentication design is deferred with the frontend; tokens must not be placed in URLs.

Frames contain `type`, `sequence`, `resource_id`, and a minimal payload. A notification frame carries its ID, not private task text. The inbox is authoritative; reconnect fetches missed records over HTTP. AI streams require current run ownership and project permission and stop on revocation, expiry, or failed permission revalidation. Bound per-client queues; disconnect slow clients with a resumable status rather than growing memory indefinitely. Send heartbeat/pong checks and drain on shutdown. Disable gzip/compression on token and streaming responses.

## MQTT delivery

A Mosquitto broker and a protocol test client run only in the `mqtt` profile. Notifications bridges durable inbox IDs to `tm/users/{user_id}/notifications` at QoS 1, with no retained private payload. Broker credentials are provisioned per test user, and ACLs allow subscription only to that user's topic; wildcard subscription is denied. Service credentials alone publish. Provisioning and revocation are a scripted development exercise, not arbitrary user-selected topic names. Reconnect synchronizes the HTTP inbox; MQTT duplication does not create new notifications.

## MCP tools

Task and Search expose Streamable HTTP MCP adapters at `/mcp`. Workspace tools follow only when the product needs them. Start with `search_tasks`, `get_task`, `create_task`, and `update_task`; mutations use the same application use cases, validation, expected versions, authorization, and idempotency as HTTP.

The in-product AI agent may call read tools directly. A proposed mutation becomes an approval record before the execution node can invoke the tool. The approved argument digest and execution key are bound to the call. Direct non-agent MCP clients remain subject to normal user scopes and permissions; the agent's confirmation requirement does not pretend to be an inherent guarantee of the MCP transport.

HTTP MCP authentication follows its [authorization specification](https://modelcontextprotocol.io/specification/2026-07-28/basic/authorization), including resource-specific audience validation and discovery. Implement the required OAuth authorization-code/PKCE flow for interactive clients in S07; do not advertise full OAuth or MCP authorization compliance before conformance tests pass. Incoming MCP tokens are never passed unchanged to a different resource server. See [security](SECURITY.md) for delegated internal calls.

Identity owns the S07 authorization-server additions: authorization-server metadata, `/oauth/authorize`, and `/oauth/token`, with authorization-code/PKCE and refresh grants for pre-registered development clients. Resource servers publish protected-resource metadata and challenges. Restrict redirect URIs to exact registered values and codes to a short-lived, single-use, client/redirect/PKCE-bound transaction. A minimal server-rendered sign-in/consent page supports this protocol exercise; it does not select or implement the task-manager frontend. Its transaction cookie requires the Origin/CSRF protections defined in the security design. Dynamic client registration is outside the initial client interoperability target; document the supported client profile.
