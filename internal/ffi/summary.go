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

// SummaryDescriptionKind mirrors zktf_message_content_summary_description_type.
type SummaryDescriptionKind uint32

const (
	SummaryDescriptionUnknown        SummaryDescriptionKind = C.CONTENT_SUMMARY_DESCRIPTION_UNKNOWN
	SummaryDescriptionChatMessage    SummaryDescriptionKind = C.CONTENT_SUMMARY_DESCRIPTION_CHAT_MESSAGE
	SummaryDescriptionChatReference  SummaryDescriptionKind = C.CONTENT_SUMMARY_DESCRIPTION_CHAT_REFERENCE
	SummaryDescriptionChatAttachment SummaryDescriptionKind = C.CONTENT_SUMMARY_DESCRIPTION_CHAT_ATTACHMENT
	SummaryDescriptionCredential     SummaryDescriptionKind = C.CONTENT_SUMMARY_DESCRIPTION_CREDENTIAL
	SummaryDescriptionPresentation   SummaryDescriptionKind = C.CONTENT_SUMMARY_DESCRIPTION_PRESENTATION
	SummaryDescriptionAsset          SummaryDescriptionKind = C.CONTENT_SUMMARY_DESCRIPTION_ASSET
	SummaryDescriptionSignature      SummaryDescriptionKind = C.CONTENT_SUMMARY_DESCRIPTION_SIGNATURE
	SummaryDescriptionVerification   SummaryDescriptionKind = C.CONTENT_SUMMARY_DESCRIPTION_VERIFICATION
	SummaryDescriptionPairing        SummaryDescriptionKind = C.CONTENT_SUMMARY_DESCRIPTION_PAIRING
)

// Descriptions returns the structured descriptions that make up the summary.
func (s *MessageContentSummary) Descriptions() []*SummaryDescription {
	c := C.zktf_message_content_summary_descriptions(s.ptr)
	if c == nil {
		return nil
	}
	defer C.zktf_collection_message_content_summary_description_destroy(c)
	n := int(C.zktf_collection_message_content_summary_description_len(c))
	out := make([]*SummaryDescription, n)
	for i := 0; i < n; i++ {
		out[i] = newSummaryDescription(C.zktf_collection_message_content_summary_description_at(c, C.size_t(i)))
	}
	return out
}

// SummaryDescription is one structured element of a content summary.
type SummaryDescription struct {
	ptr *C.zktf_message_content_summary_description
}

func newSummaryDescription(ptr *C.zktf_message_content_summary_description) *SummaryDescription {
	if ptr == nil {
		return nil
	}
	d := &SummaryDescription{ptr: ptr}
	runtime.AddCleanup(d, func(ptr *C.zktf_message_content_summary_description) {
		C.zktf_message_content_summary_description_destroy(ptr)
	}, d.ptr)
	return d
}

// Kind returns which kind of description this is.
func (d *SummaryDescription) Kind() SummaryDescriptionKind {
	return SummaryDescriptionKind(C.zktf_message_content_summary_description_type_of(d.ptr))
}

// AsChatMessage returns the chat message text (CHAT_MESSAGE descriptions).
func (d *SummaryDescription) AsChatMessage() string {
	return goStringFromBuffer(C.zktf_message_content_summary_description_as_chat_message(d.ptr))
}

// AsChatReference returns the referenced message id (CHAT_REFERENCE descriptions).
func (d *SummaryDescription) AsChatReference() []byte {
	return goBytesFromBuffer(C.zktf_message_content_summary_description_as_chat_reference(d.ptr))
}

// AsChatAttachment returns the attached object (CHAT_ATTACHMENT descriptions).
func (d *SummaryDescription) AsChatAttachment() *Object {
	return newObject(C.zktf_message_content_summary_description_as_chat_attachment(d.ptr))
}

// AsCredential returns the credential types (CREDENTIAL descriptions).
func (d *SummaryDescription) AsCredential() []string {
	return newCredentialTypeCollection(C.zktf_message_content_summary_description_as_credential(d.ptr)).Strings()
}

// AsPresentation returns the presentation types (PRESENTATION descriptions).
func (d *SummaryDescription) AsPresentation() []string {
	return presentationTypesFrom(C.zktf_message_content_summary_description_as_presentation(d.ptr))
}

// AsAsset returns the asset object (ASSET descriptions).
func (d *SummaryDescription) AsAsset() *Object {
	return newObject(C.zktf_message_content_summary_description_as_asset(d.ptr))
}

// AsSignature returns the signature bytes (SIGNATURE descriptions).
func (d *SummaryDescription) AsSignature() []byte {
	return goBytesFromBuffer(C.zktf_message_content_summary_description_as_signature(d.ptr))
}

// AsVerification returns the verified credential types (VERIFICATION descriptions).
func (d *SummaryDescription) AsVerification() []string {
	return newCredentialTypeCollection(C.zktf_message_content_summary_description_as_verification(d.ptr)).Strings()
}

// AsPairing returns the pairing roles bitmask (PAIRING descriptions).
func (d *SummaryDescription) AsPairing() uint64 {
	return uint64(C.zktf_message_content_summary_description_as_pairing(d.ptr))
}

// NotificationSend sends a push notification to the given address carrying the
// content summary, via callback.
func (a *Account) NotificationSend(to *SigningPublicKey, summary *MessageContentSummary, timeout time.Duration) error {
	fut := C.zktf_account_notification_send(a.ptr, to.ptr, summary.ptr)

	return AwaitStatus(fut, timeout)
}
