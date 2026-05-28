package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

import "unsafe"

// goBytesFromBuffer copies a zktf_bytes_buffer into a Go slice and destroys the
// buffer. A nil buffer yields nil.
func goBytesFromBuffer(buf *C.zktf_bytes_buffer) []byte {
	if buf == nil {
		return nil
	}
	defer C.zktf_bytes_buffer_destroy(buf)
	return C.GoBytes(
		unsafe.Pointer(C.zktf_bytes_buffer_buf(buf)),
		C.int(C.zktf_bytes_buffer_len(buf)),
	)
}

// didAddressesFrom copies a caller-owned zktf_collection_did_address into Go
// wrappers and destroys the collection.
func didAddressesFrom(c *C.zktf_collection_did_address) []*DIDAddress {
	if c == nil {
		return nil
	}
	defer C.zktf_collection_did_address_destroy(c)
	n := int(C.zktf_collection_did_address_len(c))
	out := make([]*DIDAddress, n)
	for i := 0; i < n; i++ {
		out[i] = newDIDAddress(C.zktf_collection_did_address_at(c, C.size_t(i)))
	}
	return out
}

// groupsFrom copies a caller-owned zktf_collection_group into Go wrappers and
// destroys the collection.
func groupsFrom(c *C.zktf_collection_group) []*Group {
	if c == nil {
		return nil
	}
	defer C.zktf_collection_group_destroy(c)
	n := int(C.zktf_collection_group_len(c))
	out := make([]*Group, n)
	for i := 0; i < n; i++ {
		out[i] = newGroup(C.zktf_collection_group_at(c, C.size_t(i)))
	}
	return out
}

// cryptoKeyPackageCollection builds a zktf_collection_crypto_key_package from Go
// key packages. Caller must destroy with destroyCryptoKeyPackages.
func cryptoKeyPackageCollection(packages []*CryptoKeyPackage) *C.zktf_collection_crypto_key_package {
	collection := C.zktf_collection_crypto_key_package_init()
	for _, p := range packages {
		C.zktf_collection_crypto_key_package_append(collection, p.ptr)
	}
	return collection
}

func destroyCryptoKeyPackages(c *C.zktf_collection_crypto_key_package) {
	C.zktf_collection_crypto_key_package_destroy(c)
}

// tokensFrom copies a caller-owned zktf_collection_token into Go wrappers and
// destroys the collection.
func tokensFrom(c *C.zktf_collection_token) []*Token {
	if c == nil {
		return nil
	}
	defer C.zktf_collection_token_destroy(c)
	n := int(C.zktf_collection_token_len(c))
	out := make([]*Token, n)
	for i := 0; i < n; i++ {
		out[i] = newToken(C.zktf_collection_token_at(c, C.size_t(i)))
	}
	return out
}

// objectsFrom copies a caller-owned zktf_collection_object into Go wrappers and
// destroys the collection.
func objectsFrom(c *C.zktf_collection_object) []*Object {
	if c == nil {
		return nil
	}
	defer C.zktf_collection_object_destroy(c)
	n := int(C.zktf_collection_object_len(c))
	out := make([]*Object, n)
	for i := 0; i < n; i++ {
		out[i] = newObject(C.zktf_collection_object_at(c, C.size_t(i)))
	}
	return out
}

// exchangePublicKeysFrom copies a zktf_collection_exchange_public_key into Go
// wrappers and destroys the collection.
func exchangePublicKeysFrom(c *C.zktf_collection_exchange_public_key) []*ExchangePublicKey {
	if c == nil {
		return nil
	}
	defer C.zktf_collection_exchange_public_key_destroy(c)
	n := int(C.zktf_collection_exchange_public_key_len(c))
	out := make([]*ExchangePublicKey, n)
	for i := 0; i < n; i++ {
		out[i] = newExchangePublicKey(C.zktf_collection_exchange_public_key_at(c, C.size_t(i)))
	}
	return out
}

// verifiableCredentialsFrom copies a zktf_collection_verifiable_credential into
// Go wrappers and destroys the collection.
func verifiableCredentialsFrom(c *C.zktf_collection_verifiable_credential) []*VerifiableCredential {
	if c == nil {
		return nil
	}
	defer C.zktf_collection_verifiable_credential_destroy(c)
	n := int(C.zktf_collection_verifiable_credential_len(c))
	out := make([]*VerifiableCredential, n)
	for i := 0; i < n; i++ {
		// the collection owns its elements; wrap a borrowed pointer without a
		// finalizer is unsafe, so we treat elements as owned copies here.
		out[i] = newVerifiableCredential(C.zktf_collection_verifiable_credential_at(c, C.size_t(i)))
	}
	return out
}

// verifiablePresentationsFrom copies a zktf_collection_verifiable_presentation
// into Go wrappers and destroys the collection.
func verifiablePresentationsFrom(c *C.zktf_collection_verifiable_presentation) []*VerifiablePresentation {
	if c == nil {
		return nil
	}
	defer C.zktf_collection_verifiable_presentation_destroy(c)
	n := int(C.zktf_collection_verifiable_presentation_len(c))
	out := make([]*VerifiablePresentation, n)
	for i := 0; i < n; i++ {
		out[i] = newVerifiablePresentation(C.zktf_collection_verifiable_presentation_at(c, C.size_t(i)))
	}
	return out
}

// presentationTypesFrom copies a zktf_collection_presentation_type into Go
// strings and destroys the collection.
func presentationTypesFrom(c *C.zktf_collection_presentation_type) []string {
	if c == nil {
		return nil
	}
	defer C.zktf_collection_presentation_type_destroy(c)
	n := int(C.zktf_collection_presentation_type_len(c))
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = C.GoString(C.zktf_collection_presentation_type_at(c, C.size_t(i)))
	}
	return out
}

// messageIDsFrom copies a zktf_collection_message_id into Go slices and destroys
// the collection.
func messageIDsFrom(c *C.zktf_collection_message_id) [][]byte {
	if c == nil {
		return nil
	}
	defer C.zktf_collection_message_id_destroy(c)
	n := int(C.zktf_collection_message_id_len(c))
	out := make([][]byte, n)
	for i := 0; i < n; i++ {
		out[i] = C.GoBytes(unsafe.Pointer(C.zktf_collection_message_id_at(c, C.size_t(i))), messageIDLen)
	}
	return out
}

// signingPublicKeyCollection builds a zktf_collection_signing_public_key from Go
// keys for passing to C. The caller must destroy it with destroySigningKeys.
func signingPublicKeyCollection(keys []*SigningPublicKey) *C.zktf_collection_signing_public_key {
	collection := C.zktf_collection_signing_public_key_init()
	for _, k := range keys {
		C.zktf_collection_signing_public_key_append(collection, k.ptr)
	}
	return collection
}

func destroySigningKeys(c *C.zktf_collection_signing_public_key) {
	C.zktf_collection_signing_public_key_destroy(c)
}

// signingPublicKeysFrom copies a caller-owned C collection into Go keys and
// destroys the collection.
func signingPublicKeysFrom(c *C.zktf_collection_signing_public_key) []*SigningPublicKey {
	if c == nil {
		return nil
	}
	defer C.zktf_collection_signing_public_key_destroy(c)
	n := int(C.zktf_collection_signing_public_key_len(c))
	out := make([]*SigningPublicKey, n)
	for i := 0; i < n; i++ {
		out[i] = newSigningPublicKey(C.zktf_collection_signing_public_key_at(c, C.size_t(i)))
	}
	return out
}
