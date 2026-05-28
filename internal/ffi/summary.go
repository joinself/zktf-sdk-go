package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

import (
	"runtime"
	"time"
	"unsafe"
)

// MessageContentSummary wraps a zktf_message_content_summary handle — a compact
// summary of message content, suitable for inclusion in a push notification.
type MessageContentSummary struct {
	ptr *C.zktf_message_content_summary
}

func newMessageContentSummary(ptr *C.zktf_message_content_summary) *MessageContentSummary {
	if ptr == nil {
		return nil
	}
	s := &MessageContentSummary{ptr: ptr}
	runtime.AddCleanup(s, func(ptr *C.zktf_message_content_summary) {
		C.zktf_message_content_summary_destroy(ptr)
	}, s.ptr)
	return s
}

// SummaryOf builds a summary of a piece of message content.
func SummaryOf(content *Content) (*MessageContentSummary, error) {
	var out *C.zktf_message_content_summary
	if err := status(C.zktf_message_content_summary_of(content.ptr, &out)); err != nil {
		return nil, err
	}
	return newMessageContentSummary(out), nil
}

// ID returns the id of the underlying content.
func (s *MessageContentSummary) ID() []byte {
	return C.GoBytes(unsafe.Pointer(C.zktf_message_content_summary_id(s.ptr)), messageIDLen)
}

// TypeOf returns the type of the summarized content.
func (s *MessageContentSummary) TypeOf() ContentType {
	return ContentType(C.zktf_message_content_summary_type_of(s.ptr))
}

// NotificationSend sends a push notification to the given address carrying the
// content summary, via callback.
func (a *Account) NotificationSend(to *SigningPublicKey, summary *MessageContentSummary, timeout time.Duration) error {
	fut := C.zktf_account_notification_send(a.ptr, to.ptr, summary.ptr)

	return AwaitStatus(fut, timeout)
}
