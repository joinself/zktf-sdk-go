// Package crypto provides MLS key packages and welcomes — primitives for
// establishing encrypted group sessions out of band (e.g. via discovery).
package crypto

import (
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// KeyPackage is an MLS key package distributed (e.g. inside a discovery
// request) so a counterparty can establish an encrypted session with you.
type KeyPackage struct {
	h *ffi.CryptoKeyPackage
}

// Welcome is an MLS welcome message inviting a member into a group.
type Welcome struct {
	h *ffi.CryptoWelcome
}

func init() {
	ffi.CryptoKeyPackageOf = func(o any) *ffi.CryptoKeyPackage {
		return o.(*KeyPackage).h
	}
	ffi.ToCryptoKeyPackage = func(h *ffi.CryptoKeyPackage) any {
		return &KeyPackage{h: h}
	}
	ffi.CryptoWelcomeOf = func(o any) *ffi.CryptoWelcome {
		return o.(*Welcome).h
	}
	ffi.ToCryptoWelcome = func(h *ffi.CryptoWelcome) any {
		return &Welcome{h: h}
	}
}

// FromAddress returns the signing address the key package is for.
func (k *KeyPackage) FromAddress() *signing.PublicKey {
	return ffi.ToSigningPublicKey(k.h.FromAddress()).(*signing.PublicKey)
}
