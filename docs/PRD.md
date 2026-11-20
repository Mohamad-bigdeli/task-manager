# Task Manager product requirements

## Purpose and success

Build a usable personal and team task manager while learning how to design, secure, test, observe, and operate distributed software. Product success means a user can organize work and collaborate without exposing private resources. Engineering success means the developer can explain tradeoffs and demonstrate correctness during failures using repeatable experiments.

The audience is the developer and a small group of local test users. No public launch, billing, enterprise compliance, mobile application, or frontend framework is required. Kubernetes and CI/CD are final learning stages, rather than prerequisites for the first usable slice.

## Vocabulary

| Term | Meaning |
| --- | --- |
| Account | A person's global identity and sign-in credentials |
| Workspace | A boundary for membership and collaboration; either personal or team |
| Personal workspace | A private workspace with one owner and no invitations |
| Project | A collection of tasks inside exactly one workspace |
| Workspace-visible project | Visible to active workspace members; never internet-public |
| Private project | Visible to owners/admins and explicitly added project members |
| Role | Workspace capabilities assigned to a member |
| Project ACL | Explicit list of members permitted to access a private project |
| Task | A work item belonging to exactly one project |
| Notification | A durable in-app record for one recipient |
| Document | User-provided UTF-8 text or Markdown associated with one project |
| AI run | One resumable execution of a user request through the agent graph |

## Releases

**Release 1** delivers local email/password registration, sessions, personal and team workspaces, invitations, projects, tasks, and in-app notifications. The implementation sequence may expose smaller usable slices before the release is complete.

**Release 2** adds full-text search, live notification delivery, practical gRPC and MQTT adapters, measured caching, and failure/performance exercises.

**Release 3** adds project documents and the AI assistant: permission-aware retrieval, tools, approval before mutations, memory, and evaluations.

**Release 4** completes hardening, recovery exercises, Kubernetes, and CI/CD. OAuth social login, email delivery, file attachments beyond text/Markdown, multi-agent orchestration, and outside integrations are optional extensions. Frontend work remains unplanned.

## Functional requirements

| ID | Requirement | Acceptance |
| --- | --- | --- |
| FR-01 | Register, log in, refresh, and log out | Email is unique after trim/lowercase normalization; a bad login is generic; logout invalidates the refresh session |
| FR-02 | Provision a personal workspace | An account-created event provisions exactly one personal workspace despite redelivery; onboarding exposes pending/ready state |
| FR-03 | Create and manage team workspaces | Creator becomes owner; owner can transfer ownership to an active member; the sole owner cannot leave or be removed |
| FR-04 | Invite members | An owner/admin creates a single-use invitation for an email and non-owner role; it expires after 72 hours; only the matching authenticated account can accept it |
| FR-05 | Manage membership | Owner/admin can remove or change eligible members; removal prevents subsequent access even if search or AI projections lag |
| FR-06 | Manage projects | Authorized writers create/rename/archive projects and choose workspace-visible or private visibility; archived projects remain readable and reject task mutations |
| FR-07 | Manage tasks | Authorized writers create/read/update/soft-delete tasks; list supports status, assignee, priority, due date, and stable cursor pagination |
| FR-08 | Assign tasks | A task has zero or one assignee, who must currently be an active workspace member with project access |
| FR-09 | Notify users | Invitation, assignment, and due-reminder notifications appear in a recipient's inbox; read state is durable and belongs to that recipient |
| FR-10 | Search authorized work | Search returns only currently permitted project/task matches; a stale projection never establishes permission |
| FR-11 | Use AI to understand work | Users can ask about authorized tasks and documents; factual retrieved claims identify their sources |
| FR-12 | Use AI to modify work | The assistant proposes concrete tool arguments, waits for user confirmation, then rechecks permission and resource version before execution |
| FR-13 | Manage documents and memory | Project writers upload/delete text documents; users can inspect/delete persistent personal AI memory; deleted content is excluded from retrieval immediately |

Registration does not require an email provider in local development. Invitations are delivered to existing accounts through in-app notifications; for unregistered accounts, the inviter receives a development-only link for manual use. Automated acceptance is permitted only in test fixtures. Email verification and recovery must be added before exposing registration beyond the trusted local environment.

## Role and ACL rules

| Capability | Owner | Admin | Member | Viewer |
| --- | --- | --- | --- | --- |
| Read workspace-visible projects | Yes | Yes | Yes | Yes |
| Read every private project | Yes | Yes | Only if in project ACL | Only if in project ACL |
| Create projects and edit accessible tasks | Yes | Yes | Yes | No |
| Manage a project's visibility and ACL | Yes | Yes | If project creator | No |
| Archive a project | Yes | Yes | If project creator | No |
| Invite/remove Member or Viewer | Yes | Yes | No | No |
| Grant/revoke Admin | Yes | No | No | No |
| Transfer ownership or archive workspace | Yes | No | No | No |

Each workspace has exactly one owner. Admins cannot modify the owner or another admin's role. Only active workspace members can appear in a project ACL. Project creators cannot elevate a Viewer to write access. A project creator who loses membership or project access loses all creator capabilities. Removed assignees remain visible as historical identifiers but receive no new notifications; edits must clear or replace an invalid assignment.

Authorization is deny-by-default. Inaccessible resource IDs return the same not-found behavior as unknown IDs. Project ACL and current workspace role determine permission; task authorship grants no exception. Personal workspaces reject invitations, ownership transfer, and role changes.

## Task behavior

- Statuses: `todo`, `in_progress`, `in_review`, `done`. Any transition is allowed initially; changes are audited.
- Priorities: `low`, `medium`, `high`, `urgent`; default `medium`.
- Required title: 1–200 characters after trimming. Description: plain text, maximum 10,000 characters. HTML is not accepted.
- Due time is optional and stored as an absolute UTC instant. The account stores an IANA timezone for presentation and future scheduling. Past due times are allowed.
- New tasks default to `todo` and unassigned. Updates require the current version to detect conflicting edits.
- A due reminder is generated at the due time for an assigned, unfinished task. A task without an assignee has no due reminder. Completion, deletion, reassignment, and due changes invalidate older pending reminders.
- No recurring tasks, custom workflows, comments, subtasks, task-level sharing, or task moves across projects in Release 1.

## Quality requirements

Security and data isolation are release gates. Critical workflows must survive process restarts and duplicate messages. Notifications and projections are eventually consistent; task writes are acknowledged after the authoritative transaction commits, without waiting for notifications or AI.

The base Compose profile targets approximately 8 GiB of container memory. AI and full observability run in separate profiles and may require stopping other profiles. This is a resource budget to validate, not a measured result.

For the first performance baseline, seed 10 workspaces and 10,000 tasks, run a five-minute mixed read/write workload at 20 requests/second, and record hardware and p50/p95/p99 latency, errors, throughput, and memory. A provisional target is p95 under 500 ms for non-AI task APIs with less than 1% unexpected errors. Authorization failures are tested separately. Change targets only with recorded evidence.

## End-to-end acceptance story

Two users register and receive separate personal workspaces. One creates a team, invites the other, and creates workspace-visible and private projects. A member creates and assigns a task, and the assignee receives one notification despite replaying the event. A Viewer cannot edit it. A removed member cannot retrieve it through HTTP, gRPC-mediated flows, search, MCP, or AI. The AI proposes an update, waits for approval, and either executes once with current permission or reports why the proposal is stale.
