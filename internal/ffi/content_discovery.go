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

// DiscoveryRequest wraps a zktf_message_content_discovery_request handle (a QR
// onboarding / out-of-band discovery request).
type DiscoveryRequest struct {
	ptr *C.zktf_message_content_discovery_request
}

func newDiscoveryRequest(ptr *C.zktf_message_content_discovery_request) *DiscoveryRequest {
	if ptr == nil {
		return nil
	}
	r := &DiscoveryRequest{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_message_content_discovery_request) {
		C.zktf_message_content_discovery_request_destroy(ptr)
	}, r.ptr)
	return r
}

// DiscoveryRequestFromContent decodes message content as a discovery request.
func DiscoveryRequestFromContent(content *Content) (*DiscoveryRequest, error) {
	var ptr *C.zktf_message_content_discovery_request
	if err := status(C.zktf_message_content_as_discovery_request(content.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newDiscoveryRequest(ptr), nil
}

// DocumentAddress returns the discovery requester's document address, or nil.
func (r *DiscoveryRequest) DocumentAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_message_content_discovery_request_document_address(r.ptr))
}

// FromAddress returns the inbox address the discovery request was issued by, or nil.
func (r *DiscoveryRequest) FromAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_message_content_discovery_request_from_address(r.ptr))
}

// KeyPackage returns the key package used to establish an inbound session, or nil.
func (r *DiscoveryRequest) KeyPackage() *CryptoKeyPackage {
	return newCryptoKeyPackage(C.zktf_message_content_discovery_request_key_package(r.ptr))
}

// Expires returns the unix timestamp (seconds) the request expires.
func (r *DiscoveryRequest) Expires() int64 {
	return int64(C.zktf_message_content_discovery_request_expires(r.ptr))
}

// DiscoveryRequestBuilder builds a discovery request.
type DiscoveryRequestBuilder struct {
	ptr *C.zktf_message_content_discovery_request_builder
}

// NewDiscoveryRequestBuilder initializes a discovery request builder.
func NewDiscoveryRequestBuilder() *DiscoveryRequestBuilder {
	ptr := C.zktf_message_content_discovery_request_builder_init()
	b := &DiscoveryRequestBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_discovery_request_builder) {
		C.zktf_message_content_discovery_request_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// DocumentAddress sets an optional document address.
func (b *DiscoveryRequestBuilder) DocumentAddress(address *SigningPublicKey) *DiscoveryRequestBuilder {
	C.zktf_message_content_discovery_request_builder_document_address(b.ptr, address.ptr)
	return b
}

// FromAddress sets the inbox address the discovery request is issued by.
func (b *DiscoveryRequestBuilder) FromAddress(address *SigningPublicKey) *DiscoveryRequestBuilder {
	C.zktf_message_content_discovery_request_builder_from_address(b.ptr, address.ptr)
	return b
}

// KeyPackage attaches a key package the receiver can use to establish a session.
func (b *DiscoveryRequestBuilder) KeyPackage(kp *CryptoKeyPackage) *DiscoveryRequestBuilder {
	C.zktf_message_content_discovery_request_builder_key_package(b.ptr, kp.ptr)
	return b
}

// Expires sets the unix timestamp (seconds) the request expires.
func (b *DiscoveryRequestBuilder) Expires(unix int64) *DiscoveryRequestBuilder {
	C.zktf_message_content_discovery_request_builder_expires(b.ptr, C.int64_t(unix))
	return b
}

// Finish finalizes the discovery request, ready to send.
func (b *DiscoveryRequestBuilder) Finish() (*Content, error) {
	var out *C.zktf_message_content
	if err := status(C.zktf_message_content_discovery_request_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newContent(out), nil
}

// DiscoveryResponse wraps a zktf_message_content_discovery_response handle.
type DiscoveryResponse struct {
	ptr *C.zktf_message_content_discovery_response
}

func newDiscoveryResponse(ptr *C.zktf_message_content_discovery_response) *DiscoveryResponse {
	if ptr == nil {
		return nil
	}
	r := &DiscoveryResponse{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_message_content_discovery_response) {
		C.zktf_message_content_discovery_response_destroy(ptr)
	}, r.ptr)
	return r
}

// DiscoveryResponseFromContent decodes message content as a discovery response.
func DiscoveryResponseFromContent(content *Content) (*DiscoveryResponse, error) {
	var ptr *C.zktf_message_content_discovery_response
	if err := status(C.zktf_message_content_as_discovery_response(content.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newDiscoveryResponse(ptr), nil
}

// ResponseTo returns the id of the request being responded to.
func (r *DiscoveryResponse) ResponseTo() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.zktf_message_content_discovery_response_response_to(r.ptr)),
		messageIDLen,
	)
}

// Status returns the response status.
func (r *DiscoveryResponse) Status() ResponseStatus {
	return ResponseStatus(C.zktf_message_content_discovery_response_status(r.ptr))
}

// ErrorMessage returns the response error message, or "".
func (r *DiscoveryResponse) ErrorMessage() string {
	return C.GoString(C.zktf_message_content_discovery_response_error_message(r.ptr))
}

// DiscoveryResponseBuilder builds a discovery response.
type DiscoveryResponseBuilder struct {
	ptr *C.zktf_message_content_discovery_response_builder
}

// NewDiscoveryResponseBuilder initializes a discovery response builder.
func NewDiscoveryResponseBuilder() *DiscoveryResponseBuilder {
	ptr := C.zktf_message_content_discovery_response_builder_init()
	b := &DiscoveryResponseBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_discovery_response_builder) {
		C.zktf_message_content_discovery_response_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// ResponseTo sets the id of the request being responded to.
func (b *DiscoveryResponseBuilder) ResponseTo(requestID []byte) *DiscoveryResponseBuilder {
	buf, _ := cbytes(requestID)
	defer free(unsafe.Pointer(buf))
	C.zktf_message_content_discovery_response_builder_response_to(b.ptr, buf)
	return b
}

// Status sets the response status.
func (b *DiscoveryResponseBuilder) Status(s ResponseStatus) *DiscoveryResponseBuilder {
	C.zktf_message_content_discovery_response_builder_status(b.ptr, C.enum_zktf_message_response_status(s))
	return b
}

// ErrorMessage sets the response error message.
func (b *DiscoveryResponseBuilder) ErrorMessage(msg string) *DiscoveryResponseBuilder {
	cmsg := cstring(msg)
	defer free(unsafe.Pointer(cmsg))
	C.zktf_message_content_discovery_response_builder_error_message(b.ptr, cmsg)
	return b
}

// Finish finalizes the discovery response, ready to send.
func (b *DiscoveryResponseBuilder) Finish() (*Content, error) {
	var out *C.zktf_message_content
	if err := status(C.zktf_message_content_discovery_response_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newContent(out), nil
}
