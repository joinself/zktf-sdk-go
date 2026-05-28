package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

import (
	"runtime"
	"time"
)

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

// SigningKeys returns all signing keys in the document (no role/time filter).
func (d *IdentityDocument) SigningKeys() []*SigningPublicKey {
	return signingPublicKeysFrom(C.zktf_identity_document_signing_keys(d.ptr, nil))
}

// ExchangeKeys returns all exchange keys in the document (no role/time filter).
func (d *IdentityDocument) ExchangeKeys() []*ExchangePublicKey {
	return exchangePublicKeysFrom(C.zktf_identity_document_exchange_keys(d.ptr, nil))
}

// SigningKeyValid reports whether the signing key is currently valid in the
// document.
func (d *IdentityDocument) SigningKeyValid(key *SigningPublicKey) bool {
	return bool(C.zktf_identity_document_signing_key_valid(d.ptr, key.ptr, nil))
}

// ExchangeKeyValid reports whether the exchange key is currently valid in the
// document.
func (d *IdentityDocument) ExchangeKeyValid(key *ExchangePublicKey) bool {
	return bool(C.zktf_identity_document_exchange_key_valid(d.ptr, key.ptr, nil))
}

// IdentityResolve resolves the identity document for an address via callback,
// returning once the result has been delivered.
func (a *Account) IdentityResolve(address *DIDAddress, timeout time.Duration) (*IdentityDocument, error) {
	fut := C.zktf_account_identity_resolve(a.ptr, address.ptr)

	return AwaitIdentityDocument(fut, timeout)
}
