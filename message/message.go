// Package message provides message and content types for the zktf SDK.
package message

import (
	"time"

	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/object"
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

// Descriptions returns the structured descriptions that make up the summary.
func (s *ContentSummary) Descriptions() []*ContentSummaryDescription {
	ds := s.h.Descriptions()
	out := make([]*ContentSummaryDescription, len(ds))
	for i, d := range ds {
		out[i] = &ContentSummaryDescription{h: d}
	}
	return out
}

// ContentSummaryDescriptionKind identifies the kind of a ContentSummaryDescription.
type ContentSummaryDescriptionKind uint32

const (
	SummaryUnknown        ContentSummaryDescriptionKind = ContentSummaryDescriptionKind(ffi.SummaryDescriptionUnknown)
	SummaryChatMessage    ContentSummaryDescriptionKind = ContentSummaryDescriptionKind(ffi.SummaryDescriptionChatMessage)
	SummaryChatReference  ContentSummaryDescriptionKind = ContentSummaryDescriptionKind(ffi.SummaryDescriptionChatReference)
	SummaryChatAttachment ContentSummaryDescriptionKind = ContentSummaryDescriptionKind(ffi.SummaryDescriptionChatAttachment)
	SummaryCredential     ContentSummaryDescriptionKind = ContentSummaryDescriptionKind(ffi.SummaryDescriptionCredential)
	SummaryPresentation   ContentSummaryDescriptionKind = ContentSummaryDescriptionKind(ffi.SummaryDescriptionPresentation)
	SummaryAsset          ContentSummaryDescriptionKind = ContentSummaryDescriptionKind(ffi.SummaryDescriptionAsset)
	SummarySignature      ContentSummaryDescriptionKind = ContentSummaryDescriptionKind(ffi.SummaryDescriptionSignature)
	SummaryVerification   ContentSummaryDescriptionKind = ContentSummaryDescriptionKind(ffi.SummaryDescriptionVerification)
	SummaryPairing        ContentSummaryDescriptionKind = ContentSummaryDescriptionKind(ffi.SummaryDescriptionPairing)
)

// ContentSummaryDescription is one structured element of a ContentSummary.
type ContentSummaryDescription struct {
	h *ffi.SummaryDescription
}

// Kind returns which kind of description this is.
func (d *ContentSummaryDescription) Kind() ContentSummaryDescriptionKind {
	return ContentSummaryDescriptionKind(d.h.Kind())
}

// ChatMessage returns the chat message text (SummaryChatMessage descriptions).
func (d *ContentSummaryDescription) ChatMessage() string { return d.h.AsChatMessage() }

// ChatReference returns the referenced message id (SummaryChatReference descriptions).
func (d *ContentSummaryDescription) ChatReference() []byte { return d.h.AsChatReference() }

// ChatAttachment returns the attached object (SummaryChatAttachment descriptions).
func (d *ContentSummaryDescription) ChatAttachment() *object.Object {
	o := d.h.AsChatAttachment()
	if o == nil {
		return nil
	}
	return ffi.ToObject(o).(*object.Object)
}

// Credential returns the credential types (SummaryCredential descriptions).
func (d *ContentSummaryDescription) Credential() []string { return d.h.AsCredential() }

// Presentation returns the presentation types (SummaryPresentation descriptions).
func (d *ContentSummaryDescription) Presentation() []string { return d.h.AsPresentation() }

// Asset returns the asset object (SummaryAsset descriptions).
func (d *ContentSummaryDescription) Asset() *object.Object {
	o := d.h.AsAsset()
	if o == nil {
		return nil
	}
	return ffi.ToObject(o).(*object.Object)
}

// Signature returns the signature bytes (SummarySignature descriptions).
func (d *ContentSummaryDescription) Signature() []byte { return d.h.AsSignature() }

// Verification returns the verified credential types (SummaryVerification descriptions).
func (d *ContentSummaryDescription) Verification() []string { return d.h.AsVerification() }

// Pairing returns the pairing roles bitmask (SummaryPairing descriptions).
func (d *ContentSummaryDescription) Pairing() uint64 { return d.h.AsPairing() }

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

// Metadata returns the opaque metadata payload attached to the message, if one
// was present. The payload is internal to the network and is not interpreted by
// this SDK. The boolean is false when no metadata is present.
func (m *Message) Metadata() ([]byte, bool) { return m.h.Metadata() }
