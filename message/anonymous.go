package message

import "github.com/joinself/zktf-sdk-go/internal/ffi"

// AnonymousMessage is a decoded anonymous message — an unencrypted,
// self-describing envelope (e.g. a pairing/QR code) carrying message content.
type AnonymousMessage struct {
	h *ffi.AnonymousMessage
}

// AnonymousMessageDecode decodes a base64 URL encoded anonymous message, such
// as a pairing code, into its message wrapper.
func AnonymousMessageDecode(code string) (*AnonymousMessage, error) {
	m, err := ffi.AnonymousMessageDecodeFromString(code)
	if err != nil {
		return nil, err
	}

	return &AnonymousMessage{h: m}, nil
}

// ID returns the message id.
func (m *AnonymousMessage) ID() []byte { return m.h.ID() }

// Content returns the message content.
func (m *AnonymousMessage) Content() *Content { return &Content{h: m.h.Content()} }
