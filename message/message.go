// Package message provides message and content types for the zktf SDK.
package message

import (
	"time"

	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// ContentType identifies the type of a message content payload.
type ContentType uint32

const (
	ContentUnknown           ContentType = ContentType(ffi.ContentUnknown)
	ContentCustom            ContentType = ContentType(ffi.ContentCustom)
	ContentChat              ContentType = ContentType(ffi.ContentChat)
	ContentReceipt           ContentType = ContentType(ffi.ContentReceipt)
	ContentCredential        ContentType = ContentType(ffi.ContentCredential)
	ContentIntroduction      ContentType = ContentType(ffi.ContentIntroduction)
	ContentDiscoveryRequest  ContentType = ContentType(ffi.ContentDiscoveryRequest)
	ContentDiscoveryResponse ContentType = ContentType(ffi.ContentDiscoveryResponse)
	ContentExchangeRequest   ContentType = ContentType(ffi.ContentExchangeRequest)
	ContentExchangeResponse  ContentType = ContentType(ffi.ContentExchangeResponse)
)

// Content is a message payload, either decoded from a Message or produced by a
// content builder ready to send.
type Content struct {
	h *ffi.Content
}

// ContentSummary is a compact summary of message content carried in a push
// notification.
type ContentSummary struct {
	h *ffi.MessageContentSummary
}

// Message is a message delivered to an account's OnMessage callback.
type Message struct {
	h *ffi.Message
}

func init() {
	ffi.ContentOf = func(o any) *ffi.Content { return o.(*Content).h }
	ffi.ToContent = func(h *ffi.Content) any { return &Content{h: h} }

	ffi.MessageContentSummaryOf = func(o any) *ffi.MessageContentSummary {
		return o.(*ContentSummary).h
	}
	ffi.ToMessageContentSummary = func(h *ffi.MessageContentSummary) any {
		return &ContentSummary{h: h}
	}

	ffi.MessageOf = func(o any) *ffi.Message { return o.(*Message).h }
	ffi.ToMessage = func(h *ffi.Message) any { return &Message{h: h} }
}

// Type returns the type of content.
func (c *Content) Type() ContentType { return ContentType(c.h.TypeOf()) }

// ID returns the content's id.
func (c *Content) ID() []byte { return c.h.ID() }

// Summary builds a compact summary suitable for inclusion in a push notification.
func (c *Content) Summary() (*ContentSummary, error) {
	s, err := ffi.SummaryOf(c.h)
	if err != nil {
		return nil, err
	}

	return &ContentSummary{h: s}, nil
}

// ID returns the id of the underlying content.
func (s *ContentSummary) ID() []byte { return s.h.ID() }

// Type returns the type of the summarized content.
func (s *ContentSummary) Type() ContentType { return ContentType(s.h.TypeOf()) }

// ID returns the message id.
func (m *Message) ID() []byte { return m.h.ID() }

// From returns the sender's address.
func (m *Message) From() *signing.PublicKey {
	return ffi.ToSigningPublicKey(m.h.FromAddress()).(*signing.PublicKey)
}

// To returns the recipient's address.
func (m *Message) To() *signing.PublicKey {
	return ffi.ToSigningPublicKey(m.h.ToAddress()).(*signing.PublicKey)
}

// Timestamp returns when the message was sent.
func (m *Message) Timestamp() time.Time { return time.Unix(m.h.Timestamp(), 0) }

// Content decodes and returns the message content.
func (m *Message) Content() *Content { return &Content{h: m.h.Content()} }
