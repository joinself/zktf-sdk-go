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

// Chat wraps a zktf_message_content_chat handle.
type Chat struct {
	ptr *C.zktf_message_content_chat
}

func newChat(ptr *C.zktf_message_content_chat) *Chat {
	if ptr == nil {
		return nil
	}
	c := &Chat{ptr: ptr}
	runtime.AddCleanup(c, func(ptr *C.zktf_message_content_chat) {
		C.zktf_message_content_chat_destroy(ptr)
	}, c.ptr)
	return c
}

// ChatFromContent decodes message content as a chat message.
func ChatFromContent(content *Content) (*Chat, error) {
	var ptr *C.zktf_message_content_chat
	if err := status(C.zktf_message_content_as_chat(content.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newChat(ptr), nil
}

// Message returns the chat message text.
func (c *Chat) Message() string {
	return C.GoString(C.zktf_message_content_chat_message(c.ptr))
}

// Referencing returns the id of the message this chat references, or nil.
func (c *Chat) Referencing() []byte {
	ptr := C.zktf_message_content_chat_referencing(c.ptr)
	if ptr == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(ptr), messageIDLen)
}

// ChatBuilder wraps a zktf_message_content_chat_builder handle.
type ChatBuilder struct {
	ptr *C.zktf_message_content_chat_builder
}

// NewChatBuilder initializes a new chat message builder.
func NewChatBuilder() *ChatBuilder {
	ptr := C.zktf_message_content_chat_builder_init()
	b := &ChatBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_chat_builder) {
		C.zktf_message_content_chat_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// Message sets the chat message text.
func (b *ChatBuilder) Message(message string) *ChatBuilder {
	cmsg := cstring(message)
	defer free(unsafe.Pointer(cmsg))
	C.zktf_message_content_chat_builder_message(b.ptr, cmsg)
	return b
}

// Reference sets the id of a message this chat references.
func (b *ChatBuilder) Reference(messageID []byte) *ChatBuilder {
	buf, _ := cbytes(messageID)
	defer free(unsafe.Pointer(buf))
	C.zktf_message_content_chat_builder_reference(b.ptr, buf)
	return b
}

// Finish finalizes the chat content, ready to send.
func (b *ChatBuilder) Finish() (*Content, error) {
	var ptr *C.zktf_message_content
	if err := status(C.zktf_message_content_chat_builder_finish(b.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newContent(ptr), nil
}
