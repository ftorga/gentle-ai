# Delta for SDD Orchestrator Assets

## ADDED Requirements

### Requirement: Downstream SDD audit and bounded review

SDD status MUST consume CompactState/CompactReceipt evidence as downstream audit/binding and MUST NOT become an authority. Apply MUST produce one candidate-bound handoff; verification MUST audit it without routine re-derivation. Ordinary bounded 4R MUST run once after apply; R3 MUST be the omission-confirmation lens. Normal guidance MUST NOT invoke Judgment Day, duplicate budget, or alter the 400-line ask-on-risk guard.

#### Scenario: SDD audits compact authority

- GIVEN an activated SDD change
- WHEN status or verification evaluates evidence
- THEN it audits compact-authority binding without replacing it

#### Scenario: One bounded review budget

- GIVEN ordinary bounded 4R ran after apply
- WHEN R3 examines a possible omission
- THEN no additional review budget or normal Judgment Day runs
