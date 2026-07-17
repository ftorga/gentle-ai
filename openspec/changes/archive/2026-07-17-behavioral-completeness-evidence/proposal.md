# Proposal: Behavioral Completeness Evidence

## Intent and Outcome

Maintainers lack concise proof that a CLI-reviewed change covers its directly observable behavior. Add a candidate-bound contract: behavioral changes carry sufficient proof; behavior-preserving work records a concise non-applicability declaration visible in the final receipt. Missing evidence warns maintainers; R3-confirmed material omissions block delivery until corrected.

## Scope

### In Scope
- CLI `review finalize` / `review validate` only; directly observable user, API, data, safety, permission, state-transition, or failure outcomes.
- One proof reference may satisfy multiple obligations when it sufficiently proves each.
- Compact-authority-first ownership in `CompactState` / `CompactReceipt`, with finalize/validate enforcing it; SDD status remains downstream audit/binding.
- Receipt-visible non-applicability, candidate freshness/replay rejection, and sequential-receipt continuity.

### Out of Scope
- Broader delivery entry points, a new SDD phase/lens, exhaustive matrices, mandatory tests per condition, or unconditional model calls.
- Replacing upstream authority, causality, low-risk zero-lens behavior, receipt freshness/replay, integration contracts, or sequential receipts.
- Treating a warning alone as a delivery block, or allowing follow-ups to bypass an R3-confirmed omission.

## Capabilities

### New Capabilities
- `behavioral-evidence-contract`: candidate-bound behavioral proof or receipt-visible non-applicability.

### Modified Capabilities
- `review-findings-ledger`: finalize/validate and receipts validate behavioral-evidence identity and applicability.
- `sdd-orchestrator-assets`: consume the contract as downstream audit/binding without duplicating authority.

## Approach and Constraints

“Sufficient” means the cheapest concrete proof for a material outcome, not branch-count coverage. The contract complements upstream mechanisms: authority and causality stay upstream; low-risk zero-lens behavior stays unchanged; receipt freshness/replay, integration, and sequential receipt contracts bind evidence rather than reimplement it.

Use one canonical 4R set after apply: this is ordinary bounded high-risk review because review authority/security-like authorization boundaries change. R3 Reliability may confirm an omission. Do not create duplicate budgets. Reserve Judgment Day for genuinely conflicting severe findings.

## Compatibility and Acceptance Direction

Keep legacy receipt compatibility explicit while v3 receipts expose digest/applicability. Acceptance proves: warning-only missing evidence is review-facing; non-applicability is receipt-visible; shared proof refs work; stale/replayed or out-of-sequence evidence cannot authorize delivery; R3-confirmed material omissions block until corrected.

## Affected Areas

| Area | Impact | Description |
|---|---|---|
| `internal/reviewtransaction/` | Modified | Compact state, receipts, freshness and validation |
| `internal/cli/review_facade.go` | Modified | CLI finalize/validate integration |
| `internal/sddstatus/` | Modified | Downstream audit/binding only |
| Review assets/tests | Modified | R3 and receipt regressions |

## Delivery, Risks, and Rollback

Provisional forecast: one focused CLI-review work unit/PR, target ≤400 authored changed lines; `sdd-tasks` makes the final split. With `ask-on-risk`, exceeding the budget requires a user decision; no silent exception.

| Risk | Mitigation |
|---|---|
| Authority duplication | Keep CompactState/Receipt canonical |
| Stale evidence authorizes delivery | Fail closed on candidate/receipt mismatch |

Rollback: revert the isolated work unit; existing authority and receipt validation remain intact.
