// Package trust provides the trusted-issuer registry.
package trust

import (
	"time"

	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
)

// Registry is a set of issuers and their per-credential-type authority windows.
type Registry struct {
	h *ffi.TrustedIssuerRegistry
}

func init() {
	ffi.TrustedIssuerRegistryOf = func(o any) *ffi.TrustedIssuerRegistry {
		return o.(*Registry).h
	}
	ffi.ToTrustedIssuerRegistry = func(h *ffi.TrustedIssuerRegistry) any {
		return &Registry{h: h}
	}
}

// NewRegistry initializes an empty registry.
func NewRegistry() *Registry { return &Registry{h: ffi.NewTrustedIssuerRegistry()} }

// DefaultProduction returns the registry of self's production issuers.
func DefaultProduction() *Registry {
	return &Registry{h: ffi.DefaultProductionTrustedIssuerRegistry()}
}

// DefaultSandbox returns the registry of self's sandbox issuers.
func DefaultSandbox() *Registry {
	return &Registry{h: ffi.DefaultSandboxTrustedIssuerRegistry()}
}

// DefaultCredentialTypes returns the default credential types issued by self.
func DefaultCredentialTypes() []string { return ffi.DefaultIssuedCredentialTypes() }

// DefaultIssuerEpoch returns the default epoch from when self-issued credentials are valid.
func DefaultIssuerEpoch() time.Time { return time.Unix(ffi.DefaultIssuerEpoch(), 0) }

// AddIssuer adds an issuer. Returns false if it was already present.
func (r *Registry) AddIssuer(issuer *credential.Address) bool {
	return r.h.IssuerAdd(ffi.DIDAddressOf(issuer))
}

// RemoveIssuer removes an issuer. Returns false if it was not present.
func (r *Registry) RemoveIssuer(issuer *credential.Address) bool {
	return r.h.IssuerRemove(ffi.DIDAddressOf(issuer))
}

// GrantAuthority grants the issuer authority over a credential type from
// `granted` until `revoked` (zero `revoked` means no end).
func (r *Registry) GrantAuthority(issuer *credential.Address, credentialType string, granted, revoked time.Time) error {
	var revokedUnix int64
	if !revoked.IsZero() {
		revokedUnix = revoked.Unix()
	}

	return r.h.AuthorityGrant(ffi.DIDAddressOf(issuer), credentialType, granted.Unix(), revokedUnix)
}

// RevokeAuthority marks the issuer's authority over a credential type revoked
// from the given time.
func (r *Registry) RevokeAuthority(issuer *credential.Address, credentialType string, revoked time.Time) error {
	return r.h.AuthorityRevoke(ffi.DIDAddressOf(issuer), credentialType, revoked.Unix())
}

// AuthorityFor returns the credential types the issuer is authorized for.
func (r *Registry) AuthorityFor(issuer *credential.Address) ([]string, error) {
	return r.h.AuthorityFor(ffi.DIDAddressOf(issuer))
}

// AuthorityAt reports whether the issuer was authorized for the credential
// type at the given time.
func (r *Registry) AuthorityAt(issuer *credential.Address, credentialType string, issued time.Time) bool {
	return r.h.AuthorityAt(ffi.DIDAddressOf(issuer), credentialType, issued.Unix())
}

// Issuers returns all issuers in the registry.
func (r *Registry) Issuers() []*credential.Address {
	as := r.h.Issuers()
	out := make([]*credential.Address, len(as))

	for i, a := range as {
		out[i] = ffi.ToDIDAddress(a).(*credential.Address)
	}

	return out
}
