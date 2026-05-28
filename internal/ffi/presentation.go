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

// PresentationTypeCollection wraps a zktf_collection_presentation_type handle.
type PresentationTypeCollection struct {
	ptr *C.zktf_collection_presentation_type
}

func newPresentationTypeCollection(ptr *C.zktf_collection_presentation_type) *PresentationTypeCollection {
	if ptr == nil {
		return nil
	}
	c := &PresentationTypeCollection{ptr: ptr}
	runtime.AddCleanup(c, func(ptr *C.zktf_collection_presentation_type) {
		C.zktf_collection_presentation_type_destroy(ptr)
	}, c.ptr)
	return c
}

// NewPresentationTypes builds a presentation type collection from type strings.
func NewPresentationTypes(types []string) *PresentationTypeCollection {
	ptr := C.zktf_collection_presentation_type_init()
	for _, t := range types {
		ct := cstring(t)
		C.zktf_collection_presentation_type_append(ptr, ct)
		free(unsafe.Pointer(ct))
	}
	return newPresentationTypeCollection(ptr)
}

// Strings returns the type strings in the collection.
func (c *PresentationTypeCollection) Strings() []string {
	n := int(C.zktf_collection_presentation_type_len(c.ptr))
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = C.GoString(C.zktf_collection_presentation_type_at(c.ptr, C.size_t(i)))
	}
	return out
}

// Presentation wraps an unsigned zktf_presentation handle.
type Presentation struct {
	ptr *C.zktf_presentation
}

func newPresentation(ptr *C.zktf_presentation) *Presentation {
	if ptr == nil {
		return nil
	}
	p := &Presentation{ptr: ptr}
	runtime.AddCleanup(p, func(ptr *C.zktf_presentation) {
		C.zktf_presentation_destroy(ptr)
	}, p.ptr)
	return p
}

// PresentationBuilder wraps a zktf_presentation_builder handle.
type PresentationBuilder struct {
	ptr *C.zktf_presentation_builder
}

// NewPresentationBuilder initializes a new presentation builder.
func NewPresentationBuilder() *PresentationBuilder {
	ptr := C.zktf_presentation_builder_init()
	b := &PresentationBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_presentation_builder) {
		C.zktf_presentation_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// PresentationType sets the presentation's types.
func (b *PresentationBuilder) PresentationType(types *PresentationTypeCollection) *PresentationBuilder {
	C.zktf_presentation_builder_presentation_type(b.ptr, types.ptr)
	return b
}

// CredentialAdd adds a verifiable credential to the presentation.
func (b *PresentationBuilder) CredentialAdd(credential *VerifiableCredential) *PresentationBuilder {
	C.zktf_presentation_builder_credential_add(b.ptr, credential.ptr)
	return b
}

// Holder sets the holder/bearer address.
func (b *PresentationBuilder) Holder(holder *DIDAddress) *PresentationBuilder {
	C.zktf_presentation_builder_holder(b.ptr, holder.ptr)
	return b
}

// Finish finalizes the unsigned presentation.
func (b *PresentationBuilder) Finish() (*Presentation, error) {
	var ptr *C.zktf_presentation
	if err := status(C.zktf_presentation_builder_finish(b.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newPresentation(ptr), nil
}

// VerifiablePresentation wraps a signed zktf_verifiable_presentation handle.
type VerifiablePresentation struct {
	ptr *C.zktf_verifiable_presentation
}

func newVerifiablePresentation(ptr *C.zktf_verifiable_presentation) *VerifiablePresentation {
	if ptr == nil {
		return nil
	}
	p := &VerifiablePresentation{ptr: ptr}
	runtime.AddCleanup(p, func(ptr *C.zktf_verifiable_presentation) {
		C.zktf_verifiable_presentation_destroy(ptr)
	}, p.ptr)
	return p
}

// VerifiablePresentationDecode decodes a JSON-encoded verifiable presentation.
func VerifiablePresentationDecode(data []byte) (*VerifiablePresentation, error) {
	buf, length := cbytes(data)
	defer free(unsafe.Pointer(buf))

	var ptr *C.zktf_verifiable_presentation
	if err := status(C.zktf_verifiable_presentation_decode(&ptr, buf, length)); err != nil {
		return nil, err
	}
	return newVerifiablePresentation(ptr), nil
}

// Validate returns an error if the presentation is invalid.
func (p *VerifiablePresentation) Validate() error {
	return status(C.zktf_verifiable_presentation_validate(p.ptr))
}

// Types returns the presentation's type strings.
func (p *VerifiablePresentation) Types() []string {
	return presentationTypesFrom(C.zktf_verifiable_presentation_type_of(p.ptr))
}

// Holder returns the holder address, or nil.
func (p *VerifiablePresentation) Holder() *DIDAddress {
	return newDIDAddress(C.zktf_verifiable_presentation_holder(p.ptr))
}

// Credentials returns the credentials contained in the presentation.
func (p *VerifiablePresentation) Credentials() []*VerifiableCredential {
	return verifiableCredentialsFrom(C.zktf_verifiable_presentation_credentials(p.ptr))
}

// Encode returns the JSON-encoded presentation.
func (p *VerifiablePresentation) Encode() ([]byte, error) {
	var buf *C.zktf_bytes_buffer
	if err := status(C.zktf_verifiable_presentation_encode(p.ptr, &buf)); err != nil {
		return nil, err
	}
	return goBytesFromBuffer(buf), nil
}
