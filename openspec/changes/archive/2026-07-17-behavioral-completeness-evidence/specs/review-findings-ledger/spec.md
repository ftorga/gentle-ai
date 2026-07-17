# Delta for Review Findings Ledger

## ADDED Requirements

### Requirement: Finalize and validate completeness enforcement

`review finalize` and `review validate` MUST enforce the compact authority's evidence binding without replacing existing authority, causality, low-risk, integration, or sequential-receipt contracts. Missing activated evidence MUST warn maintainers. An R3-confirmed material omission MUST block delivery until corrected; speculative or non-material concerns MUST NOT block solely as evidence warnings.

#### Scenario: Missing evidence warns

- GIVEN activated evidence is missing
- WHEN finalize or validate runs
- THEN maintainers receive a visible warning

#### Scenario: R3 confirms a material omission

- GIVEN R3 confirms an unproved material obligation
- WHEN delivery is evaluated
- THEN delivery is blocked until corrected

#### Scenario: Existing lifecycle remains authoritative

- GIVEN evidence is evaluated during finalize or validate
- WHEN an upstream authority or causality rule applies
- THEN its existing outcome remains controlling
