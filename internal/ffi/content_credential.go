package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import "runtime"

// CredentialContent wraps a zktf_message_content_credential handle.
type CredentialContent struct {
	ptr *C.zktf_message_content_credential
}

func newCredentialContent(ptr *C.zktf_message_content_credential) *CredentialContent {
	if ptr == nil {
		return nil
	}
	c := &CredentialContent{ptr: ptr}
	runtime.AddCleanup(c, func(ptr *C.zktf_message_content_credential) {
		C.zktf_message_content_credential_destroy(ptr)
	}, c.ptr)
	return c
}

// CredentialContentFromContent decodes message content as a credential payload.
func CredentialContentFromContent(content *Content) (*CredentialContent, error) {
	var ptr *C.zktf_message_content_credential
	if err := status(C.zktf_message_content_as_credential(content.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newCredentialContent(ptr), nil
}

// VerifiablePresentations returns the presentations carried in the content.
func (c *CredentialContent) VerifiablePresentations() []*VerifiablePresentation {
	return verifiablePresentationsFrom(
		C.zktf_message_content_credential_verifiable_presentations(c.ptr),
	)
}

// VerifiableCredentials returns the credentials carried in the content.
func (c *CredentialContent) VerifiableCredentials() []*VerifiableCredential {
	return verifiableCredentialsFrom(
		C.zktf_message_content_credential_verifiable_credentials(c.ptr),
	)
}

// Assets returns supporting object assets carried in the content.
func (c *CredentialContent) Assets() []*Object {
	return objectsFrom(C.zktf_message_content_credential_assets(c.ptr))
}

// CredentialContentBuilder wraps a zktf_message_content_credential_builder.
type CredentialContentBuilder struct {
	ptr *C.zktf_message_content_credential_builder
}

// NewCredentialContentBuilder initializes a new credential-content builder.
func NewCredentialContentBuilder() *CredentialContentBuilder {
	ptr := C.zktf_message_content_credential_builder_init()
	b := &CredentialContentBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_credential_builder) {
		C.zktf_message_content_credential_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// VerifiablePresentation adds a presentation to the credential content.
func (b *CredentialContentBuilder) VerifiablePresentation(p *VerifiablePresentation) *CredentialContentBuilder {
	C.zktf_message_content_credential_builder_verifiable_presentation(b.ptr, p.ptr)
	return b
}

// VerifiableCredential adds a credential to the credential content.
func (b *CredentialContentBuilder) VerifiableCredential(c *VerifiableCredential) *CredentialContentBuilder {
	C.zktf_message_content_credential_builder_verifiable_credential(b.ptr, c.ptr)
	return b
}

// Asset attaches a supporting object asset.
func (b *CredentialContentBuilder) Asset(o *Object) *CredentialContentBuilder {
	C.zktf_message_content_credential_builder_asset(b.ptr, o.ptr)
	return b
}

// Finish finalizes the credential content, ready to send.
func (b *CredentialContentBuilder) Finish() (*Content, error) {
	var ptr *C.zktf_message_content
	if err := status(C.zktf_message_content_credential_builder_finish(b.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newContent(ptr), nil
}
