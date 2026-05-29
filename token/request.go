package token

import (
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/exchange"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// Request is a validated token request ready to be issued into a token by
// Account.TokenIssue.
type Request struct {
	h *ffi.TokenRequest
}

func init() {
	ffi.TokenRequestOf = func(o any) *ffi.TokenRequest { return o.(*Request).h }
	ffi.ToTokenRequest = func(h *ffi.TokenRequest) any { return &Request{h: h} }
}

// PushBuilder builds a push token request.
type PushBuilder struct {
	h *ffi.PushTokenBuilder
}

// NewPushBuilder creates a push token builder.
func NewPushBuilder() *PushBuilder {
	return &PushBuilder{h: ffi.NewPushTokenBuilder()}
}

// ForAddress sets the local group address the token authorizes notifications for.
func (b *PushBuilder) ForAddress(address *signing.PublicKey) *PushBuilder {
	b.h.ForAddress(ffi.SigningPublicKeyOf(address))
	return b
}

// ProviderAddress sets the exchange public key of the push provider.
func (b *PushBuilder) ProviderAddress(address *exchange.PublicKey) *PushBuilder {
	b.h.ProviderAddress(ffi.ExchangePublicKeyOf(address))
	return b
}

// Delegatable allows the bearer to further delegate the issued token.
func (b *PushBuilder) Delegatable(delegatable bool) *PushBuilder {
	b.h.Delegatable(delegatable)
	return b
}

// Finish validates the configured fields and returns a token request.
func (b *PushBuilder) Finish() (*Request, error) {
	r, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return ffi.ToTokenRequest(r).(*Request), nil
}
