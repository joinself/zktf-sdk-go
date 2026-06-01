// Package pairwise provides pairwise identity / relationship / introduction
// read types. (Pairwise lookup, store, and validation are mobile-only in the
// native SDK and are not exposed here.)
package pairwise

import (
	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
)

// Status is the connection state of a pairwise relationship.
type Status uint32

const (
	StatusPending     Status = Status(ffi.PairwiseStatusPending)
	StatusNegotiating Status = Status(ffi.PairwiseStatusNegotiating)
	StatusEstablished Status = Status(ffi.PairwiseStatusEstablished)
)

// Identity is a pairwise identity (the counterparty's document address plus
// optional biometric anchor hash).
type Identity struct {
	h *ffi.PairwiseIdentity
}

// Introduction is the validated result of an introduction message.
type Introduction struct {
	h *ffi.PairwiseIntroduction
}

func init() {
	ffi.PairwiseIdentityOf = func(o any) *ffi.PairwiseIdentity { return o.(*Identity).h }
	ffi.ToPairwiseIdentity = func(h *ffi.PairwiseIdentity) any { return &Identity{h: h} }

	ffi.PairwiseIntroductionOf = func(o any) *ffi.PairwiseIntroduction { return o.(*Introduction).h }
	ffi.ToPairwiseIntroduction = func(h *ffi.PairwiseIntroduction) any { return &Introduction{h: h} }
}

// DecodeIdentity decodes an encoded pairwise identity.
func DecodeIdentity(data []byte) (*Identity, error) {
	i, err := ffi.PairwiseIdentityDecode(data)
	if err != nil {
		return nil, err
	}

	return &Identity{h: i}, nil
}

// DocumentAddress returns the counterparty's document DID address.
func (i *Identity) DocumentAddress() *identity.Address {
	return ffi.ToDIDAddress(i.h.DocumentAddress()).(*identity.Address)
}

// BiometricAnchorHash returns the 20-byte biometric anchor hash, or nil.
func (i *Identity) BiometricAnchorHash() []byte { return i.h.BiometricAnchorHash() }

// Encode returns the encoded bytes of the identity.
func (i *Identity) Encode() []byte { return i.h.Encode() }

// Relationship is a pairwise relationship.
type Relationship struct {
	h *ffi.PairwiseRelationship
}

// AsIdentity returns the identity this account presents to the counterparty.
func (r *Relationship) AsIdentity() *Identity { return &Identity{h: r.h.AsIdentity()} }

// WithIdentity returns the counterparty's identity.
func (r *Relationship) WithIdentity() *Identity { return &Identity{h: r.h.WithIdentity()} }

// Status returns the connection status.
func (r *Relationship) Status() Status { return Status(r.h.Status()) }

// DocumentAddress returns the introduced party's document DID address.
func (i *Introduction) DocumentAddress() *identity.Address {
	return ffi.ToDIDAddress(i.h.DocumentAddress()).(*identity.Address)
}

// Presentations returns the presentations shared by the sender.
func (i *Introduction) Presentations() []*credential.VerifiablePresentation {
	ps := i.h.Presentations()
	out := make([]*credential.VerifiablePresentation, len(ps))

	for j, p := range ps {
		out[j] = ffi.ToVerifiablePresentation(p).(*credential.VerifiablePresentation)
	}

	return out
}
