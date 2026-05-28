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

// PairwiseStatus mirrors zktf_pairwise_status.
type PairwiseStatus uint32

const (
	PairwiseStatusPending     PairwiseStatus = C.CONNECTION_STATUS_PENDING
	PairwiseStatusNegotiating PairwiseStatus = C.CONNECTION_STATUS_NEGOTIATING
	PairwiseStatusEstablished PairwiseStatus = C.CONNECTION_STATUS_ESTABLISHED
)

const biometricAnchorHashLen = 20

// PairwiseIdentity wraps a zktf_pairwise_identity handle.
type PairwiseIdentity struct {
	ptr *C.zktf_pairwise_identity
}

func newPairwiseIdentity(ptr *C.zktf_pairwise_identity) *PairwiseIdentity {
	if ptr == nil {
		return nil
	}
	i := &PairwiseIdentity{ptr: ptr}
	runtime.AddCleanup(i, func(ptr *C.zktf_pairwise_identity) {
		C.zktf_pairwise_identity_destroy(ptr)
	}, i.ptr)
	return i
}

// PairwiseIdentityDecode decodes an encoded pairwise identity.
func PairwiseIdentityDecode(data []byte) (*PairwiseIdentity, error) {
	buf, length := cbytes(data)
	defer free(unsafe.Pointer(buf))

	var ptr *C.zktf_pairwise_identity
	if err := status(C.zktf_pairwise_identity_decode(buf, length, &ptr)); err != nil {
		return nil, err
	}
	return newPairwiseIdentity(ptr), nil
}

// DocumentAddress returns the counterparty's document DID address.
func (i *PairwiseIdentity) DocumentAddress() *DIDAddress {
	return newDIDAddress(C.zktf_pairwise_identity_document_address(i.ptr))
}

// BiometricAnchorHash returns the 20-byte biometric anchor hash, or nil.
func (i *PairwiseIdentity) BiometricAnchorHash() []byte {
	p := C.zktf_pairwise_identity_biometric_anchor_hash(i.ptr)
	if p == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(p), biometricAnchorHashLen)
}

// Encode returns the encoded bytes of the identity.
func (i *PairwiseIdentity) Encode() []byte {
	return goBytesFromBuffer(C.zktf_pairwise_identity_encode(i.ptr))
}

// PairwiseRelationship wraps a zktf_pairwise_relationship handle.
type PairwiseRelationship struct {
	ptr *C.zktf_pairwise_relationship
}

func newPairwiseRelationship(ptr *C.zktf_pairwise_relationship) *PairwiseRelationship {
	if ptr == nil {
		return nil
	}
	r := &PairwiseRelationship{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_pairwise_relationship) {
		C.zktf_pairwise_relationship_destroy(ptr)
	}, r.ptr)
	return r
}

// AsIdentity returns the identity this account presents to the counterparty.
func (r *PairwiseRelationship) AsIdentity() *PairwiseIdentity {
	return newPairwiseIdentity(C.zktf_pairwise_relationship_as_identity(r.ptr))
}

// WithIdentity returns the counterparty's identity.
func (r *PairwiseRelationship) WithIdentity() *PairwiseIdentity {
	return newPairwiseIdentity(C.zktf_pairwise_relationship_with_identity(r.ptr))
}

// Status returns the connection status.
func (r *PairwiseRelationship) Status() PairwiseStatus {
	return PairwiseStatus(C.zktf_pairwise_relationship_status(r.ptr))
}

// PairwiseIntroduction wraps a zktf_pairwise_introduction handle (the result of
// validating an introduction message).
type PairwiseIntroduction struct {
	ptr *C.zktf_pairwise_introduction
}

func newPairwiseIntroduction(ptr *C.zktf_pairwise_introduction) *PairwiseIntroduction {
	if ptr == nil {
		return nil
	}
	i := &PairwiseIntroduction{ptr: ptr}
	runtime.AddCleanup(i, func(ptr *C.zktf_pairwise_introduction) {
		C.zktf_pairwise_introduction_destroy(ptr)
	}, i.ptr)
	return i
}

// DocumentAddress returns the introduced party's document DID address.
func (i *PairwiseIntroduction) DocumentAddress() *DIDAddress {
	return newDIDAddress(C.zktf_pairwise_introduction_document_address(i.ptr))
}

// Presentations returns the presentations shared by the sender.
func (i *PairwiseIntroduction) Presentations() []*VerifiablePresentation {
	return verifiablePresentationsFrom(C.zktf_pairwise_introduction_presentations(i.ptr))
}

// pairwise_lookup/store/validate_introduction are mobile-only in the native
// SDK; they are intentionally not wrapped on this server-side build. A server
// receives an introduction inside a message.Introduction and reads the embedded
// PairwiseIntroduction's document address + presentations directly.
