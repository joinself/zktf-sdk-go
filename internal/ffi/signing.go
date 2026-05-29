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

// signingPublicKeyEncodedLen is the length of a hex-encoded signing address.
const signingPublicKeyEncodedLen = 66

// signingPublicKeyBytesLen is the length of a raw signing public key.
const signingPublicKeyBytesLen = 33

// SigningPublicKey wraps a zktf_signing_public_key handle. The C pointer is held
// in the unexported ptr field so it never leaks into an exported signature.
type SigningPublicKey struct {
	ptr *C.zktf_signing_public_key
}

// newSigningPublicKey adopts a C handle and attaches a finalizer that frees it.
func newSigningPublicKey(ptr *C.zktf_signing_public_key) *SigningPublicKey {
	if ptr == nil {
		return nil
	}
	k := &SigningPublicKey{ptr: ptr}
	runtime.AddCleanup(k, func(ptr *C.zktf_signing_public_key) {
		C.zktf_signing_public_key_destroy(ptr)
	}, k.ptr)
	return k
}

// SigningPublicKeyFromAddress decodes a hex address into a public key.
func SigningPublicKeyFromAddress(hex string) (*SigningPublicKey, error) {
	buf, length := cbytes([]byte(hex))
	defer free(unsafe.Pointer(buf))

	var ptr *C.zktf_signing_public_key
	if err := status(C.zktf_signing_public_key_decode(&ptr, buf, length)); err != nil {
		return nil, err
	}
	return newSigningPublicKey(ptr), nil
}

// SigningPublicKeyFromBytes constructs a public key from its raw bytes.
func SigningPublicKeyFromBytes(data []byte) (*SigningPublicKey, error) {
	buf, length := cbytes(data)
	defer free(unsafe.Pointer(buf))

	var ptr *C.zktf_signing_public_key
	if err := status(C.zktf_signing_public_key_from_bytes(&ptr, buf, length)); err != nil {
		return nil, err
	}
	return newSigningPublicKey(ptr), nil
}

// String returns the hex encoded address.
func (k *SigningPublicKey) String() string {
	buf := C.malloc(signingPublicKeyEncodedLen)
	defer C.free(buf)

	if err := status(C.zktf_signing_public_key_encode(
		k.ptr, (*C.uint8_t)(buf), signingPublicKeyEncodedLen,
	)); err != nil {
		return ""
	}
	return string(C.GoBytes(buf, signingPublicKeyEncodedLen))
}

// Bytes returns the raw bytes of the public key.
func (k *SigningPublicKey) Bytes() []byte {
	buf := C.malloc(signingPublicKeyBytesLen)
	defer C.free(buf)

	if err := status(C.zktf_signing_public_key_as_bytes(
		k.ptr, (*C.uint8_t)(buf), signingPublicKeyBytesLen,
	)); err != nil {
		return nil
	}
	return C.GoBytes(buf, signingPublicKeyBytesLen)
}

// Matches reports whether two public keys are equal.
// Verify reports whether signature is a valid signature of message by this key.
func (k *SigningPublicKey) Verify(message, signature []byte) bool {
	msgBuf, msgLen := cbytes(message)
	defer free(unsafe.Pointer(msgBuf))

	sigBuf, sigLen := cbytes(signature)
	defer free(unsafe.Pointer(sigBuf))

	return bool(C.zktf_signing_public_key_verify(k.ptr, msgBuf, msgLen, sigBuf, sigLen))
}

func (k *SigningPublicKey) Matches(other *SigningPublicKey) bool {
	if other == nil {
		return false
	}
	return bool(C.zktf_signing_public_key_matches(k.ptr, other.ptr))
}
