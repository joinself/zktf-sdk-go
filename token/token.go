// Package token provides zktf tokens.
package token

import (
	"time"

	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// Kind identifies the purpose of a token.
type Kind uint8

const (
	KindUnknown        Kind = Kind(ffi.TokenKindUnknown)
	KindAuthentication Kind = Kind(ffi.TokenKindAuthentication)
	KindSend           Kind = Kind(ffi.TokenKindSend)
	KindPush           Kind = Kind(ffi.TokenKindPush)
	KindSubscription   Kind = Kind(ffi.TokenKindSubscription)
	KindDelegation     Kind = Kind(ffi.TokenKindDelegation)
	KindIdentity       Kind = Kind(ffi.TokenKindIdentity)
)

// Token is a decoded zktf token.
type Token struct {
	h *ffi.Token
}

func init() {
	ffi.TokenOf = func(o any) *ffi.Token {
		return o.(*Token).h
	}
	ffi.ToToken = func(h *ffi.Token) any {
		return &Token{h: h}
	}
}

// Decode decodes an encoded token.
func Decode(data []byte) (*Token, error) {
	tk, err := ffi.TokenDecode(data)
	if err != nil {
		return nil, err
	}

	return &Token{h: tk}, nil
}

// Kind returns the token kind.
func (t *Token) Kind() Kind { return Kind(t.h.Kind()) }

// Issuer returns the address that issued the token, or nil.
func (t *Token) Issuer() *signing.PublicKey { return wrapSigningKey(t.h.Issuer()) }

// Bearer returns the address the token is intended for, or nil.
func (t *Token) Bearer() *signing.PublicKey { return wrapSigningKey(t.h.Bearer()) }

// Application returns the application the token is scoped to, or nil.
func (t *Token) Application() *signing.PublicKey { return wrapSigningKey(t.h.Application()) }

// Issued returns when the token was issued.
func (t *Token) Issued() time.Time { return time.Unix(t.h.Issued(), 0) }

// Expires returns when the token expires.
func (t *Token) Expires() time.Time { return time.Unix(t.h.Expires(), 0) }

// Nonce returns the token's 20-byte nonce.
func (t *Token) Nonce() []byte { return t.h.Nonce() }

// Encode returns the encoded token bytes.
func (t *Token) Encode() ([]byte, error) { return t.h.Encode() }

func wrapSigningKey(k *ffi.SigningPublicKey) *signing.PublicKey {
	if k == nil {
		return nil
	}

	return ffi.ToSigningPublicKey(k).(*signing.PublicKey)
}
