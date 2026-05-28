package message

import (
	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
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

// Finish finalizes the verification request.
func (b *VerificationRequestBuilder) Finish() (*VerificationRequest, error) {
	a, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &VerificationRequest{h: a}, nil
}

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
