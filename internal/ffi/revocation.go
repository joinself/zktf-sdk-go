package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import (
	"runtime"
	"time"
	"unsafe"
)

const revocationHashLen = 32

// RevocationStatement wraps a zktf_revocation_statement handle.
type RevocationStatement struct {
	ptr *C.zktf_revocation_statement
}

func newRevocationStatement(ptr *C.zktf_revocation_statement) *RevocationStatement {
	if ptr == nil {
		return nil
	}
	s := &RevocationStatement{ptr: ptr}
	runtime.AddCleanup(s, func(ptr *C.zktf_revocation_statement) {
		C.zktf_revocation_statement_destroy(ptr)
	}, s.ptr)
	return s
}

// RevocationStatementDecode decodes an encoded revocation statement.
func RevocationStatementDecode(data []byte) (*RevocationStatement, error) {
	buf, length := cbytes(data)
	defer free(unsafe.Pointer(buf))

	var ptr *C.zktf_revocation_statement
	if err := status(C.zktf_revocation_statement_decode(buf, length, &ptr)); err != nil {
		return nil, err
	}
	return newRevocationStatement(ptr), nil
}

// Encode returns the encoded bytes of the statement.
func (s *RevocationStatement) Encode() ([]byte, error) {
	var buf *C.zktf_bytes_buffer
	if err := status(C.zktf_revocation_statement_encode(s.ptr, &buf)); err != nil {
		return nil, err
	}
	return goBytesFromBuffer(buf), nil
}

// Issuer returns the issuer's signing public key.
func (s *RevocationStatement) Issuer() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_revocation_statement_issuer(s.ptr))
}

// Sequence returns the statement's sequence number.
func (s *RevocationStatement) Sequence() uint64 {
	return uint64(C.zktf_revocation_statement_sequence(s.ptr))
}

// Timestamp returns the statement's unix timestamp (seconds).
func (s *RevocationStatement) Timestamp() int64 {
	return int64(C.zktf_revocation_statement_timestamp(s.ptr))
}

// SignedBy reports whether the given signer's signature appears on the statement.
func (s *RevocationStatement) SignedBy(signer *SigningPublicKey) bool {
	return bool(C.zktf_revocation_statement_signed_by(s.ptr, signer.ptr))
}

// RevokedAt returns the revocation timestamp for the given revocation hash, or
// false if the hash is not in the statement.
func (s *RevocationStatement) RevokedAt(hash []byte) (int64, bool) {
	buf, length := cbytes(hash)
	defer free(unsafe.Pointer(buf))
	var ts C.int64_t
	ok := bool(C.zktf_revocation_statement_revoked_at(s.ptr, buf, C.uintptr_t(length), &ts))
	return int64(ts), ok
}

// Revocations returns the per-credential revocations in the statement.
func (s *RevocationStatement) Revocations() []*RevocationEntry {
	c := C.zktf_revocation_statement_revocations(s.ptr)
	if c == nil {
		return nil
	}
	defer C.zktf_collection_revocation_statement_revocation_destroy(c)
	n := int(C.zktf_collection_revocation_statement_revocation_len(c))
	out := make([]*RevocationEntry, n)
	for i := 0; i < n; i++ {
		out[i] = newRevocationEntry(C.zktf_collection_revocation_statement_revocation_at(c, C.size_t(i)))
	}
	return out
}

// Signers returns the signer entries in the statement.
func (s *RevocationStatement) Signers() []*RevocationSigner {
	c := C.zktf_revocation_statement_signers(s.ptr)
	if c == nil {
		return nil
	}
	defer C.zktf_collection_revocation_statement_signer_destroy(c)
	n := int(C.zktf_collection_revocation_statement_signer_len(c))
	out := make([]*RevocationSigner, n)
	for i := 0; i < n; i++ {
		out[i] = newRevocationSigner(C.zktf_collection_revocation_statement_signer_at(c, C.size_t(i)))
	}
	return out
}

// RevocationEntry is a single per-credential revocation entry in a statement.
type RevocationEntry struct {
	ptr *C.zktf_revocation_statement_revocation
}

func newRevocationEntry(ptr *C.zktf_revocation_statement_revocation) *RevocationEntry {
	if ptr == nil {
		return nil
	}
	e := &RevocationEntry{ptr: ptr}
	runtime.AddCleanup(e, func(ptr *C.zktf_revocation_statement_revocation) {
		C.zktf_revocation_statement_revocation_destroy(ptr)
	}, e.ptr)
	return e
}

// Hash returns the 32-byte revocation hash.
func (e *RevocationEntry) Hash() []byte {
	return C.GoBytes(unsafe.Pointer(C.zktf_revocation_statement_revocation_hash(e.ptr)), revocationHashLen)
}

// Timestamp returns the revocation timestamp.
func (e *RevocationEntry) Timestamp() int64 {
	return int64(C.zktf_revocation_statement_revocation_timestamp(e.ptr))
}

// RevocationSigner is a single signer entry on a revocation statement.
type RevocationSigner struct {
	ptr *C.zktf_revocation_statement_signer
}

func newRevocationSigner(ptr *C.zktf_revocation_statement_signer) *RevocationSigner {
	if ptr == nil {
		return nil
	}
	s := &RevocationSigner{ptr: ptr}
	runtime.AddCleanup(s, func(ptr *C.zktf_revocation_statement_signer) {
		C.zktf_revocation_statement_signer_destroy(ptr)
	}, s.ptr)
	return s
}

// Address returns the signer's signing public key.
func (s *RevocationSigner) Address() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_revocation_statement_signer_address(s.ptr))
}

// Issued returns the unix timestamp the signature was issued.
func (s *RevocationSigner) Issued() int64 {
	return int64(C.zktf_revocation_statement_signer_issued(s.ptr))
}

// RevocationStatementBuilder builds a revocation statement.
type RevocationStatementBuilder struct {
	ptr *C.zktf_revocation_statement_builder
}

// NewRevocationStatementBuilder initializes a revocation statement builder.
func NewRevocationStatementBuilder() *RevocationStatementBuilder {
	ptr := C.zktf_revocation_statement_builder_init()
	b := &RevocationStatementBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_revocation_statement_builder) {
		C.zktf_revocation_statement_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// Issuer sets the statement's issuer.
func (b *RevocationStatementBuilder) Issuer(issuer *SigningPublicKey) *RevocationStatementBuilder {
	C.zktf_revocation_statement_builder_issuer(b.ptr, issuer.ptr)
	return b
}

// Sequence sets the statement's sequence number.
func (b *RevocationStatementBuilder) Sequence(seq uint64) *RevocationStatementBuilder {
	C.zktf_revocation_statement_builder_sequence(b.ptr, C.uint64_t(seq))
	return b
}

// Timestamp sets the statement's timestamp.
func (b *RevocationStatementBuilder) Timestamp(unix int64) *RevocationStatementBuilder {
	C.zktf_revocation_statement_builder_timestamp(b.ptr, C.int64_t(unix))
	return b
}

// Revoke revokes a verifiable credential at the given timestamp.
func (b *RevocationStatementBuilder) Revoke(credential *VerifiableCredential, revokedAtUnix int64) *RevocationStatementBuilder {
	C.zktf_revocation_statement_builder_revoke(b.ptr, credential.ptr, C.int64_t(revokedAtUnix))
	return b
}

// RevokeBy revokes a credential identified by its 32-byte revocation hash.
func (b *RevocationStatementBuilder) RevokeBy(hash []byte, revokedAtUnix int64) *RevocationStatementBuilder {
	buf, length := cbytes(hash)
	defer free(unsafe.Pointer(buf))
	C.zktf_revocation_statement_builder_revoke_by(b.ptr, buf, length, C.int64_t(revokedAtUnix))
	return b
}

// SignWith records the signing key and issuance time for the statement.
func (b *RevocationStatementBuilder) SignWith(signer *SigningPublicKey, issuedAtUnix int64) *RevocationStatementBuilder {
	C.zktf_revocation_statement_builder_sign_with(b.ptr, signer.ptr, C.int64_t(issuedAtUnix))
	return b
}

// Finish finalizes the revocation statement.
func (b *RevocationStatementBuilder) Finish() (*RevocationStatement, error) {
	var out *C.zktf_revocation_statement
	if err := status(C.zktf_revocation_statement_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newRevocationStatement(out), nil
}

// RevocationSign signs an unsigned revocation statement with the account's keys.
func (a *Account) RevocationSign(statement *RevocationStatement) error {
	return status(C.zktf_account_revocation_sign(a.ptr, statement.ptr))
}

// RevocationRevoke publishes a signed revocation statement via callback.
func (a *Account) RevocationRevoke(statement *RevocationStatement, timeout time.Duration) error {
	fut := C.zktf_account_revocation_revoke(a.ptr, statement.ptr)

	return AwaitStatus(fut, timeout)
}
