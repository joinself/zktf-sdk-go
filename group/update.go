package group

import (
	"github.com/joinself/zktf-sdk-go/crypto"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// UpdateBuilder builds a group update (add/remove members), to be published via
// account.Account.GroupUpdate.
type UpdateBuilder struct {
	h *ffi.GroupUpdateBuilder
}

// UpdateRequest is a built, signed group update ready to publish.
type UpdateRequest struct {
	h *ffi.GroupUpdateRequest
}

func init() {
	ffi.GroupUpdateRequestOf = func(o any) *ffi.GroupUpdateRequest { return o.(*UpdateRequest).h }
	ffi.ToGroupUpdateRequest = func(h *ffi.GroupUpdateRequest) any { return &UpdateRequest{h: h} }
}

// NewUpdateBuilder starts building an update for the given group.
func NewUpdateBuilder(g *Group) *UpdateBuilder {
	return &UpdateBuilder{h: ffi.NewGroupUpdateBuilder(g.h)}
}

// AddMembers stages new members to be added via their MLS key packages.
func (b *UpdateBuilder) AddMembers(packages []*crypto.KeyPackage) *UpdateBuilder {
	ps := make([]*ffi.CryptoKeyPackage, len(packages))
	for i, p := range packages {
		ps[i] = ffi.CryptoKeyPackageOf(p)
	}

	b.h.AddMembers(ps)

	return b
}

// RemoveMembers stages members to be removed by their addresses.
func (b *UpdateBuilder) RemoveMembers(members []*signing.PublicKey) *UpdateBuilder {
	ms := make([]*ffi.SigningPublicKey, len(members))
	for i, m := range members {
		ms[i] = ffi.SigningPublicKeyOf(m)
	}

	b.h.RemoveMembers(ms)

	return b
}

// AsProposal marks the update to be sent as a proposal (not auto-committed).
func (b *UpdateBuilder) AsProposal() *UpdateBuilder {
	b.h.AsProposal()
	return b
}

// Finish validates the staged changes and produces an update request.
func (b *UpdateBuilder) Finish() (*UpdateRequest, error) {
	r, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &UpdateRequest{h: r}, nil
}
