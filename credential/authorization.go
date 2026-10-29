package credential

import (
	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

type ReasonKind int

const (
	ReasonRejected       ReasonKind = ReasonKind(ffi.AuthorizationReasonRejected)
	ReasonQuorumNotMet   ReasonKind = ReasonKind(ffi.AuthorizationReasonQuorumNotMet)
	ReasonPersonUnproven ReasonKind = ReasonKind(ffi.AuthorizationReasonPersonUnproven)
)

type Gap int

const (
	GapNone            Gap = Gap(ffi.AuthorizationGapNone)
	GapAnchorMissing   Gap = Gap(ffi.AuthorizationGapAnchorMissing)
	GapAnchorMismatch  Gap = Gap(ffi.AuthorizationGapAnchorMismatch)
	GapLivenessMissing Gap = Gap(ffi.AuthorizationGapLivenessMissing)
	GapEvidenceRevoked Gap = Gap(ffi.AuthorizationGapEvidenceRevoked)
)

type RequirementKind int

const (
	RequirementSignatures RequirementKind = RequirementKind(ffi.AuthorizationRequirementSignatures)
	RequirementAnchor     RequirementKind = RequirementKind(ffi.AuthorizationRequirementAnchor)
	RequirementLiveness   RequirementKind = RequirementKind(ffi.AuthorizationRequirementLiveness)
)

type Reason struct {
	Kind      ReasonKind
	Err       error
	Document  *identity.Address
	Role      identity.KeyRole
	Weight    uint64
	Threshold uint64
	Signer    *signing.PublicKey
	Identity  *identity.Address
	Gap       Gap
}

type Requirement struct {
	Kind       RequirementKind
	Document   *identity.Address
	Role       identity.KeyRole
	Weight     uint64
	Candidates []*signing.PublicKey
	Identity   *identity.Address
	Commitment []byte
	Hash       []byte
}

type Verdict struct {
	Authorized bool
	Reasons    []Reason
}

type OperationVerdict struct {
	Verdict
	Document *identity.Document
	Missing  []Requirement
}

type Person struct {
	Identity   *identity.Address
	Commitment []byte
	Anchor     *VerifiableCredential
	Evidence   *VerifiableCredential
}

type Link struct {
	From       *identity.Address
	To         *identity.Address
	Key        *signing.PublicKey
	Commitment []byte
	Roles      identity.KeyRole
}

type AuthorizationChain struct {
	h *ffi.AuthorizationChain
}

func (g *Graph) ValidOperationFor(document *identity.Address, operation *identity.Operation) (*OperationVerdict, error) {
	var address *ffi.DIDAddress
	if document != nil {
		address = ffi.DIDAddressOf(document)
	}

	v, err := g.h.ValidOperationFor(address, ffi.IdentityOperationOf(operation))
	if err != nil {
		return nil, err
	}

	out := &OperationVerdict{
		Verdict: verdictFrom(v.Verdict),
		Missing: requirementsFrom(v.Missing),
	}
	if v.Document != nil {
		out.Document = ffi.ToIdentityDocument(v.Document).(*identity.Document)
	}

	return out, nil
}

func (g *Graph) AuthorizationFor(hash []byte) (*AuthorizationChain, error) {
	c, err := g.h.AuthorizationFor(hash)
	if err != nil {
		return nil, err
	}

	return &AuthorizationChain{h: c}, nil
}

func (c *AuthorizationChain) People() []Person {
	people := c.h.People()
	out := make([]Person, len(people))
	for i, p := range people {
		out[i] = Person{
			Identity:   address(p.Identity),
			Commitment: p.Commitment,
			Anchor:     &VerifiableCredential{h: p.Anchor},
			Evidence:   &VerifiableCredential{h: p.Evidence},
		}
	}

	return out
}

func (c *AuthorizationChain) Links() []Link {
	links := c.h.Links()
	out := make([]Link, len(links))
	for i, l := range links {
		out[i] = Link{
			From:       address(l.From),
			To:         address(l.To),
			Key:        publicKey(l.Key),
			Commitment: l.Commitment,
			Roles:      identity.KeyRole(l.Roles),
		}
	}

	return out
}

func (c *AuthorizationChain) AuthorizedBy(party *identity.Address) bool {
	return c.h.AuthorizedBy(ffi.DIDAddressOf(party))
}

func (c *AuthorizationChain) Verdict(document *identity.Document, role identity.KeyRole) (Verdict, error) {
	v, err := c.h.Verdict(ffi.IdentityDocumentOf(document), uint64(role))
	if err != nil {
		return Verdict{}, err
	}

	return verdictFrom(v), nil
}

func (c *AuthorizationChain) Missing(document *identity.Document, role identity.KeyRole) ([]Requirement, error) {
	r, err := c.h.Missing(ffi.IdentityDocumentOf(document), uint64(role))
	if err != nil {
		return nil, err
	}

	return requirementsFrom(r), nil
}

func verdictFrom(v ffi.AuthorizationVerdict) Verdict {
	reasons := make([]Reason, len(v.Reasons))
	for i, r := range v.Reasons {
		reasons[i] = Reason{
			Kind:      ReasonKind(r.Kind),
			Err:       r.Error,
			Document:  address(r.Document),
			Role:      identity.KeyRole(r.Role),
			Weight:    r.Weight,
			Threshold: r.Threshold,
			Signer:    publicKey(r.Signer),
			Identity:  address(r.Identity),
			Gap:       Gap(r.Gap),
		}
	}

	return Verdict{Authorized: v.Authorized, Reasons: reasons}
}

func requirementsFrom(requirements []ffi.AuthorizationRequirement) []Requirement {
	out := make([]Requirement, len(requirements))
	for i, r := range requirements {
		candidates := make([]*signing.PublicKey, len(r.Candidates))
		for j, k := range r.Candidates {
			candidates[j] = publicKey(k)
		}
		out[i] = Requirement{
			Kind:       RequirementKind(r.Kind),
			Document:   address(r.Document),
			Role:       identity.KeyRole(r.Role),
			Weight:     r.Weight,
			Candidates: candidates,
			Identity:   address(r.Identity),
			Commitment: r.Commitment,
			Hash:       r.Hash,
		}
	}

	return out
}

func address(h *ffi.DIDAddress) *identity.Address {
	if h == nil {
		return nil
	}

	return ffi.ToDIDAddress(h).(*identity.Address)
}

func publicKey(h *ffi.SigningPublicKey) *signing.PublicKey {
	if h == nil {
		return nil
	}

	return ffi.ToSigningPublicKey(h).(*signing.PublicKey)
}
