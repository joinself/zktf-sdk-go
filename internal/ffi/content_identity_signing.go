package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import "runtime"

// IdentitySigningAction is a request to sign an identity-document operation.
type IdentitySigningAction struct {
	ptr *C.zktf_message_content_identity_signing_action
}

func newIdentitySigningAction(ptr *C.zktf_message_content_identity_signing_action) *IdentitySigningAction {
	if ptr == nil {
		return nil
	}
	a := &IdentitySigningAction{ptr: ptr}
	runtime.AddCleanup(a, func(ptr *C.zktf_message_content_identity_signing_action) {
		C.zktf_message_content_identity_signing_action_destroy(ptr)
	}, a.ptr)
	return a
}

// DocumentAddress returns the document address the operation targets.
func (a *IdentitySigningAction) DocumentAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_message_content_identity_signing_action_document_address(a.ptr))
}

// Operation returns the hashgraph operation to sign.
func (a *IdentitySigningAction) Operation() *IdentityOperation {
	return newIdentityOperation(C.zktf_message_content_identity_signing_action_operation(a.ptr))
}

// AsAction wraps this identity-signing action into a generic Action.
func (a *IdentitySigningAction) AsAction() *Action {
	return newAction(C.zktf_message_content_action_signing(a.ptr))
}

// IdentitySigningActionBuilder builds an identity-signing action.
type IdentitySigningActionBuilder struct {
	ptr *C.zktf_message_content_identity_signing_action_builder
}

// NewIdentitySigningActionBuilder initializes the builder.
func NewIdentitySigningActionBuilder() *IdentitySigningActionBuilder {
	ptr := C.zktf_message_content_identity_signing_action_builder_init()
	b := &IdentitySigningActionBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_identity_signing_action_builder) {
		C.zktf_message_content_identity_signing_action_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// DocumentAddress sets the document address the operation targets.
func (b *IdentitySigningActionBuilder) DocumentAddress(address *SigningPublicKey) *IdentitySigningActionBuilder {
	C.zktf_message_content_identity_signing_action_builder_document_address(b.ptr, address.ptr)
	return b
}

// Operation sets the operation to sign.
func (b *IdentitySigningActionBuilder) Operation(operation *IdentityOperation) *IdentitySigningActionBuilder {
	C.zktf_message_content_identity_signing_action_builder_operation(b.ptr, operation.ptr)
	return b
}

// Finish finalizes the action.
func (b *IdentitySigningActionBuilder) Finish() (*IdentitySigningAction, error) {
	var out *C.zktf_message_content_identity_signing_action
	if err := status(C.zktf_message_content_identity_signing_action_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newIdentitySigningAction(out), nil
}

// IdentitySigningResult is the response to an identity-signing request.
type IdentitySigningResult struct {
	ptr *C.zktf_message_content_identity_signing_result
}

func newIdentitySigningResult(ptr *C.zktf_message_content_identity_signing_result) *IdentitySigningResult {
	if ptr == nil {
		return nil
	}
	r := &IdentitySigningResult{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_message_content_identity_signing_result) {
		C.zktf_message_content_identity_signing_result_destroy(ptr)
	}, r.ptr)
	return r
}

// DocumentAddress returns the document address the result is for.
func (r *IdentitySigningResult) DocumentAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_message_content_identity_signing_result_document_address(r.ptr))
}

// Operation returns the signed operation.
func (r *IdentitySigningResult) Operation() *IdentityOperation {
	return newIdentityOperation(C.zktf_message_content_identity_signing_result_operation(r.ptr))
}

// Presentations returns the presentations attached to the result.
func (r *IdentitySigningResult) Presentations() []*VerifiablePresentation {
	return verifiablePresentationsFrom(C.zktf_message_content_identity_signing_result_presentations(r.ptr))
}

// Assets returns the assets attached to the result.
func (r *IdentitySigningResult) Assets() []*Object {
	return objectsFrom(C.zktf_message_content_identity_signing_result_assets(r.ptr))
}

// IdentitySigningResultBuilder builds an identity-signing result.
type IdentitySigningResultBuilder struct {
	ptr *C.zktf_message_content_identity_signing_result_builder
}

// NewIdentitySigningResultBuilder initializes the builder.
func NewIdentitySigningResultBuilder() *IdentitySigningResultBuilder {
	ptr := C.zktf_message_content_identity_signing_result_builder_init()
	b := &IdentitySigningResultBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_identity_signing_result_builder) {
		C.zktf_message_content_identity_signing_result_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// DocumentAddress sets the document address the result is for.
func (b *IdentitySigningResultBuilder) DocumentAddress(address *SigningPublicKey) *IdentitySigningResultBuilder {
	C.zktf_message_content_identity_signing_result_builder_document_address(b.ptr, address.ptr)
	return b
}

// Operation sets the signed operation.
func (b *IdentitySigningResultBuilder) Operation(operation *IdentityOperation) *IdentitySigningResultBuilder {
	C.zktf_message_content_identity_signing_result_builder_operation(b.ptr, operation.ptr)
	return b
}

// Presentation adds a presentation to the result.
func (b *IdentitySigningResultBuilder) Presentation(p *VerifiablePresentation) *IdentitySigningResultBuilder {
	C.zktf_message_content_identity_signing_result_builder_presentation(b.ptr, p.ptr)
	return b
}

// Asset attaches a supporting object asset.
func (b *IdentitySigningResultBuilder) Asset(o *Object) *IdentitySigningResultBuilder {
	C.zktf_message_content_identity_signing_result_builder_asset(b.ptr, o.ptr)
	return b
}

// Finish finalizes the result.
func (b *IdentitySigningResultBuilder) Finish() (*IdentitySigningResult, error) {
	var out *C.zktf_message_content_identity_signing_result
	if err := status(C.zktf_message_content_identity_signing_result_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newIdentitySigningResult(out), nil
}
