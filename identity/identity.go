// Package identity provides resolved identity documents and identity operations.
package identity

import (
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/exchange"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// Document is the resolved key state of an identity.
type Document struct {
	h *ffi.IdentityDocument
}

// Operation is a hashgraph operation describing a change to an identity document.
type Operation struct {
	h *ffi.IdentityOperation
}

func init() {
	ffi.IdentityDocumentOf = func(o any) *ffi.IdentityDocument { return o.(*Document).h }
	ffi.ToIdentityDocument = func(h *ffi.IdentityDocument) any { return &Document{h: h} }

	ffi.IdentityOperationOf = func(o any) *ffi.IdentityOperation { return o.(*Operation).h }
	ffi.ToIdentityOperation = func(h *ffi.IdentityOperation) any { return &Operation{h: h} }
}

// SigningKeys returns all signing keys in the document.
func (d *Document) SigningKeys() []*signing.PublicKey {
	ks := d.h.SigningKeys()
	out := make([]*signing.PublicKey, len(ks))

	for i, k := range ks {
		out[i] = ffi.ToSigningPublicKey(k).(*signing.PublicKey)
	}

	return out
}

// ExchangeKeys returns all exchange keys in the document.
func (d *Document) ExchangeKeys() []*exchange.PublicKey {
	ks := d.h.ExchangeKeys()
	out := make([]*exchange.PublicKey, len(ks))

	for i, k := range ks {
		out[i] = ffi.ToExchangePublicKey(k).(*exchange.PublicKey)
	}

	return out
}

// SigningKeyValid reports whether the signing key is currently valid.
func (d *Document) SigningKeyValid(key *signing.PublicKey) bool {
	return d.h.SigningKeyValid(ffi.SigningPublicKeyOf(key))
}

// ExchangeKeyValid reports whether the exchange key is currently valid.
func (d *Document) ExchangeKeyValid(key *exchange.PublicKey) bool {
	return d.h.ExchangeKeyValid(ffi.ExchangePublicKeyOf(key))
}
