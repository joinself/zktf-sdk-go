package message

import (
	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/object"
)

// IdentitySigningRequest asks a counterparty to sign an identity-document operation.
type IdentitySigningRequest struct {
	h *ffi.IdentitySigningAction
}

// IdentitySigningRequestBuilder builds an identity-signing request.
type IdentitySigningRequestBuilder struct {
	h *ffi.IdentitySigningActionBuilder
}

// IdentitySigningResponse is the response to an identity-signing request.
type IdentitySigningResponse struct {
	h *ffi.IdentitySigningResult
}

// IdentitySigningResponseBuilder builds an identity-signing response.
type IdentitySigningResponseBuilder struct {
	h *ffi.IdentitySigningResultBuilder
}

func init() {
	ffi.IdentitySigningActionOf = func(o any) *ffi.IdentitySigningAction {
		return o.(*IdentitySigningRequest).h
	}
	ffi.ToIdentitySigningAction = func(h *ffi.IdentitySigningAction) any {
		return &IdentitySigningRequest{h: h}
	}

	ffi.IdentitySigningResultOf = func(o any) *ffi.IdentitySigningResult {
		return o.(*IdentitySigningResponse).h
	}
	ffi.ToIdentitySigningResult = func(h *ffi.IdentitySigningResult) any {
		return &IdentitySigningResponse{h: h}
	}
}

// DocumentAddress returns the document address the operation targets.
func (r *IdentitySigningRequest) DocumentAddress() *signing.PublicKey {
	return ffi.ToSigningPublicKey(r.h.DocumentAddress()).(*signing.PublicKey)
}

// Operation returns the hashgraph operation to be signed.
func (r *IdentitySigningRequest) Operation() *identity.Operation {
	return ffi.ToIdentityOperation(r.h.Operation()).(*identity.Operation)
}

// AsAction wraps this request into a generic Action.
func (r *IdentitySigningRequest) AsAction() *Action { return &Action{h: r.h.AsAction()} }

// NewIdentitySigningRequest starts building an identity-signing request.
func NewIdentitySigningRequest() *IdentitySigningRequestBuilder {
	return &IdentitySigningRequestBuilder{h: ffi.NewIdentitySigningActionBuilder()}
}

// DocumentAddress sets the document address.
func (b *IdentitySigningRequestBuilder) DocumentAddress(address *signing.PublicKey) *IdentitySigningRequestBuilder {
	b.h.DocumentAddress(ffi.SigningPublicKeyOf(address))
	return b
}

// Operation sets the operation to sign.
func (b *IdentitySigningRequestBuilder) Operation(op *identity.Operation) *IdentitySigningRequestBuilder {
	b.h.Operation(ffi.IdentityOperationOf(op))
	return b
}

// Finish finalizes the request.
func (b *IdentitySigningRequestBuilder) Finish() (*IdentitySigningRequest, error) {
	a, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &IdentitySigningRequest{h: a}, nil
}

// DocumentAddress returns the document address.
func (r *IdentitySigningResponse) DocumentAddress() *signing.PublicKey {
	return ffi.ToSigningPublicKey(r.h.DocumentAddress()).(*signing.PublicKey)
}

// Operation returns the signed operation.
func (r *IdentitySigningResponse) Operation() *identity.Operation {
	return ffi.ToIdentityOperation(r.h.Operation()).(*identity.Operation)
}

// Presentations returns the verifiable presentations attached.
func (r *IdentitySigningResponse) Presentations() []*credential.VerifiablePresentation {
	ps := r.h.Presentations()
	out := make([]*credential.VerifiablePresentation, len(ps))

	for i, p := range ps {
		out[i] = ffi.ToVerifiablePresentation(p).(*credential.VerifiablePresentation)
	}

	return out
}

// Assets returns the supporting object assets.
func (r *IdentitySigningResponse) Assets() []*object.Object {
	os := r.h.Assets()
	out := make([]*object.Object, len(os))

	for i, o := range os {
		out[i] = ffi.ToObject(o).(*object.Object)
	}

	return out
}

// NewIdentitySigningResponse starts building an identity-signing response.
func NewIdentitySigningResponse() *IdentitySigningResponseBuilder {
	return &IdentitySigningResponseBuilder{h: ffi.NewIdentitySigningResultBuilder()}
}

// DocumentAddress sets the document address.
func (b *IdentitySigningResponseBuilder) DocumentAddress(address *signing.PublicKey) *IdentitySigningResponseBuilder {
	b.h.DocumentAddress(ffi.SigningPublicKeyOf(address))
	return b
}

// Operation sets the signed operation.
func (b *IdentitySigningResponseBuilder) Operation(op *identity.Operation) *IdentitySigningResponseBuilder {
	b.h.Operation(ffi.IdentityOperationOf(op))
	return b
}

// Presentation adds a verifiable presentation.
func (b *IdentitySigningResponseBuilder) Presentation(p *credential.VerifiablePresentation) *IdentitySigningResponseBuilder {
	b.h.Presentation(ffi.VerifiablePresentationOf(p))
	return b
}

// Asset attaches a supporting object asset.
func (b *IdentitySigningResponseBuilder) Asset(o *object.Object) *IdentitySigningResponseBuilder {
	b.h.Asset(ffi.ObjectOf(o))
	return b
}

// Finish finalizes the response.
func (b *IdentitySigningResponseBuilder) Finish() (*IdentitySigningResponse, error) {
	r, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &IdentitySigningResponse{h: r}, nil
}
