package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

// CredentialIssue signs a credential (with pending signers queued via
// CredentialBuilder.SignWith) into a verifiable credential. The signature is
// applied in place; the same credential is returned once signed.
func (a *Account) CredentialIssue(credential *VerifiableCredential) (*VerifiableCredential, error) {
	if err := status(C.zktf_account_credential_sign(a.ptr, credential.ptr)); err != nil {
		return nil, err
	}

	return credential, nil
}

// CredentialStore stores a verifiable credential in the account's local store.
func (a *Account) CredentialStore(credential *VerifiableCredential) error {
	return status(C.zktf_account_credential_store(a.ptr, credential.ptr))
}

// CredentialLookup returns credentials in the account's local store that
// satisfy the given predicate tree.
func (a *Account) CredentialLookup(tree *PredicateTree) ([]*VerifiableCredential, error) {
	var c *C.zktf_collection_verifiable_credential

	if err := status(C.zktf_account_credential_lookup(a.ptr, tree.ptr, &c)); err != nil {
		return nil, err
	}

	return verifiableCredentialsFrom(c), nil
}

// CredentialSharedWith returns credentials the account has shared with the
// given address that satisfy the predicate tree.
func (a *Account) CredentialSharedWith(with *SigningPublicKey, tree *PredicateTree) ([]*VerifiableCredential, error) {
	var c *C.zktf_collection_verifiable_credential

	if err := status(C.zktf_account_credential_shared_with(a.ptr, with.ptr, tree.ptr, &c)); err != nil {
		return nil, err
	}

	return verifiableCredentialsFrom(c), nil
}
