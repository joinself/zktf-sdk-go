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

// NewAnonymousMessage wraps content as an anonymous message, addressed to no
// particular inbox, ready to encode for out-of-band delivery (e.g. a QR code).
func NewAnonymousMessage(content *Content) *AnonymousMessage {
	defer runtime.KeepAlive(content)
	return newAnonymousMessage(C.zktf_anonymous_message_init(content.ptr))
}

// EncodeAsString encodes the message as a base64 URL encoded string.
func (m *AnonymousMessage) EncodeAsString() (string, error) {
	defer runtime.KeepAlive(m)
	var buf *C.zktf_string_buffer
	if err := status(C.zktf_anonymous_message_encode_as_string(m.ptr, &buf)); err != nil {
		return "", err
	}
	return goStringFromBuffer(buf), nil
}

type QREncoding int

const (
	QRSVG     QREncoding = C.QR_SVG
	QRUnicode QREncoding = C.QR_UNICODE
)

func (m *AnonymousMessage) EncodeAsQR(encoding QREncoding) ([]byte, error) {
	defer runtime.KeepAlive(m)
	var buf *C.zktf_bytes_buffer
	if err := status(C.zktf_anonymous_message_encode_as_qr(m.ptr, &buf, C.enum_zktf_qr_code_encoding(encoding))); err != nil {
		return nil, err
	}
	return goBytesFromBuffer(buf), nil
}

// ID returns the id of the message.
func (m *AnonymousMessage) ID() []byte {
	defer runtime.KeepAlive(m)
	return C.GoBytes(unsafe.Pointer(C.zktf_anonymous_message_id(m.ptr)), messageIDLen)
}

// Content returns the message content.
func (m *AnonymousMessage) Content() *Content {
	defer runtime.KeepAlive(m)
	return newContent(C.zktf_anonymous_message_message_content(m.ptr))
}

type MessageFlag uint64

const (
	MessageFlagNone          MessageFlag = C.MESSAGE_FLAG_NONE
	MessageFlagTargetSandbox MessageFlag = C.MESSAGE_FLAG_TARGET_SANDBOX
)

func AnonymousMessageDecode(data []byte) (*AnonymousMessage, error) {
	buf, length := cbytes(data)
	defer free(unsafe.Pointer(buf))

	var ptr *C.zktf_anonymous_message
	if err := status(C.zktf_anonymous_message_decode(&ptr, buf, length)); err != nil {
		return nil, err
	}
	return newAnonymousMessage(ptr), nil
}

func (m *AnonymousMessage) Encode() ([]byte, error) {
	defer runtime.KeepAlive(m)
	var buf *C.zktf_bytes_buffer
	if err := status(C.zktf_anonymous_message_encode(m.ptr, &buf)); err != nil {
		return nil, err
	}
	return goBytesFromBuffer(buf), nil
}

func (m *AnonymousMessage) Flags() MessageFlag {
	defer runtime.KeepAlive(m)
	return MessageFlag(C.zktf_anonymous_message_flags(m.ptr))
}

func (m *AnonymousMessage) SetFlags(flags MessageFlag) {
	defer runtime.KeepAlive(m)
	C.zktf_anonymous_message_set_flags(m.ptr, C.zktf_message_flag(flags))
}

func (m *AnonymousMessage) HasFlags(flags MessageFlag) bool {
	defer runtime.KeepAlive(m)
	return bool(C.zktf_anonymous_message_has_flags(m.ptr, C.zktf_message_flag(flags)))
}
