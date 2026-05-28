package message

import "github.com/joinself/zktf-sdk-go/internal/ffi"

// Chat is a decoded chat message.
type Chat struct {
	h *ffi.Chat
}

// ChatBuilder builds chat message content.
type ChatBuilder struct {
	h *ffi.ChatBuilder
}

// ChatDecode decodes message content as a chat message.
func ChatDecode(content *Content) (*Chat, error) {
	c, err := ffi.ChatFromContent(content.h)
	if err != nil {
		return nil, err
	}

	return &Chat{h: c}, nil
}

// Message returns the chat message text.
func (c *Chat) Message() string { return c.h.Message() }

// Referencing returns the id of the message this chat references, or nil.
func (c *Chat) Referencing() []byte { return c.h.Referencing() }

// NewChat starts building a chat message.
func NewChat() *ChatBuilder { return &ChatBuilder{h: ffi.NewChatBuilder()} }

// Message sets the chat message text.
func (b *ChatBuilder) Message(message string) *ChatBuilder {
	b.h.Message(message)
	return b
}

// Reference sets the id of a message this chat references.
func (b *ChatBuilder) Reference(messageID []byte) *ChatBuilder {
	b.h.Reference(messageID)
	return b
}

// Finish finalizes the chat content, ready to send.
func (b *ChatBuilder) Finish() (*Content, error) {
	c, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Content{h: c}, nil
}
