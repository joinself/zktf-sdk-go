// Package exchange provides exchange (encryption) public keys.
package exchange

import "github.com/joinself/zktf-sdk-go/internal/ffi"

// PublicKey is an exchange public key.
type PublicKey struct {
	h *ffi.ExchangePublicKey
}

func init() {
	ffi.ExchangePublicKeyOf = func(o any) *ffi.ExchangePublicKey {
		return o.(*PublicKey).h
	}

	ffi.ToExchangePublicKey = func(h *ffi.ExchangePublicKey) any {
		return &PublicKey{h: h}
	}
}

// FromAddress decodes a hex address into an exchange key.
func FromAddress(hex string) (*PublicKey, error) {
	k, err := ffi.ExchangePublicKeyFromAddress(hex)
	if err != nil {
		return nil, err
	}

	return &PublicKey{h: k}, nil
}

// FromBytes constructs an exchange key from raw bytes.
func FromBytes(data []byte) (*PublicKey, error) {
	k, err := ffi.ExchangePublicKeyFromBytes(data)
	if err != nil {
		return nil, err
	}

	return &PublicKey{h: k}, nil
}

// String returns the hex encoded address.
func (p *PublicKey) String() string { return p.h.String() }

// Bytes returns the raw bytes of the key.
func (p *PublicKey) Bytes() []byte { return p.h.Bytes() }
