package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

// CredentialIssue signs an unsigned credential into a verifiable credential.
func (a *Account) CredentialIssue(credential *Credential) (*VerifiableCredential, error) {
	var out *C.zktf_verifiable_credential

	if err := status(C.zktf_account_credential_issue(a.ptr, credential.ptr, &out)); err != nil {
		return nil, err
	}

	return newVerifiableCredential(out), nil
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

// PresentationIssue signs an unsigned presentation into a verifiable presentation.
func (a *Account) PresentationIssue(presentation *Presentation) (*VerifiablePresentation, error) {
	var out *C.zktf_verifiable_presentation

	if err := status(C.zktf_account_presentation_issue(a.ptr, presentation.ptr, &out)); err != nil {
		return nil, err
	}

	return newVerifiablePresentation(out), nil
}
