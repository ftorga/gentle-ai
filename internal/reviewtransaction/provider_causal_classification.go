package reviewtransaction

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type ProviderCausalClassification string

const (
	ProviderCandidateCausal    ProviderCausalClassification = "candidate-causal"
	ProviderProvenNonCandidate ProviderCausalClassification = "proven-non-candidate"
	ProviderUnknown            ProviderCausalClassification = "unknown"
)

var providerDiffHunk = regexp.MustCompile(`(?m)^@@ -\d+(?:,\d+)? \+(\d+)(?:,(\d+))? @@`)

type ProviderCausalEvidence struct {
	FindingID string   `json:"finding_id"`
	Location  string   `json:"location"`
	ProofRefs []string `json:"proof_refs"`
}
type ProviderCausalFinding struct {
	FindingID      string                       `json:"finding_id"`
	Location       string                       `json:"location"`
	ProofRefs      []string                     `json:"proof_refs"`
	Classification ProviderCausalClassification `json:"classification"`
	EvidenceDigest string                       `json:"evidence_digest"`
}
type ProviderCausalCarrier struct {
	SubjectHash       string                  `json:"subject_hash"`
	CandidateIdentity CandidateIdentity       `json:"candidate_identity"`
	Findings          []ProviderCausalFinding `json:"findings"`
	AggregateDigest   string                  `json:"aggregate_digest"`

	ArtifactBinding NewLineageArtifactBinding `json:"artifact_binding"`
}

var (
	ErrProviderCausalCarrierMissing  = errors.New("new-lineage provider causal carrier is missing")                     // refusal:by-design operator-knowledge: missing persisted authority requires fresh capture
	ErrProviderCausalCarrierConflict = errors.New("new-lineage provider causal classifications conflict across lenses") // refusal:by-design operator-knowledge: conflicting authority requires fresh capture
)

type ProviderCausalFailure struct {
	Kind, Operation string
	Cause           error
}

func (e *ProviderCausalFailure) Error() string {
	return fmt.Sprintf("provider causal %s failure during %s: %v", e.Kind, e.Operation, e.Cause)
}
func (e *ProviderCausalFailure) Unwrap() error { return e.Cause }
func (c ProviderCausalCarrier) Validate() error {
	if !validSHA256(c.SubjectHash) {
		return errors.New("provider causal carrier requires a canonical subject hash") // refusal:by-design world-action: immutable carrier bytes with an invalid subject cannot be repaired safely by an operator command
	}
	for i, f := range c.Findings {
		if f.FindingID == "" || (i > 0 && f.FindingID <= c.Findings[i-1].FindingID) {
			return errors.New("provider causal carrier findings must be unique and canonical") // refusal:by-design world-action: immutable carrier findings cannot be canonicalized safely by an operator command
		}
		if f.Location != strings.TrimSpace(f.Location) {
			return errors.New("provider causal carrier finding location must be canonical") // refusal:by-design world-action: immutable carrier locations cannot be canonicalized safely by an operator command
		}
		if f.Classification != ProviderCandidateCausal && f.Classification != ProviderProvenNonCandidate && f.Classification != ProviderUnknown {
			return errors.New("provider causal carrier classification is unsupported") // refusal:by-design world-action: an unsupported provider classification requires corrected source evidence, not an operator command
		}
		if f.Classification == ProviderCandidateCausal && len(f.ProofRefs) == 0 {
			return errors.New("provider causal carrier candidate-causal findings require proof refs") // refusal:by-design world-action: missing immutable proof cannot be invented by an operator command
		}
		if f.EvidenceDigest != providerFindingDigest(f) {
			return errors.New("provider causal carrier finding digest does not match its content") // refusal:by-design world-action: an immutable digest mismatch cannot be repaired safely by an operator command
		}
	}
	if c.AggregateDigest != providerAggregateDigest(c) {
		return errors.New("provider causal carrier aggregate digest does not match its findings") // refusal:by-design world-action: an immutable aggregate mismatch cannot be repaired safely by an operator command
	}
	if c.ArtifactBinding.Subject.SubjectHash != "" && c.ArtifactBinding.Subject.SubjectHash != c.SubjectHash {
		return errors.New("provider causal carrier artifact binding subject hash does not match carrier") // refusal:by-design world-action: an immutable binding mismatch requires fresh capture
	}
	return nil
}
func canonicalProviderProofRefs(refs []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, r := range refs {
		r = strings.TrimSpace(r)
		if r != "" && !seen[r] {
			seen[r] = true
			out = append(out, r)
		}
	}
	sort.Strings(out)
	return out
}

func (authority NewLineageAuthority) ProviderCausalCarrierDigest() (string, error) {
	type entry struct{ Lens, SubjectHash, AggregateDigest string }
	entries := make([]entry, 0, len(authority.CapturedResults))
	for _, captured := range authority.CapturedResults {
		if captured.Provider.SubjectHash == "" {
			continue
		}
		if err := captured.Provider.Validate(); err != nil {
			return "", err
		}
		if captured.Provider.SubjectHash != captured.SubjectHash || captured.Provider.CandidateIdentity != authority.CandidateIdentity {
			// refusal:by-design world-action: mismatched or corrupt persisted authority must be restored and re-reviewed; no caller command can reinterpret it
			return "", errors.New("provider causal authority aggregate has an invalid carrier binding")
		}
		entries = append(entries, entry{captured.Lens, captured.SubjectHash, captured.Provider.AggregateDigest})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Lens < entries[j].Lens })
	payload, _ := json.Marshal(struct {
		Candidate CandidateIdentity `json:"candidate_identity"`
		Carriers  []entry           `json:"carriers"`
	}{authority.CandidateIdentity, entries})
	sum := sha256.Sum256(append([]byte("gentle-ai.provider-causal-authority/v1\x00"), payload...))
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}
func (authority NewLineageAuthority) ProviderCausalReceiptDigest(revision string) (string, error) {
	if !validSHA256(revision) {
		// refusal:by-design world-action: an invalid authority revision must be restored and re-reviewed; no caller command can reinterpret it
		return "", errors.New("provider causal aggregate requires a valid authority revision")
	}
	carrier, err := authority.ProviderCausalCarrierDigest()
	if err != nil {
		return "", err
	}
	payload, _ := json.Marshal(struct{ AuthorityRevision, CarrierDigest string }{revision, carrier})
	sum := sha256.Sum256(append([]byte("gentle-ai.provider-causal-aggregate/v1\x00"), payload...))
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}
func canonicalProviderClaims(in []ProviderCausalEvidence) ([]ProviderCausalEvidence, error) {
	out := append([]ProviderCausalEvidence(nil), in...)
	for i := range out {
		out[i].FindingID = strings.TrimSpace(out[i].FindingID)
		out[i].Location = strings.TrimSpace(out[i].Location)
		out[i].ProofRefs = canonicalProviderProofRefs(out[i].ProofRefs)
		if out[i].FindingID == "" {
			return nil, errors.New("provider causal evidence requires a finding id") // refusal:by-design world-action: incomplete provider evidence requires a new source result, not an operator command
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].FindingID < out[j].FindingID })
	for i := 1; i < len(out); i++ {
		if out[i].FindingID == out[i-1].FindingID && (out[i].Location != out[i-1].Location || !equalStrings(out[i].ProofRefs, out[i-1].ProofRefs)) {
			return nil, errors.New("provider causal evidence contains conflicting duplicate finding id") // refusal:by-design world-action: conflicting provider evidence requires a new source result, not an operator command
		}
	}
	dedup := out[:0]
	for _, v := range out {
		if len(dedup) == 0 || dedup[len(dedup)-1].FindingID != v.FindingID {
			dedup = append(dedup, v)
		}
	}
	return dedup, nil
}
func providerFindingDigest(f ProviderCausalFinding) string {
	b, _ := json.Marshal([]any{f.FindingID, f.Location, f.ProofRefs, f.Classification})
	return fmt.Sprintf("sha256:%x", sha256.Sum256(append([]byte("gentle-ai.provider-causal-finding/v1\x00"), b...)))
}
func providerAggregateDigest(c ProviderCausalCarrier) string {
	ids := make([]string, len(c.Findings))
	for i, f := range c.Findings {
		ids[i] = f.EvidenceDigest
	}
	b, _ := json.Marshal([]any{c.SubjectHash, c.CandidateIdentity, c.ArtifactBinding, ids})
	return fmt.Sprintf("sha256:%x", sha256.Sum256(append([]byte("gentle-ai.provider-causal-aggregate/v1\x00"), b...)))
}
func ProviderCausalAggregateDigest(a NewLineageAuthority) string {
	parts := make([]string, 0, len(a.CapturedResults))
	for _, captured := range a.CapturedResults {
		if captured.Provider.SubjectHash != "" {
			parts = append(parts, captured.Provider.AggregateDigest)
		}
	}
	sort.Strings(parts)
	b, _ := json.Marshal([]any{a.CandidateIdentity, parts})
	sum := sha256.Sum256(append([]byte("gentle-ai.provider-causal-authority/v1\x00"), b...))
	return fmt.Sprintf("sha256:%x", sum)
}

func DeriveProviderCausalCarrier(ctx context.Context, repo, subject string, candidate CandidateIdentity, claims []ProviderCausalEvidence) (ProviderCausalCarrier, error) {
	if !validSHA256(subject) || !validGitTree(candidate.BaseTree) || !validGitTree(candidate.CandidateTree) {
		return ProviderCausalCarrier{}, errors.New("provider causal derivation requires a valid subject and frozen candidate trees") // refusal:by-design world-action: caller-supplied frozen identity must be valid before derivation and cannot be repaired by an operator command
	}
	claims, err := canonicalProviderClaims(claims)
	if err != nil {
		return ProviderCausalCarrier{}, &ProviderCausalFailure{"capture", "canonicalize provider claims", err}
	}
	c := ProviderCausalCarrier{SubjectHash: subject, CandidateIdentity: candidate, Findings: make([]ProviderCausalFinding, 0, len(claims))}
	for _, claim := range claims {
		class := ProviderUnknown
		if claim.Location != "" {
			proof, err := providerProofRefsValid(ctx, repo, candidate, claim.ProofRefs)
			if err != nil {
				return ProviderCausalCarrier{}, err
			}
			changed, same, err := providerCandidateLineChanged(ctx, repo, candidate, claim.Location)
			if err != nil {
				return ProviderCausalCarrier{}, err
			}
			if changed && proof {
				class = ProviderCandidateCausal
			} else if same {
				class = ProviderProvenNonCandidate
			}
		}
		f := ProviderCausalFinding{FindingID: claim.FindingID, Location: claim.Location, ProofRefs: claim.ProofRefs, Classification: class}
		f.EvidenceDigest = providerFindingDigest(f)
		c.Findings = append(c.Findings, f)
	}
	c.AggregateDigest = providerAggregateDigest(c)
	return c, nil
}
func providerProofRefsValid(ctx context.Context, repo string, c CandidateIdentity, refs []string) (bool, error) {
	if len(refs) == 0 {
		return false, nil
	}
	for _, r := range refs {
		p, e := parseFindingLocation(r)
		if e != nil || p.StartLine < 1 {
			return false, nil
		}
		b, x := runGit(ctx, repo, nil, nil, "show", c.CandidateTree+":"+p.Path)
		if x != nil {
			return false, &ProviderCausalFailure{"infrastructure", "read frozen proof ref", x}
		}
		lineCount := len(strings.Split(strings.TrimSuffix(string(b), "\n"), "\n"))
		if len(b) == 0 {
			lineCount = 0
		}
		if p.StartLine > lineCount || p.EndLine > lineCount {
			return false, nil
		}
	}
	return true, nil
}
func providerCandidateLineChanged(ctx context.Context, repo string, c CandidateIdentity, loc string) (bool, bool, error) {
	p, e := parseFindingLocation(loc)
	if e != nil {
		return false, false, e
	}
	b, e := runGit(ctx, repo, nil, nil, "diff", "--unified=0", "--no-renames", "--no-ext-diff", "--no-textconv", c.BaseTree, c.CandidateTree, "--", literalPathspec(p.Path))
	if e != nil {
		return false, false, &ProviderCausalFailure{"infrastructure", "read frozen candidate line", e}
	}
	for _, m := range providerDiffHunk.FindAllSubmatch(b, -1) {
		s, _ := strconv.Atoi(string(m[1]))
		n := 1
		if len(m[2]) > 0 {
			n, _ = strconv.Atoi(string(m[2]))
		}
		if n > 0 && p.StartLine < s+n && s <= p.EndLine {
			return true, false, nil
		}
	}
	a, e := runGit(ctx, repo, nil, nil, "show", c.BaseTree+":"+p.Path)
	if e != nil {
		return false, false, &ProviderCausalFailure{"infrastructure", "read frozen base path", e}
	}
	candidateBytes, e := runGit(ctx, repo, nil, nil, "show", c.CandidateTree+":"+p.Path)
	if e != nil {
		return false, false, &ProviderCausalFailure{"infrastructure", "read frozen candidate path", e}
	}
	return false, p.StartLine > 0 && string(a) == string(candidateBytes), nil
}
