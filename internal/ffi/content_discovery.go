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
	defer runtime.KeepAlive(content)
	var ptr *C.zktf_message_content_discovery_request
	if err := status(C.zktf_message_content_as_discovery_request(content.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newDiscoveryRequest(ptr), nil
}

// DocumentAddress returns the discovery requester's document address, or nil.
func (r *DiscoveryRequest) DocumentAddress() *SigningPublicKey {
	defer runtime.KeepAlive(r)
	return newSigningPublicKey(C.zktf_message_content_discovery_request_document_address(r.ptr))
}

// FromAddress returns the inbox address the discovery request was issued by, or nil.
func (r *DiscoveryRequest) FromAddress() *SigningPublicKey {
	defer runtime.KeepAlive(r)
	return newSigningPublicKey(C.zktf_message_content_discovery_request_from_address(r.ptr))
}

// KeyPackage returns the key package used to establish an inbound session, or nil.
func (r *DiscoveryRequest) KeyPackage() *CryptoKeyPackage {
	defer runtime.KeepAlive(r)
	return newCryptoKeyPackage(C.zktf_message_content_discovery_request_key_package(r.ptr))
}

// Expires returns the unix timestamp (seconds) the request expires.
func (r *DiscoveryRequest) Expires() int64 {
	defer runtime.KeepAlive(r)
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
	defer runtime.KeepAlive(b)
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_discovery_request_builder) {
		C.zktf_message_content_discovery_request_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// DocumentAddress sets an optional document address.
func (b *DiscoveryRequestBuilder) DocumentAddress(address *SigningPublicKey) *DiscoveryRequestBuilder {
	defer runtime.KeepAlive(b)
	defer runtime.KeepAlive(address)
	C.zktf_message_content_discovery_request_builder_document_address(b.ptr, address.ptr)
	return b
}

// FromAddress sets the inbox address the discovery request is issued by.
func (b *DiscoveryRequestBuilder) FromAddress(address *SigningPublicKey) *DiscoveryRequestBuilder {
	defer runtime.KeepAlive(b)
	defer runtime.KeepAlive(address)
	C.zktf_message_content_discovery_request_builder_from_address(b.ptr, address.ptr)
	return b
}

// KeyPackage attaches a key package the receiver can use to establish a session.
func (b *DiscoveryRequestBuilder) KeyPackage(kp *CryptoKeyPackage) *DiscoveryRequestBuilder {
	defer runtime.KeepAlive(b)
	defer runtime.KeepAlive(kp)
	C.zktf_message_content_discovery_request_builder_key_package(b.ptr, kp.ptr)
	return b
}

// Expires sets the unix timestamp (seconds) the request expires.
func (b *DiscoveryRequestBuilder) Expires(unix int64) *DiscoveryRequestBuilder {
	defer runtime.KeepAlive(b)
	C.zktf_message_content_discovery_request_builder_expires(b.ptr, C.int64_t(unix))
	return b
}

// Finish finalizes the discovery request, ready to send.
func (b *DiscoveryRequestBuilder) Finish() (*Content, error) {
	defer runtime.KeepAlive(b)
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
	defer runtime.KeepAlive(content)
	var ptr *C.zktf_message_content_discovery_response
	if err := status(C.zktf_message_content_as_discovery_response(content.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newDiscoveryResponse(ptr), nil
}

// ResponseTo returns the id of the request being responded to.
func (r *DiscoveryResponse) ResponseTo() []byte {
	defer runtime.KeepAlive(r)
	return C.GoBytes(
		unsafe.Pointer(C.zktf_message_content_discovery_response_response_to(r.ptr)),
		messageIDLen,
	)
}

// Status returns the response status.
func (r *DiscoveryResponse) Status() ResponseStatus {
	defer runtime.KeepAlive(r)
	return ResponseStatus(C.zktf_message_content_discovery_response_status(r.ptr))
}

// ErrorMessage returns the response error message, or "".
func (r *DiscoveryResponse) ErrorMessage() string {
	defer runtime.KeepAlive(r)
	return C.GoString(C.zktf_message_content_discovery_response_error_message(r.ptr))
}

// Tokens returns the tokens issued to the recipient.
func (r *DiscoveryResponse) Tokens() ([]*Token, error) {
	defer runtime.KeepAlive(r)
	var c *C.zktf_collection_token
	if err := status(C.zktf_message_content_discovery_response_tokens(r.ptr, &c)); err != nil {
		return nil, err
	}
	return tokensFrom(c), nil
}

// DiscoveryResponseBuilder builds a discovery response.
type DiscoveryResponseBuilder struct {
	ptr *C.zktf_message_content_discovery_response_builder
}

// NewDiscoveryResponseBuilder initializes a discovery response builder.
func NewDiscoveryResponseBuilder() *DiscoveryResponseBuilder {
	ptr := C.zktf_message_content_discovery_response_builder_init()
	b := &DiscoveryResponseBuilder{ptr: ptr}
	defer runtime.KeepAlive(b)
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_discovery_response_builder) {
		C.zktf_message_content_discovery_response_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// ResponseTo sets the id of the request being responded to.
func (b *DiscoveryResponseBuilder) ResponseTo(requestID []byte) *DiscoveryResponseBuilder {
	defer runtime.KeepAlive(b)
	buf, _ := cbytes(requestID)
	defer free(unsafe.Pointer(buf))
	C.zktf_message_content_discovery_response_builder_response_to(b.ptr, buf)
	return b
}

// Status sets the response status.
func (b *DiscoveryResponseBuilder) Status(s ResponseStatus) *DiscoveryResponseBuilder {
	defer runtime.KeepAlive(b)
	C.zktf_message_content_discovery_response_builder_status(b.ptr, C.enum_zktf_message_response_status(s))
	return b
}

// ErrorMessage sets the response error message.
func (b *DiscoveryResponseBuilder) ErrorMessage(msg string) *DiscoveryResponseBuilder {
	defer runtime.KeepAlive(b)
	cmsg := cstring(msg)
	defer free(unsafe.Pointer(cmsg))
	C.zktf_message_content_discovery_response_builder_error_message(b.ptr, cmsg)
	return b
}

// Token attaches a token for the recipient.
func (b *DiscoveryResponseBuilder) Token(t *Token) *DiscoveryResponseBuilder {
	defer runtime.KeepAlive(b)
	defer runtime.KeepAlive(t)
	C.zktf_message_content_discovery_response_builder_token(b.ptr, t.ptr)
	return b
}

// Finish finalizes the discovery response, ready to send.
func (b *DiscoveryResponseBuilder) Finish() (*Content, error) {
	defer runtime.KeepAlive(b)
	var out *C.zktf_message_content
	if err := status(C.zktf_message_content_discovery_response_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newContent(out), nil
}
