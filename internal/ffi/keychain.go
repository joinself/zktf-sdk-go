package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// signatureLen is the byte length of a signature produced by KeychainSign.
const signatureLen = 64

// KeychainSigningCreate generates a new signing key in the keychain and returns
// its public address.
func (a *Account) KeychainSigningCreate() (*SigningPublicKey, error) {
	var out *C.zktf_signing_public_key

	if err := status(C.zktf_account_keychain_signing_create(a.ptr, &out)); err != nil {
		return nil, err
	}

	return newSigningPublicKey(out), nil
}

// KeychainExchangeCreate generates a new exchange key in the keychain and
// returns its public address.
func (a *Account) KeychainExchangeCreate() (*ExchangePublicKey, error) {
	var out *C.zktf_exchange_public_key

	if err := status(C.zktf_account_keychain_exchange_create(a.ptr, &out)); err != nil {
		return nil, err
	}

	return newExchangePublicKey(out), nil
}

// KeychainSign signs payload with the keychain key identified by address.
func (a *Account) KeychainSign(address *SigningPublicKey, payload []byte) ([]byte, error) {
	payloadBuf, payloadLen := cbytes(payload)
	defer free(unsafe.Pointer(payloadBuf))

	sig := make([]byte, signatureLen)

	if err := status(C.zktf_account_keychain_sign(
		a.ptr,
		address.ptr,
		payloadBuf, payloadLen,
		(*C.uint8_t)(unsafe.Pointer(&sig[0])), C.size_t(signatureLen),
	)); err != nil {
		return nil, err
	}

	return sig, nil
}

// KeychainLookup builds a query for keychain signing keys.
type KeychainLookup struct {
	ptr *C.zktf_keychain_lookup
}

// NewKeychainLookup initializes a keychain lookup query. With no filters
// applied it matches every signing key in the keychain.
func NewKeychainLookup() *KeychainLookup {
	ptr := C.zktf_keychain_lookup_init()
	l := &KeychainLookup{ptr: ptr}
	runtime.AddCleanup(l, func(ptr *C.zktf_keychain_lookup) {
		C.zktf_keychain_lookup_destroy(ptr)
	}, l.ptr)

	return l
}

// ByIdentity restricts the lookup to keys associated with identity.
func (l *KeychainLookup) ByIdentity(identity *SigningPublicKey) *KeychainLookup {
	C.zktf_keychain_lookup_by_identity(l.ptr, identity.ptr)
	return l
}

// WithRoles restricts the lookup to keys carrying every role in roles. Only
// applies in combination with ByIdentity.
func (l *KeychainLookup) WithRoles(roles IdentityKeyRole) *KeychainLookup {
	C.zktf_keychain_lookup_with_roles(l.ptr, C.zktf_identity_key_role(roles))
	return l
}

// KeychainLookup resolves the signing keys held in the keychain that satisfy
// the lookup query.
func (a *Account) KeychainLookup(lookup *KeychainLookup) ([]*SigningPublicKey, error) {
	var c *C.zktf_collection_signing_public_key

	if err := status(C.zktf_account_keychain_lookup(a.ptr, lookup.ptr, &c)); err != nil {
		return nil, err
	}

	return signingPublicKeysFrom(c), nil
}
