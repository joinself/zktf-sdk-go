package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

import (
	"runtime"
	"time"
	"unsafe"
)

// commitmentLen is the length of an identity document commitment hash.
const commitmentLen = 32

// IdentityOperation wraps a zktf_identity_operation handle (a hashgraph
// operation describing a change to an identity document — e.g. adding/removing
// keys, recovery). The slice exposes it as an opaque type; Wave 12 adds the
// operation builder and accessors.
type IdentityOperation struct {
	ptr *C.zktf_identity_operation
}

func newIdentityOperation(ptr *C.zktf_identity_operation) *IdentityOperation {
	if ptr == nil {
		return nil
	}
	o := &IdentityOperation{ptr: ptr}
	runtime.AddCleanup(o, func(ptr *C.zktf_identity_operation) {
		C.zktf_identity_operation_destroy(ptr)
	}, o.ptr)
	return o
}

// IdentityDocument wraps a zktf_identity_document handle (the resolved key state
// of an identity).
type IdentityDocument struct {
	ptr *C.zktf_identity_document
}

func newIdentityDocument(ptr *C.zktf_identity_document) *IdentityDocument {
	if ptr == nil {
		return nil
	}
	d := &IdentityDocument{ptr: ptr}
	runtime.AddCleanup(d, func(ptr *C.zktf_identity_document) {
		C.zktf_identity_document_destroy(ptr)
	}, d.ptr)
	return d
}

// Commitment returns the document's commitment hash, or nil if none is set.
func (d *IdentityDocument) Commitment() []byte {
	p := C.zktf_identity_document_commitment(d.ptr)
	if p == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(p), commitmentLen)
}

// Create returns an operation builder seeded from the current document state.
func (d *IdentityDocument) Create() *IdentityOperationBuilder {
	return newIdentityOperationBuilder(C.zktf_identity_document_create(d.ptr))
}

// SigningKeys returns the signing keys in the document. Pass nil for lookup to
// list every signing key in the latest snapshot.
func (d *IdentityDocument) SigningKeys(lookup *IdentityKeyLookup) []*SigningPublicKey {
	return signingPublicKeysFrom(C.zktf_identity_document_signing_keys(d.ptr, lookupPtr(lookup)))
}

// ExchangeKeys returns the exchange keys in the document. Pass nil for lookup
// to list every exchange key in the latest snapshot.
func (d *IdentityDocument) ExchangeKeys(lookup *IdentityKeyLookup) []*ExchangePublicKey {
	return exchangePublicKeysFrom(C.zktf_identity_document_exchange_keys(d.ptr, lookupPtr(lookup)))
}

// Descriptions returns the key descriptions in the document. Pass nil for
// lookup to list every description in the latest snapshot.
func (d *IdentityDocument) Descriptions(lookup *IdentityKeyLookup) []*IdentityOperationDescription {
	c := C.zktf_identity_document_descriptions(d.ptr, lookupPtr(lookup))
	if c == nil {
		return nil
	}
	defer C.zktf_collection_identity_operation_description_destroy(c)
	n := int(C.zktf_collection_identity_operation_description_len(c))
	out := make([]*IdentityOperationDescription, n)
	for i := 0; i < n; i++ {
		out[i] = newIdentityOperationDescription(C.zktf_collection_identity_operation_description_at(c, C.size_t(i)))
	}
	return out
}

// SigningKeyHasRoles reports whether the signing key holds every role in roles.
// Pass nil for lookup to evaluate against the latest snapshot.
func (d *IdentityDocument) SigningKeyHasRoles(key *SigningPublicKey, roles IdentityKeyRole, lookup *IdentityKeyLookup) bool {
	return bool(C.zktf_identity_document_signing_key_has_roles(d.ptr, key.ptr, C.zktf_identity_key_role(roles), lookupPtr(lookup)))
}

// ExchangeKeyHasRoles reports whether the exchange key holds every role in
// roles. Pass nil for lookup to evaluate against the latest snapshot.
func (d *IdentityDocument) ExchangeKeyHasRoles(key *ExchangePublicKey, roles IdentityKeyRole, lookup *IdentityKeyLookup) bool {
	return bool(C.zktf_identity_document_exchange_key_has_roles(d.ptr, key.ptr, C.zktf_identity_key_role(roles), lookupPtr(lookup)))
}

// SigningKeyValid reports whether the signing key is valid in the document.
// Pass nil for lookup to evaluate against the latest snapshot.
func (d *IdentityDocument) SigningKeyValid(key *SigningPublicKey, lookup *IdentityKeyLookup) bool {
	return bool(C.zktf_identity_document_signing_key_valid(d.ptr, key.ptr, lookupPtr(lookup)))
}

// ExchangeKeyValid reports whether the exchange key is valid in the document.
// Pass nil for lookup to evaluate against the latest snapshot.
func (d *IdentityDocument) ExchangeKeyValid(key *ExchangePublicKey, lookup *IdentityKeyLookup) bool {
	return bool(C.zktf_identity_document_exchange_key_valid(d.ptr, key.ptr, lookupPtr(lookup)))
}

// ThresholdMet reports whether signers collectively satisfy role's threshold.
// Pass nil for lookup to evaluate against the latest snapshot.
func (d *IdentityDocument) ThresholdMet(role IdentityKeyRole, signers []*SigningPublicKey, lookup *IdentityKeyLookup) bool {
	c := signingPublicKeyCollection(signers)
	defer destroySigningKeys(c)
	return bool(C.zktf_identity_document_threshold_met(d.ptr, C.zktf_identity_key_role(role), c, lookupPtr(lookup)))
}

// IdentityKeyLookup filters a document key query by snapshot time and/or roles.
// With no options applied, queries run against the latest snapshot with no role
// filter.
type IdentityKeyLookup struct {
	ptr *C.zktf_identity_key_lookup
}

// NewIdentityKeyLookup initializes a key lookup.
func NewIdentityKeyLookup() *IdentityKeyLookup {
	ptr := C.zktf_identity_key_lookup_init()
	l := &IdentityKeyLookup{ptr: ptr}
	runtime.AddCleanup(l, func(ptr *C.zktf_identity_key_lookup) {
		C.zktf_identity_key_lookup_destroy(ptr)
	}, l.ptr)
	return l
}

// AtTime evaluates the query against the document as it existed at the given
// unix timestamp (seconds).
func (l *IdentityKeyLookup) AtTime(unix int64) {
	C.zktf_identity_key_lookup_at_time(l.ptr, C.int64_t(unix))
}

// WithRoles restricts the result to keys holding every role in roles.
func (l *IdentityKeyLookup) WithRoles(roles IdentityKeyRole) {
	C.zktf_identity_key_lookup_with_roles(l.ptr, C.zktf_identity_key_role(roles))
}

// lookupPtr returns the C pointer for l, or nil if l is nil.
func lookupPtr(l *IdentityKeyLookup) *C.zktf_identity_key_lookup {
	if l == nil {
		return nil
	}
	return l.ptr
}

// IdentityResolve resolves the identity document for an address via callback,
// returning once the result has been delivered.
func (a *Account) IdentityResolve(address *DIDAddress, timeout time.Duration) (*IdentityDocument, error) {
	fut := C.zktf_account_identity_resolve(a.ptr, address.ptr)

	return AwaitIdentityDocument(fut, timeout)
}
