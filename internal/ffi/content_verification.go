package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import "runtime"

// VerificationAction is a credential-verification request (issuer/verifier
// asking for proof a credential should be (re-)issued or verified).
type VerificationAction struct {
	ptr *C.zktf_message_content_credential_verification_action
}

func newVerificationAction(ptr *C.zktf_message_content_credential_verification_action) *VerificationAction {
	if ptr == nil {
		return nil
	}
	a := &VerificationAction{ptr: ptr}
	runtime.AddCleanup(a, func(ptr *C.zktf_message_content_credential_verification_action) {
		C.zktf_message_content_credential_verification_action_destroy(ptr)
	}, a.ptr)
	return a
}

// CredentialTypes returns the requested credential types.
func (a *VerificationAction) CredentialTypes() []string {
	c := C.zktf_message_content_credential_verification_action_credential_type(a.ptr)
	if c == nil {
		return nil
	}
	wrapped := newCredentialTypeCollection(c)
	return wrapped.Strings()
}

// Proof returns the verifiable presentations attached as proof.
func (a *VerificationAction) Proof() []*VerifiablePresentation {
	return verifiablePresentationsFrom(
		C.zktf_message_content_credential_verification_action_proof(a.ptr),
	)
}

// AsAction wraps this verification action into a generic Action (consuming it).
func (a *VerificationAction) AsAction() *Action {
	return newAction(C.zktf_message_content_action_verification(a.ptr))
}

// VerificationActionBuilder builds a credential-verification action.
type VerificationActionBuilder struct {
	ptr *C.zktf_message_content_credential_verification_action_builder
}

// NewVerificationActionBuilder initializes a verification action builder.
func NewVerificationActionBuilder() *VerificationActionBuilder {
	ptr := C.zktf_message_content_credential_verification_action_builder_init()
	b := &VerificationActionBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_credential_verification_action_builder) {
		C.zktf_message_content_credential_verification_action_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// CredentialType sets the requested credential types.
func (b *VerificationActionBuilder) CredentialType(types *CredentialTypeCollection) *VerificationActionBuilder {
	C.zktf_message_content_credential_verification_action_builder_credential_type(b.ptr, types.ptr)
	return b
}

// Proof attaches a verifiable presentation as proof.
func (b *VerificationActionBuilder) Proof(p *VerifiablePresentation) *VerificationActionBuilder {
	C.zktf_message_content_credential_verification_action_builder_proof(b.ptr, p.ptr)
	return b
}

// Finish finalizes the verification action.
func (b *VerificationActionBuilder) Finish() (*VerificationAction, error) {
	var out *C.zktf_message_content_credential_verification_action
	if err := status(C.zktf_message_content_credential_verification_action_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newVerificationAction(out), nil
}

// VerificationResult is the response to a credential-verification request.
type VerificationResult struct {
	ptr *C.zktf_message_content_credential_verification_result
}

func newVerificationResult(ptr *C.zktf_message_content_credential_verification_result) *VerificationResult {
	if ptr == nil {
		return nil
	}
	r := &VerificationResult{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_message_content_credential_verification_result) {
		C.zktf_message_content_credential_verification_result_destroy(ptr)
	}, r.ptr)
	return r
}

// Credentials returns the verifiable credentials carried in the result.
func (r *VerificationResult) Credentials() []*VerifiableCredential {
	return verifiableCredentialsFrom(
		C.zktf_message_content_credential_verification_result_credentials(r.ptr),
	)
}

// VerificationResultBuilder builds a credential-verification result.
type VerificationResultBuilder struct {
	ptr *C.zktf_message_content_credential_verification_result_builder
}

// NewVerificationResultBuilder initializes a verification result builder.
func NewVerificationResultBuilder() *VerificationResultBuilder {
	ptr := C.zktf_message_content_credential_verification_result_builder_init()
	b := &VerificationResultBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_credential_verification_result_builder) {
		C.zktf_message_content_credential_verification_result_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// Credential adds a verifiable credential to the result.
func (b *VerificationResultBuilder) Credential(c *VerifiableCredential) *VerificationResultBuilder {
	C.zktf_message_content_credential_verification_result_builder_credential(b.ptr, c.ptr)
	return b
}

// Finish finalizes the verification result.
func (b *VerificationResultBuilder) Finish() (*VerificationResult, error) {
	var out *C.zktf_message_content_credential_verification_result
	if err := status(C.zktf_message_content_credential_verification_result_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newVerificationResult(out), nil
}
