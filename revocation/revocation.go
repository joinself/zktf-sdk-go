// Package revocation provides revocation statements.
package revocation

import (
	"time"

	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// Statement is a revocation statement.
type Statement struct {
	h *ffi.RevocationStatement
}

func init() {
	ffi.RevocationStatementOf = func(o any) *ffi.RevocationStatement { return o.(*Statement).h }
	ffi.ToRevocationStatement = func(h *ffi.RevocationStatement) any { return &Statement{h: h} }
}

// Decode decodes an encoded revocation statement.
func Decode(data []byte) (*Statement, error) {
	s, err := ffi.RevocationStatementDecode(data)
	if err != nil {
		return nil, err
	}

	return &Statement{h: s}, nil
}

// Encode returns the encoded bytes of the statement.
func (s *Statement) Encode() ([]byte, error) { return s.h.Encode() }

// Issuer returns the statement's issuer.
func (s *Statement) Issuer() *signing.PublicKey {
	return ffi.ToSigningPublicKey(s.h.Issuer()).(*signing.PublicKey)
}

// Sequence returns the statement's sequence number.
func (s *Statement) Sequence() uint64 { return s.h.Sequence() }

// Timestamp returns the statement's timestamp.
func (s *Statement) Timestamp() time.Time { return time.Unix(s.h.Timestamp(), 0) }

// SignedBy reports whether the given signer signed the statement.
func (s *Statement) SignedBy(signer *signing.PublicKey) bool {
	return s.h.SignedBy(ffi.SigningPublicKeyOf(signer))
}

// RevokedAt returns the revocation timestamp for the given hash, or false if
// the hash is not in the statement.
func (s *Statement) RevokedAt(hash []byte) (time.Time, bool) {
	ts, ok := s.h.RevokedAt(hash)
	if !ok {
		return time.Time{}, false
	}

	return time.Unix(ts, 0), true
}

// Revocations returns the per-credential revocations in the statement.
func (s *Statement) Revocations() []*Entry {
	es := s.h.Revocations()
	out := make([]*Entry, len(es))

	for i, e := range es {
		out[i] = &Entry{h: e}
	}

	return out
}

// Signers returns the signer entries on the statement.
func (s *Statement) Signers() []*Signer {
	ss := s.h.Signers()
	out := make([]*Signer, len(ss))

	for i, x := range ss {
		out[i] = &Signer{h: x}
	}

	return out
}

// Entry is a single per-credential revocation in a statement.
type Entry struct {
	h *ffi.RevocationEntry
}

// Hash returns the 32-byte revocation hash.
func (e *Entry) Hash() []byte { return e.h.Hash() }

// Timestamp returns when the revocation occurred.
func (e *Entry) Timestamp() time.Time { return time.Unix(e.h.Timestamp(), 0) }

// Signer is a single signer entry on a statement.
type Signer struct {
	h *ffi.RevocationSigner
}

// Address returns the signer's signing public key.
func (s *Signer) Address() *signing.PublicKey {
	return ffi.ToSigningPublicKey(s.h.Address()).(*signing.PublicKey)
}

// Issued returns when the signer signed the statement.
func (s *Signer) Issued() time.Time { return time.Unix(s.h.Issued(), 0) }

// Builder builds a revocation statement.
type Builder struct {
	h *ffi.RevocationStatementBuilder
}

// NewBuilder starts building a revocation statement.
func NewBuilder() *Builder { return &Builder{h: ffi.NewRevocationStatementBuilder()} }

// Issuer sets the statement's issuer.
func (b *Builder) Issuer(issuer *signing.PublicKey) *Builder {
	b.h.Issuer(ffi.SigningPublicKeyOf(issuer))
	return b
}

// Sequence sets the statement's sequence number.
func (b *Builder) Sequence(seq uint64) *Builder {
	b.h.Sequence(seq)
	return b
}

// Timestamp sets the statement's timestamp.
func (b *Builder) Timestamp(t time.Time) *Builder {
	b.h.Timestamp(t.Unix())
	return b
}

// Revoke revokes a verifiable credential at the given time.
func (b *Builder) Revoke(c *credential.Verifiable, revokedAt time.Time) *Builder {
	b.h.Revoke(ffi.VerifiableCredentialOf(c), revokedAt.Unix())
	return b
}

// RevokeBy revokes a credential identified by its revocation hash.
func (b *Builder) RevokeBy(hash []byte, revokedAt time.Time) *Builder {
	b.h.RevokeBy(hash, revokedAt.Unix())
	return b
}

// SignWith records the signing key and issuance time.
func (b *Builder) SignWith(signer *signing.PublicKey, issuedAt time.Time) *Builder {
	b.h.SignWith(ffi.SigningPublicKeyOf(signer), issuedAt.Unix())
	return b
}

// Finish finalizes the unsigned statement.
func (b *Builder) Finish() (*Statement, error) {
	s, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Statement{h: s}, nil
}
