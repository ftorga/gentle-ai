# Behavioral Evidence Contract Specification

## Purpose

Define candidate-bound proof for directly observable behavior in CLI `review finalize` and `review validate`. Non-goals: a new phase/lens, matrix, mandatory test per condition, or unconditional model call.

## Requirements

### Requirement: Activated and non-applicable evidence

The system SHALL activate evidence only for directly observable behavior. An activated candidate MUST record material obligations and sufficient proof. Behavior-preserving work MUST record a concise non-applicability basis and MUST NOT require obligations. One proof reference MAY satisfy multiple obligations when sufficient for each.

#### Scenario: Observable behavior activates evidence

- GIVEN a candidate changes a directly observable outcome
- WHEN evidence is prepared
- THEN material obligations and sufficient proof are recorded

#### Scenario: Shared proof is reused

- GIVEN one proof sufficiently demonstrates two obligations
- WHEN evidence is recorded
- THEN both obligations MAY reference that proof

#### Scenario: Non-applicability is concise

- GIVEN a behavior-preserving candidate
- WHEN evidence is prepared
- THEN it records a concise visible basis and no obligations

### Requirement: Candidate binding and receipt compatibility

Evidence MUST bind the current candidate and path identity. A stale or mismatched binding MUST NOT authorize delivery. A corrected candidate MUST replace or rebind its evidence before authorization. Receipts MUST visibly report digest and applicability when supported, preserve legacy v2 readability, and preserve current v3 and sequential-receipt continuity.

#### Scenario: Stale candidate binding is rejected

- GIVEN evidence bound to a different candidate or paths
- WHEN finalize or validate evaluates it
- THEN delivery is denied

#### Scenario: Correction replaces evidence

- GIVEN a correction changes the candidate
- WHEN delivery is retried
- THEN current evidence replaces or rebinds prior evidence

#### Scenario: Receipt exposes applicability

- GIVEN evidence is activated or non-applicable
- WHEN the final receipt is emitted
- THEN digest and applicability are visible

#### Scenario: Legacy and current receipts remain valid

- GIVEN a valid v2 receipt or a valid current v3 receipt
- WHEN it is read or validated
- THEN its applicable compatibility rules remain unchanged

#### Scenario: Sequential receipts preserve continuity

- GIVEN a successor receipt follows a valid prior receipt
- WHEN it is validated
- THEN evidence binding does not break sequential continuity
