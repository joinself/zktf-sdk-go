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

// Evidence returns the objects attached to the action as evidence.
func (a *VerificationAction) Evidence() []*VerificationEvidence {
	return verificationEvidenceFrom(
		C.zktf_message_content_credential_verification_action_evidence(a.ptr),
	)
}

// Parameters returns the typed parameters attached to the action.
func (a *VerificationAction) Parameters() []*VerificationParameter {
	return verificationParametersFrom(
		C.zktf_message_content_credential_verification_action_parameters(a.ptr),
	)
}

// AsAction wraps this verification action into a generic Action (consuming it).
func (a *VerificationAction) AsAction() *Action {
	return newAction(C.zktf_message_content_action_verification(a.ptr))
}

// VerificationEvidence is an object attached to a verification action as
// supporting evidence, tagged with an evidence type.
type VerificationEvidence struct {
	ptr *C.zktf_credential_verification_evidence
}

func newVerificationEvidence(ptr *C.zktf_credential_verification_evidence) *VerificationEvidence {
	if ptr == nil {
		return nil
	}
	e := &VerificationEvidence{ptr: ptr}
	runtime.AddCleanup(e, func(ptr *C.zktf_credential_verification_evidence) {
		C.zktf_credential_verification_evidence_destroy(ptr)
	}, e.ptr)
	return e
}

// EvidenceType returns the evidence type tag.
func (e *VerificationEvidence) EvidenceType() string {
	return C.GoString(C.zktf_credential_verification_evidence_evidence_type(e.ptr))
}

// Object returns the object forming the evidence.
func (e *VerificationEvidence) Object() *Object {
	return newObject(C.zktf_credential_verification_evidence_object(e.ptr))
}

// VerificationParameter is a typed key/value parameter attached to a
// verification action.
type VerificationParameter struct {
	ptr *C.zktf_credential_verification_parameter
}

func newVerificationParameter(ptr *C.zktf_credential_verification_parameter) *VerificationParameter {
	if ptr == nil {
		return nil
	}
	p := &VerificationParameter{ptr: ptr}
	runtime.AddCleanup(p, func(ptr *C.zktf_credential_verification_parameter) {
		C.zktf_credential_verification_parameter_destroy(ptr)
	}, p.ptr)
	return p
}

// Key returns the parameter key.
func (p *VerificationParameter) Key() string {
	return C.GoString(C.zktf_credential_verification_parameter_parameter_key(p.ptr))
}

// Value decodes the parameter value into a native Go type. See
// ParameterValue.Value for the supported types.
func (p *VerificationParameter) Value() any {
	pv := newParameterValue(C.zktf_credential_verification_parameter_value(p.ptr))
	if pv == nil {
		return nil
	}
	return pv.Value()
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

// Evidence attaches an object as supporting evidence under a named type.
func (b *VerificationActionBuilder) Evidence(evidenceType string, object *Object) *VerificationActionBuilder {
	ct := cstring(evidenceType)
	C.zktf_message_content_credential_verification_action_builder_evidence(b.ptr, ct, object.ptr)
	free(unsafe.Pointer(ct))
	return b
}

// Parameter attaches a typed key/value parameter.
func (b *VerificationActionBuilder) Parameter(key string, value *ParameterValue) *VerificationActionBuilder {
	ck := cstring(key)
	C.zktf_message_content_credential_verification_action_builder_parameter(b.ptr, ck, value.ptr)
	free(unsafe.Pointer(ck))
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
