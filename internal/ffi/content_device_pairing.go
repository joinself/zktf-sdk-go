package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import "runtime"

// DevicePairingAction is a request to pair another device (signing key) into
// an identity document with a given role bitmask.
type DevicePairingAction struct {
	ptr *C.zktf_message_content_device_pairing_action
}

func newDevicePairingAction(ptr *C.zktf_message_content_device_pairing_action) *DevicePairingAction {
	if ptr == nil {
		return nil
	}
	a := &DevicePairingAction{ptr: ptr}
	runtime.AddCleanup(a, func(ptr *C.zktf_message_content_device_pairing_action) {
		C.zktf_message_content_device_pairing_action_destroy(ptr)
	}, a.ptr)
	return a
}

// Address returns the signing address to pair.
func (a *DevicePairingAction) Address() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_message_content_device_pairing_action_address(a.ptr))
}

// Roles returns the requested role bitmask for the paired key.
func (a *DevicePairingAction) Roles() uint64 {
	return uint64(C.zktf_message_content_device_pairing_action_roles(a.ptr))
}

// AsAction wraps this device-pairing action into a generic Action.
func (a *DevicePairingAction) AsAction() *Action {
	return newAction(C.zktf_message_content_action_pairing(a.ptr))
}

// DevicePairingActionBuilder builds a device-pairing action.
type DevicePairingActionBuilder struct {
	ptr *C.zktf_message_content_device_pairing_action_builder
}

// NewDevicePairingActionBuilder initializes the builder.
func NewDevicePairingActionBuilder() *DevicePairingActionBuilder {
	ptr := C.zktf_message_content_device_pairing_action_builder_init()
	b := &DevicePairingActionBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_device_pairing_action_builder) {
		C.zktf_message_content_device_pairing_action_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// Address sets the signing address to pair.
func (b *DevicePairingActionBuilder) Address(address *SigningPublicKey) *DevicePairingActionBuilder {
	C.zktf_message_content_device_pairing_action_builder_address(b.ptr, address.ptr)
	return b
}

// Roles sets the requested role bitmask.
func (b *DevicePairingActionBuilder) Roles(roles uint64) *DevicePairingActionBuilder {
	C.zktf_message_content_device_pairing_action_builder_roles(b.ptr, C.uint64_t(roles))
	return b
}

// Finish finalizes the action.
func (b *DevicePairingActionBuilder) Finish() (*DevicePairingAction, error) {
	var out *C.zktf_message_content_device_pairing_action
	if err := status(C.zktf_message_content_device_pairing_action_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newDevicePairingAction(out), nil
}

// DevicePairingResult is the response to a device-pairing request.
type DevicePairingResult struct {
	ptr *C.zktf_message_content_device_pairing_result
}

func newDevicePairingResult(ptr *C.zktf_message_content_device_pairing_result) *DevicePairingResult {
	if ptr == nil {
		return nil
	}
	r := &DevicePairingResult{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_message_content_device_pairing_result) {
		C.zktf_message_content_device_pairing_result_destroy(ptr)
	}, r.ptr)
	return r
}

// DocumentAddress returns the document address the result is for.
func (r *DevicePairingResult) DocumentAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_message_content_device_pairing_result_document_address(r.ptr))
}

// Operation returns the signed operation that paired the device.
func (r *DevicePairingResult) Operation() *IdentityOperation {
	return newIdentityOperation(C.zktf_message_content_device_pairing_result_operation(r.ptr))
}

// Presentations returns the presentations attached to the result.
func (r *DevicePairingResult) Presentations() []*VerifiablePresentation {
	return verifiablePresentationsFrom(C.zktf_message_content_device_pairing_result_presentations(r.ptr))
}

// Assets returns the assets attached to the result.
func (r *DevicePairingResult) Assets() []*Object {
	return objectsFrom(C.zktf_message_content_device_pairing_result_assets(r.ptr))
}

// DevicePairingResultBuilder builds a device-pairing result.
type DevicePairingResultBuilder struct {
	ptr *C.zktf_message_content_device_pairing_result_builder
}

// NewDevicePairingResultBuilder initializes the builder.
func NewDevicePairingResultBuilder() *DevicePairingResultBuilder {
	ptr := C.zktf_message_content_device_pairing_result_builder_init()
	b := &DevicePairingResultBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_device_pairing_result_builder) {
		C.zktf_message_content_device_pairing_result_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// DocumentAddress sets the document address the result is for.
func (b *DevicePairingResultBuilder) DocumentAddress(address *SigningPublicKey) *DevicePairingResultBuilder {
	C.zktf_message_content_device_pairing_result_builder_document_address(b.ptr, address.ptr)
	return b
}

// Operation sets the signed operation.
func (b *DevicePairingResultBuilder) Operation(operation *IdentityOperation) *DevicePairingResultBuilder {
	C.zktf_message_content_device_pairing_result_builder_operation(b.ptr, operation.ptr)
	return b
}

// Presentation adds a presentation to the result.
func (b *DevicePairingResultBuilder) Presentation(p *VerifiablePresentation) *DevicePairingResultBuilder {
	C.zktf_message_content_device_pairing_result_builder_presentation(b.ptr, p.ptr)
	return b
}

// Asset attaches a supporting object asset.
func (b *DevicePairingResultBuilder) Asset(o *Object) *DevicePairingResultBuilder {
	C.zktf_message_content_device_pairing_result_builder_asset(b.ptr, o.ptr)
	return b
}

// Finish finalizes the result.
func (b *DevicePairingResultBuilder) Finish() (*DevicePairingResult, error) {
	var out *C.zktf_message_content_device_pairing_result
	if err := status(C.zktf_message_content_device_pairing_result_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newDevicePairingResult(out), nil
}
