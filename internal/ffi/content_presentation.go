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

// presentationChallengeLen is the byte length of a presentation challenge.
const presentationChallengeLen = 32

// PresentationAction is a credential-presentation request (verifier → holder).
type PresentationAction struct {
	ptr *C.zktf_message_content_credential_presentation_action
}

func newPresentationAction(ptr *C.zktf_message_content_credential_presentation_action) *PresentationAction {
	if ptr == nil {
		return nil
	}
	a := &PresentationAction{ptr: ptr}
	runtime.AddCleanup(a, func(ptr *C.zktf_message_content_credential_presentation_action) {
		C.zktf_message_content_credential_presentation_action_destroy(ptr)
	}, a.ptr)
	return a
}

// PresentationTypes returns the requested presentation types.
func (a *PresentationAction) PresentationTypes() []string {
	return presentationTypesFrom(
		C.zktf_message_content_credential_presentation_action_presentation_type(a.ptr),
	)
}

// Holder returns the expected holder address, or nil.
func (a *PresentationAction) Holder() (*DIDAddress, error) {
	var out *C.zktf_did_address
	if err := status(C.zktf_message_content_credential_presentation_action_holder(a.ptr, &out)); err != nil {
		return nil, err
	}
	return newDIDAddress(out), nil
}

// Challenge returns the random challenge bytes the verifier expects to be
// signed back, or nil.
func (a *PresentationAction) Challenge() []byte {
	p := C.zktf_message_content_credential_presentation_action_challenge(a.ptr)
	if p == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(p), presentationChallengeLen)
}

// Predicates returns the predicate tree describing the request's credential
// constraints, or nil.
func (a *PresentationAction) Predicates() *PredicateTree {
	return newPredicateTree(C.zktf_message_content_credential_presentation_action_predicates(a.ptr))
}

// Proof returns the verifiable presentations attached as proof.
func (a *PresentationAction) Proof() []*VerifiablePresentation {
	return verifiablePresentationsFrom(C.zktf_message_content_credential_presentation_action_proof(a.ptr))
}

// Term returns the term the requester would like to access credentials under,
// or nil.
func (a *PresentationAction) Term() *CredentialTerm {
	return newCredentialTerm(C.zktf_message_content_credential_presentation_action_term(a.ptr))
}

// BiometricAnchor returns the 20-byte biometric anchor hash, or nil.
func (a *PresentationAction) BiometricAnchor() []byte {
	p := C.zktf_message_content_credential_presentation_action_biometric_anchor(a.ptr)
	if p == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(p), biometricAnchorHashLen)
}

// AsAction wraps this presentation action into a generic Action (consuming it).
func (a *PresentationAction) AsAction() *Action {
	return newAction(C.zktf_message_content_action_presentation(a.ptr))
}

// PresentationActionBuilder builds a credential presentation action.
type PresentationActionBuilder struct {
	ptr *C.zktf_message_content_credential_presentation_action_builder
}

// NewPresentationActionBuilder initializes a presentation action builder.
func NewPresentationActionBuilder() *PresentationActionBuilder {
	ptr := C.zktf_message_content_credential_presentation_action_builder_init()
	b := &PresentationActionBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_credential_presentation_action_builder) {
		C.zktf_message_content_credential_presentation_action_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// PresentationType sets the requested presentation types.
func (b *PresentationActionBuilder) PresentationType(types *PresentationTypeCollection) *PresentationActionBuilder {
	C.zktf_message_content_credential_presentation_action_builder_presentation_type(b.ptr, types.ptr)
	return b
}

// Holder sets the expected holder address.
func (b *PresentationActionBuilder) Holder(holder *DIDAddress) *PresentationActionBuilder {
	C.zktf_message_content_credential_presentation_action_builder_holder(b.ptr, holder.ptr)
	return b
}

// Challenge sets the random challenge the verifier expects to be signed back.
func (b *PresentationActionBuilder) Challenge(challenge []byte) *PresentationActionBuilder {
	buf, _ := cbytes(challenge)
	defer free(unsafe.Pointer(buf))
	C.zktf_message_content_credential_presentation_action_builder_challenge(b.ptr, buf)
	return b
}

// Predicates sets the predicate tree describing the request's credential constraints.
func (b *PresentationActionBuilder) Predicates(tree *PredicateTree) *PresentationActionBuilder {
	C.zktf_message_content_credential_presentation_action_builder_predicates(b.ptr, tree.ptr)
	return b
}

// Proof attaches a verifiable presentation as proof.
func (b *PresentationActionBuilder) Proof(p *VerifiablePresentation) *PresentationActionBuilder {
	C.zktf_message_content_credential_presentation_action_builder_proof(b.ptr, p.ptr)
	return b
}

// Term sets the term the requester would like to access the credentials under.
func (b *PresentationActionBuilder) Term(term *CredentialTerm) *PresentationActionBuilder {
	C.zktf_message_content_credential_presentation_action_builder_term(b.ptr, term.ptr)
	return b
}

// BiometricAnchor sets the biometric anchor hash.
func (b *PresentationActionBuilder) BiometricAnchor(anchor []byte) *PresentationActionBuilder {
	buf, _ := cbytes(anchor)
	defer free(unsafe.Pointer(buf))
	C.zktf_message_content_credential_presentation_action_builder_biometric_anchor(b.ptr, buf)
	return b
}

// Finish finalizes the presentation action.
func (b *PresentationActionBuilder) Finish() (*PresentationAction, error) {
	var out *C.zktf_message_content_credential_presentation_action
	if err := status(C.zktf_message_content_credential_presentation_action_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newPresentationAction(out), nil
}

// PresentationResult is the response to a credential-presentation request.
type PresentationResult struct {
	ptr *C.zktf_message_content_credential_presentation_result
}

func newPresentationResult(ptr *C.zktf_message_content_credential_presentation_result) *PresentationResult {
	if ptr == nil {
		return nil
	}
	r := &PresentationResult{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_message_content_credential_presentation_result) {
		C.zktf_message_content_credential_presentation_result_destroy(ptr)
	}, r.ptr)
	return r
}

// Presentations returns the verifiable presentations contained in the result.
func (r *PresentationResult) Presentations() []*VerifiablePresentation {
	return verifiablePresentationsFrom(
		C.zktf_message_content_credential_presentation_result_presentations(r.ptr),
	)
}

// PresentationResultBuilder builds a credential-presentation result.
type PresentationResultBuilder struct {
	ptr *C.zktf_message_content_credential_presentation_result_builder
}

// NewPresentationResultBuilder initializes a presentation result builder.
func NewPresentationResultBuilder() *PresentationResultBuilder {
	ptr := C.zktf_message_content_credential_presentation_result_builder_init()
	b := &PresentationResultBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_credential_presentation_result_builder) {
		C.zktf_message_content_credential_presentation_result_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// Presentation adds a verifiable presentation to the result.
func (b *PresentationResultBuilder) Presentation(p *VerifiablePresentation) *PresentationResultBuilder {
	C.zktf_message_content_credential_presentation_result_builder_presentation(b.ptr, p.ptr)
	return b
}

// Finish finalizes the presentation result.
func (b *PresentationResultBuilder) Finish() (*PresentationResult, error) {
	var out *C.zktf_message_content_credential_presentation_result
	if err := status(C.zktf_message_content_credential_presentation_result_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newPresentationResult(out), nil
}
