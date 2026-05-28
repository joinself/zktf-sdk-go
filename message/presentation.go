package message

import (
	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/credential/predicate"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
)

// PresentationRequest is a credential-presentation request (verifier → holder).
type PresentationRequest struct {
	h *ffi.PresentationAction
}

// PresentationRequestBuilder builds a credential-presentation request.
type PresentationRequestBuilder struct {
	h *ffi.PresentationActionBuilder
}

// PresentationResponse is the response to a credential-presentation request.
type PresentationResponse struct {
	h *ffi.PresentationResult
}

// PresentationResponseBuilder builds a credential-presentation response.
type PresentationResponseBuilder struct {
	h *ffi.PresentationResultBuilder
}

func init() {
	ffi.PresentationActionOf = func(o any) *ffi.PresentationAction {
		return o.(*PresentationRequest).h
	}
	ffi.ToPresentationAction = func(h *ffi.PresentationAction) any {
		return &PresentationRequest{h: h}
	}

	ffi.PresentationResultOf = func(o any) *ffi.PresentationResult {
		return o.(*PresentationResponse).h
	}
	ffi.ToPresentationResult = func(h *ffi.PresentationResult) any {
		return &PresentationResponse{h: h}
	}
}

// PresentationTypes returns the requested presentation types.
func (r *PresentationRequest) PresentationTypes() []string { return r.h.PresentationTypes() }

// Holder returns the expected holder address, or an error if not set.
func (r *PresentationRequest) Holder() (*credential.Address, error) {
	h, err := r.h.Holder()
	if err != nil {
		return nil, err
	}

	return ffi.ToDIDAddress(h).(*credential.Address), nil
}

// Challenge returns the random challenge bytes, or nil.
func (r *PresentationRequest) Challenge() []byte { return r.h.Challenge() }

// Predicates returns the predicate tree describing credential constraints, or nil.
func (r *PresentationRequest) Predicates() *predicate.Tree {
	t := r.h.Predicates()
	if t == nil {
		return nil
	}

	return ffi.ToPredicateTree(t).(*predicate.Tree)
}

// Proof returns the verifiable presentations attached as proof.
func (r *PresentationRequest) Proof() []*credential.VerifiablePresentation {
	ps := r.h.Proof()
	out := make([]*credential.VerifiablePresentation, len(ps))

	for i, p := range ps {
		out[i] = ffi.ToVerifiablePresentation(p).(*credential.VerifiablePresentation)
	}

	return out
}

// Term returns the term the requester would like to access credentials under,
// or nil.
func (r *PresentationRequest) Term() *credential.Term {
	t := r.h.Term()
	if t == nil {
		return nil
	}

	return ffi.ToCredentialTerm(t).(*credential.Term)
}

// BiometricAnchor returns the 20-byte biometric anchor hash, or nil.
func (r *PresentationRequest) BiometricAnchor() []byte { return r.h.BiometricAnchor() }

// AsAction wraps this request into a generic Action for inclusion in an
// ExchangeRequest.
func (r *PresentationRequest) AsAction() *Action { return &Action{h: r.h.AsAction()} }

// NewPresentationRequest starts building a presentation request.
func NewPresentationRequest() *PresentationRequestBuilder {
	return &PresentationRequestBuilder{h: ffi.NewPresentationActionBuilder()}
}

// PresentationType sets the requested presentation types.
func (b *PresentationRequestBuilder) PresentationType(types ...string) *PresentationRequestBuilder {
	b.h.PresentationType(ffi.NewPresentationTypes(types))
	return b
}

// Holder sets the expected holder address.
func (b *PresentationRequestBuilder) Holder(holder *credential.Address) *PresentationRequestBuilder {
	b.h.Holder(ffi.DIDAddressOf(holder))
	return b
}

// Challenge sets the random challenge.
func (b *PresentationRequestBuilder) Challenge(challenge []byte) *PresentationRequestBuilder {
	b.h.Challenge(challenge)
	return b
}

// Predicates sets the predicate tree describing credential constraints.
func (b *PresentationRequestBuilder) Predicates(tree *predicate.Tree) *PresentationRequestBuilder {
	b.h.Predicates(ffi.PredicateTreeOf(tree))
	return b
}

// Proof attaches a verifiable presentation as proof.
func (b *PresentationRequestBuilder) Proof(p *credential.VerifiablePresentation) *PresentationRequestBuilder {
	b.h.Proof(ffi.VerifiablePresentationOf(p))
	return b
}

// Term sets the access term.
func (b *PresentationRequestBuilder) Term(term *credential.Term) *PresentationRequestBuilder {
	b.h.Term(ffi.CredentialTermOf(term))
	return b
}

// BiometricAnchor sets the biometric anchor hash.
func (b *PresentationRequestBuilder) BiometricAnchor(anchor []byte) *PresentationRequestBuilder {
	b.h.BiometricAnchor(anchor)
	return b
}

// Finish finalizes the presentation request.
func (b *PresentationRequestBuilder) Finish() (*PresentationRequest, error) {
	a, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &PresentationRequest{h: a}, nil
}

// Presentations returns the verifiable presentations carried in the response.
func (r *PresentationResponse) Presentations() []*credential.VerifiablePresentation {
	ps := r.h.Presentations()
	out := make([]*credential.VerifiablePresentation, len(ps))

	for i, p := range ps {
		out[i] = ffi.ToVerifiablePresentation(p).(*credential.VerifiablePresentation)
	}

	return out
}

// NewPresentationResponse starts building a presentation response.
func NewPresentationResponse() *PresentationResponseBuilder {
	return &PresentationResponseBuilder{h: ffi.NewPresentationResultBuilder()}
}

// Presentation adds a verifiable presentation to the response.
func (b *PresentationResponseBuilder) Presentation(p *credential.VerifiablePresentation) *PresentationResponseBuilder {
	b.h.Presentation(ffi.VerifiablePresentationOf(p))
	return b
}

// Finish finalizes the presentation response.
func (b *PresentationResponseBuilder) Finish() (*PresentationResponse, error) {
	r, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &PresentationResponse{h: r}, nil
}
