package credential

import (
	"time"

	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// Graph is a verified credential graph for a holder, derived by validating a
// set of presentations against a trusted-issuer registry. Build it via
// account.Account.CredentialGraphCreate.
type Graph struct {
	h *ffi.CredentialGraph
}

// RevocationProof carries the revocation entry and statement metadata that
// invalidated a credential.
type RevocationProof struct {
	h *ffi.RevocationProof
}

func init() {
	ffi.CredentialGraphOf = func(o any) *ffi.CredentialGraph { return o.(*Graph).h }
	ffi.ToCredentialGraph = func(h *ffi.CredentialGraph) any { return &Graph{h: h} }

	ffi.RevocationProofOf = func(o any) *ffi.RevocationProof { return o.(*RevocationProof).h }
	ffi.ToRevocationProof = func(h *ffi.RevocationProof) any { return &RevocationProof{h: h} }
}

// ValidCredentialsFor returns the holder's currently-valid credentials.
func (g *Graph) ValidCredentialsFor(holder *identity.Address) ([]*Verifiable, error) {
	cs, err := g.h.ValidCredentialsFor(ffi.DIDAddressOf(holder))
	if err != nil {
		return nil, err
	}

	out := make([]*Verifiable, len(cs))
	for i, c := range cs {
		out[i] = &Verifiable{h: c}
	}

	return out, nil
}

// RevokedCredentialsFor returns the holder's revoked credentials.
func (g *Graph) RevokedCredentialsFor(holder *identity.Address) ([]*Verifiable, error) {
	cs, err := g.h.RevokedCredentialsFor(ffi.DIDAddressOf(holder))
	if err != nil {
		return nil, err
	}

	out := make([]*Verifiable, len(cs))
	for i, c := range cs {
		out[i] = &Verifiable{h: c}
	}

	return out, nil
}

// ValidDocumentFor reports whether the document at the address is currently
// valid (no effective recovery / deactivation).
func (g *Graph) ValidDocumentFor(document *identity.Address) bool {
	return g.h.ValidDocumentFor(ffi.DIDAddressOf(document))
}

// BiometricAnchorHashFor returns the holder's 20-byte biometric anchor hash,
// or nil if not present.
func (g *Graph) BiometricAnchorHashFor(holder *identity.Address) []byte {
	return g.h.BiometricAnchorHashFor(ffi.DIDAddressOf(holder))
}

// RevocationProof returns the revocation proof for the given hash, or nil.
func (g *Graph) RevocationProof(revocationHash []byte) *RevocationProof {
	p := g.h.RevocationProofFor(revocationHash)
	if p == nil {
		return nil
	}

	return &RevocationProof{h: p}
}

// Issuer returns the issuer's signing public key.
func (p *RevocationProof) Issuer() *signing.PublicKey {
	return ffi.ToSigningPublicKey(p.h.Issuer()).(*signing.PublicKey)
}

// Sequence returns the originating statement's sequence number.
func (p *RevocationProof) Sequence() uint64 { return p.h.Sequence() }

// Timestamp returns the originating statement's time.
func (p *RevocationProof) Timestamp() time.Time { return time.Unix(p.h.Timestamp(), 0) }

// RevokedAt returns when the revocation took effect.
func (p *RevocationProof) RevokedAt() time.Time { return time.Unix(p.h.RevokedAt(), 0) }

// RevocationHash returns the 32-byte hash of the revoked entity.
func (p *RevocationProof) RevocationHash() []byte { return p.h.RevocationHash() }
