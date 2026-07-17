# Apply Progress: Behavioral Completeness Evidence

## Delivery Context

- Artifact store: BOTH
- Execution mode: interactive
- Delivery strategy: ask-on-risk
- Chain strategy: feature-branch-chain
- Current work unit: 3 / Phase 3 — Correction, R3, and Downstream Audit
- Review budget: 400 authored changed lines per chain slice; no size exception

## Cumulative Completed Work

- [x] 1.1 Resolved the contract gates: missing evidence remains warning metadata; `review-integration/v1` remains closed.
- [x] 1.2–1.3 Added compact-authority behavioral evidence, receipt v3 visibility, v2 reads, candidate/path checks, and focused verification.
- [x] 2.1–2.3 Added RED-first CLI evidence input, warning reporting, receipt visibility, and focused finalize/validate verification.
- [x] 3.1 Reject correction completion that would retain prior behavioral evidence without an atomic replacement; preserve state on rejection.
- [x] 3.2 Proved R3 candidate-causal CRITICAL omissions use the existing correction ledger while WARNING observations remain non-blocking; proved v3 receipt binding is audit-only and stable across sequential binding.
- [x] 3.3 Confirmed no shared assets are required by current specs; integration tests and `go vet ./...` completed.

## TDD Cycle Evidence

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|---|---|---|---|---|---|---|---|
| 1.1–1.3 | `internal/reviewtransaction/{behavioral_evidence,compact_store}_test.go` | Unit/temp-dir repository | focused compact suite | contract symbols absent | 443 package tests passed | activated, non-applicable, v2/v3, stale/replay | gofmt |
| 2.1–2.3 | `internal/cli/review_facade_test.go` | Unit/temp-dir Git repository | 47 focused CLI tests passed | facade warning field absent | 52 focused CLI tests passed | malformed/stale, applicability, warning paths | gofmt |
| 3.1 | `internal/reviewtransaction/compact_store_test.go` | Unit/temp-dir repository | 90 focused compact tests passed | `TestCompactCorrectionRejectsRetainedBehavioralEvidenceWithoutReplacement` failed: correction accepted retained evidence | 2 replacement/rejection tests passed | replacement binding and legacy retained-evidence rejection | gofmt; no further refactor needed |
| 3.2 | `internal/reviewtransaction/compact_store_test.go`, `internal/sddstatus/review_binding_test.go` | Unit/temp-dir repository | 90 compact + 41 SDD-status focused tests passed | test-first coverage exposed invalid empty-lens fixture setup; no production gap after a valid R3 fixture was established | R3 table cases and v3 binding test passed | CRITICAL correction vs WARNING follow-up; initial bind vs sequential retry | gofmt; no production refactor needed |
| 3.3 | Existing focused fixtures | Integration boundary | 447 reviewtransaction + 151 sddstatus tests passed | N/A — audit-only behavior was covered by 3.2 | `go test ./...` completed with approved baseline only; `go vet ./...` passed | v2 reads, v3 binding, sequential retry | no shared assets required |

## Work Unit Evidence

| Evidence | Result |
|---|---|
| Work Unit 1 focused/package tests | `go test ./internal/reviewtransaction -run 'BehavioralEvidence|Compact' -count=1` → 90 passed; package → 443 passed |
| Work Unit 1 runtime harness | temp-dir Git-repository correction scenario → passed |
| Work Unit 2 focused/package tests | `go test ./internal/cli -run 'ReviewFacade|Review' -count=1` → 52 passed; package → 510 passed with 2 approved baseline failures |
| Work Unit 2 runtime harness | temp-dir Git-repository finalize/validate JSON scenario → passed |
| Work Unit 3 focused test | `go test ./internal/reviewtransaction ./internal/sddstatus -run 'Binding|Correction|Gate|BehavioralEvidence' -count=1` → 170 passed |
| Work Unit 3 package boundary | `go test ./internal/reviewtransaction -count=1` → 447 passed; `go test ./internal/sddstatus -count=1` → 151 passed |
| Work Unit 3 runtime harness | N/A — the correction ledger and downstream receipt audit are exercised through temp-dir Git-repository fixtures in the focused tests; no separate runtime boundary exists |
| Integration boundary | `go test ./...` → 5073 passed, 5 approved baseline failures, 15 skipped in 57 packages; `go vet ./...` → passed |
| Work Unit 3 rollback boundary | Revert Work Unit 3 hunks in `internal/reviewtransaction/compact.go`, `internal/reviewtransaction/compact_store_test.go`, `internal/sddstatus/review_binding_test.go`, `internal/sddstatus/bounded_review_test.go`, and these task/progress records. This restores prior correction behavior without removing Units 1–2. |

## Approved Baseline

`go test ./...` retains five pre-existing environment-sensitive CodeGraph failures: `internal/cli` has two failures caused by macOS `/var` versus `/private/var` path canonicalization and MCP init timing; `internal/components/communitytool` has three failures caused by the same canonicalization and MCP 100ms timing. They are outside this candidate and were not modified.

## Scope and Next Action

All 9/9 planned tasks are complete. The final planned implementation slice is ready for bounded post-apply review. No broader delivery entry points, integration-contract changes, duplicate 4R/Judgment Day budget, commits, PR operations, or GitHub updates were performed.

## Scoped Verify Remediation

The Verify CRITICAL was corrected without starting a review, lens, Judgment Day, or budget. `RunReviewFacadeFinalize` now reserves supplied behavioral evidence for the correction branch and passes it through `CompleteCorrectionWithBehavioralEvidence`; ordinary initial binding remains unchanged. The CLI regression starts with initial bound evidence, proves a corrected v3 receipt contains the replacement digest and applicability, and proves rejected replacement input leaves the correction record unchanged.

### TDD Cycle Evidence

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|---|---|---|---|---|---|---|---|
| Verify CRITICAL remediation | `internal/cli/review_facade_test.go` | Unit/temp-dir Git repository | `go test ./internal/cli -run 'ReviewFacade|Review' -count=1` → 52 passed | replacement test failed with `behavioral evidence is already bound to the compact authority` | focused regression → 3 passed; focused CLI → 55 passed | valid replacement receipt plus rejected replacement/no mutation | gofmt; no further refactor needed |

### Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused CLI regression | `go test ./internal/cli -run '^TestReviewFacadeCorrectionReplacesBehavioralEvidenceAtomically$' -count=1` → 3 passed |
| Focused CLI suite | `go test ./internal/cli -run 'ReviewFacade|Review' -count=1` → 55 passed |
| Focused compact authority suite | `go test ./internal/reviewtransaction -run 'Correction|BehavioralEvidence' -count=1` → 30 passed |
| Runtime harness | temp-dir Git repository: `review finalize` receives initial evidence, correction validation, replacement evidence, and final evidence → v3 receipt verified; rejected replacement keeps revision and state unchanged |
| Full suite | `go test ./...` → 5,076 passed, 5 approved baseline failures, 15 skipped across 57 packages |
| Static analysis | `go vet ./...` → passed |
| Rollback boundary | Revert only the remediation hunks in `internal/cli/review_facade.go` and `internal/cli/review_facade_test.go`; the existing behavioral-evidence contract and normal initial-binding path remain intact |

```yaml
schema: gentle-ai.remediation-result/v1
lineage_id: review-7f7c4fe5640e2f5e
generation: unavailable-in-verify-artifact
fix_batch: unavailable-in-verify-artifact
failed_evidence_revision: sha256:626a7aff3557b10be60717270f99942a3301a575b511a7e06c153e5ed308e92e
status: implementation-complete-metadata-pending
```

```json
{
  "schema": "gentle-ai.remediation-evidence/v1",
  "lineage_id": "review-7f7c4fe5640e2f5e",
  "generation": "unavailable-in-verify-artifact",
  "fix_batch": "unavailable-in-verify-artifact",
  "failed_evidence_revision": "sha256:626a7aff3557b10be60717270f99942a3301a575b511a7e06c153e5ed308e92e",
  "tests": {
    "cli_regression": "3 passed",
    "cli_focused": "55 passed",
    "reviewtransaction_focused": "30 passed",
    "full_suite": "5076 passed; 5 approved baseline failures; 15 skipped",
    "vet": "passed"
  },
  "authored_native_changed_lines": 146,
  "within_native_correction_budget": true
}
```

The supplied verification artifact identifies the lineage and failed evidence revision but does not record the required generation or mode-specific fix batch. They are intentionally not inferred.
