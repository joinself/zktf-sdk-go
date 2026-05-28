package message

import "github.com/joinself/zktf-sdk-go/internal/ffi"

// Receipt is a decoded delivery/read receipt.
type Receipt struct {
	h *ffi.Receipt
}

// ReceiptBuilder builds receipt message content.
type ReceiptBuilder struct {
	h *ffi.ReceiptBuilder
}

// ReceiptDecode decodes message content as a receipt.
func ReceiptDecode(content *Content) (*Receipt, error) {
	r, err := ffi.ReceiptFromContent(content.h)
	if err != nil {
		return nil, err
	}

	return &Receipt{h: r}, nil
}

// Delivered returns the ids of messages marked delivered.
func (r *Receipt) Delivered() [][]byte { return r.h.Delivered() }

// Read returns the ids of messages marked read.
func (r *Receipt) Read() [][]byte { return r.h.Read() }

// NewReceipt starts building a receipt.
func NewReceipt() *ReceiptBuilder { return &ReceiptBuilder{h: ffi.NewReceiptBuilder()} }

// Delivered marks a message id as delivered.
func (b *ReceiptBuilder) Delivered(messageID []byte) *ReceiptBuilder {
	b.h.Delivered(messageID)
	return b
}

// Read marks a message id as read.
func (b *ReceiptBuilder) Read(messageID []byte) *ReceiptBuilder {
	b.h.Read(messageID)
	return b
}

// Finish finalizes the receipt content, ready to send.
func (b *ReceiptBuilder) Finish() (*Content, error) {
	c, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Content{h: c}, nil
}
