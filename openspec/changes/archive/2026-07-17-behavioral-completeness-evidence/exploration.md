## Exploration: behavioral-completeness-evidence

### Current State
The codebase already has the core machinery this change needs. `internal/reviewtransaction/behavioral_evidence.go` defines a canonical, candidate-bound contract with `activated` and `non_applicable` modes, proof refs, canonical ordering, and a digest. `internal/reviewtransaction/compact.go` already stores that contract on `CompactState`, validates it against the live candidate snapshot, and emits v3 receipts with behavioral evidence digest/applicability. `internal/cli/review_facade.go` is the real delivery seam: it owns `review finalize` / `review validate` and already dispatches through compact authority and native gate evaluation. `internal/sddstatus/review_binding.go` is downstream-only OpenSpec binding, not the authority boundary.

So the intended change is not “invent behavioral evidence”; it is to finish wiring the existing contract into the right delivery/review paths and make the SDD artifacts reflect that contract consistently.

### Affected Areas
- `internal/reviewtransaction/behavioral_evidence.go` — canonical contract and validation rules.
- `internal/reviewtransaction/compact.go` — state/receipt binding and terminal validation.
- `internal/reviewtransaction/compact_store_test.go` — current evidence-binding regression coverage.
- `internal/cli/review_facade.go` / `internal/cli/review.go` — finalize/validate surface and review dispatch.
- `internal/sddstatus/review_gate.go` / `internal/sddstatus/status.go` — bounded review routing, not authority ownership.
- `openspec/changes/behavioral-completeness-evidence/*` — proposal/spec/design/tasks/apply-progress artifacts for the current change.

### Approaches
1. **Compact-authority-first integration** — keep the contract owned by `CompactState`/`CompactReceipt`, then have finalize/validate consume it and let SDD status/binding only observe it.
   - Pros: Matches current architecture; smallest seam; preserves authority/casuality boundaries; low regression risk.
   - Cons: Requires careful cross-file wiring and compatibility validation for v2/v3 receipts.
   - Effort: Medium

2. **SDD-status-led integration** — surface behavioral completeness primarily through `internal/sddstatus` and let review transaction follow later.
   - Pros: Feels close to SDD workflow.
   - Cons: Wrong authority boundary; duplicates lifecycle logic; risks drifting from compact-v2/facade truth.
   - Effort: High

3. **Facade-only shim** — expose evidence only at `review finalize`/`review validate` and avoid touching compact state/receipt.
   - Pros: Smaller diff.
   - Cons: Breaks candidate binding and sequential receipts; leaves SDD binding unable to audit freshness.
   - Effort: Low initially, but incomplete.

### Recommendation
Use **Compact-authority-first integration**. The smallest viable seam is the existing compact review authority: bind evidence to `CompactState`, carry digest/applicability in `CompactReceipt`, and have `review finalize` / `review validate` enforce freshness and terminal binding. That gives the candidate-bound contract without inventing a new review phase or a new authority layer.

### Risks
- v2/v3 receipt compatibility must stay explicit; missing digest/applicability must continue to fail closed for v3 while legacy reads remain valid.
- If the change expands into SDD-status ownership, it will duplicate authority and grow beyond the 400-line review budget.
- The live gate/review-bind path must preserve sequential receipts and cannot treat `sddstatus.ReviewBinding` as authority.

### Ready for Proposal
Yes — but the proposal should stay narrow: define the exact compact-authority seam, confirm whether any new SDD artifact text is required, and keep review-budget planning centered on the compact finalize/validate path.
