package identity

import (
	"time"

	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/exchange"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// KeyRole is a bitmask of roles a key may hold in an identity document.
type KeyRole uint64

const (
	KeyRoleVerification   KeyRole = KeyRole(ffi.KeyRoleVerification)
	KeyRoleAssertion      KeyRole = KeyRole(ffi.KeyRoleAssertion)
	KeyRoleAuthentication KeyRole = KeyRole(ffi.KeyRoleAuthentication)
	KeyRoleDelegation     KeyRole = KeyRole(ffi.KeyRoleDelegation)
	KeyRoleInvocation     KeyRole = KeyRole(ffi.KeyRoleInvocation)
	KeyRoleKeyAgreement   KeyRole = KeyRole(ffi.KeyRoleKeyAgreement)
	KeyRoleMessaging      KeyRole = KeyRole(ffi.KeyRoleMessaging)
)

// KeypairType identifies the cryptographic type of a key.
type KeypairType uint32

const (
	KeypairSigning  KeypairType = KeypairType(ffi.KeypairSigning)
	KeypairExchange KeypairType = KeypairType(ffi.KeypairExchange)
)

// AddressMethod identifies the DID method an address uses.
type AddressMethod uint32

const (
	AddressMethodZktf AddressMethod = AddressMethod(ffi.AddressMethodZktf)
	AddressMethodKey  AddressMethod = AddressMethod(ffi.AddressMethodKey)
)

// OperationActionKind identifies an action inside an operation.
type OperationActionKind uint32

const (
	OperationActionGrant      OperationActionKind = OperationActionKind(ffi.OperationActionGrant)
	OperationActionModify     OperationActionKind = OperationActionKind(ffi.OperationActionModify)
	OperationActionRevoke     OperationActionKind = OperationActionKind(ffi.OperationActionRevoke)
	OperationActionRecover    OperationActionKind = OperationActionKind(ffi.OperationActionRecover)
	OperationActionDeactivate OperationActionKind = OperationActionKind(ffi.OperationActionDeactivate)
)

// DescriptionKind identifies the kind of description on an action.
type DescriptionKind uint32

const (
	DescriptionNone      DescriptionKind = DescriptionKind(ffi.DescriptionKindNone)
	DescriptionEmbedded  DescriptionKind = DescriptionKind(ffi.DescriptionKindEmbedded)
	DescriptionReference DescriptionKind = DescriptionKind(ffi.DescriptionKindReference)
)

// OperationBuilder builds an identity-document operation.
type OperationBuilder struct {
	h *ffi.IdentityOperationBuilder
}

// NewOperationBuilder starts building an identity operation.
func NewOperationBuilder() *OperationBuilder {
	return &OperationBuilder{h: ffi.NewIdentityOperationBuilder()}
}

// ID sets the document address the operation targets.
func (b *OperationBuilder) ID(id *signing.PublicKey) *OperationBuilder {
	b.h.ID(ffi.SigningPublicKeyOf(id))
	return b
}

// Previous sets the hash of the previous operation.
func (b *OperationBuilder) Previous(hash []byte) *OperationBuilder {
	b.h.Previous(hash)
	return b
}

// Sequence sets the operation sequence number.
func (b *OperationBuilder) Sequence(seq uint32) *OperationBuilder {
	b.h.Sequence(seq)
	return b
}

// Timestamp sets the operation timestamp.
func (b *OperationBuilder) Timestamp(t time.Time) *OperationBuilder {
	b.h.Timestamp(t.Unix())
	return b
}

// Anchor attaches a biometric anchor + nonce.
func (b *OperationBuilder) Anchor(anchor, nonce []byte) *OperationBuilder {
	b.h.Anchor(anchor, nonce)
	return b
}

// Commitment attaches a commitment value.
func (b *OperationBuilder) Commitment(commitment []byte) *OperationBuilder {
	b.h.Commitment(commitment)
	return b
}

// Threshold sets the threshold required to satisfy a role.
func (b *OperationBuilder) Threshold(role KeyRole, threshold uint64) *OperationBuilder {
	b.h.Threshold(ffi.IdentityKeyRole(role), threshold)
	return b
}

// Weight assigns a weight to a key for a role.
func (b *OperationBuilder) Weight(key *signing.PublicKey, role KeyRole, weight uint64) *OperationBuilder {
	b.h.Weight(ffi.SigningPublicKeyOf(key), ffi.IdentityKeyRole(role), weight)
	return b
}

// GrantSigning grants a signing key roles (embedded).
func (b *OperationBuilder) GrantSigning(key *signing.PublicKey, roles KeyRole) *OperationBuilder {
	b.h.SigningGrantEmbedded(ffi.SigningPublicKeyOf(key), ffi.IdentityKeyRole(roles))
	return b
}

// GrantSigningReferenced grants a signing key via a referenced description.
func (b *OperationBuilder) GrantSigningReferenced(method uint16, controller, key *signing.PublicKey, commitment []byte, roles KeyRole) *OperationBuilder {
	b.h.SigningGrantReferenced(method, ffi.SigningPublicKeyOf(controller), ffi.SigningPublicKeyOf(key), commitment, ffi.IdentityKeyRole(roles))
	return b
}

// ModifySigning modifies the roles of a signing key.
func (b *OperationBuilder) ModifySigning(key *signing.PublicKey, roles KeyRole) *OperationBuilder {
	b.h.SigningModify(ffi.SigningPublicKeyOf(key), ffi.IdentityKeyRole(roles))
	return b
}

// RevokeSigning revokes a signing key effective from the given time.
func (b *OperationBuilder) RevokeSigning(key *signing.PublicKey, effectiveFrom time.Time) *OperationBuilder {
	b.h.SigningRevoke(ffi.SigningPublicKeyOf(key), effectiveFrom.Unix())
	return b
}

// GrantExchange grants an exchange key roles (embedded).
func (b *OperationBuilder) GrantExchange(key *exchange.PublicKey, roles KeyRole) *OperationBuilder {
	b.h.ExchangeGrantEmbedded(ffi.ExchangePublicKeyOf(key), ffi.IdentityKeyRole(roles))
	return b
}

// ModifyExchange modifies the roles of an exchange key.
func (b *OperationBuilder) ModifyExchange(key *exchange.PublicKey, roles KeyRole) *OperationBuilder {
	b.h.ExchangeModify(ffi.ExchangePublicKeyOf(key), ffi.IdentityKeyRole(roles))
	return b
}

// RevokeExchange revokes an exchange key effective from the given time.
func (b *OperationBuilder) RevokeExchange(key *exchange.PublicKey, effectiveFrom time.Time) *OperationBuilder {
	b.h.ExchangeRevoke(ffi.ExchangePublicKeyOf(key), effectiveFrom.Unix())
	return b
}

// Recover stages a recovery effective from the given time.
func (b *OperationBuilder) Recover(effectiveFrom time.Time) *OperationBuilder {
	b.h.Recover(effectiveFrom.Unix())
	return b
}

// Deactivate stages a deactivation effective from the given time.
func (b *OperationBuilder) Deactivate(effectiveFrom time.Time) *OperationBuilder {
	b.h.Deactivate(effectiveFrom.Unix())
	return b
}

// SignWith records the signing key.
func (b *OperationBuilder) SignWith(signer *signing.PublicKey) *OperationBuilder {
	b.h.SignWith(ffi.SigningPublicKeyOf(signer))
	return b
}

// Finish finalizes the operation. Sign it via account.Account.IdentitySign.
func (b *OperationBuilder) Finish() (*Operation, error) {
	o, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Operation{h: o}, nil
}

// DecodeOperation decodes an encoded operation for the given document address.
func DecodeOperation(documentAddress *signing.PublicKey, data []byte) (*Operation, error) {
	o, err := ffi.IdentityOperationDecode(ffi.SigningPublicKeyOf(documentAddress), data)
	if err != nil {
		return nil, err
	}

	return &Operation{h: o}, nil
}

// Encode returns the encoded bytes of the operation.
func (o *Operation) Encode() ([]byte, error) { return o.h.Encode() }

// Sequence returns the operation sequence number.
func (o *Operation) Sequence() uint32 { return o.h.Sequence() }

// Hash returns the 32-byte operation hash.
func (o *Operation) Hash() []byte { return o.h.Hash() }

// SignedBy reports whether the operation has been signed by the given key.
func (o *Operation) SignedBy(signer *signing.PublicKey) bool {
	return o.h.SignedBy(ffi.SigningPublicKeyOf(signer))
}

// Merge merges signatures from another operation into this one.
func (o *Operation) Merge(other *Operation) error { return o.h.Merge(other.h) }

// Actions returns the actions inside the operation.
func (o *Operation) Actions() []*Action {
	as := o.h.Actions()
	out := make([]*Action, len(as))

	for i, a := range as {
		out[i] = &Action{h: a}
	}

	return out
}

// Action is one of the actions inside an identity operation.
type Action struct {
	h *ffi.OperationAction
}

// Kind returns the kind of action.
func (a *Action) Kind() OperationActionKind { return OperationActionKind(a.h.Kind()) }

// Roles returns the roles assigned by a grant or modify action.
func (a *Action) Roles() KeyRole { return KeyRole(a.h.Roles()) }

// EffectiveFrom returns when the action takes effect.
func (a *Action) EffectiveFrom() time.Time { return time.Unix(a.h.EffectiveFrom(), 0) }

// DescriptionKind returns the kind of description on the action.
func (a *Action) DescriptionKind() DescriptionKind {
	return DescriptionKind(a.h.DescriptionKind())
}

// AsEmbedded returns the embedded description, or nil.
func (a *Action) AsEmbedded() *DescEmbedded {
	d := a.h.DescriptionEmbedded()
	if d == nil {
		return nil
	}

	return &DescEmbedded{h: d}
}

// AsReference returns the reference description, or nil.
func (a *Action) AsReference() *DescReference {
	d := a.h.DescriptionReference()
	if d == nil {
		return nil
	}

	return &DescReference{h: d}
}

// DescEmbedded describes a key embedded in an action.
type DescEmbedded struct {
	h *ffi.OperationDescriptionEmbedded
}

// AddressType returns the type of address (signing or exchange).
func (d *DescEmbedded) AddressType() KeypairType { return KeypairType(d.h.AddressType()) }

// AsSigning returns the address as a signing key, or nil.
func (d *DescEmbedded) AsSigning() *signing.PublicKey {
	k := d.h.AddressAsSigning()
	if k == nil {
		return nil
	}

	return ffi.ToSigningPublicKey(k).(*signing.PublicKey)
}

// AsExchange returns the address as an exchange key, or nil.
func (d *DescEmbedded) AsExchange() *exchange.PublicKey {
	k := d.h.AddressAsExchange()
	if k == nil {
		return nil
	}

	return ffi.ToExchangePublicKey(k).(*exchange.PublicKey)
}

// Controller returns the controller, or nil.
func (d *DescEmbedded) Controller() *signing.PublicKey {
	k := d.h.Controller()
	if k == nil {
		return nil
	}

	return ffi.ToSigningPublicKey(k).(*signing.PublicKey)
}

// DescReference describes a key referenced by another method.
type DescReference struct {
	h *ffi.OperationDescriptionReference
}

// Method returns the DID method.
func (d *DescReference) Method() AddressMethod { return AddressMethod(d.h.Method()) }

// AddressType returns the type of address (signing or exchange).
func (d *DescReference) AddressType() KeypairType { return KeypairType(d.h.AddressType()) }

// AsSigning returns the address as a signing key, or nil.
func (d *DescReference) AsSigning() *signing.PublicKey {
	k := d.h.AddressAsSigning()
	if k == nil {
		return nil
	}

	return ffi.ToSigningPublicKey(k).(*signing.PublicKey)
}

// AsExchange returns the address as an exchange key, or nil.
func (d *DescReference) AsExchange() *exchange.PublicKey {
	k := d.h.AddressAsExchange()
	if k == nil {
		return nil
	}

	return ffi.ToExchangePublicKey(k).(*exchange.PublicKey)
}

// Controller returns the controller, or nil.
func (d *DescReference) Controller() *signing.PublicKey {
	k := d.h.Controller()
	if k == nil {
		return nil
	}

	return ffi.ToSigningPublicKey(k).(*signing.PublicKey)
}

// Lookup is the constructed lookup query.
type Lookup struct {
	h *ffi.IdentityLookup
}

func init() {
	ffi.IdentityLookupOf = func(o any) *ffi.IdentityLookup { return o.(*Lookup).h }
	ffi.ToIdentityLookup = func(h *ffi.IdentityLookup) any { return &Lookup{h: h} }
}

// LookupOption is one of the variadic filters accepted by IdentityLookup.
type LookupOption func(*lookupOpts)

type lookupOpts struct {
	key *signing.PublicKey
}

// ByKey restricts the lookup to identities the given key is associated with.
func ByKey(key *signing.PublicKey) LookupOption {
	return func(o *lookupOpts) {
		o.key = key
	}
}

// BuildLookup applies the given options and returns a Lookup ready to pass to
// the FFI. For use by other zktf-sdk-go packages.
func BuildLookup(options ...LookupOption) *Lookup {
	var o lookupOpts
	for _, opt := range options {
		opt(&o)
	}

	l := &Lookup{h: ffi.NewIdentityLookup()}

	if o.key != nil {
		l.h.ByKey(ffi.SigningPublicKeyOf(o.key))
	}

	return l
}
