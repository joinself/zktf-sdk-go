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

// AnonymousMessage wraps a zktf_anonymous_message handle. An anonymous message
// is an unencrypted, self-describing envelope (e.g. a pairing/QR code) that
// carries message content addressed to no particular inbox.
type AnonymousMessage struct {
	ptr *C.zktf_anonymous_message
}

func newAnonymousMessage(ptr *C.zktf_anonymous_message) *AnonymousMessage {
	if ptr == nil {
		return nil
	}
	m := &AnonymousMessage{ptr: ptr}
	runtime.AddCleanup(m, func(ptr *C.zktf_anonymous_message) {
		C.zktf_anonymous_message_destroy(ptr)
	}, m.ptr)
	return m
}

// AnonymousMessageDecodeFromString decodes a base64 URL encoded anonymous
// message (e.g. a pairing code).
func AnonymousMessageDecodeFromString(encoded string) (*AnonymousMessage, error) {
	cencoded := cstring(encoded)
	defer free(unsafe.Pointer(cencoded))

	var ptr *C.zktf_anonymous_message
	if err := status(C.zktf_anonymous_message_decode_from_string(&ptr, cencoded)); err != nil {
		return nil, err
	}
	return newAnonymousMessage(ptr), nil
}

// ID returns the 16 byte id of the message.
func (m *AnonymousMessage) ID() []byte {
	return C.GoBytes(unsafe.Pointer(C.zktf_anonymous_message_id(m.ptr)), messageIDLen)
}

// Content returns the message content.
func (m *AnonymousMessage) Content() *Content {
	return newContent(C.zktf_anonymous_message_message_content(m.ptr))
}
