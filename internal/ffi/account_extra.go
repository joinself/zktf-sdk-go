package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// SetupPairingCode sets the account up for pairing with an application identity
// and returns the pairing code. It fails if the account has already been paired.
func (a *Account) SetupPairingCode() (string, error) {
	var buf *C.zktf_string_buffer

	if err := status(C.zktf_account_setup_pairing_code(a.ptr, &buf)); err != nil {
		return "", err
	}

	return goStringFromBuffer(buf), nil
}

// PresentationSign signs a presentation with any available keys it requires.
func (a *Account) PresentationSign(vp *VerifiablePresentation) error {
	return status(C.zktf_account_presentation_sign(a.ptr, vp.ptr))
}

// PresentationStore stores a presentation on the account for later retrieval.
func (a *Account) PresentationStore(vp *VerifiablePresentation) error {
	return status(C.zktf_account_presentation_store(a.ptr, vp.ptr))
}

// PresentationLookup returns presentations stored on the account that satisfy
// the predicate tree. A nil tree returns every stored presentation.
func (a *Account) PresentationLookup(tree *PredicateTree) ([]*VerifiablePresentation, error) {
	var treePtr *C.zktf_credential_predicate_tree
	if tree != nil {
		treePtr = tree.ptr
	}

	var c *C.zktf_collection_verifiable_presentation
	if err := status(C.zktf_account_presentation_lookup(a.ptr, treePtr, &c)); err != nil {
		return nil, err
	}

	return verifiablePresentationsFrom(c), nil
}

// ObjectStore stores an object in the account's local data store.
func (a *Account) ObjectStore(obj *Object) error {
	return status(C.zktf_account_object_store(a.ptr, obj.ptr))
}

// ObjectRetrieve loads a locally stored object by its id.
func (a *Account) ObjectRetrieve(objectID []byte) (*Object, error) {
	idBuf, _ := cbytes(objectID)
	defer free(unsafe.Pointer(idBuf))

	var out *C.zktf_object
	if err := status(C.zktf_account_object_retrieve(a.ptr, idBuf, &out)); err != nil {
		return nil, err
	}

	return newObject(out), nil
}

// CredentialExchangeTrack records that a credential was exchanged with an address.
func (a *Account) CredentialExchangeTrack(with *SigningPublicKey, vc *VerifiableCredential) error {
	return status(C.zktf_account_credential_exchange_track(a.ptr, with.ptr, vc.ptr, nil))
}

// CredentialExchangeLog returns the credential exchange log, optionally
// restricted to exchanges with an address and to credentials satisfying a
// predicate tree. Either filter may be nil.
func (a *Account) CredentialExchangeLog(with *SigningPublicKey, tree *PredicateTree) ([]*CredentialExchange, error) {
	var withPtr *C.zktf_signing_public_key
	if with != nil {
		withPtr = with.ptr
	}

	var treePtr *C.zktf_credential_predicate_tree
	if tree != nil {
		treePtr = tree.ptr
	}

	var c *C.zktf_collection_credential_exchange
	if err := status(C.zktf_account_credential_exchange_log(a.ptr, withPtr, treePtr, &c)); err != nil {
		return nil, err
	}

	return credentialExchangesFrom(c), nil
}

// TokenIssue issues a fresh token from a validated request.
func (a *Account) TokenIssue(req *TokenRequest) (*Token, error) {
	var out *C.zktf_token

	if err := status(C.zktf_account_token_issue(a.ptr, req.ptr, &out)); err != nil {
		return nil, err
	}

	return newToken(out), nil
}

// TokenStore stores a token. The issuer, bearer and local owner are derived
// from the token itself.
func (a *Account) TokenStore(tk *Token) error {
	return status(C.zktf_account_token_store(a.ptr, tk.ptr))
}

// CredentialExchange records a single credential exchanged with an address.
type CredentialExchange struct {
	ptr *C.zktf_credential_exchange
}

func newCredentialExchange(ptr *C.zktf_credential_exchange) *CredentialExchange {
	if ptr == nil {
		return nil
	}
	e := &CredentialExchange{ptr: ptr}
	runtime.AddCleanup(e, func(ptr *C.zktf_credential_exchange) {
		C.zktf_credential_exchange_destroy(ptr)
	}, e.ptr)

	return e
}

// WithAddress returns the address the credential was exchanged with.
func (e *CredentialExchange) WithAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_credential_exchange_with_address(e.ptr))
}

// ContentHash returns the hash of the exchanged credential content.
func (e *CredentialExchange) ContentHash() []byte {
	return goBytesFromBuffer(C.zktf_credential_exchange_content_hash(e.ptr))
}

// SharedAt returns the unix timestamp (seconds) the credential was shared.
func (e *CredentialExchange) SharedAt() int64 {
	return int64(C.zktf_credential_exchange_shared_at(e.ptr))
}

// credentialExchangesFrom copies a caller-owned collection into Go wrappers and
// destroys the collection.
func credentialExchangesFrom(c *C.zktf_collection_credential_exchange) []*CredentialExchange {
	if c == nil {
		return nil
	}
	defer C.zktf_collection_credential_exchange_destroy(c)

	n := int(C.zktf_collection_credential_exchange_len(c))
	out := make([]*CredentialExchange, n)
	for i := 0; i < n; i++ {
		out[i] = newCredentialExchange(C.zktf_collection_credential_exchange_at(c, C.size_t(i)))
	}

	return out
}

// TokenRequest is a validated request ready to be issued into a token.
type TokenRequest struct {
	ptr *C.zktf_token_request
}

func newTokenRequest(ptr *C.zktf_token_request) *TokenRequest {
	if ptr == nil {
		return nil
	}
	r := &TokenRequest{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_token_request) {
		C.zktf_token_request_destroy(ptr)
	}, r.ptr)

	return r
}

// PushTokenBuilder builds a push token request.
type PushTokenBuilder struct {
	ptr *C.zktf_push_token_builder
}

// NewPushTokenBuilder initializes a push token builder.
func NewPushTokenBuilder() *PushTokenBuilder {
	ptr := C.zktf_push_token_builder_init()
	b := &PushTokenBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_push_token_builder) {
		C.zktf_push_token_builder_destroy(ptr)
	}, b.ptr)

	return b
}

// ForAddress sets the local group address the token authorizes notifications for.
func (b *PushTokenBuilder) ForAddress(address *SigningPublicKey) *PushTokenBuilder {
	C.zktf_push_token_builder_for_address(b.ptr, address.ptr)
	return b
}

// ProviderAddress sets the exchange public key of the push provider.
func (b *PushTokenBuilder) ProviderAddress(address *ExchangePublicKey) *PushTokenBuilder {
	C.zktf_push_token_builder_provider_address(b.ptr, address.ptr)
	return b
}

// Delegatable allows the bearer to further delegate the issued token.
func (b *PushTokenBuilder) Delegatable(delegatable bool) *PushTokenBuilder {
	C.zktf_push_token_builder_delegatable(b.ptr, C.bool(delegatable))
	return b
}

// Finish validates the configured fields and returns a token request.
func (b *PushTokenBuilder) Finish() (*TokenRequest, error) {
	var out *C.zktf_token_request
	if err := status(C.zktf_push_token_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}

	return newTokenRequest(out), nil
}
