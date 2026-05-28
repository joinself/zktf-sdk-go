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

// TrustedIssuerRegistry wraps a zktf_trusted_issuer_registry handle: the set of
// issuers (and their per-credential-type authority windows) a verifier trusts.
type TrustedIssuerRegistry struct {
	ptr *C.zktf_trusted_issuer_registry
}

func newTrustedIssuerRegistry(ptr *C.zktf_trusted_issuer_registry) *TrustedIssuerRegistry {
	if ptr == nil {
		return nil
	}
	r := &TrustedIssuerRegistry{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_trusted_issuer_registry) {
		C.zktf_trusted_issuer_registry_destroy(ptr)
	}, r.ptr)
	return r
}

// NewTrustedIssuerRegistry initializes an empty registry.
func NewTrustedIssuerRegistry() *TrustedIssuerRegistry {
	return newTrustedIssuerRegistry(C.zktf_trusted_issuer_registry_init())
}

// DefaultProductionTrustedIssuerRegistry returns the registry of self's
// production issuers.
func DefaultProductionTrustedIssuerRegistry() *TrustedIssuerRegistry {
	return newTrustedIssuerRegistry(C.zktf_trusted_issuer_registry_default_production())
}

// DefaultSandboxTrustedIssuerRegistry returns the registry of self's sandbox
// issuers.
func DefaultSandboxTrustedIssuerRegistry() *TrustedIssuerRegistry {
	return newTrustedIssuerRegistry(C.zktf_trusted_issuer_registry_default_sandbox())
}

// DefaultIssuedCredentialTypes returns the default credential types issued by self.
func DefaultIssuedCredentialTypes() []string {
	c := C.zktf_trusted_issuer_registry_default_credential_types()
	if c == nil {
		return nil
	}
	defer C.zktf_collection_string_buffer_destroy(c)
	n := int(C.zktf_collection_string_buffer_len(c))
	out := make([]string, n)
	for i := 0; i < n; i++ {
		buf := C.zktf_collection_string_buffer_at(c, C.size_t(i))
		if buf == nil {
			continue
		}
		out[i] = C.GoString(C.zktf_string_buffer_ptr(buf))
	}
	return out
}

// DefaultIssuerEpoch returns the default epoch (unix seconds) from when
// credentials issued by self are valid.
func DefaultIssuerEpoch() int64 {
	return int64(C.zktf_trusted_issuer_registry_default_issuer_epoch())
}

// IssuerAdd adds an issuer to the registry. Returns false if it was already present.
func (r *TrustedIssuerRegistry) IssuerAdd(issuer *DIDAddress) bool {
	return bool(C.zktf_trusted_issuer_registry_issuer_add(r.ptr, issuer.ptr))
}

// IssuerRemove removes an issuer. Returns false if it was not present.
func (r *TrustedIssuerRegistry) IssuerRemove(issuer *DIDAddress) bool {
	return bool(C.zktf_trusted_issuer_registry_issuer_remove(r.ptr, issuer.ptr))
}

// AuthorityGrant grants the issuer authority over a credential type from a
// time, optionally bounded by a revocation time (pass 0 for no revocation).
func (r *TrustedIssuerRegistry) AuthorityGrant(issuer *DIDAddress, credentialType string, grantedUnix int64, revokedUnix int64) error {
	cType := cstring(credentialType)
	defer free(unsafe.Pointer(cType))
	var revokedArg *C.int64_t
	rev := C.int64_t(revokedUnix)
	if revokedUnix != 0 {
		revokedArg = &rev
	}
	return status(C.zktf_trusted_issuer_registry_authority_grant(
		r.ptr, issuer.ptr, cType, C.int64_t(grantedUnix), revokedArg,
	))
}

// AuthorityRevoke marks the issuer's authority over a credential type revoked
// from the given time.
func (r *TrustedIssuerRegistry) AuthorityRevoke(issuer *DIDAddress, credentialType string, revokedUnix int64) error {
	cType := cstring(credentialType)
	defer free(unsafe.Pointer(cType))
	return status(C.zktf_trusted_issuer_registry_authority_revoke(
		r.ptr, issuer.ptr, cType, C.int64_t(revokedUnix),
	))
}

// AuthorityFor returns the credential types the issuer is authorized for.
func (r *TrustedIssuerRegistry) AuthorityFor(issuer *DIDAddress) ([]string, error) {
	var c *C.zktf_collection_string_buffer
	if err := status(C.zktf_trusted_issuer_registry_authority_for(r.ptr, issuer.ptr, &c)); err != nil {
		return nil, err
	}
	defer C.zktf_collection_string_buffer_destroy(c)
	n := int(C.zktf_collection_string_buffer_len(c))
	out := make([]string, n)
	for i := 0; i < n; i++ {
		buf := C.zktf_collection_string_buffer_at(c, C.size_t(i))
		if buf == nil {
			continue
		}
		out[i] = C.GoString(C.zktf_string_buffer_ptr(buf))
	}
	return out, nil
}

// AuthorityAt reports whether the issuer was authorized for the credential
// type at the given timestamp.
func (r *TrustedIssuerRegistry) AuthorityAt(issuer *DIDAddress, credentialType string, issuedUnix int64) bool {
	cType := cstring(credentialType)
	defer free(unsafe.Pointer(cType))
	return bool(C.zktf_trusted_issuer_registry_authority_at(
		r.ptr, issuer.ptr, cType, C.int64_t(issuedUnix),
	))
}

// Issuers returns the DID addresses of all issuers in the registry.
func (r *TrustedIssuerRegistry) Issuers() []*DIDAddress {
	return didAddressesFrom(C.zktf_trusted_issuer_registry_issuers(r.ptr))
}

// CredentialGraph wraps a zktf_credential_graph handle — the verified state of a
// holder's credentials, derived from a set of presentations against a registry.
type CredentialGraph struct {
	ptr *C.zktf_credential_graph
}

func newCredentialGraph(ptr *C.zktf_credential_graph) *CredentialGraph {
	if ptr == nil {
		return nil
	}
	g := &CredentialGraph{ptr: ptr}
	runtime.AddCleanup(g, func(ptr *C.zktf_credential_graph) {
		C.zktf_credential_graph_destroy(ptr)
	}, g.ptr)
	return g
}

// ValidCredentialsFor returns the holder's currently-valid credentials.
func (g *CredentialGraph) ValidCredentialsFor(holder *DIDAddress) ([]*VerifiableCredential, error) {
	var c *C.zktf_collection_verifiable_credential
	if err := status(C.zktf_credential_graph_valid_credentials_for(g.ptr, holder.ptr, &c)); err != nil {
		return nil, err
	}
	return verifiableCredentialsFrom(c), nil
}

// RevokedCredentialsFor returns the holder's revoked credentials.
func (g *CredentialGraph) RevokedCredentialsFor(holder *DIDAddress) ([]*VerifiableCredential, error) {
	var c *C.zktf_collection_verifiable_credential
	if err := status(C.zktf_credential_graph_revoked_credentials_for(g.ptr, holder.ptr, &c)); err != nil {
		return nil, err
	}
	return verifiableCredentialsFrom(c), nil
}

// ValidDocumentFor reports whether the document at the address is currently
// valid (no recovery / deactivation effective).
func (g *CredentialGraph) ValidDocumentFor(document *DIDAddress) bool {
	return bool(C.zktf_credential_graph_valid_document_for(g.ptr, document.ptr))
}

// BiometricAnchorHashFor returns the holder's 20-byte biometric anchor hash, or nil.
func (g *CredentialGraph) BiometricAnchorHashFor(holder *DIDAddress) []byte {
	return goBytesFromBuffer(C.zktf_credential_graph_biometric_anchor_hash_for(g.ptr, holder.ptr))
}

// RevocationProofFor returns the revocation proof for the given hash, or nil
// if no revocation has been recorded.
func (g *CredentialGraph) RevocationProofFor(revocationHash []byte) *RevocationProof {
	buf, length := cbytes(revocationHash)
	defer free(unsafe.Pointer(buf))
	return newRevocationProof(C.zktf_credential_graph_revocation_proof_for(g.ptr, buf, C.uintptr_t(length)))
}

// ValidAuthenticationFor reports whether the given pairwise identity has signed
// the supplied challenge with currently-valid keys.
func (g *CredentialGraph) ValidAuthenticationFor(identity *PairwiseIdentity, challenge []byte) bool {
	buf, _ := cbytes(challenge)
	defer free(unsafe.Pointer(buf))
	return bool(C.zktf_credential_graph_valid_authentication_for(g.ptr, identity.ptr, buf))
}

// CredentialGraphCreate builds a credential graph for a holder by validating
// the given presentations against the trusted-issuer registry, via callback.
func (a *Account) CredentialGraphCreate(registry *TrustedIssuerRegistry, presentations []*VerifiablePresentation, timeout time.Duration) (*CredentialGraph, error) {
	in := verifiablePresentationCollection(presentations)
	defer C.zktf_collection_verifiable_presentation_destroy(in)

	fut := C.zktf_account_credential_graph_create(a.ptr, registry.ptr, in)

	return AwaitCredentialGraph(fut, timeout)
}

// verifiablePresentationCollection builds a zktf_collection_verifiable_presentation
// from Go presentations. Caller destroys.
func verifiablePresentationCollection(presentations []*VerifiablePresentation) *C.zktf_collection_verifiable_presentation {
	c := C.zktf_collection_verifiable_presentation_init()
	for _, p := range presentations {
		C.zktf_collection_verifiable_presentation_append(c, p.ptr)
	}
	return c
}

// RevocationProof wraps a zktf_revocation_proof handle — a single revocation
// entry plus the signers and statement metadata that produced it.
type RevocationProof struct {
	ptr *C.zktf_revocation_proof
}

func newRevocationProof(ptr *C.zktf_revocation_proof) *RevocationProof {
	if ptr == nil {
		return nil
	}
	p := &RevocationProof{ptr: ptr}
	runtime.AddCleanup(p, func(ptr *C.zktf_revocation_proof) {
		C.zktf_revocation_proof_destroy(ptr)
	}, p.ptr)
	return p
}

// Issuer returns the issuer's signing public key.
func (p *RevocationProof) Issuer() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_revocation_proof_issuer(p.ptr))
}

// Sequence returns the originating statement's sequence number.
func (p *RevocationProof) Sequence() uint64 {
	return uint64(C.zktf_revocation_proof_sequence(p.ptr))
}

// Timestamp returns the originating statement's timestamp (unix seconds).
func (p *RevocationProof) Timestamp() int64 {
	return int64(C.zktf_revocation_proof_timestamp(p.ptr))
}

// RevokedAt returns when the revocation took effect (unix seconds).
func (p *RevocationProof) RevokedAt() int64 {
	return int64(C.zktf_revocation_proof_revoked(p.ptr))
}

// RevocationHash returns the 32-byte hash of the revoked entity.
func (p *RevocationProof) RevocationHash() []byte {
	return C.GoBytes(unsafe.Pointer(C.zktf_revocation_proof_revocation_hash(p.ptr)), revocationHashLen)
}

// Signers returns the signer entries of the originating statement.
func (p *RevocationProof) Signers() []*RevocationSigner {
	c := C.zktf_revocation_proof_signers(p.ptr)
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
