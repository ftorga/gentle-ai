# Tasks: Behavioral Completeness Evidence

## Review Workload Forecast

Estimated authored changed lines: 520–680 (tests/fixtures included; generated assets excluded from risk count).
Decision needed before apply: Yes
Chained PRs recommended: Yes
Chain strategy: feature-branch-chain
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Focused test command | Runtime harness | Rollback boundary |
|---|---|---|---|---|
| 1 | Contract/authority | `go test ./internal/reviewtransaction -run 'BehavioralEvidence|Compact'` | temp-dir store + v2/v3 fixtures | `behavioral_evidence.go`, compact files/tests |
| 2 | CLI delivery | `go test ./internal/cli -run 'ReviewFacade|Review'` | finalize/validate JSON scenario | `review_facade.go`, `review.go`, tests |
| 3 | correction/consumers | `go test ./internal/sddstatus ./internal/reviewtransaction -run 'Binding|Correction|Gate'` | N/A: ledger/binding fixtures | correction, binding, shared asset files/tests |

## Phase 1: Contract and Authority (strict TDD)

- [x] 1.1 Decision gate: resolve missing evidence as receipt applicability `missing` versus warning metadata only, and whether `review-integration/v1` may carry warnings without a version bump; record the choice in fixtures/contracts, never infer compatibility.
  - Decision: missing evidence remains warning metadata only; `review-integration/v1` is a closed negotiated contract, so no warning fields are added without a version bump. This slice does not modify that contract.
- [x] 1.2 RED, then minimal GREEN: test and implement canonical schema/order/digest, shared proof refs, activated/non-applicable validation, candidate/path binding, v2 read compatibility, current v3 writes, receipt visibility, stale/replay/sequential rejection, and atomic store replacement in `internal/reviewtransaction/{behavioral_evidence,compact,compact_store,compact_gate}*_test.go` and corresponding symbols.
- [x] 1.3 Focused verification, then `go test ./internal/reviewtransaction`; preserve legacy v2 fixtures and add exact v3 fixtures.

## Phase 2: CLI Finalize/Validate (strict TDD)

- [x] 2.1 RED: cover `--behavioral-evidence <file|->`, canonical input errors, missing warnings, valid activated/non-applicable input, stale binding, unchanged authority after rejection, receipt JSON visibility, and v2 reads in `internal/cli/review_facade_test.go`/`review_test.go`.
- [x] 2.2 GREEN: wire `RunReviewFacadeFinalize`/`RunReviewFacadeValidate` and argument parsing in `internal/cli/review_facade.go`/`review.go`; keep `CompactState`/`CompactReceipt`/`EvaluateCompactGate` authoritative.
- [x] 2.3 Focused CLI verification, then `go test ./internal/cli`.

## Phase 3: Correction, R3, and Downstream Audit (strict TDD)

- [x] 3.1 RED, then GREEN: test atomic correction replacement/rebinding, rejection of retained old evidence, and no partial mutation in compact correction symbols; verify stale/replay behavior.
- [x] 3.2 RED, then GREEN: test R3-confirmed material omission blocks through the existing correction ledger; warnings/non-material uncertainty remain non-blocking follow-ups; test `internal/sddstatus/review_binding.go` audit-only receipt binding and sequential continuity.
- [x] 3.3 Confirmed current specs require no additional shared assets: `review_binding.go` continues to audit authoritative v3 receipts without state mutation or authorization duplication. Ran `go test ./...` at this integration boundary, then `go vet ./...`.

Ordinary bounded 4R runs once after apply; R3 owns omission confirmation. Judgment Day is escalation only for genuinely conflicting severe findings.
