package credential

import (
	"time"

	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// Exchange records a single credential shared with an address, as returned by
// Account.CredentialExchangeLog.
type Exchange struct {
	h *ffi.CredentialExchange
}

func init() {
	ffi.CredentialExchangeOf = func(o any) *ffi.CredentialExchange { return o.(*Exchange).h }
	ffi.ToCredentialExchange = func(h *ffi.CredentialExchange) any { return &Exchange{h: h} }
}

// WithAddress returns the address the credential was exchanged with.
func (e *Exchange) WithAddress() *signing.PublicKey {
	return ffi.ToSigningPublicKey(e.h.WithAddress()).(*signing.PublicKey)
}

// ContentHash returns the hash of the exchanged credential content.
func (e *Exchange) ContentHash() []byte {
	return e.h.ContentHash()
}

// SharedAt returns the time the credential was shared.
func (e *Exchange) SharedAt() time.Time {
	return time.Unix(e.h.SharedAt(), 0)
}
