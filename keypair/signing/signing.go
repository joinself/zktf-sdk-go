// Package signing provides signing public keys (zktf addresses).
package signing

import "github.com/joinself/zktf-sdk-go/internal/ffi"

// PublicKey is a signing public key / address.
type PublicKey struct {
	h *ffi.SigningPublicKey
}

func init() {
	ffi.SigningPublicKeyOf = func(o any) *ffi.SigningPublicKey {
		return o.(*PublicKey).h
	}

	ffi.ToSigningPublicKey = func(h *ffi.SigningPublicKey) any {
		return &PublicKey{h: h}
	}
}

// FromAddress decodes a hex address into a public key.
func FromAddress(hex string) (*PublicKey, error) {
	k, err := ffi.SigningPublicKeyFromAddress(hex)
	if err != nil {
		return nil, err
	}

	return &PublicKey{h: k}, nil
}

// FromBytes constructs a public key from its raw bytes.
func FromBytes(data []byte) (*PublicKey, error) {
	k, err := ffi.SigningPublicKeyFromBytes(data)
	if err != nil {
		return nil, err
	}

	return &PublicKey{h: k}, nil
}

// String returns the hex encoded address.
func (p *PublicKey) String() string { return p.h.String() }

// Bytes returns the raw bytes of the public key.
func (p *PublicKey) Bytes() []byte { return p.h.Bytes() }

// Verify reports whether signature is a valid signature of message by this key.
func (p *PublicKey) Verify(message, signature []byte) bool {
	return p.h.Verify(message, signature)
}

// Matches reports whether two public keys are equal.
func (p *PublicKey) Matches(other *PublicKey) bool {
	if other == nil {
		return false
	}

	return p.h.Matches(other.h)
}
