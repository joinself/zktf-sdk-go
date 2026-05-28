package message

import (
	"time"

	"github.com/joinself/zktf-sdk-go/internal/ffi"
)

// ActionKind identifies the per-kind type of a polymorphic Action.
type ActionKind uint32

const (
	ActionUnknown                ActionKind = ActionKind(ffi.ActionKindUnknown)
	ActionCredentialPresentation ActionKind = ActionKind(ffi.ActionKindCredentialPresentation)
	ActionCredentialVerification ActionKind = ActionKind(ffi.ActionKindCredentialVerification)
	ActionIdentitySigning        ActionKind = ActionKind(ffi.ActionKindIdentitySigning)
	ActionDevicePairing          ActionKind = ActionKind(ffi.ActionKindDevicePairing)
)

// OutcomeKind identifies the per-kind type of a polymorphic Outcome.
type OutcomeKind uint32

const (
	OutcomeUnknown                OutcomeKind = OutcomeKind(ffi.OutcomeKindUnknown)
	OutcomeCredentialPresentation OutcomeKind = OutcomeKind(ffi.OutcomeKindCredentialPresentation)
	OutcomeCredentialVerification OutcomeKind = OutcomeKind(ffi.OutcomeKindCredentialVerification)
	OutcomeIdentitySigning        OutcomeKind = OutcomeKind(ffi.OutcomeKindIdentitySigning)
	OutcomeDevicePairing          OutcomeKind = OutcomeKind(ffi.OutcomeKindDevicePairing)
)

// ResponseStatus is the per-action or overall status of an exchange response.
type ResponseStatus uint32

const (
	StatusUnknown       ResponseStatus = ResponseStatus(ffi.ResponseStatusUnknown)
	StatusOK            ResponseStatus = ResponseStatus(ffi.ResponseStatusOK)
	StatusAccepted      ResponseStatus = ResponseStatus(ffi.ResponseStatusAccepted)
	StatusCreated       ResponseStatus = ResponseStatus(ffi.ResponseStatusCreated)
	StatusBadRequest    ResponseStatus = ResponseStatus(ffi.ResponseStatusBadRequest)
	StatusUnauthorized  ResponseStatus = ResponseStatus(ffi.ResponseStatusUnauthorized)
	StatusForbidden     ResponseStatus = ResponseStatus(ffi.ResponseStatusForbidden)
	StatusNotFound      ResponseStatus = ResponseStatus(ffi.ResponseStatusNotFound)
	StatusNotAcceptable ResponseStatus = ResponseStatus(ffi.ResponseStatusNotAcceptable)
	StatusConflict      ResponseStatus = ResponseStatus(ffi.ResponseStatusConflict)
)

// Action is the generic polymorphic wrapper for the per-kind action types
// inside an exchange request. Downcast via the As* methods.
type Action struct {
	h *ffi.Action
}

// Outcome is the generic polymorphic wrapper for per-kind result types inside
// an exchange response.
type Outcome struct {
	h *ffi.Outcome
}

// OutcomeBuilder builds a generic outcome carrying one of the per-kind results.
type OutcomeBuilder struct {
	h *ffi.OutcomeBuilder
}

// ExchangeRequest is a decoded exchange request.
type ExchangeRequest struct {
	h *ffi.ExchangeRequest
}

// ExchangeRequestBuilder builds an exchange request.
type ExchangeRequestBuilder struct {
	h *ffi.ExchangeRequestBuilder
}

// ExchangeResponse is a decoded exchange response.
type ExchangeResponse struct {
	h *ffi.ExchangeResponse
}

// ExchangeResponseBuilder builds an exchange response.
type ExchangeResponseBuilder struct {
	h *ffi.ExchangeResponseBuilder
}

func init() {
	ffi.ActionOf = func(o any) *ffi.Action { return o.(*Action).h }
	ffi.ToAction = func(h *ffi.Action) any { return &Action{h: h} }

	ffi.OutcomeOf = func(o any) *ffi.Outcome { return o.(*Outcome).h }
	ffi.ToOutcome = func(h *ffi.Outcome) any { return &Outcome{h: h} }

	ffi.ExchangeRequestOf = func(o any) *ffi.ExchangeRequest { return o.(*ExchangeRequest).h }
	ffi.ToExchangeRequest = func(h *ffi.ExchangeRequest) any { return &ExchangeRequest{h: h} }

	ffi.ExchangeResponseOf = func(o any) *ffi.ExchangeResponse { return o.(*ExchangeResponse).h }
	ffi.ToExchangeResponse = func(h *ffi.ExchangeResponse) any { return &ExchangeResponse{h: h} }
}

// Kind returns the kind of action.
func (a *Action) Kind() ActionKind { return ActionKind(a.h.Kind()) }

// ID returns the action's id bytes.
func (a *Action) ID() []byte { return a.h.ID() }

// AsPresentation downcasts the action to a credential presentation request.
func (a *Action) AsPresentation() (*PresentationRequest, error) {
	r, err := a.h.AsPresentation()
	if err != nil {
		return nil, err
	}

	return &PresentationRequest{h: r}, nil
}

// AsVerification downcasts the action to a credential verification request.
func (a *Action) AsVerification() (*VerificationRequest, error) {
	r, err := a.h.AsVerification()
	if err != nil {
		return nil, err
	}

	return &VerificationRequest{h: r}, nil
}

// AsIdentitySigning downcasts the action to an identity-signing request.
func (a *Action) AsIdentitySigning() (*IdentitySigningRequest, error) {
	r, err := a.h.AsIdentitySigning()
	if err != nil {
		return nil, err
	}

	return &IdentitySigningRequest{h: r}, nil
}

// AsDevicePairing downcasts the action to a device-pairing request.
func (a *Action) AsDevicePairing() (*DevicePairingRequest, error) {
	r, err := a.h.AsDevicePairing()
	if err != nil {
		return nil, err
	}

	return &DevicePairingRequest{h: r}, nil
}

// Kind returns the kind of outcome.
func (o *Outcome) Kind() OutcomeKind { return OutcomeKind(o.h.Kind()) }

// ActionID returns the id of the action this outcome refers to.
func (o *Outcome) ActionID() []byte { return o.h.ActionID() }

// Status returns the per-action response status.
func (o *Outcome) Status() ResponseStatus { return ResponseStatus(o.h.Status()) }

// ErrorMessage returns the per-action error message, or "".
func (o *Outcome) ErrorMessage() string { return o.h.ErrorMessage() }

// AsPresentation downcasts the outcome to a credential presentation result.
func (o *Outcome) AsPresentation() (*PresentationResponse, error) {
	r, err := o.h.AsPresentation()
	if err != nil {
		return nil, err
	}

	return &PresentationResponse{h: r}, nil
}

// AsVerification downcasts the outcome to a credential verification result.
func (o *Outcome) AsVerification() (*VerificationResponse, error) {
	r, err := o.h.AsVerification()
	if err != nil {
		return nil, err
	}

	return &VerificationResponse{h: r}, nil
}

// AsIdentitySigning downcasts the outcome to an identity-signing result.
func (o *Outcome) AsIdentitySigning() (*IdentitySigningResponse, error) {
	r, err := o.h.AsIdentitySigning()
	if err != nil {
		return nil, err
	}

	return &IdentitySigningResponse{h: r}, nil
}

// AsDevicePairing downcasts the outcome to a device-pairing result.
func (o *Outcome) AsDevicePairing() (*DevicePairingResponse, error) {
	r, err := o.h.AsDevicePairing()
	if err != nil {
		return nil, err
	}

	return &DevicePairingResponse{h: r}, nil
}

// NewOutcome starts building an outcome.
func NewOutcome() *OutcomeBuilder { return &OutcomeBuilder{h: ffi.NewOutcomeBuilder()} }

// ActionID sets the id of the action this outcome refers to.
func (b *OutcomeBuilder) ActionID(id []byte) *OutcomeBuilder {
	b.h.ActionID(id)
	return b
}

// Status sets the per-action response status.
func (b *OutcomeBuilder) Status(s ResponseStatus) *OutcomeBuilder {
	b.h.Status(ffi.ResponseStatus(s))
	return b
}

// ErrorMessage sets a per-action error message.
func (b *OutcomeBuilder) ErrorMessage(msg string) *OutcomeBuilder {
	b.h.ErrorMessage(msg)
	return b
}

// ResultPresentation attaches a credential-presentation response.
func (b *OutcomeBuilder) ResultPresentation(r *PresentationResponse) *OutcomeBuilder {
	b.h.ResultPresentation(r.h)
	return b
}

// ResultVerification attaches a credential-verification response.
func (b *OutcomeBuilder) ResultVerification(r *VerificationResponse) *OutcomeBuilder {
	b.h.ResultVerification(r.h)
	return b
}

// ResultSigning attaches an identity-signing response.
func (b *OutcomeBuilder) ResultSigning(r *IdentitySigningResponse) *OutcomeBuilder {
	b.h.ResultSigning(r.h)
	return b
}

// ResultPairing attaches a device-pairing response.
func (b *OutcomeBuilder) ResultPairing(r *DevicePairingResponse) *OutcomeBuilder {
	b.h.ResultPairing(r.h)
	return b
}

// Finish finalizes the outcome.
func (b *OutcomeBuilder) Finish() (*Outcome, error) {
	o, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Outcome{h: o}, nil
}

// ExchangeRequestDecode decodes message content as an exchange request.
func ExchangeRequestDecode(content *Content) (*ExchangeRequest, error) {
	r, err := ffi.ExchangeRequestFromContent(content.h)
	if err != nil {
		return nil, err
	}

	return &ExchangeRequest{h: r}, nil
}

// ID returns the request id.
func (r *ExchangeRequest) ID() []byte { return r.h.ID() }

// Purpose returns the request purpose string.
func (r *ExchangeRequest) Purpose() string { return r.h.Purpose() }

// Expires returns when the request expires.
func (r *ExchangeRequest) Expires() time.Time { return time.Unix(r.h.Expires(), 0) }

// Flags returns the request flags bitfield.
func (r *ExchangeRequest) Flags() uint64 { return r.h.Flags() }

// Actions returns the actions contained in the request.
func (r *ExchangeRequest) Actions() ([]*Action, error) {
	as, err := r.h.Actions()
	if err != nil {
		return nil, err
	}

	out := make([]*Action, len(as))
	for i, a := range as {
		out[i] = &Action{h: a}
	}

	return out, nil
}

// NewExchangeRequest starts building an exchange request.
func NewExchangeRequest() *ExchangeRequestBuilder {
	return &ExchangeRequestBuilder{h: ffi.NewExchangeRequestBuilder()}
}

// ID sets the request id.
func (b *ExchangeRequestBuilder) ID(id []byte) *ExchangeRequestBuilder {
	b.h.ID(id)
	return b
}

// Purpose sets the request purpose string.
func (b *ExchangeRequestBuilder) Purpose(p string) *ExchangeRequestBuilder {
	b.h.Purpose(p)
	return b
}

// Expires sets when the request expires.
func (b *ExchangeRequestBuilder) Expires(t time.Time) *ExchangeRequestBuilder {
	b.h.Expires(t.Unix())
	return b
}

// Flags sets the request flags bitfield.
func (b *ExchangeRequestBuilder) Flags(flags uint64) *ExchangeRequestBuilder {
	b.h.Flags(flags)
	return b
}

// Action appends an action to the request.
func (b *ExchangeRequestBuilder) Action(a *Action) *ExchangeRequestBuilder {
	b.h.Action(a.h)
	return b
}

// Finish finalizes the exchange request, ready to send.
func (b *ExchangeRequestBuilder) Finish() (*Content, error) {
	c, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Content{h: c}, nil
}

// ExchangeResponseDecode decodes message content as an exchange response.
func ExchangeResponseDecode(content *Content) (*ExchangeResponse, error) {
	r, err := ffi.ExchangeResponseFromContent(content.h)
	if err != nil {
		return nil, err
	}

	return &ExchangeResponse{h: r}, nil
}

// ID returns the response id.
func (r *ExchangeResponse) ID() []byte { return r.h.ID() }

// ResponseTo returns the id of the request being responded to.
func (r *ExchangeResponse) ResponseTo() []byte { return r.h.ResponseTo() }

// Status returns the overall response status.
func (r *ExchangeResponse) Status() ResponseStatus { return ResponseStatus(r.h.Status()) }

// ErrorMessage returns the overall error message, or "".
func (r *ExchangeResponse) ErrorMessage() string { return r.h.ErrorMessage() }

// Outcomes returns the per-action outcomes in the response.
func (r *ExchangeResponse) Outcomes() ([]*Outcome, error) {
	os, err := r.h.Outcomes()
	if err != nil {
		return nil, err
	}

	out := make([]*Outcome, len(os))
	for i, o := range os {
		out[i] = &Outcome{h: o}
	}

	return out, nil
}

// NewExchangeResponse starts building an exchange response.
func NewExchangeResponse() *ExchangeResponseBuilder {
	return &ExchangeResponseBuilder{h: ffi.NewExchangeResponseBuilder()}
}

// ID sets the response id.
func (b *ExchangeResponseBuilder) ID(id []byte) *ExchangeResponseBuilder {
	b.h.ID(id)
	return b
}

// ResponseTo sets the id of the request being responded to.
func (b *ExchangeResponseBuilder) ResponseTo(requestID []byte) *ExchangeResponseBuilder {
	b.h.ResponseTo(requestID)
	return b
}

// Status sets the overall response status.
func (b *ExchangeResponseBuilder) Status(s ResponseStatus) *ExchangeResponseBuilder {
	b.h.Status(ffi.ResponseStatus(s))
	return b
}

// ErrorMessage sets the overall error message.
func (b *ExchangeResponseBuilder) ErrorMessage(msg string) *ExchangeResponseBuilder {
	b.h.ErrorMessage(msg)
	return b
}

// Outcome appends an outcome to the response.
func (b *ExchangeResponseBuilder) Outcome(o *Outcome) *ExchangeResponseBuilder {
	b.h.Outcome(o.h)
	return b
}

// Finish finalizes the exchange response, ready to send.
func (b *ExchangeResponseBuilder) Finish() (*Content, error) {
	c, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Content{h: c}, nil
}
