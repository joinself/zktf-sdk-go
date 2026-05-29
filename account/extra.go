package account

import (
	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/credential/predicate"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/object"
	"github.com/joinself/zktf-sdk-go/token"
)

// SetupPairingCode sets the account up for pairing with an application identity
// and returns the pairing code. It fails if the account is already paired.
func (a *Account) SetupPairingCode() (string, error) {
	return a.h.SetupPairingCode()
}

// PresentationSign signs a presentation with any account keys it requires.
func (a *Account) PresentationSign(p *credential.VerifiablePresentation) error {
	return a.h.PresentationSign(ffi.VerifiablePresentationOf(p))
}

// PresentationStore stores a presentation on the account for later retrieval.
func (a *Account) PresentationStore(p *credential.VerifiablePresentation) error {
	return a.h.PresentationStore(ffi.VerifiablePresentationOf(p))
}

// PresentationLookup returns presentations stored on the account that satisfy
// the predicate tree. A nil tree returns every stored presentation.
func (a *Account) PresentationLookup(tree *predicate.Tree) ([]*credential.VerifiablePresentation, error) {
	var t *ffi.PredicateTree
	if tree != nil {
		t = ffi.PredicateTreeOf(tree)
	}

	vps, err := a.h.PresentationLookup(t)
	if err != nil {
		return nil, err
	}

	out := make([]*credential.VerifiablePresentation, len(vps))
	for i, vp := range vps {
		out[i] = ffi.ToVerifiablePresentation(vp).(*credential.VerifiablePresentation)
	}

	return out, nil
}

// ObjectStore stores an object in the account's local data store.
func (a *Account) ObjectStore(obj *object.Object) error {
	return a.h.ObjectStore(ffi.ObjectOf(obj))
}

// ObjectRetrieve loads a locally stored object by its id.
func (a *Account) ObjectRetrieve(objectID []byte) (*object.Object, error) {
	o, err := a.h.ObjectRetrieve(objectID)
	if err != nil {
		return nil, err
	}

	return ffi.ToObject(o).(*object.Object), nil
}

// CredentialExchangeTrack records that a credential was exchanged with an address.
func (a *Account) CredentialExchangeTrack(with *signing.PublicKey, c *credential.Verifiable) error {
	return a.h.CredentialExchangeTrack(ffi.SigningPublicKeyOf(with), ffi.VerifiableCredentialOf(c))
}

// CredentialExchangeLog returns the credential exchange log, optionally
// restricted to exchanges with an address and to credentials satisfying a
// predicate tree. Either filter may be nil.
func (a *Account) CredentialExchangeLog(with *signing.PublicKey, tree *predicate.Tree) ([]*credential.Exchange, error) {
	var w *ffi.SigningPublicKey
	if with != nil {
		w = ffi.SigningPublicKeyOf(with)
	}

	var t *ffi.PredicateTree
	if tree != nil {
		t = ffi.PredicateTreeOf(tree)
	}

	es, err := a.h.CredentialExchangeLog(w, t)
	if err != nil {
		return nil, err
	}

	out := make([]*credential.Exchange, len(es))
	for i, e := range es {
		out[i] = ffi.ToCredentialExchange(e).(*credential.Exchange)
	}

	return out, nil
}

// TokenIssue issues a fresh token from a validated request.
func (a *Account) TokenIssue(req *token.Request) (*token.Token, error) {
	tk, err := a.h.TokenIssue(ffi.TokenRequestOf(req))
	if err != nil {
		return nil, err
	}

	return ffi.ToToken(tk).(*token.Token), nil
}

// TokenStore stores a token. Its issuer, bearer and local owner are derived
// from the token itself.
func (a *Account) TokenStore(tk *token.Token) error {
	return a.h.TokenStore(ffi.TokenOf(tk))
}
