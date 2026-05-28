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

// Receipt wraps a zktf_message_content_receipt handle.
type Receipt struct {
	ptr *C.zktf_message_content_receipt
}

func newReceipt(ptr *C.zktf_message_content_receipt) *Receipt {
	if ptr == nil {
		return nil
	}
	r := &Receipt{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_message_content_receipt) {
		C.zktf_message_content_receipt_destroy(ptr)
	}, r.ptr)
	return r
}

// ReceiptFromContent decodes message content as a receipt.
func ReceiptFromContent(content *Content) (*Receipt, error) {
	var ptr *C.zktf_message_content_receipt
	if err := status(C.zktf_message_content_as_receipt(content.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newReceipt(ptr), nil
}

// Delivered returns the ids of messages marked delivered.
func (r *Receipt) Delivered() [][]byte {
	return messageIDsFrom(C.zktf_message_content_receipt_delivered(r.ptr))
}

// Read returns the ids of messages marked read.
func (r *Receipt) Read() [][]byte {
	return messageIDsFrom(C.zktf_message_content_receipt_read(r.ptr))
}

// ReceiptBuilder wraps a zktf_message_content_receipt_builder handle.
type ReceiptBuilder struct {
	ptr *C.zktf_message_content_receipt_builder
}

// NewReceiptBuilder initializes a new receipt content builder.
func NewReceiptBuilder() *ReceiptBuilder {
	ptr := C.zktf_message_content_receipt_builder_init()
	b := &ReceiptBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_receipt_builder) {
		C.zktf_message_content_receipt_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// Delivered marks a message id as delivered.
func (b *ReceiptBuilder) Delivered(messageID []byte) *ReceiptBuilder {
	buf, _ := cbytes(messageID)
	defer free(unsafe.Pointer(buf))
	C.zktf_message_content_receipt_builder_delivered(b.ptr, buf)
	return b
}

// Read marks a message id as read.
func (b *ReceiptBuilder) Read(messageID []byte) *ReceiptBuilder {
	buf, _ := cbytes(messageID)
	defer free(unsafe.Pointer(buf))
	C.zktf_message_content_receipt_builder_read(b.ptr, buf)
	return b
}

// Finish finalizes the receipt content, ready to send.
func (b *ReceiptBuilder) Finish() (*Content, error) {
	var ptr *C.zktf_message_content
	if err := status(C.zktf_message_content_receipt_builder_finish(b.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newContent(ptr), nil
}
