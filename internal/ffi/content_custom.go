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

// Custom wraps a zktf_message_content_custom handle.
type Custom struct {
	ptr *C.zktf_message_content_custom
}

func newCustom(ptr *C.zktf_message_content_custom) *Custom {
	if ptr == nil {
		return nil
	}
	c := &Custom{ptr: ptr}
	runtime.AddCleanup(c, func(ptr *C.zktf_message_content_custom) {
		C.zktf_message_content_custom_destroy(ptr)
	}, c.ptr)
	return c
}

// CustomFromContent decodes message content as a custom payload.
func CustomFromContent(content *Content) (*Custom, error) {
	var ptr *C.zktf_message_content_custom
	if err := status(C.zktf_message_content_as_custom(content.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newCustom(ptr), nil
}

// Payload returns the custom payload bytes.
func (c *Custom) Payload() []byte {
	return goBytesFromBuffer(C.zktf_message_content_custom_payload(c.ptr))
}

// CustomBuilder wraps a zktf_message_content_custom_builder handle.
type CustomBuilder struct {
	ptr *C.zktf_message_content_custom_builder
}

// NewCustomBuilder initializes a new custom content builder.
func NewCustomBuilder() *CustomBuilder {
	ptr := C.zktf_message_content_custom_builder_init()
	b := &CustomBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_custom_builder) {
		C.zktf_message_content_custom_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// Payload sets the custom payload bytes.
func (b *CustomBuilder) Payload(payload []byte) *CustomBuilder {
	buf, length := cbytes(payload)
	defer free(unsafe.Pointer(buf))
	C.zktf_message_content_custom_builder_payload(b.ptr, buf, length)
	return b
}

// Finish finalizes the custom content, ready to send.
func (b *CustomBuilder) Finish() (*Content, error) {
	var ptr *C.zktf_message_content
	if err := status(C.zktf_message_content_custom_builder_finish(b.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newContent(ptr), nil
}
