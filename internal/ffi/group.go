package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

import (
	"runtime"
	"time"
)

// Group wraps a zktf_group handle (a member of an established encrypted group
// session). The slice exposes Address() and Members(); deeper accessors can be
// added without changing the wrapping pattern.
type Group struct {
	ptr *C.zktf_group
}

func newGroup(ptr *C.zktf_group) *Group {
	if ptr == nil {
		return nil
	}
	g := &Group{ptr: ptr}
	runtime.AddCleanup(g, func(ptr *C.zktf_group) {
		C.zktf_group_destroy(ptr)
	}, g.ptr)
	return g
}

// Address returns the group address.
func (g *Group) Address() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_group_address(g.ptr))
}

// Members returns the addresses of the group members.
func (g *Group) Members() []*SigningPublicKey {
	return signingPublicKeysFrom(C.zktf_group_members(g.ptr))
}

// MemberAs returns the address this account presents within the group.
func (g *Group) MemberAs() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_group_member_as(g.ptr))
}

// GroupNegotiateOutOfBand creates a key package for establishing an encrypted
// session with this account out of band (e.g. for inclusion in a discovery
// request). expiresUnix of 0 means no expiry.
func (a *Account) GroupNegotiateOutOfBand(as *SigningPublicKey, expiresUnix int64) (*CryptoKeyPackage, error) {
	var out *C.zktf_crypto_key_package

	if err := status(C.zktf_account_group_negotiate_out_of_band(a.ptr, as.ptr, C.int64_t(expiresUnix), &out)); err != nil {
		return nil, err
	}

	return newCryptoKeyPackage(out), nil
}

// GroupEstablish uses a received key package to establish an encrypted group
// session via callback.
func (a *Account) GroupEstablish(as *SigningPublicKey, keyPackage *CryptoKeyPackage, timeout time.Duration) (*Group, error) {
	fut := C.zktf_account_group_establish(a.ptr, as.ptr, keyPackage.ptr)

	return AwaitGroup(fut, timeout)
}

// GroupAccept accepts a received welcome to join an encrypted group session
// via callback.
func (a *Account) GroupAccept(as *SigningPublicKey, welcome *CryptoWelcome, timeout time.Duration) (*Group, error) {
	fut := C.zktf_account_group_accept(a.ptr, as.ptr, welcome.ptr)

	return AwaitGroup(fut, timeout)
}

// GroupLookup builds a query for groups on an account.
type GroupLookup struct {
	ptr *C.zktf_group_lookup
}

// NewGroupLookup initializes a group lookup query.
func NewGroupLookup() *GroupLookup {
	ptr := C.zktf_group_lookup_init()
	l := &GroupLookup{ptr: ptr}
	runtime.AddCleanup(l, func(ptr *C.zktf_group_lookup) {
		C.zktf_group_lookup_destroy(ptr)
	}, l.ptr)
	return l
}

// ByAddress restricts the lookup to a group at the given address.
func (l *GroupLookup) ByAddress(address *SigningPublicKey) *GroupLookup {
	C.zktf_group_lookup_by_address(l.ptr, address.ptr)
	return l
}

// ByMember restricts the lookup to groups including the given member.
func (l *GroupLookup) ByMember(member *SigningPublicKey) *GroupLookup {
	C.zktf_group_lookup_by_member(l.ptr, member.ptr)
	return l
}

// GroupUpdateBuilder builds a group update (add/remove members).
type GroupUpdateBuilder struct {
	ptr *C.zktf_group_update_builder
}

// NewGroupUpdateBuilder initializes an update builder for the given group.
func NewGroupUpdateBuilder(g *Group) *GroupUpdateBuilder {
	ptr := C.zktf_group_update_builder_init(g.ptr)
	b := &GroupUpdateBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_group_update_builder) {
		C.zktf_group_update_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// AddMembers stages every member in `packages` to be added.
func (b *GroupUpdateBuilder) AddMembers(packages []*CryptoKeyPackage) *GroupUpdateBuilder {
	c := cryptoKeyPackageCollection(packages)
	defer destroyCryptoKeyPackages(c)
	C.zktf_group_update_builder_add_members(b.ptr, c)
	return b
}

// RemoveMembers stages every member in `members` to be removed.
func (b *GroupUpdateBuilder) RemoveMembers(members []*SigningPublicKey) *GroupUpdateBuilder {
	c := signingPublicKeyCollection(members)
	defer destroySigningKeys(c)
	C.zktf_group_update_builder_remove_members(b.ptr, c)
	return b
}

// AsProposal marks the update to be sent as a proposal (not auto-committed).
func (b *GroupUpdateBuilder) AsProposal() *GroupUpdateBuilder {
	C.zktf_group_update_builder_as_proposal(b.ptr)
	return b
}

// Finish validates the staged changes and produces a request.
func (b *GroupUpdateBuilder) Finish() (*GroupUpdateRequest, error) {
	var out *C.zktf_group_update_request
	if err := status(C.zktf_group_update_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	r := &GroupUpdateRequest{ptr: out}
	runtime.AddCleanup(r, func(ptr *C.zktf_group_update_request) {
		C.zktf_group_update_request_destroy(ptr)
	}, r.ptr)
	return r, nil
}

// GroupUpdateRequest is a signed group update ready to publish via Account.GroupUpdate.
type GroupUpdateRequest struct {
	ptr *C.zktf_group_update_request
}

// GroupLookup returns groups matching the lookup query.
func (a *Account) GroupLookup(l *GroupLookup) ([]*Group, error) {
	var c *C.zktf_collection_group
	if err := status(C.zktf_account_group_lookup(a.ptr, l.ptr, &c)); err != nil {
		return nil, err
	}
	return groupsFrom(c), nil
}

// GroupUpdate publishes a built group update.
func (a *Account) GroupUpdate(r *GroupUpdateRequest) error {
	return status(C.zktf_account_group_update(a.ptr, r.ptr))
}

// GroupLeave leaves a group.
func (a *Account) GroupLeave(g *Group) error {
	return status(C.zktf_account_group_leave(a.ptr, g.ptr))
}
