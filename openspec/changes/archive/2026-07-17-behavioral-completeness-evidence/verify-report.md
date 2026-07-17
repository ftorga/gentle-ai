```yaml
schema: gentle-ai.verify-result/v1
evidence_revision: sha256:775e4f15a5743fbb36da44921eea3c15f5a0723edf09fef823c6e4d72290fc4a
verdict: pass_with_warnings
blockers: 0
critical_findings: 0
requirements: 4/4
scenarios: 13/13
test_command: go test ./...
test_exit_code: 1
test_output_hash: sha256:4170190cf89f4221a9fd28ec94cb2e4edbec591676b87ab78e3f9995556953a3
build_command: go vet ./...
build_exit_code: 0
build_output_hash: sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

# Verification Report: Behavioral Completeness Evidence

**Mode**: Strict TDD  
**Authority consumed**: approved native post-apply lineage `review-7f7c4fe5640e2f5e`, generation `1`; `gentle-ai review validate --gate post-apply` previously returned `allow`.

## Completeness

| Metric | Value |
|---|---:|
| Tasks total | 9 |
| Tasks complete | 9 |
| Tasks incomplete | 0 |
| Requirements | 4/4 |
| Scenarios | 13/13 |

## Runtime Evidence

| Scope | Command | Exit | Output SHA-256 | Result |
|---|---|---:|---|---|
| CLI correction regression | `go test ./internal/cli -run '^TestReviewFacadeCorrectionReplacesBehavioralEvidenceAtomically$' -count=1` | 0 | `sha256:f18f79b7576bc1b62f409a19eb94c256929f1103deea5f75c5bc0879371bda8f` | PASS |
| Focused CLI | `go test ./internal/cli -run 'ReviewFacade|Review' -count=1` | 0 | `sha256:decc257bce5de125e31ebe562900d58899250ea31990a1a9c7b1ec2e3a454c52` | PASS |
| Focused compact authority and SDD | `go test ./internal/reviewtransaction ./internal/sddstatus -run 'Binding|Correction|Gate|BehavioralEvidence' -count=1` | 0 | `sha256:48efa99cd4a9bd9bf852cb29f0150fd998e444cb169960eaa618b26420fa6440` | PASS |
| Full suite | `go test ./...` | 1 | `sha256:4170190cf89f4221a9fd28ec94cb2e4edbec591676b87ab78e3f9995556953a3` | Known baseline warning |
| Static analysis | `go vet ./...` | 0 | `sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` | PASS |

The focused CLI regression passed. Its temp-directory Git harness proves a corrected v3 receipt carries replacement evidence and invalid replacement input preserves the correction state.

### Known Baseline Warning

The full suite reproduced four of the five approved pre-existing environment-sensitive CodeGraph failures; the fifth did not reproduce in this run. They are outside this candidate:

- `internal/cli`: two macOS `/var` versus `/private/var` canonicalization/MCP-init failures.
- `internal/components/communitytool`: path canonicalization and an MCP deadline phase assertion failure.

This is a known baseline warning, not a candidate-caused blocker. All candidate-focused packages and the correction regression passed.

## Spec Compliance Matrix

| Requirement | Scenario | Runtime coverage | Result |
|---|---|---|---|
| Activated and non-applicable evidence | Observable behavior activates evidence | `behavioral_evidence_test.go` in focused compact suite | ✅ COMPLIANT |
| Activated and non-applicable evidence | Shared proof is reused | `behavioral_evidence_test.go` in focused compact suite | ✅ COMPLIANT |
| Activated and non-applicable evidence | Non-applicability is concise | behavioral evidence and CLI facade tests | ✅ COMPLIANT |
| Candidate binding and receipt compatibility | Stale candidate binding is rejected | compact/store and CLI facade focused tests | ✅ COMPLIANT |
| Candidate binding and receipt compatibility | Correction replaces evidence | compact replacement test and CLI correction regression | ✅ COMPLIANT |
| Candidate binding and receipt compatibility | Receipt exposes applicability | CLI regression verifies the v3 receipt | ✅ COMPLIANT |
| Candidate binding and receipt compatibility | Legacy and current receipts remain valid | compact/store/gate focused tests | ✅ COMPLIANT |
| Candidate binding and receipt compatibility | Sequential receipts preserve continuity | SDD binding focused tests | ✅ COMPLIANT |
| Finalize and validate completeness enforcement | Missing evidence warns | CLI facade focused tests | ✅ COMPLIANT |
| Finalize and validate completeness enforcement | R3 confirms a material omission | `TestCompactReviewR3MaterialOmissionUsesExistingCorrectionLedger` | ✅ COMPLIANT |
| Finalize and validate completeness enforcement | Existing lifecycle remains authoritative | compact gate and CLI facade focused tests | ✅ COMPLIANT |
| Downstream SDD audit and bounded review | SDD audits compact authority | SDD binding focused tests | ✅ COMPLIANT |
| Downstream SDD audit and bounded review | One bounded review budget | bounded-review/compact focused tests and approved authority | ✅ COMPLIANT |

**Compliance summary**: 13/13 scenarios compliant.

## Correctness and Design Coherence

| Check | Result | Evidence |
|---|---|---|
| Compact authority remains canonical | ✅ | Evidence is validated/stored in `CompactState`; receipts derive digest and applicability from it. |
| Correction replacement is atomic | ✅ | The facade reserves supplied behavioral evidence for the correction branch and calls `CompleteCorrectionWithBehavioralEvidence`. |
| Initial binding behavior is unchanged | ✅ | Initial binding occurs only outside `StateCorrectionRequired`. |
| Rejected correction input is immutable | ✅ | CLI and compact tests assert unchanged state/revision after rejection. |
| SDD remains downstream audit only | ✅ | Focused SDD binding coverage passes; no authorization path was added. |
| v2/v3 and sequential continuity | ✅ | Focused compact/store/gate and SDD binding coverage passes. |

## TDD Compliance

| Check | Result | Details |
|---|---|---|
| TDD evidence reported | ✅ | Current apply-progress contains task-by-task evidence, including remediation. |
| All tasks have tests | ✅ | 5/5 TDD rows name existing test coverage. |
| RED confirmed | ✅ | RED descriptions identify former failing/absent behavior; named tests exist. |
| GREEN confirmed | ✅ | 5/5 rows have current passing focused runtime evidence. |
| Triangulation adequate | ✅ | Activated/non-applicable, stale/replay, replacement/rejection, R3 severity, and sequential cases vary behavior. |
| Safety net | ✅ | Every row records a focused safety-net command. |

**TDD compliance**: 6/6 checks passed.

## Test Layers, Coverage, and Assertion Quality

| Layer | Tests | Files | Result |
|---|---:|---:|---|
| Unit / temp-directory repository | focused compact, SDD, and CLI suites | 5+ | PASS |
| Integration boundary | full Go suite | 57 packages | Known baseline warning only |
| E2E | 0 | 0 | Not applicable |

Focused coverage is informational and reflects narrow behavioral selections:

| Scope | Coverage | Output SHA-256 |
|---|---:|---|
| CLI correction regression | 7.2% | `sha256:5440f48138e7acc94fabefaf504dd198799476b0df39f5d5cbdddadda5f71d3e` |
| Compact authority focused suite | 56.8% | `sha256:3b0047f320bce56cd3256740266ffa26b137a3762f44fbb862c4449e2f0f3262` |
| SDD status focused suite | 51.4% | `sha256:3b0047f320bce56cd3256740266ffa26b137a3762f44fbb862c4449e2f0f3262` |

**Assertion quality**: ✅ Inspected correction, compact, and SDD tests invoke production boundaries and assert state, receipt, digest, or immutability outcomes. No tautology, ghost-loop, type-only-only, or smoke-only assertion was found.

## Issues Found

**CRITICAL**: None.

**WARNING**:
- `go test ./...` exits 1 because of the approved, unrelated CodeGraph baseline.
- Focused coverage is below a broad-suite threshold because it intentionally executes behavior-specific selections; informational only.

**SUGGESTION**: Stabilize the approved CodeGraph path canonicalization and MCP timing baseline independently.

## Verdict

**PASS WITH WARNINGS** — all 4 requirements and 13 scenarios have passing focused runtime coverage, the corrected CLI path calls `CompleteCorrectionWithBehavioralEvidence`, the CLI regression passes, and static analysis passes. The only full-suite failure is the approved unrelated baseline.
