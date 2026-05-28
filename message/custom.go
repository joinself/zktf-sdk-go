package message

import "github.com/joinself/zktf-sdk-go/internal/ffi"

// Custom is a decoded custom message payload.
type Custom struct {
	h *ffi.Custom
}

// CustomBuilder builds custom message content.
type CustomBuilder struct {
	h *ffi.CustomBuilder
}

// CustomDecode decodes message content as a custom payload.
func CustomDecode(content *Content) (*Custom, error) {
	c, err := ffi.CustomFromContent(content.h)
	if err != nil {
		return nil, err
	}

	return &Custom{h: c}, nil
}

// Payload returns the custom payload bytes.
func (c *Custom) Payload() []byte { return c.h.Payload() }

// NewCustom starts building custom content.
func NewCustom() *CustomBuilder { return &CustomBuilder{h: ffi.NewCustomBuilder()} }

// Payload sets the custom payload bytes.
func (b *CustomBuilder) Payload(payload []byte) *CustomBuilder {
	b.h.Payload(payload)
	return b
}

// Finish finalizes the custom content, ready to send.
func (b *CustomBuilder) Finish() (*Content, error) {
	c, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Content{h: c}, nil
}
