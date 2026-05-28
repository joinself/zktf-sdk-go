package message

import (
	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/object"
)

// DevicePairingRequest asks a counterparty to pair a device into an identity
// document with a given role bitmask.
type DevicePairingRequest struct {
	h *ffi.DevicePairingAction
}

// DevicePairingRequestBuilder builds a device-pairing request.
type DevicePairingRequestBuilder struct {
	h *ffi.DevicePairingActionBuilder
}

// DevicePairingResponse is the response to a device-pairing request.
type DevicePairingResponse struct {
	h *ffi.DevicePairingResult
}

// DevicePairingResponseBuilder builds a device-pairing response.
type DevicePairingResponseBuilder struct {
	h *ffi.DevicePairingResultBuilder
}

func init() {
	ffi.DevicePairingActionOf = func(o any) *ffi.DevicePairingAction {
		return o.(*DevicePairingRequest).h
	}
	ffi.ToDevicePairingAction = func(h *ffi.DevicePairingAction) any {
		return &DevicePairingRequest{h: h}
	}

	ffi.DevicePairingResultOf = func(o any) *ffi.DevicePairingResult {
		return o.(*DevicePairingResponse).h
	}
	ffi.ToDevicePairingResult = func(h *ffi.DevicePairingResult) any {
		return &DevicePairingResponse{h: h}
	}
}

// Address returns the signing address to pair.
func (r *DevicePairingRequest) Address() *signing.PublicKey {
	return ffi.ToSigningPublicKey(r.h.Address()).(*signing.PublicKey)
}

// Roles returns the requested role bitmask.
func (r *DevicePairingRequest) Roles() uint64 { return r.h.Roles() }

// AsAction wraps this request into a generic Action.
func (r *DevicePairingRequest) AsAction() *Action { return &Action{h: r.h.AsAction()} }

// NewDevicePairingRequest starts building a device-pairing request.
func NewDevicePairingRequest() *DevicePairingRequestBuilder {
	return &DevicePairingRequestBuilder{h: ffi.NewDevicePairingActionBuilder()}
}

// Address sets the signing address to pair.
func (b *DevicePairingRequestBuilder) Address(address *signing.PublicKey) *DevicePairingRequestBuilder {
	b.h.Address(ffi.SigningPublicKeyOf(address))
	return b
}

// Roles sets the role bitmask.
func (b *DevicePairingRequestBuilder) Roles(roles uint64) *DevicePairingRequestBuilder {
	b.h.Roles(roles)
	return b
}

// Finish finalizes the request.
func (b *DevicePairingRequestBuilder) Finish() (*DevicePairingRequest, error) {
	a, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &DevicePairingRequest{h: a}, nil
}

// DocumentAddress returns the document address.
func (r *DevicePairingResponse) DocumentAddress() *signing.PublicKey {
	return ffi.ToSigningPublicKey(r.h.DocumentAddress()).(*signing.PublicKey)
}

// Operation returns the signed operation that paired the device.
func (r *DevicePairingResponse) Operation() *identity.Operation {
	return ffi.ToIdentityOperation(r.h.Operation()).(*identity.Operation)
}

// Presentations returns the presentations attached.
func (r *DevicePairingResponse) Presentations() []*credential.VerifiablePresentation {
	ps := r.h.Presentations()
	out := make([]*credential.VerifiablePresentation, len(ps))

	for i, p := range ps {
		out[i] = ffi.ToVerifiablePresentation(p).(*credential.VerifiablePresentation)
	}

	return out
}

// Assets returns the supporting object assets.
func (r *DevicePairingResponse) Assets() []*object.Object {
	os := r.h.Assets()
	out := make([]*object.Object, len(os))

	for i, o := range os {
		out[i] = ffi.ToObject(o).(*object.Object)
	}

	return out
}

// NewDevicePairingResponse starts building a device-pairing response.
func NewDevicePairingResponse() *DevicePairingResponseBuilder {
	return &DevicePairingResponseBuilder{h: ffi.NewDevicePairingResultBuilder()}
}

// DocumentAddress sets the document address.
func (b *DevicePairingResponseBuilder) DocumentAddress(address *signing.PublicKey) *DevicePairingResponseBuilder {
	b.h.DocumentAddress(ffi.SigningPublicKeyOf(address))
	return b
}

// Operation sets the signed operation.
func (b *DevicePairingResponseBuilder) Operation(op *identity.Operation) *DevicePairingResponseBuilder {
	b.h.Operation(ffi.IdentityOperationOf(op))
	return b
}

// Presentation adds a verifiable presentation.
func (b *DevicePairingResponseBuilder) Presentation(p *credential.VerifiablePresentation) *DevicePairingResponseBuilder {
	b.h.Presentation(ffi.VerifiablePresentationOf(p))
	return b
}

// Asset attaches a supporting object asset.
func (b *DevicePairingResponseBuilder) Asset(o *object.Object) *DevicePairingResponseBuilder {
	b.h.Asset(ffi.ObjectOf(o))
	return b
}

// Finish finalizes the response.
func (b *DevicePairingResponseBuilder) Finish() (*DevicePairingResponse, error) {
	r, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &DevicePairingResponse{h: r}, nil
}
