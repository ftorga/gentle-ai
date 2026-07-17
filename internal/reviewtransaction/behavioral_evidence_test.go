package reviewtransaction

import (
	"reflect"
	"strings"
	"testing"
)

func TestCanonicalBehavioralEvidence(t *testing.T) {
	activated := BehavioralEvidence{
		Schema: BehavioralEvidenceSchema, Applicability: BehavioralEvidenceActivated,
		CandidateTree: tree("a"), PathsDigest: hash("b"),
		Obligations: []BehavioralObligation{
			{ID: "outcome-b", OutcomeOrInvariant: "rejected input preserves stored value", ProofRefs: []string{"go test ./internal/reviewtransaction -run Rejection"}, Disposition: "proved"},
			{ID: "outcome-a", OutcomeOrInvariant: "accepted input persists value", ProofRefs: []string{"go test ./internal/reviewtransaction -run Acceptance", "go test ./internal/reviewtransaction -run Rejection"}, Disposition: "proved"},
		},
	}

	tests := []struct {
		name     string
		evidence BehavioralEvidence
		wantErr  string
	}{
		{name: "activated canonicalizes obligations and shared proof references", evidence: activated},
		{name: "non applicable requires concise basis", evidence: BehavioralEvidence{Schema: BehavioralEvidenceSchema, Applicability: BehavioralEvidenceNonApplicable, Basis: "documentation only", CandidateTree: tree("a"), PathsDigest: hash("b"), Obligations: []BehavioralObligation{}}},
		{name: "activated requires proof references", evidence: BehavioralEvidence{Schema: BehavioralEvidenceSchema, Applicability: BehavioralEvidenceActivated, CandidateTree: tree("a"), PathsDigest: hash("b"), Obligations: []BehavioralObligation{{ID: "missing-proof", OutcomeOrInvariant: "outcome", Disposition: "proved"}}}, wantErr: "proof"},
		{name: "duplicate obligations are rejected", evidence: BehavioralEvidence{Schema: BehavioralEvidenceSchema, Applicability: BehavioralEvidenceActivated, CandidateTree: tree("a"), PathsDigest: hash("b"), Obligations: []BehavioralObligation{{ID: "duplicate", OutcomeOrInvariant: "one", ProofRefs: []string{"test"}, Disposition: "proved"}, {ID: "duplicate", OutcomeOrInvariant: "two", ProofRefs: []string{"test"}, Disposition: "proved"}}}, wantErr: "duplicate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CanonicalBehavioralEvidence(tt.evidence)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("CanonicalBehavioralEvidence() error = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("CanonicalBehavioralEvidence() error = %v", err)
			}
			if len(got.Obligations) > 1 && got.Obligations[0].ID != "outcome-a" {
				t.Fatalf("obligation order = %#v", got.Obligations)
			}
			if len(got.Obligations) == 2 && !reflect.DeepEqual(got.Obligations[0].ProofRefs, []string{"go test ./internal/reviewtransaction -run Acceptance", "go test ./internal/reviewtransaction -run Rejection"}) {
				t.Fatalf("proof references = %#v", got.Obligations[0].ProofRefs)
			}
			if digest := BehavioralEvidenceDigest(got); !validSHA256(digest) || digest != BehavioralEvidenceDigest(got) {
				t.Fatalf("digest = %q", digest)
			}
		})
	}
}
