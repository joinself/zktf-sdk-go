package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import "runtime"

// RevocationSigningAction is a request asking a device to co-sign a
// revocation statement.
type RevocationSigningAction struct {
	ptr *C.zktf_message_content_revocation_signing_action
}

func newRevocationSigningAction(ptr *C.zktf_message_content_revocation_signing_action) *RevocationSigningAction {
	if ptr == nil {
		return nil
	}
	a := &RevocationSigningAction{ptr: ptr}
	runtime.AddCleanup(a, func(ptr *C.zktf_message_content_revocation_signing_action) {
		C.zktf_message_content_revocation_signing_action_destroy(ptr)
	}, a.ptr)
	return a
}

// Statement returns the revocation statement this device is asked to co-sign.
func (a *RevocationSigningAction) Statement() *RevocationStatement {
	return newRevocationStatement(C.zktf_message_content_revocation_signing_action_statement(a.ptr))
}

// AsAction wraps this revocation-signing action into a generic Action.
func (a *RevocationSigningAction) AsAction() *Action {
	return newAction(C.zktf_message_content_action_revocation_signing(a.ptr))
}

// RevocationSigningActionBuilder builds a revocation-signing action.
type RevocationSigningActionBuilder struct {
	ptr *C.zktf_message_content_revocation_signing_action_builder
}

// NewRevocationSigningActionBuilder initializes the builder.
func NewRevocationSigningActionBuilder() *RevocationSigningActionBuilder {
	ptr := C.zktf_message_content_revocation_signing_action_builder_init()
	b := &RevocationSigningActionBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_revocation_signing_action_builder) {
		C.zktf_message_content_revocation_signing_action_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// Statement sets the revocation statement to be co-signed.
func (b *RevocationSigningActionBuilder) Statement(statement *RevocationStatement) *RevocationSigningActionBuilder {
	C.zktf_message_content_revocation_signing_action_builder_statement(b.ptr, statement.ptr)
	return b
}

// Finish finalizes the action.
func (b *RevocationSigningActionBuilder) Finish() (*RevocationSigningAction, error) {
	var out *C.zktf_message_content_revocation_signing_action
	if err := status(C.zktf_message_content_revocation_signing_action_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newRevocationSigningAction(out), nil
}

// RevocationSigningResult is the response to a revocation-signing request.
type RevocationSigningResult struct {
	ptr *C.zktf_message_content_revocation_signing_result
}

func newRevocationSigningResult(ptr *C.zktf_message_content_revocation_signing_result) *RevocationSigningResult {
	if ptr == nil {
		return nil
	}
	r := &RevocationSigningResult{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_message_content_revocation_signing_result) {
		C.zktf_message_content_revocation_signing_result_destroy(ptr)
	}, r.ptr)
	return r
}

// Statement returns the co-signed revocation statement.
func (r *RevocationSigningResult) Statement() *RevocationStatement {
	return newRevocationStatement(C.zktf_message_content_revocation_signing_result_statement(r.ptr))
}

// RevocationSigningResultBuilder builds a revocation-signing result.
type RevocationSigningResultBuilder struct {
	ptr *C.zktf_message_content_revocation_signing_result_builder
}

// NewRevocationSigningResultBuilder initializes the builder.
func NewRevocationSigningResultBuilder() *RevocationSigningResultBuilder {
	ptr := C.zktf_message_content_revocation_signing_result_builder_init()
	b := &RevocationSigningResultBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_revocation_signing_result_builder) {
		C.zktf_message_content_revocation_signing_result_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// Statement sets the co-signed revocation statement.
func (b *RevocationSigningResultBuilder) Statement(statement *RevocationStatement) *RevocationSigningResultBuilder {
	C.zktf_message_content_revocation_signing_result_builder_statement(b.ptr, statement.ptr)
	return b
}

// Finish finalizes the result.
func (b *RevocationSigningResultBuilder) Finish() (*RevocationSigningResult, error) {
	var out *C.zktf_message_content_revocation_signing_result
	if err := status(C.zktf_message_content_revocation_signing_result_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newRevocationSigningResult(out), nil
}
