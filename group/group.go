// Package group provides established encrypted group sessions.
package group

import (
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// Group is a member of an established encrypted group session.
type Group struct {
	h *ffi.Group
}

// Lookup is the constructed lookup query.
type Lookup struct {
	h *ffi.GroupLookup
}

func init() {
	ffi.GroupOf = func(o any) *ffi.Group { return o.(*Group).h }
	ffi.ToGroup = func(h *ffi.Group) any { return &Group{h: h} }

	ffi.GroupLookupOf = func(o any) *ffi.GroupLookup { return o.(*Lookup).h }
	ffi.ToGroupLookup = func(h *ffi.GroupLookup) any { return &Lookup{h: h} }
}

// Address returns the group address.
func (g *Group) Address() *signing.PublicKey {
	return ffi.ToSigningPublicKey(g.h.Address()).(*signing.PublicKey)
}

// Members returns the addresses of the group members.
func (g *Group) Members() []*signing.PublicKey {
	ks := g.h.Members()
	out := make([]*signing.PublicKey, len(ks))

	for i, k := range ks {
		out[i] = ffi.ToSigningPublicKey(k).(*signing.PublicKey)
	}

	return out
}

// MemberAs returns the address this account presents within the group.
func (g *Group) MemberAs() *signing.PublicKey {
	return ffi.ToSigningPublicKey(g.h.MemberAs()).(*signing.PublicKey)
}

// LookupOption is one of the variadic filters accepted by GroupLookup.
type LookupOption func(*lookupOpts)

type lookupOpts struct {
	address *signing.PublicKey
	member  *signing.PublicKey
}

// ByAddress restricts the lookup to a group at the given address.
func ByAddress(address *signing.PublicKey) LookupOption {
	return func(o *lookupOpts) {
		o.address = address
	}
}

// ByMember restricts the lookup to groups including the given member.
func ByMember(member *signing.PublicKey) LookupOption {
	return func(o *lookupOpts) {
		o.member = member
	}
}

// BuildLookup applies the given options and returns a Lookup ready to pass to
// the FFI. For use by other zktf-sdk-go packages.
func BuildLookup(options ...LookupOption) *Lookup {
	var o lookupOpts
	for _, opt := range options {
		opt(&o)
	}

	l := &Lookup{h: ffi.NewGroupLookup()}

	if o.address != nil {
		l.h.ByAddress(ffi.SigningPublicKeyOf(o.address))
	}

	if o.member != nil {
		l.h.ByMember(ffi.SigningPublicKeyOf(o.member))
	}

	return l
}
