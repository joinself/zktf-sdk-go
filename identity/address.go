package identity

import (
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// Address is a decentralized identifier (DID) address.
type Address struct {
	h *ffi.DIDAddress
}

func init() {
	ffi.DIDAddressOf = func(o any) *ffi.DIDAddress { return o.(*Address).h }
	ffi.ToDIDAddress = func(h *ffi.DIDAddress) any { return &Address{h: h} }
}

// AddressKey builds a key-method DID address from a signing key.
func AddressKey(key *signing.PublicKey) *Address {
	return &Address{h: ffi.DIDAddressKey(ffi.SigningPublicKeyOf(key))}
}

// ParseAddress decodes a DID string into an address.
func ParseAddress(did string) (*Address, error) {
	a, err := ffi.DIDAddressDecode(did)
	if err != nil {
		return nil, err
	}

	return &Address{h: a}, nil
}

// Key returns the signing public key embedded in the address.
func (a *Address) Key() *signing.PublicKey {
	return ffi.ToSigningPublicKey(a.h.Address()).(*signing.PublicKey)
}

// String returns the encoded DID string.
func (a *Address) String() string { return a.h.String() }
