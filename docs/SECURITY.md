# Task Manager security design

## Trust boundaries and threats

Protected assets include credentials, signing keys, workspace content, invitation tokens, model API keys, and user memory. Treat clients, retrieved documents, model output, and message payloads as untrusted inputs. Container networking alone does not authenticate a service.

| Threat | Required control | Verification |
| --- | --- | --- |
| Credential theft or token forgery | Password hashing, short-lived access tokens, refresh rotation, strict JWT validation | Invalid signature/issuer/audience/algorithm tests and refresh reuse test |
| Cross-workspace access | Scoped queries, current membership, project ACL on every adapter | Two-tenant adversarial fixtures, including lists and batch APIs |
| Stale authorization | Canonical checks before data disclosure or execution | Remove membership while projections and tokens remain unchanged |
| SQL/NoSQL injection and malicious text | Parameterized queries; allowlisted filter operators; plain-text fields; escaped rendering | Injection fixtures and rejection of client-supplied Mongo operators |
| CSRF and cross-origin abuse | Explicit credential mode and origin policy | Cross-origin mutation, preflight, and credential tests |
| Retry abuse and duplicate effects | Idempotency ledger and uniqueness constraints | Concurrent duplicate requests and crash recovery |
| Prompt injection or confused deputy | Treat retrieval as data, tool allowlists, delegated scopes, confirmed argument binding | Malicious documents cannot expand authority or approve a mutation |
| Resource exhaustion | Body/time/concurrency limits, distributed rate limits, bounded queues | Slow-client, large-input, pool exhaustion, and model budget tests |
| Secret leakage | Mounted secrets, response/log redaction, scoped credentials | Log capture tests and secret scanning in S11 |

## Passwords and sessions

Use Argon2id with a random per-password salt and a well-maintained implementation. Start from the [OWASP password-storage guidance](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html), then record the chosen memory/time parameters and measured login cost on the development host. Permit passphrases of 12–128 characters; never silently truncate them. Store hashes outside public account DTOs and rehash on login when parameters change.

Initial access-token lifetime is five minutes. Refresh tokens are random opaque secrets, stored only as hashes, valid for at most seven days, and rotated atomically on each refresh. Reuse of a consumed token revokes its family and emits a security log. Logout revokes the refresh family. An already issued access token may remain valid until expiry; membership removal is checked independently on each protected request. Do not promise instantaneous global logout with offline JWT verification.

Initial product API clients transmit access tokens with `Authorization: Bearer` and refresh tokens in JSON bodies over TLS, except explicit loopback-only local HTTP. They store credentials outside logs and URLs. No product-session cookies or browser token storage are selected before frontend design. S07's minimal OAuth authorization page uses a separate short-lived transaction cookie with Secure, HttpOnly, SameSite, Origin, and CSRF protections; local HTTP relaxation is limited to loopback development.

## RS256 and service identity

Identity signs access tokens with RS256. Validate a fixed algorithm allowlist, signature, issuer, audience, expiration, not-before, token type, and subject. Reject `none`, algorithm confusion, malformed keys, and untrusted `jku`/`x5u` URLs. `kid` selects only from the configured issuer's JWKS. This policy follows [JWT best current practices](https://www.rfc-editor.org/rfc/rfc8725.html).

Product tokens use the explicit audience `task-manager-api`; first-party services participating in that audience still enforce their own scopes and resource permissions. Claims include `sub`, `iss`, `aud`, `iat`, `exp`, `jti`, and session ID. Workspace roles are not authoritative token claims. Publish public keys through JWKS; private keys are readable only by Identity. Rotate by publishing the new public key before signing with it and retaining the old public key until all valid old tokens and clock-skew allowance expire. Unknown key IDs trigger a bounded refresh, not an arbitrary network fetch per request.

Internal HTTP/gRPC uses mTLS for caller identity, with development certificates generated locally and private keys excluded from Git. Bind service identities to allowlisted operations. For user-triggered calls, also require a signed user credential; a caller cannot choose an actor through an unsigned header.

In S02, a product token may accompany mTLS when checking the same first-party API audience. For MCP's separate resource audiences in S07, Identity adds `POST /internal/delegations`: an authenticated service presents the subject token and requested internal scope; Identity validates the original token, actor, caller, and permitted audience transition, then issues a short-lived delegated assertion. Preserve user subject and calling-service identity, narrow scopes, and cap expiry at the original credential's expiry. An MCP resource token is not forwarded unchanged to Workspace or Task. This is a first-party delegation contract, not a claim to implement a general token-exchange standard.

## Authorization and revocation

Use the role matrix in [PRD](PRD.md). Workspace checks active membership and workspace state; Task combines that result with its own project ACL and archive state. Internal service credentials authenticate the caller but never grant all users' data by default. Workers have purpose-scoped access for their task type.

Do not cache positive permission decisions initially. Payload caches may be used only after a fresh permission check. Lists, counts, search ranking/snippets, source citations, exports, and streams are data disclosures and require the same policy as direct reads. Search overfetches candidate IDs, authorizes/hydrates them through Task, then builds the response without exposing raw projection snippets or hidden totals.

AI rechecks permissions before retrieval enters the model, before returning a stored response, and before every tool execution. Conversation ownership does not permanently grant access to previously retrieved project content. Store source dependencies for grounded messages; deny access to those messages while dependencies are no longer readable. A conversation with revoked dependencies becomes blocked for generation until the user starts a new conversation with permitted context. Stream a run in one selected project scope and revalidate at least every five seconds, stopping immediately when a revocation event arrives. Data already sent to a client or model cannot be recalled; an operation authorized before a concurrent revocation may finish.

MQTT carries only opaque notification IDs. WebSocket notification frames do the same; fetching private task text requires a fresh API check. Revoke test MQTT credentials when a user/session is disabled. A user's historical inbox may retain a generic notification while links to removed projects return not found.

## CORS and CSRF

CORS defaults to no browser origins. When a local browser client is introduced, use an explicit origin allowlist; never combine wildcard origin with credentialed requests. CORS is not authentication. Trust proxy headers only from configured proxies; direct loopback deployments do not trust arbitrary forwarded IPs.

Bearer credentials sent explicitly in headers are the initial product mode, so cookie-based CSRF middleware is not enabled blindly. Add a controlled learning exercise comparing that mode with a cookie session. If browser cookies are adopted, define Secure/HttpOnly/SameSite policy, validate Origin, add a session-bound CSRF token for mutations, and reject state changes over GET. OAuth authorization endpoints separately require their flow protections, including PKCE and state, even before the product adopts cookies.

## Limits and sanitization

Start with 100 requests/minute per authenticated user, 10 login attempts/minute per IP, and an additional account/IP login bucket. Model runs start at five/minute/user and at most two concurrent runs/user. Return `429` with a retry interval. These are configurable learning defaults, not production capacity claims.

The distributed limiter uses Redis atomic operations. If it is unavailable, login, token issuance, AI runs, and other costly/security-sensitive routes return retryable unavailability. Ordinary reads use a bounded local emergency limiter. Data-cache outages do not bypass those controls. Password-hash work has a separate concurrency bound.

Validate identifiers, enum values, length, timestamps, expected versions, and filter fields. Parameterize SQL. Reject raw MongoDB query objects and unsupported content types. Document ingestion accepts only bounded UTF-8 text/Markdown; file names never become filesystem paths. The agent has no general shell, arbitrary URL fetcher, or unrestricted filesystem tool. Adding such capabilities requires an ADR and targeted tests for SSRF and sandbox escape.

## Secrets and external model data

Use environment variables for nonsecret configuration and mounted files for passwords, keys, and provider tokens. Commit only placeholder examples. Each service has its own database/broker credentials; migrator credentials are separate. Secrets must be rotatable without rebuilding images.

Only authorized, relevant excerpts are sent to the configured model API. The developer selects a provider and model explicitly before live evaluation; record data-retention settings, price source, and token budgets in the evaluation report. Use synthetic data by default. Live external calls are opt-in in tests and never part of deterministic unit/integration gates.

Log authentication outcomes, role/ACL changes, approval decisions, and execution IDs. Never log passwords, full tokens, invitation links, raw document content, or model prompts/responses by default. Audit records identify the actor, action, resource IDs, outcome, and time without copying private text.
