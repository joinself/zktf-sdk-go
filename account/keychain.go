package account

import (
	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/exchange"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// KeychainSigningCreate generates a new signing key in the keychain and returns
// its public address.
func (a *Account) KeychainSigningCreate() (*signing.PublicKey, error) {
	k, err := a.h.KeychainSigningCreate()
	if err != nil {
		return nil, err
	}

	return ffi.ToSigningPublicKey(k).(*signing.PublicKey), nil
}

// KeychainExchangeCreate generates a new exchange key in the keychain and
// returns its public address.
func (a *Account) KeychainExchangeCreate() (*exchange.PublicKey, error) {
	k, err := a.h.KeychainExchangeCreate()
	if err != nil {
		return nil, err
	}

	return ffi.ToExchangePublicKey(k).(*exchange.PublicKey), nil
}

// KeychainSign signs payload with the keychain key identified by address.
func (a *Account) KeychainSign(address *signing.PublicKey, payload []byte) ([]byte, error) {
	return a.h.KeychainSign(ffi.SigningPublicKeyOf(address), payload)
}

// KeychainLookupOption is one of the variadic filters accepted by
// KeychainLookup.
type KeychainLookupOption func(*keychainLookupOpts)

type keychainLookupOpts struct {
	identity *signing.PublicKey
	roles    identity.KeyRole
	hasRoles bool
}

// ByIdentity restricts the lookup to keys associated with the given identity.
func ByIdentity(key *signing.PublicKey) KeychainLookupOption {
	return func(o *keychainLookupOpts) {
		o.identity = key
	}
}

// WithRoles restricts the lookup to keys carrying every role in roles. Only
// applies in combination with ByIdentity.
func WithRoles(roles identity.KeyRole) KeychainLookupOption {
	return func(o *keychainLookupOpts) {
		o.roles = roles
		o.hasRoles = true
	}
}

// KeychainLookup returns the signing keys held in the keychain that satisfy the
// given filters. With no filters it returns every signing key.
func (a *Account) KeychainLookup(options ...KeychainLookupOption) ([]*signing.PublicKey, error) {
	var o keychainLookupOpts
	for _, opt := range options {
		opt(&o)
	}

	l := ffi.NewKeychainLookup()
	if o.identity != nil {
		l.ByIdentity(ffi.SigningPublicKeyOf(o.identity))
	}
	if o.hasRoles {
		l.WithRoles(ffi.IdentityKeyRole(o.roles))
	}

	ks, err := a.h.KeychainLookup(l)
	if err != nil {
		return nil, err
	}

	out := make([]*signing.PublicKey, len(ks))
	for i, k := range ks {
		out[i] = ffi.ToSigningPublicKey(k).(*signing.PublicKey)
	}

	return out, nil
}
