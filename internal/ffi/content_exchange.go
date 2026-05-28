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

// ActionKind mirrors zktf_message_content_action_kind.
type ActionKind uint32

const (
	ActionKindUnknown                ActionKind = C.ACTION_KIND_UNKNOWN
	ActionKindCredentialPresentation ActionKind = C.ACTION_KIND_CREDENTIAL_PRESENTATION
	ActionKindCredentialVerification ActionKind = C.ACTION_KIND_CREDENTIAL_VERIFICATION
	ActionKindIdentitySigning        ActionKind = C.ACTION_KIND_IDENTITY_SIGNING
	ActionKindDevicePairing          ActionKind = C.ACTION_KIND_DEVICE_PAIRING
)

// OutcomeKind mirrors zktf_message_content_outcome_kind.
type OutcomeKind uint32

const (
	OutcomeKindUnknown                OutcomeKind = C.OUTCOME_KIND_UNKNOWN
	OutcomeKindCredentialPresentation OutcomeKind = C.OUTCOME_KIND_CREDENTIAL_PRESENTATION
	OutcomeKindCredentialVerification OutcomeKind = C.OUTCOME_KIND_CREDENTIAL_VERIFICATION
	OutcomeKindIdentitySigning        OutcomeKind = C.OUTCOME_KIND_IDENTITY_SIGNING
	OutcomeKindDevicePairing          OutcomeKind = C.OUTCOME_KIND_DEVICE_PAIRING
)

// ResponseStatus mirrors zktf_message_response_status.
type ResponseStatus uint32

const (
	ResponseStatusUnknown       ResponseStatus = C.RESPONSE_STATUS_UNKNOWN
	ResponseStatusOK            ResponseStatus = C.RESPONSE_STATUS_OK
	ResponseStatusAccepted      ResponseStatus = C.RESPONSE_STATUS_ACCEPTED
	ResponseStatusCreated       ResponseStatus = C.RESPONSE_STATUS_CREATED
	ResponseStatusBadRequest    ResponseStatus = C.RESPONSE_STATUS_BAD_REQUEST
	ResponseStatusUnauthorized  ResponseStatus = C.RESPONSE_STATUS_UNAUTHORIZED
	ResponseStatusForbidden     ResponseStatus = C.RESPONSE_STATUS_FORBIDDEN
	ResponseStatusNotFound      ResponseStatus = C.RESPONSE_STATUS_NOT_FOUND
	ResponseStatusNotAcceptable ResponseStatus = C.RESPONSE_STATUS_NOT_ACCEPTABLE
	ResponseStatusConflict      ResponseStatus = C.RESPONSE_STATUS_CONFLICT
)

// Action is the generic polymorphic wrapper for the per-kind action types
// inside an exchange request.
type Action struct {
	ptr *C.zktf_message_content_action
}

func newAction(ptr *C.zktf_message_content_action) *Action {
	if ptr == nil {
		return nil
	}
	a := &Action{ptr: ptr}
	runtime.AddCleanup(a, func(ptr *C.zktf_message_content_action) {
		C.zktf_message_content_action_destroy(ptr)
	}, a.ptr)
	return a
}

// Kind returns the kind of action.
func (a *Action) Kind() ActionKind {
	return ActionKind(C.zktf_message_content_action_kind_of(a.ptr))
}

// ID returns the action's id bytes.
func (a *Action) ID() []byte {
	n := C.zktf_message_content_action_id_len(a.ptr)
	return C.GoBytes(unsafe.Pointer(C.zktf_message_content_action_id(a.ptr)), C.int(n))
}

// AsPresentation downcasts the action to a credential presentation action.
func (a *Action) AsPresentation() (*PresentationAction, error) {
	var out *C.zktf_message_content_credential_presentation_action
	if err := status(C.zktf_message_content_action_as_credential_presentation(a.ptr, &out)); err != nil {
		return nil, err
	}
	return newPresentationAction(out), nil
}

// AsVerification downcasts the action to a credential verification action.
func (a *Action) AsVerification() (*VerificationAction, error) {
	var out *C.zktf_message_content_credential_verification_action
	if err := status(C.zktf_message_content_action_as_credential_verification(a.ptr, &out)); err != nil {
		return nil, err
	}
	return newVerificationAction(out), nil
}

// AsIdentitySigning downcasts the action to an identity-signing action.
func (a *Action) AsIdentitySigning() (*IdentitySigningAction, error) {
	var out *C.zktf_message_content_identity_signing_action
	if err := status(C.zktf_message_content_action_as_identity_signing(a.ptr, &out)); err != nil {
		return nil, err
	}
	return newIdentitySigningAction(out), nil
}

// AsDevicePairing downcasts the action to a device-pairing action.
func (a *Action) AsDevicePairing() (*DevicePairingAction, error) {
	var out *C.zktf_message_content_device_pairing_action
	if err := status(C.zktf_message_content_action_as_device_pairing(a.ptr, &out)); err != nil {
		return nil, err
	}
	return newDevicePairingAction(out), nil
}

// Outcome is the generic polymorphic wrapper for per-kind result types inside
// an exchange response.
type Outcome struct {
	ptr *C.zktf_message_content_outcome
}

func newOutcome(ptr *C.zktf_message_content_outcome) *Outcome {
	if ptr == nil {
		return nil
	}
	o := &Outcome{ptr: ptr}
	runtime.AddCleanup(o, func(ptr *C.zktf_message_content_outcome) {
		C.zktf_message_content_outcome_destroy(ptr)
	}, o.ptr)
	return o
}

// Kind returns the kind of outcome.
func (o *Outcome) Kind() OutcomeKind {
	return OutcomeKind(C.zktf_message_content_outcome_kind_of(o.ptr))
}

// ActionID returns the id of the action this outcome refers to.
func (o *Outcome) ActionID() []byte {
	n := C.zktf_message_content_outcome_action_id_len(o.ptr)
	return C.GoBytes(unsafe.Pointer(C.zktf_message_content_outcome_action_id(o.ptr)), C.int(n))
}

// Status returns the response status carried by the outcome.
func (o *Outcome) Status() ResponseStatus {
	return ResponseStatus(C.zktf_message_content_outcome_status(o.ptr))
}

// ErrorMessage returns the error message carried by the outcome, or "".
func (o *Outcome) ErrorMessage() string {
	return C.GoString(C.zktf_message_content_outcome_error_message(o.ptr))
}

// AsPresentation downcasts the outcome to a credential presentation result.
func (o *Outcome) AsPresentation() (*PresentationResult, error) {
	var out *C.zktf_message_content_credential_presentation_result
	if err := status(C.zktf_message_content_outcome_as_credential_presentation(o.ptr, &out)); err != nil {
		return nil, err
	}
	return newPresentationResult(out), nil
}

// AsVerification downcasts the outcome to a credential verification result.
func (o *Outcome) AsVerification() (*VerificationResult, error) {
	var out *C.zktf_message_content_credential_verification_result
	if err := status(C.zktf_message_content_outcome_as_credential_verification(o.ptr, &out)); err != nil {
		return nil, err
	}
	return newVerificationResult(out), nil
}

// AsIdentitySigning downcasts the outcome to an identity-signing result.
func (o *Outcome) AsIdentitySigning() (*IdentitySigningResult, error) {
	var out *C.zktf_message_content_identity_signing_result
	if err := status(C.zktf_message_content_outcome_as_identity_signing(o.ptr, &out)); err != nil {
		return nil, err
	}
	return newIdentitySigningResult(out), nil
}

// AsDevicePairing downcasts the outcome to a device-pairing result.
func (o *Outcome) AsDevicePairing() (*DevicePairingResult, error) {
	var out *C.zktf_message_content_device_pairing_result
	if err := status(C.zktf_message_content_outcome_as_device_pairing(o.ptr, &out)); err != nil {
		return nil, err
	}
	return newDevicePairingResult(out), nil
}

// OutcomeBuilder builds a generic outcome carrying one of the per-kind results.
type OutcomeBuilder struct {
	ptr *C.zktf_message_content_outcome_builder
}

// NewOutcomeBuilder initializes an outcome builder.
func NewOutcomeBuilder() *OutcomeBuilder {
	ptr := C.zktf_message_content_outcome_builder_init()
	b := &OutcomeBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_outcome_builder) {
		C.zktf_message_content_outcome_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// ActionID sets the id of the action this outcome refers to.
func (b *OutcomeBuilder) ActionID(id []byte) *OutcomeBuilder {
	buf, length := cbytes(id)
	defer free(unsafe.Pointer(buf))
	C.zktf_message_content_outcome_builder_action_id(b.ptr, buf, length)
	return b
}

// Status sets the response status.
func (b *OutcomeBuilder) Status(s ResponseStatus) *OutcomeBuilder {
	C.zktf_message_content_outcome_builder_status(b.ptr, C.enum_zktf_message_response_status(s))
	return b
}

// ErrorMessage sets a human-readable error message.
func (b *OutcomeBuilder) ErrorMessage(msg string) *OutcomeBuilder {
	cmsg := cstring(msg)
	defer free(unsafe.Pointer(cmsg))
	C.zktf_message_content_outcome_builder_error_message(b.ptr, cmsg)
	return b
}

// ResultPresentation attaches a credential presentation result.
func (b *OutcomeBuilder) ResultPresentation(r *PresentationResult) *OutcomeBuilder {
	C.zktf_message_content_outcome_builder_result_presentation(b.ptr, r.ptr)
	return b
}

// ResultVerification attaches a credential verification result.
func (b *OutcomeBuilder) ResultVerification(r *VerificationResult) *OutcomeBuilder {
	C.zktf_message_content_outcome_builder_result_verification(b.ptr, r.ptr)
	return b
}

// ResultSigning attaches an identity-signing result.
func (b *OutcomeBuilder) ResultSigning(r *IdentitySigningResult) *OutcomeBuilder {
	C.zktf_message_content_outcome_builder_result_signing(b.ptr, r.ptr)
	return b
}

// ResultPairing attaches a device-pairing result.
func (b *OutcomeBuilder) ResultPairing(r *DevicePairingResult) *OutcomeBuilder {
	C.zktf_message_content_outcome_builder_result_pairing(b.ptr, r.ptr)
	return b
}

// Finish finalizes the outcome.
func (b *OutcomeBuilder) Finish() (*Outcome, error) {
	var out *C.zktf_message_content_outcome
	if err := status(C.zktf_message_content_outcome_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newOutcome(out), nil
}

// ExchangeRequest wraps a zktf_message_content_exchange_request handle.
type ExchangeRequest struct {
	ptr *C.zktf_message_content_exchange_request
}

func newExchangeRequest(ptr *C.zktf_message_content_exchange_request) *ExchangeRequest {
	if ptr == nil {
		return nil
	}
	r := &ExchangeRequest{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_message_content_exchange_request) {
		C.zktf_message_content_exchange_request_destroy(ptr)
	}, r.ptr)
	return r
}

// ExchangeRequestFromContent decodes message content as an exchange request.
func ExchangeRequestFromContent(content *Content) (*ExchangeRequest, error) {
	var ptr *C.zktf_message_content_exchange_request
	if err := status(C.zktf_message_content_as_exchange_request(content.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newExchangeRequest(ptr), nil
}

// ID returns the request id.
func (r *ExchangeRequest) ID() []byte {
	return C.GoBytes(unsafe.Pointer(C.zktf_message_content_exchange_request_id(r.ptr)), messageIDLen)
}

// Purpose returns the request purpose string.
func (r *ExchangeRequest) Purpose() string {
	return C.GoString(C.zktf_message_content_exchange_request_purpose(r.ptr))
}

// Expires returns the unix timestamp (seconds) the request expires.
func (r *ExchangeRequest) Expires() int64 {
	return int64(C.zktf_message_content_exchange_request_expires(r.ptr))
}

// Flags returns the request flags bitfield.
func (r *ExchangeRequest) Flags() uint64 {
	return uint64(C.zktf_message_content_exchange_request_flags(r.ptr))
}

// Actions returns the actions contained in the request.
func (r *ExchangeRequest) Actions() ([]*Action, error) {
	var c *C.zktf_collection_message_content_action
	if err := status(C.zktf_message_content_exchange_request_actions(r.ptr, &c)); err != nil {
		return nil, err
	}
	defer C.zktf_collection_message_content_action_destroy(c)
	n := int(C.zktf_collection_message_content_action_len(c))
	out := make([]*Action, n)
	for i := 0; i < n; i++ {
		out[i] = newAction(C.zktf_collection_message_content_action_at(c, C.size_t(i)))
	}
	return out, nil
}

// ExchangeRequestBuilder builds an exchange request.
type ExchangeRequestBuilder struct {
	ptr *C.zktf_message_content_exchange_request_builder
}

// NewExchangeRequestBuilder initializes an exchange request builder.
func NewExchangeRequestBuilder() *ExchangeRequestBuilder {
	ptr := C.zktf_message_content_exchange_request_builder_init()
	b := &ExchangeRequestBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_exchange_request_builder) {
		C.zktf_message_content_exchange_request_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// ID sets the request id.
func (b *ExchangeRequestBuilder) ID(id []byte) *ExchangeRequestBuilder {
	buf, length := cbytes(id)
	defer free(unsafe.Pointer(buf))
	C.zktf_message_content_exchange_request_builder_id(b.ptr, buf, length)
	return b
}

// Purpose sets the request purpose string.
func (b *ExchangeRequestBuilder) Purpose(p string) *ExchangeRequestBuilder {
	cp := cstring(p)
	defer free(unsafe.Pointer(cp))
	C.zktf_message_content_exchange_request_builder_purpose(b.ptr, cp)
	return b
}

// Expires sets the request expiry as a unix timestamp (seconds).
func (b *ExchangeRequestBuilder) Expires(unix int64) *ExchangeRequestBuilder {
	C.zktf_message_content_exchange_request_builder_expires(b.ptr, C.int64_t(unix))
	return b
}

// Flags sets the request flags bitfield.
func (b *ExchangeRequestBuilder) Flags(flags uint64) *ExchangeRequestBuilder {
	C.zktf_message_content_exchange_request_builder_flags(b.ptr, C.uint64_t(flags))
	return b
}

// Action appends an action to the request.
func (b *ExchangeRequestBuilder) Action(a *Action) *ExchangeRequestBuilder {
	C.zktf_message_content_exchange_request_builder_action(b.ptr, a.ptr)
	return b
}

// Finish finalizes the exchange request, ready to send.
func (b *ExchangeRequestBuilder) Finish() (*Content, error) {
	var out *C.zktf_message_content
	if err := status(C.zktf_message_content_exchange_request_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newContent(out), nil
}

// ExchangeResponse wraps a zktf_message_content_exchange_response handle.
type ExchangeResponse struct {
	ptr *C.zktf_message_content_exchange_response
}

func newExchangeResponse(ptr *C.zktf_message_content_exchange_response) *ExchangeResponse {
	if ptr == nil {
		return nil
	}
	r := &ExchangeResponse{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_message_content_exchange_response) {
		C.zktf_message_content_exchange_response_destroy(ptr)
	}, r.ptr)
	return r
}

// ExchangeResponseFromContent decodes message content as an exchange response.
func ExchangeResponseFromContent(content *Content) (*ExchangeResponse, error) {
	var ptr *C.zktf_message_content_exchange_response
	if err := status(C.zktf_message_content_as_exchange_response(content.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newExchangeResponse(ptr), nil
}

// ID returns the response id.
func (r *ExchangeResponse) ID() []byte {
	return C.GoBytes(unsafe.Pointer(C.zktf_message_content_exchange_response_id(r.ptr)), messageIDLen)
}

// ResponseTo returns the id of the request being responded to.
func (r *ExchangeResponse) ResponseTo() []byte {
	return C.GoBytes(unsafe.Pointer(C.zktf_message_content_exchange_response_response_to(r.ptr)), messageIDLen)
}

// Status returns the overall response status.
func (r *ExchangeResponse) Status() ResponseStatus {
	return ResponseStatus(C.zktf_message_content_exchange_response_status(r.ptr))
}

// ErrorMessage returns the response error message, or "".
func (r *ExchangeResponse) ErrorMessage() string {
	return C.GoString(C.zktf_message_content_exchange_response_error_message(r.ptr))
}

// Outcomes returns the per-action outcomes contained in the response.
func (r *ExchangeResponse) Outcomes() ([]*Outcome, error) {
	var c *C.zktf_collection_message_content_outcome
	if err := status(C.zktf_message_content_exchange_response_outcomes(r.ptr, &c)); err != nil {
		return nil, err
	}
	defer C.zktf_collection_message_content_outcome_destroy(c)
	n := int(C.zktf_collection_message_content_outcome_len(c))
	out := make([]*Outcome, n)
	for i := 0; i < n; i++ {
		out[i] = newOutcome(C.zktf_collection_message_content_outcome_at(c, C.size_t(i)))
	}
	return out, nil
}

// ExchangeResponseBuilder builds an exchange response.
type ExchangeResponseBuilder struct {
	ptr *C.zktf_message_content_exchange_response_builder
}

// NewExchangeResponseBuilder initializes an exchange response builder.
func NewExchangeResponseBuilder() *ExchangeResponseBuilder {
	ptr := C.zktf_message_content_exchange_response_builder_init()
	b := &ExchangeResponseBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_exchange_response_builder) {
		C.zktf_message_content_exchange_response_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// ID sets the response id.
func (b *ExchangeResponseBuilder) ID(id []byte) *ExchangeResponseBuilder {
	buf, length := cbytes(id)
	defer free(unsafe.Pointer(buf))
	C.zktf_message_content_exchange_response_builder_id(b.ptr, buf, length)
	return b
}

// ResponseTo sets the id of the request being responded to.
func (b *ExchangeResponseBuilder) ResponseTo(requestID []byte) *ExchangeResponseBuilder {
	buf, _ := cbytes(requestID)
	defer free(unsafe.Pointer(buf))
	C.zktf_message_content_exchange_response_builder_response_to(b.ptr, buf)
	return b
}

// Status sets the overall response status.
func (b *ExchangeResponseBuilder) Status(s ResponseStatus) *ExchangeResponseBuilder {
	C.zktf_message_content_exchange_response_builder_status(b.ptr, C.enum_zktf_message_response_status(s))
	return b
}

// ErrorMessage sets the response error message.
func (b *ExchangeResponseBuilder) ErrorMessage(msg string) *ExchangeResponseBuilder {
	cmsg := cstring(msg)
	defer free(unsafe.Pointer(cmsg))
	C.zktf_message_content_exchange_response_builder_error_message(b.ptr, cmsg)
	return b
}

// Outcome appends an outcome to the response.
func (b *ExchangeResponseBuilder) Outcome(o *Outcome) *ExchangeResponseBuilder {
	C.zktf_message_content_exchange_response_builder_outcome(b.ptr, o.ptr)
	return b
}

// Finish finalizes the exchange response, ready to send.
func (b *ExchangeResponseBuilder) Finish() (*Content, error) {
	var out *C.zktf_message_content
	if err := status(C.zktf_message_content_exchange_response_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newContent(out), nil
}
