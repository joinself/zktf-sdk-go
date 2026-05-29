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

// ContentType mirrors zktf_message_content_type.
type ContentType uint32

const (
	ContentUnknown           ContentType = C.CONTENT_UNKNOWN
	ContentCustom            ContentType = C.CONTENT_CUSTOM
	ContentChat              ContentType = C.CONTENT_CHAT
	ContentReceipt           ContentType = C.CONTENT_RECEIPT
	ContentCredential        ContentType = C.CONTENT_CREDENTIAL
	ContentIntroduction      ContentType = C.CONTENT_INTRODUCTION
	ContentDiscoveryRequest  ContentType = C.CONTENT_DISCOVERY_REQUEST
	ContentDiscoveryResponse ContentType = C.CONTENT_DISCOVERY_RESPONSE
	ContentExchangeRequest   ContentType = C.CONTENT_EXCHANGE_REQUEST
	ContentExchangeResponse  ContentType = C.CONTENT_EXCHANGE_RESPONSE
)

// messageIDLen is the length of a message/content id.
const messageIDLen = 16

// Content wraps a zktf_message_content handle — the decoded payload of a message
// or the output of a content builder, ready to send.
type Content struct {
	ptr *C.zktf_message_content
}

func newContent(ptr *C.zktf_message_content) *Content {
	if ptr == nil {
		return nil
	}
	c := &Content{ptr: ptr}
	runtime.AddCleanup(c, func(ptr *C.zktf_message_content) {
		C.zktf_message_content_destroy(ptr)
	}, c.ptr)
	return c
}

// TypeOf returns the type of content.
func (c *Content) TypeOf() ContentType {
	return ContentType(C.zktf_message_content_type_of(c.ptr))
}

// ID returns the content's id.
func (c *Content) ID() []byte {
	return C.GoBytes(unsafe.Pointer(C.zktf_message_content_id(c.ptr)), messageIDLen)
}

// Message wraps a zktf_message handle delivered to the on_message callback.
type Message struct {
	ptr *C.zktf_message
}

func newMessage(ptr *C.zktf_message) *Message {
	if ptr == nil {
		return nil
	}
	m := &Message{ptr: ptr}
	runtime.AddCleanup(m, func(ptr *C.zktf_message) {
		C.zktf_message_destroy(ptr)
	}, m.ptr)
	return m
}

// ID returns the message id.
func (m *Message) ID() []byte {
	return C.GoBytes(unsafe.Pointer(C.zktf_message_id(m.ptr)), messageIDLen)
}

// FromAddress returns the sender's address.
func (m *Message) FromAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_message_from_address(m.ptr))
}

// ToAddress returns the recipient's address.
func (m *Message) ToAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_message_to_address(m.ptr))
}

// Timestamp returns the unix timestamp (seconds) of the message.
func (m *Message) Timestamp() int64 {
	return int64(C.zktf_message_timestamp(m.ptr))
}

// Content decodes and returns the message content. The caller owns the result.
func (m *Message) Content() *Content {
	return newContent(C.zktf_message_message_content(m.ptr))
}

// Metadata returns the opaque metadata payload attached to the message, if any.
// The payload is internal to the network and is not interpreted by this SDK; the
// boolean is false when no metadata is present.
func (m *Message) Metadata() ([]byte, bool) {
	b := goBytesFromBuffer(C.zktf_message_message_metadata(m.ptr))
	return b, b != nil
}
