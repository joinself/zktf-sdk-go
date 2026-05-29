package message

import (
	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/object"
)

// VerificationRequest is a credential-verification request.
type VerificationRequest struct {
	h *ffi.VerificationAction
}

// VerificationRequestBuilder builds a credential-verification request.
type VerificationRequestBuilder struct {
	h *ffi.VerificationActionBuilder
}

// VerificationResponse is the response to a credential-verification request.
type VerificationResponse struct {
	h *ffi.VerificationResult
}

// VerificationResponseBuilder builds a credential-verification response.
type VerificationResponseBuilder struct {
	h *ffi.VerificationResultBuilder
}

func init() {
	ffi.VerificationActionOf = func(o any) *ffi.VerificationAction {
		return o.(*VerificationRequest).h
	}
	ffi.ToVerificationAction = func(h *ffi.VerificationAction) any {
		return &VerificationRequest{h: h}
	}

	ffi.VerificationResultOf = func(o any) *ffi.VerificationResult {
		return o.(*VerificationResponse).h
	}
	ffi.ToVerificationResult = func(h *ffi.VerificationResult) any {
		return &VerificationResponse{h: h}
	}
}

// CredentialTypes returns the requested credential types.
func (r *VerificationRequest) CredentialTypes() []string { return r.h.CredentialTypes() }

// Proof returns any verifiable presentations attached as proof.
func (r *VerificationRequest) Proof() []*credential.VerifiablePresentation {
	ps := r.h.Proof()
	out := make([]*credential.VerifiablePresentation, len(ps))

	for i, p := range ps {
		out[i] = ffi.ToVerifiablePresentation(p).(*credential.VerifiablePresentation)
	}

	return out
}

// Evidence returns the objects attached to the request as supporting evidence.
func (r *VerificationRequest) Evidence() []*VerificationEvidence {
	es := r.h.Evidence()
	out := make([]*VerificationEvidence, len(es))

	for i, e := range es {
		out[i] = &VerificationEvidence{h: e}
	}

	return out
}

// Parameters returns the typed parameters attached to the request.
func (r *VerificationRequest) Parameters() []*VerificationParameter {
	ps := r.h.Parameters()
	out := make([]*VerificationParameter, len(ps))

	for i, p := range ps {
		out[i] = &VerificationParameter{h: p}
	}

	return out
}

// AsAction wraps this request into a generic Action.
func (r *VerificationRequest) AsAction() *Action { return &Action{h: r.h.AsAction()} }

// NewVerificationRequest starts building a verification request.
func NewVerificationRequest() *VerificationRequestBuilder {
	return &VerificationRequestBuilder{h: ffi.NewVerificationActionBuilder()}
}

// CredentialType sets the requested credential types.
func (b *VerificationRequestBuilder) CredentialType(types ...string) *VerificationRequestBuilder {
	b.h.CredentialType(ffi.NewCredentialTypes(types))
	return b
}

// Proof attaches a verifiable presentation as proof.
func (b *VerificationRequestBuilder) Proof(p *credential.VerifiablePresentation) *VerificationRequestBuilder {
	b.h.Proof(ffi.VerifiablePresentationOf(p))
	return b
}

// Evidence attaches an object as supporting evidence under a named type.
func (b *VerificationRequestBuilder) Evidence(evidenceType string, obj *object.Object) *VerificationRequestBuilder {
	b.h.Evidence(evidenceType, ffi.ObjectOf(obj))
	return b
}

// Parameter attaches a typed parameter. value must be one of []byte, string,
// bool, an integer type (int, int8, int16, int32, int64), an unsigned type
// (uint, uint8, uint16, uint32, uint64), float32, float64, [][]byte or
// []string; values of any other type are ignored.
func (b *VerificationRequestBuilder) Parameter(key string, value any) *VerificationRequestBuilder {
	pv := ffi.NewParameterValue(value)
	if pv == nil {
		return b
	}

	b.h.Parameter(key, pv)
	return b
}

// Finish finalizes the verification request.
func (b *VerificationRequestBuilder) Finish() (*VerificationRequest, error) {
	a, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &VerificationRequest{h: a}, nil
}

// VerificationEvidence is an object attached to a verification request as
// supporting evidence, tagged with an evidence type.
type VerificationEvidence struct {
	h *ffi.VerificationEvidence
}

// Type returns the evidence type tag.
func (e *VerificationEvidence) Type() string { return e.h.EvidenceType() }

// Object returns the object forming the evidence.
func (e *VerificationEvidence) Object() *object.Object {
	return ffi.ToObject(e.h.Object()).(*object.Object)
}

// VerificationParameter is a typed key/value parameter attached to a
// verification request.
type VerificationParameter struct {
	h *ffi.VerificationParameter
}

// Key returns the parameter key.
func (p *VerificationParameter) Key() string { return p.h.Key() }

// Value returns the parameter value decoded to a native Go type ([]byte,
// string, bool, int64, uint64, float64, [][]byte or []string).
func (p *VerificationParameter) Value() any { return p.h.Value() }

// Credentials returns the verifiable credentials in the response.
func (r *VerificationResponse) Credentials() []*credential.Verifiable {
	cs := r.h.Credentials()
	out := make([]*credential.Verifiable, len(cs))

	for i, c := range cs {
		out[i] = ffi.ToVerifiableCredential(c).(*credential.Verifiable)
	}

	return out
}

// NewVerificationResponse starts building a verification response.
func NewVerificationResponse() *VerificationResponseBuilder {
	return &VerificationResponseBuilder{h: ffi.NewVerificationResultBuilder()}
}

// Credential adds a verifiable credential to the response.
func (b *VerificationResponseBuilder) Credential(c *credential.Verifiable) *VerificationResponseBuilder {
	b.h.Credential(ffi.VerifiableCredentialOf(c))
	return b
}

// Finish finalizes the verification response.
func (b *VerificationResponseBuilder) Finish() (*VerificationResponse, error) {
	r, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &VerificationResponse{h: r}, nil
}
