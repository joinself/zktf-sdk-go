// Package identity provides resolved identity documents and identity operations.
package identity

import (
	"time"

	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair"
	"github.com/joinself/zktf-sdk-go/keypair/exchange"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// Document is the resolved key state of an identity.
type Document struct {
	h *ffi.IdentityDocument
}

// Operation is a hashgraph operation describing a change to an identity document.
type Operation struct {
	h *ffi.IdentityOperation
}

func init() {
	ffi.IdentityDocumentOf = func(o any) *ffi.IdentityDocument { return o.(*Document).h }
	ffi.ToIdentityDocument = func(h *ffi.IdentityDocument) any { return &Document{h: h} }

	ffi.IdentityOperationOf = func(o any) *ffi.IdentityOperation { return o.(*Operation).h }
	ffi.ToIdentityOperation = func(h *ffi.IdentityOperation) any { return &Operation{h: h} }
}

// KeyOption filters a document key query by snapshot time and/or roles. With no
// options, queries run against the latest snapshot with no role filter.
type KeyOption func(*ffi.IdentityKeyLookup)

// AtTime evaluates the query against the document as it existed at t.
func AtTime(t time.Time) KeyOption {
	return func(l *ffi.IdentityKeyLookup) { l.AtTime(t.Unix()) }
}

// WithRoles restricts the result to keys holding every role in roles.
func WithRoles(roles KeyRole) KeyOption {
	return func(l *ffi.IdentityKeyLookup) { l.WithRoles(ffi.IdentityKeyRole(roles)) }
}

// buildKeyLookup returns a lookup configured by opts, or nil if opts is empty.
func buildKeyLookup(opts []KeyOption) *ffi.IdentityKeyLookup {
	if len(opts) == 0 {
		return nil
	}

	l := ffi.NewIdentityKeyLookup()
	for _, opt := range opts {
		opt(l)
	}

	return l
}

// Commitment returns the document's commitment hash, or nil if none is set.
func (d *Document) Commitment() []byte { return d.h.Commitment() }

// Create returns an operation builder seeded from the current document state.
func (d *Document) Create() *OperationBuilder {
	return &OperationBuilder{h: d.h.Create()}
}

// SigningKeys returns the signing keys in the document matching opts.
func (d *Document) SigningKeys(opts ...KeyOption) []*signing.PublicKey {
	ks := d.h.SigningKeys(buildKeyLookup(opts))
	out := make([]*signing.PublicKey, len(ks))

	for i, k := range ks {
		out[i] = ffi.ToSigningPublicKey(k).(*signing.PublicKey)
	}

	return out
}

// ExchangeKeys returns the exchange keys in the document matching opts.
func (d *Document) ExchangeKeys(opts ...KeyOption) []*exchange.PublicKey {
	ks := d.h.ExchangeKeys(buildKeyLookup(opts))
	out := make([]*exchange.PublicKey, len(ks))

	for i, k := range ks {
		out[i] = ffi.ToExchangePublicKey(k).(*exchange.PublicKey)
	}

	return out
}

// Descriptions returns the key descriptions in the document matching opts.
func (d *Document) Descriptions(opts ...KeyOption) []AddressDescription {
	ds := d.h.Descriptions(buildKeyLookup(opts))
	out := make([]AddressDescription, len(ds))

	for i, desc := range ds {
		out[i] = AddressDescription{h: desc}
	}

	return out
}

// HasRoles reports whether key holds every role in roles, evaluated under opts.
// It returns false for a key that is neither a signing nor an exchange key.
func (d *Document) HasRoles(key keypair.PublicKey, roles KeyRole, opts ...KeyOption) bool {
	l := buildKeyLookup(opts)

	switch k := key.(type) {
	case *signing.PublicKey:
		return d.h.SigningKeyHasRoles(ffi.SigningPublicKeyOf(k), ffi.IdentityKeyRole(roles), l)
	case *exchange.PublicKey:
		return d.h.ExchangeKeyHasRoles(ffi.ExchangePublicKeyOf(k), ffi.IdentityKeyRole(roles), l)
	default:
		return false
	}
}

// Valid reports whether key is valid in the document, evaluated under opts. It
// returns false for a key that is neither a signing nor an exchange key.
func (d *Document) Valid(key keypair.PublicKey, opts ...KeyOption) bool {
	l := buildKeyLookup(opts)

	switch k := key.(type) {
	case *signing.PublicKey:
		return d.h.SigningKeyValid(ffi.SigningPublicKeyOf(k), l)
	case *exchange.PublicKey:
		return d.h.ExchangeKeyValid(ffi.ExchangePublicKeyOf(k), l)
	default:
		return false
	}
}

// ThresholdMet reports whether signers collectively satisfy role's threshold,
// evaluated under opts.
func (d *Document) ThresholdMet(role KeyRole, signers []*signing.PublicKey, opts ...KeyOption) bool {
	hs := make([]*ffi.SigningPublicKey, len(signers))
	for i, s := range signers {
		hs[i] = ffi.SigningPublicKeyOf(s)
	}

	return d.h.ThresholdMet(ffi.IdentityKeyRole(role), hs, buildKeyLookup(opts))
}

// AddressDescription is a key description recorded in a document (embedded or
// referenced).
type AddressDescription struct {
	h *ffi.IdentityOperationDescription
}

// Kind returns whether the description is embedded or referenced.
func (d AddressDescription) Kind() DescriptionKind { return DescriptionKind(d.h.Kind()) }

// AsEmbedded returns the embedded description, or nil if not embedded.
func (d AddressDescription) AsEmbedded() *DescEmbedded {
	e := d.h.Embedded()
	if e == nil {
		return nil
	}

	return &DescEmbedded{h: e}
}

// AsReference returns the reference description, or nil if not a reference.
func (d AddressDescription) AsReference() *DescReference {
	r := d.h.Reference()
	if r == nil {
		return nil
	}

	return &DescReference{h: r}
}
