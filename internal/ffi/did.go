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

// DIDAddress wraps a zktf_did_address handle.
type DIDAddress struct {
	ptr *C.zktf_did_address
}

func newDIDAddress(ptr *C.zktf_did_address) *DIDAddress {
	if ptr == nil {
		return nil
	}
	a := &DIDAddress{ptr: ptr}
	runtime.AddCleanup(a, func(ptr *C.zktf_did_address) {
		C.zktf_did_address_destroy(ptr)
	}, a.ptr)
	return a
}

// DIDAddressKey builds a key-method DID address from a signing key.
func DIDAddressKey(key *SigningPublicKey) *DIDAddress {
	return newDIDAddress(C.zktf_did_address_key(key.ptr))
}

// DIDAddressDecode decodes a DID string into an address.
func DIDAddressDecode(did string) (*DIDAddress, error) {
	cdid := cstring(did)
	defer free(unsafe.Pointer(cdid))

	var ptr *C.zktf_did_address
	if err := status(C.zktf_did_address_decode(cdid, &ptr)); err != nil {
		return nil, err
	}
	return newDIDAddress(ptr), nil
}

// Address returns the signing public key embedded in the DID address.
func (a *DIDAddress) Address() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_did_address_address(a.ptr))
}

// String returns the encoded DID string.
func (a *DIDAddress) String() string {
	buf := C.zktf_did_address_encode(a.ptr)
	if buf == nil {
		return ""
	}
	defer C.zktf_string_buffer_destroy(buf)
	return C.GoString(C.zktf_string_buffer_ptr(buf))
}
