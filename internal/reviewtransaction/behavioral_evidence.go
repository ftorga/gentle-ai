package reviewtransaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

const BehavioralEvidenceSchema = "gentle-ai.behavioral-evidence/v1"

type BehavioralEvidenceApplicability string

const (
	BehavioralEvidenceActivated     BehavioralEvidenceApplicability = "activated"
	BehavioralEvidenceNonApplicable BehavioralEvidenceApplicability = "non_applicable"
)

// BehavioralEvidence is compact proof of material outcomes for one candidate.
type BehavioralEvidence struct {
	Schema        string                          `json:"schema"`
	Applicability BehavioralEvidenceApplicability `json:"applicability"`
	Basis         string                          `json:"basis,omitempty"`
	CandidateTree string                          `json:"candidate_tree"`
	PathsDigest   string                          `json:"paths_digest"`
	Obligations   []BehavioralObligation          `json:"obligations"`
}

type BehavioralObligation struct {
	ID                 string   `json:"id"`
	OutcomeOrInvariant string   `json:"outcome_or_invariant"`
	ProofRefs          []string `json:"proof_refs"`
	Disposition        string   `json:"disposition"`
}

func CanonicalBehavioralEvidence(evidence BehavioralEvidence) (BehavioralEvidence, error) {
	evidence.Schema = strings.TrimSpace(evidence.Schema)
	evidence.Basis = strings.TrimSpace(evidence.Basis)
	evidence.CandidateTree = strings.TrimSpace(evidence.CandidateTree)
	evidence.PathsDigest = strings.TrimSpace(evidence.PathsDigest)
	if evidence.Schema != BehavioralEvidenceSchema {
		return BehavioralEvidence{}, errors.New("unsupported behavioral evidence schema")
	}
	if !validGitTree(evidence.CandidateTree) || !validSHA256(evidence.PathsDigest) {
		return BehavioralEvidence{}, errors.New("behavioral evidence candidate binding is invalid")
	}
	switch evidence.Applicability {
	case BehavioralEvidenceNonApplicable:
		if evidence.Basis == "" || len(evidence.Obligations) != 0 {
			return BehavioralEvidence{}, errors.New("non-applicable behavioral evidence requires a basis and no obligations")
		}
		evidence.Obligations = []BehavioralObligation{}
		return evidence, nil
	case BehavioralEvidenceActivated:
		if len(evidence.Obligations) == 0 {
			return BehavioralEvidence{}, errors.New("activated behavioral evidence requires obligations")
		}
	default:
		return BehavioralEvidence{}, errors.New("behavioral evidence applicability is invalid")
	}

	seen := make(map[string]struct{}, len(evidence.Obligations))
	canonical := make([]BehavioralObligation, len(evidence.Obligations))
	for index, obligation := range evidence.Obligations {
		obligation.ID = strings.TrimSpace(obligation.ID)
		obligation.OutcomeOrInvariant = strings.TrimSpace(obligation.OutcomeOrInvariant)
		obligation.Disposition = strings.TrimSpace(obligation.Disposition)
		if obligation.ID == "" || obligation.OutcomeOrInvariant == "" || obligation.Disposition == "" {
			return BehavioralEvidence{}, errors.New("behavioral obligation requires id, outcome, and disposition")
		}
		if _, exists := seen[obligation.ID]; exists {
			return BehavioralEvidence{}, fmt.Errorf("duplicate behavioral obligation %q", obligation.ID)
		}
		seen[obligation.ID] = struct{}{}
		obligation.ProofRefs = canonicalBehavioralProofRefs(obligation.ProofRefs)
		if len(obligation.ProofRefs) == 0 {
			return BehavioralEvidence{}, fmt.Errorf("behavioral obligation %q requires proof references", obligation.ID)
		}
		canonical[index] = obligation
	}
	sort.Slice(canonical, func(i, j int) bool { return canonical[i].ID < canonical[j].ID })
	evidence.Obligations = canonical
	return evidence, nil
}

func BehavioralEvidenceDigest(evidence BehavioralEvidence) string {
	canonical, err := CanonicalBehavioralEvidence(evidence)
	if err != nil {
		return ""
	}
	payload, err := json.Marshal(canonical)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func ParseBehavioralEvidence(payload []byte) (BehavioralEvidence, error) {
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	var evidence BehavioralEvidence
	if err := decoder.Decode(&evidence); err != nil {
		return BehavioralEvidence{}, err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return BehavioralEvidence{}, errors.New("multiple JSON values in behavioral evidence")
	}
	return CanonicalBehavioralEvidence(evidence)
}

func canonicalBehavioralProofRefs(refs []string) []string {
	seen := make(map[string]struct{}, len(refs))
	canonical := make([]string, 0, len(refs))
	for _, ref := range refs {
		ref = strings.TrimSpace(ref)
		if ref == "" {
			continue
		}
		if _, exists := seen[ref]; !exists {
			seen[ref] = struct{}{}
			canonical = append(canonical, ref)
		}
	}
	sort.Strings(canonical)
	return canonical
}
