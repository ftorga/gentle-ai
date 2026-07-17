# Design: Behavioral Completeness Evidence

## Technical Approach

Design against `origin/main` (`9c7bac8`): compact state and receipts are v2. `CompactState`, receipt construction, store replay, and `EvaluateCompactGate` remain authoritative; `sddstatus.ReviewBinding` only audits after `GatePostApply` allows.

## Architecture Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Authority | `CompactState` retains canonical evidence; `CompactReceipt` exposes digest/applicability. | Preserves causality, integration, zero-lens, freshness/replay, and sequential receipts. |
| Completeness | Go enforces shape, canonical order/digest, candidate tree, and paths digest; R3 judges materiality. | No branch-count matrix or parallel policy engine. |
| Missing evidence | Absent optional input produces a machine-readable finalize/validate warning, not a denial. Invalid supplied input fails closed. | A warning is not authority; an R3-confirmed material omission uses the existing blocking finding/correction lifecycle. |
| Review | One high-risk ordinary bounded 4R after apply; R3 confirms omissions. | This changes review authorization boundaries. Judgment Day replaces 4R and is only legitimate for explicit, conflicting severe findings needing two blind judges. |

## Data Flow

```text
candidate + evidence -> finalize -> CompactState.Validate -> store revision
                                  -> CompactReceipt v3 -> validate/gate -> allow or denial
R3 material omission -> existing correction_required -> atomic correction + replacement evidence
                                                     -> SDD binding audits receipt only
```

`activated` has canonical material obligations and proof references (one may serve many). `non_applicable` has a concise basis and no obligations; v3 visibly carries its digest and applicability.

Correction must be atomic: extend the existing correction transition to receive the new snapshot and replacement evidence in one `CompactStore.Replace` record. It clears/rejects the prior value, validates the replacement against the new candidate/path digest, then enters validation. No subsequent mutation may rebind evidence; terminal/replay receipt equality and gate double-read require the exact stored state. A correction without a valid replacement cannot authorize the new candidate.

## Interfaces / Contracts

```go
type BehavioralEvidence struct {
  Schema, Applicability, Basis, CandidateTree, PathsDigest string
  Obligations []BehavioralObligation
}
```

`CompactStateSchema` stays v2: old states without the optional field retain their lifecycle. `CompactReceiptSchema` gains v3: v2 receipts parse/validate/gate unchanged and are never rewritten; v3 MUST contain a canonical digest and applicability matching state. New contract-aware terminal writes are v3. Unknown schemas/fields fail closed; no eager migration. Store/export/import and successors compare the derived receipt, preserving replay and continuity.

`review finalize --behavioral-evidence <file|->` accepts one canonical input and returns warning metadata when omitted. `review validate` surfaces the same authoritative missing-evidence warning from state/receipt while retaining existing gate result precedence. An R3 BLOCKER/CRITICAL for a named, candidate-causal unproved material obligation enters the normal frozen-ledger correction path and blocks; speculative/non-material observations remain `WARNING|SUGGESTION`/follow-up.

## File Changes

| File | Action | Ownership |
|---|---|---|
| `internal/reviewtransaction/behavioral_evidence.go` | Add | Canonical contract, digest, structural validation. |
| `internal/reviewtransaction/compact.go`, `compact_store.go`, `compact_gate.go` | Modify | State/receipt v3, atomic rebinding, store/gate freshness and replay. |
| `internal/cli/review_facade.go`, `review.go` | Modify | Finalize input and validate/finalize warning JSON. |
| `internal/sddstatus/review_binding.go` | Modify only if needed | Downstream receipt audit; no state mutation or authorization. |
| Focused existing `*_test.go` files | Modify | Contract, compact store/gate, facade, and binding regressions. |

## Testing Strategy

Strict TDD: table-driven RED tests beside compact/facade helpers, implementation, focused packages, then `go test ./...`. Cover canonical/shared/non-applicable, malformed input, v2 read, v3 fields, stale binding/replay/export/successors, warnings, R3 blocking, atomic replacement, CLI JSON, and binding. Use `t.TempDir()` and existing fixtures. v2 golden state/receipt JSON is exact/canonical: add v3 fixtures; do not rewrite legacy fixtures.

## Threat Matrix

| Boundary | Applicability | Response / RED test |
|---|---|---|
| Documentation-like paths | N/A | Evidence JSON is data, not executable classification. |
| Git repository selection | N/A | Existing `--cwd` derivation is unchanged. |
| Commit state | N/A | Existing snapshot/gate behavior is reused. |
| Push state | N/A | Existing gate behavior is reused. |
| PR commands | N/A | No PR command composition changes. |

## Migration / Rollout

No store migration or flag. v2 remains readable and immutable; v3 is written only by the new contract-aware path. Roll back by reverting this work unit; v2 authority remains valid.

## Work Unit and Open Questions

Provisional single slice: compact authority + CLI + focused tests, target ≤400 authored lines, rollback limited to these files. Evidence: focused results, `go test ./...`, receipt fixtures, and one finalize/validate scenario. `sdd-tasks` owns the final forecast; with `ask-on-risk`, >400 requires a decision before apply and a chain—no silent exception.

- [x] Missing evidence uses warning metadata only; receipt applicability remains limited to `activated` and `non_applicable`.
- [x] `review-integration/v1` is a closed negotiated contract. Warning fields require a version bump and are deferred from this authority slice.
