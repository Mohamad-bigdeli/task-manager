# Resumable AI with confirmed writes

Status: Accepted. Date: 2026-11-20.

## Context

The assistant must teach agents, RAG, memory, and LangGraph while operating the task manager. It may perform reads directly but must obtain confirmation before mutations. Restart and remote-write ambiguity must not cause silent duplicate changes.

## Decision

Use one Python AI service with separate API/worker processes, LangGraph checkpoints, a transactional approval/execution ledger, explicit structured proposals, and idempotent MCP writes. Keep deterministic authorization/approval nodes outside model judgment. Bind approval to actor, exact arguments, resource versions, and expiry. Recheck permission before execution.

Use MongoDB for content/memory and Qdrant for derived retrieval, with authoritative lifecycle tombstones and canonical permission checks. Use deterministic fake models in required test gates; live evaluation is opt-in with explicit model selection and budget.

## Alternatives and consequences

A read-only assistant would simplify safety but miss the agreed task-management capability. An autonomous writer contradicts the approval choice. Two AI domain services add contracts without a demonstrated independent owner. A checkpoint alone does not coordinate a remote mutation; the execution ledger and remote idempotency are additional required state.

## Verification and revisit

Test interruption, duplicated approval, altered arguments, revoked access, stale version, prompt injection, and restart after remote commit. Evaluate sources, retrieval, latency, and cost independently. Multi-agent design is an optional measured extension after the baseline. See [AI design](../AI.md).
