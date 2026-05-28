package message

import (
	"time"

	"github.com/joinself/zktf-sdk-go/crypto"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// DiscoveryRequest is a decoded discovery / out-of-band onboarding request.
type DiscoveryRequest struct {
	h *ffi.DiscoveryRequest
}

// DiscoveryRequestBuilder builds a discovery request.
type DiscoveryRequestBuilder struct {
	h *ffi.DiscoveryRequestBuilder
}

// DiscoveryResponse is a decoded discovery response.
type DiscoveryResponse struct {
	h *ffi.DiscoveryResponse
}

// DiscoveryResponseBuilder builds a discovery response.
type DiscoveryResponseBuilder struct {
	h *ffi.DiscoveryResponseBuilder
}

// DiscoveryRequestDecode decodes message content as a discovery request.
func DiscoveryRequestDecode(content *Content) (*DiscoveryRequest, error) {
	r, err := ffi.DiscoveryRequestFromContent(content.h)
	if err != nil {
		return nil, err
	}

	return &DiscoveryRequest{h: r}, nil
}

// DocumentAddress returns the requester's document address, or nil.
func (r *DiscoveryRequest) DocumentAddress() *signing.PublicKey {
	a := r.h.DocumentAddress()
	if a == nil {
		return nil
	}

	return ffi.ToSigningPublicKey(a).(*signing.PublicKey)
}

// FromAddress returns the requester's inbox address, or nil.
func (r *DiscoveryRequest) FromAddress() *signing.PublicKey {
	a := r.h.FromAddress()
	if a == nil {
		return nil
	}

	return ffi.ToSigningPublicKey(a).(*signing.PublicKey)
}

// KeyPackage returns the key package the recipient can use to establish an
// inbound session, or nil.
func (r *DiscoveryRequest) KeyPackage() *crypto.KeyPackage {
	k := r.h.KeyPackage()
	if k == nil {
		return nil
	}

	return ffi.ToCryptoKeyPackage(k).(*crypto.KeyPackage)
}

// Expires returns when the request expires.
func (r *DiscoveryRequest) Expires() time.Time { return time.Unix(r.h.Expires(), 0) }

// NewDiscoveryRequest starts building a discovery request.
func NewDiscoveryRequest() *DiscoveryRequestBuilder {
	return &DiscoveryRequestBuilder{h: ffi.NewDiscoveryRequestBuilder()}
}

// DocumentAddress sets an optional document address.
func (b *DiscoveryRequestBuilder) DocumentAddress(address *signing.PublicKey) *DiscoveryRequestBuilder {
	b.h.DocumentAddress(ffi.SigningPublicKeyOf(address))
	return b
}

// FromAddress sets the inbox address the discovery request is issued by.
func (b *DiscoveryRequestBuilder) FromAddress(address *signing.PublicKey) *DiscoveryRequestBuilder {
	b.h.FromAddress(ffi.SigningPublicKeyOf(address))
	return b
}

// KeyPackage attaches a key package the receiver can use to establish a session.
func (b *DiscoveryRequestBuilder) KeyPackage(kp *crypto.KeyPackage) *DiscoveryRequestBuilder {
	b.h.KeyPackage(ffi.CryptoKeyPackageOf(kp))
	return b
}

// Expires sets when the request expires.
func (b *DiscoveryRequestBuilder) Expires(t time.Time) *DiscoveryRequestBuilder {
	b.h.Expires(t.Unix())
	return b
}

// Finish finalizes the discovery request, ready to send.
func (b *DiscoveryRequestBuilder) Finish() (*Content, error) {
	c, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Content{h: c}, nil
}

// DiscoveryResponseDecode decodes message content as a discovery response.
func DiscoveryResponseDecode(content *Content) (*DiscoveryResponse, error) {
	r, err := ffi.DiscoveryResponseFromContent(content.h)
	if err != nil {
		return nil, err
	}

	return &DiscoveryResponse{h: r}, nil
}

// ResponseTo returns the id of the request being responded to.
func (r *DiscoveryResponse) ResponseTo() []byte { return r.h.ResponseTo() }

// Status returns the response status.
func (r *DiscoveryResponse) Status() ResponseStatus { return ResponseStatus(r.h.Status()) }

// ErrorMessage returns the response error message, or "".
func (r *DiscoveryResponse) ErrorMessage() string { return r.h.ErrorMessage() }

// NewDiscoveryResponse starts building a discovery response.
func NewDiscoveryResponse() *DiscoveryResponseBuilder {
	return &DiscoveryResponseBuilder{h: ffi.NewDiscoveryResponseBuilder()}
}

// ResponseTo sets the id of the request being responded to.
func (b *DiscoveryResponseBuilder) ResponseTo(requestID []byte) *DiscoveryResponseBuilder {
	b.h.ResponseTo(requestID)
	return b
}

// Status sets the response status.
func (b *DiscoveryResponseBuilder) Status(s ResponseStatus) *DiscoveryResponseBuilder {
	b.h.Status(ffi.ResponseStatus(s))
	return b
}

// ErrorMessage sets the response error message.
func (b *DiscoveryResponseBuilder) ErrorMessage(msg string) *DiscoveryResponseBuilder {
	b.h.ErrorMessage(msg)
	return b
}

// Finish finalizes the discovery response, ready to send.
func (b *DiscoveryResponseBuilder) Finish() (*Content, error) {
	c, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Content{h: c}, nil
}
