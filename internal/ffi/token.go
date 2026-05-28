package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// TokenKind mirrors zktf_token_kind.
type TokenKind uint8

const (
	TokenKindUnknown        TokenKind = C.TOKEN_KIND_UNKNOWN
	TokenKindAuthentication TokenKind = C.TOKEN_KIND_AUTHENTICATION
	TokenKindSend           TokenKind = C.TOKEN_KIND_SEND
	TokenKindPush           TokenKind = C.TOKEN_KIND_PUSH
	TokenKindSubscription   TokenKind = C.TOKEN_KIND_SUBSCRIPTION
	TokenKindDelegation     TokenKind = C.TOKEN_KIND_DELEGATION
	TokenKindIdentity       TokenKind = C.TOKEN_KIND_IDENTITY
)

const tokenNonceLen = 20

// Token wraps a zktf_token handle.
type Token struct {
	ptr *C.zktf_token
}

func newToken(ptr *C.zktf_token) *Token {
	if ptr == nil {
		return nil
	}
	t := &Token{ptr: ptr}
	runtime.AddCleanup(t, func(ptr *C.zktf_token) {
		C.zktf_token_destroy(ptr)
	}, t.ptr)
	return t
}

// TokenDecode decodes an encoded token.
func TokenDecode(data []byte) (*Token, error) {
	buf, length := cbytes(data)
	defer free(unsafe.Pointer(buf))

	var ptr *C.zktf_token
	if err := status(C.zktf_token_decode(buf, length, &ptr)); err != nil {
		return nil, err
	}
	return newToken(ptr), nil
}

// Kind returns the token kind.
func (t *Token) Kind() TokenKind { return TokenKind(C.zktf_token_get_kind(t.ptr)) }

// Issuer returns the address that issued the token, or nil.
func (t *Token) Issuer() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_token_issuer(t.ptr))
}

// Bearer returns the address the token is intended for, or nil.
func (t *Token) Bearer() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_token_bearer(t.ptr))
}

// Application returns the application the token is scoped to, or nil.
func (t *Token) Application() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_token_application(t.ptr))
}

// Issued returns the unix timestamp (seconds) the token was issued.
func (t *Token) Issued() int64 { return int64(C.zktf_token_issued(t.ptr)) }

// Expires returns the unix timestamp (seconds) the token expires.
func (t *Token) Expires() int64 { return int64(C.zktf_token_expires(t.ptr)) }

// Nonce returns the 20-byte random nonce in the token header.
func (t *Token) Nonce() []byte {
	buf := C.malloc(tokenNonceLen)
	defer C.free(buf)
	n := C.zktf_token_nonce(t.ptr, (*C.uint8_t)(buf), tokenNonceLen)
	return C.GoBytes(buf, C.int(n))
}

// Encode returns the encoded token bytes.
func (t *Token) Encode() ([]byte, error) {
	var buf *C.zktf_bytes_buffer
	if err := status(C.zktf_token_encode(t.ptr, &buf)); err != nil {
		return nil, err
	}
	return goBytesFromBuffer(buf), nil
}
