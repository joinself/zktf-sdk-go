package credential

import (
	"time"

	"github.com/joinself/zktf-sdk-go/internal/ffi"
)

// TrustedIssuerRegistry is a set of issuers and their per-credential-type
// authority windows. It is a credential-specific construct used to build
// credential graphs.
type TrustedIssuerRegistry struct {
	h *ffi.TrustedIssuerRegistry
}

func init() {
	ffi.TrustedIssuerRegistryOf = func(o any) *ffi.TrustedIssuerRegistry {
		return o.(*TrustedIssuerRegistry).h
	}
	ffi.ToTrustedIssuerRegistry = func(h *ffi.TrustedIssuerRegistry) any {
		return &TrustedIssuerRegistry{h: h}
	}
}

// NewTrustedIssuerRegistry initializes an empty registry.
func NewTrustedIssuerRegistry() *TrustedIssuerRegistry {
	return &TrustedIssuerRegistry{h: ffi.NewTrustedIssuerRegistry()}
}

// DefaultProductionRegistry returns the registry of self's production issuers.
func DefaultProductionRegistry() *TrustedIssuerRegistry {
	return &TrustedIssuerRegistry{h: ffi.DefaultProductionTrustedIssuerRegistry()}
}

// DefaultSandboxRegistry returns the registry of self's sandbox issuers.
func DefaultSandboxRegistry() *TrustedIssuerRegistry {
	return &TrustedIssuerRegistry{h: ffi.DefaultSandboxTrustedIssuerRegistry()}
}

// DefaultCredentialTypes returns the default credential types issued by self.
func DefaultCredentialTypes() []string { return ffi.DefaultIssuedCredentialTypes() }

// DefaultIssuerEpoch returns the default epoch from when self-issued credentials
// are valid.
func DefaultIssuerEpoch() time.Time { return time.Unix(ffi.DefaultIssuerEpoch(), 0) }

// AddIssuer adds an issuer. Returns false if it was already present.
func (r *TrustedIssuerRegistry) AddIssuer(issuer *Address) bool {
	return r.h.IssuerAdd(ffi.DIDAddressOf(issuer))
}

// RemoveIssuer removes an issuer. Returns false if it was not present.
func (r *TrustedIssuerRegistry) RemoveIssuer(issuer *Address) bool {
	return r.h.IssuerRemove(ffi.DIDAddressOf(issuer))
}

// GrantAuthority grants the issuer authority over a credential type from
// granted until revoked (a zero revoked means no end).
func (r *TrustedIssuerRegistry) GrantAuthority(issuer *Address, credentialType string, granted, revoked time.Time) error {
	var revokedUnix int64
	if !revoked.IsZero() {
		revokedUnix = revoked.Unix()
	}

	return r.h.AuthorityGrant(ffi.DIDAddressOf(issuer), credentialType, granted.Unix(), revokedUnix)
}

// RevokeAuthority marks the issuer's authority over a credential type revoked
// from the given time.
func (r *TrustedIssuerRegistry) RevokeAuthority(issuer *Address, credentialType string, revoked time.Time) error {
	return r.h.AuthorityRevoke(ffi.DIDAddressOf(issuer), credentialType, revoked.Unix())
}

// AuthorityFor returns the credential types the issuer is authorized for.
func (r *TrustedIssuerRegistry) AuthorityFor(issuer *Address) ([]string, error) {
	return r.h.AuthorityFor(ffi.DIDAddressOf(issuer))
}

// AuthorityAt reports whether the issuer was authorized for the credential type
// at the given time.
func (r *TrustedIssuerRegistry) AuthorityAt(issuer *Address, credentialType string, issued time.Time) bool {
	return r.h.AuthorityAt(ffi.DIDAddressOf(issuer), credentialType, issued.Unix())
}

// Issuers returns all issuers in the registry.
func (r *TrustedIssuerRegistry) Issuers() []*Address {
	as := r.h.Issuers()
	out := make([]*Address, len(as))

	for i, a := range as {
		out[i] = ffi.ToDIDAddress(a).(*Address)
	}

	return out
}
