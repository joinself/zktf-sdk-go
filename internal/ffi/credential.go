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

// CredentialTypeCollection wraps a zktf_collection_credential_type handle.
type CredentialTypeCollection struct {
	ptr *C.zktf_collection_credential_type
}

func newCredentialTypeCollection(ptr *C.zktf_collection_credential_type) *CredentialTypeCollection {
	if ptr == nil {
		return nil
	}
	c := &CredentialTypeCollection{ptr: ptr}
	runtime.AddCleanup(c, func(ptr *C.zktf_collection_credential_type) {
		C.zktf_collection_credential_type_destroy(ptr)
	}, c.ptr)
	return c
}

// NewCredentialTypes builds a credential type collection from the given type
// strings (e.g. "VerifiableCredential", "EmailCredential").
func NewCredentialTypes(types []string) *CredentialTypeCollection {
	ptr := C.zktf_collection_credential_type_init()
	for _, t := range types {
		ct := cstring(t)
		C.zktf_collection_credential_type_append(ptr, ct)
		free(unsafe.Pointer(ct))
	}
	return newCredentialTypeCollection(ptr)
}

// Strings returns the type strings in the collection.
func (c *CredentialTypeCollection) Strings() []string {
	n := int(C.zktf_collection_credential_type_len(c.ptr))
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = C.GoString(C.zktf_collection_credential_type_at(c.ptr, C.size_t(i)))
	}
	return out
}

// CredentialTerm describes the duration under which the requester wishes to
// access shared credentials.
type CredentialTerm struct {
	ptr *C.zktf_credential_term
}

func newCredentialTerm(ptr *C.zktf_credential_term) *CredentialTerm {
	if ptr == nil {
		return nil
	}
	t := &CredentialTerm{ptr: ptr}
	runtime.AddCleanup(t, func(ptr *C.zktf_credential_term) {
		C.zktf_credential_term_destroy(ptr)
	}, t.ptr)
	return t
}

// NewCredentialTerm creates a credential term with the given duration in seconds.
func NewCredentialTerm(durationSeconds uint64) *CredentialTerm {
	return newCredentialTerm(C.zktf_credential_term_create(C.uint64_t(durationSeconds)))
}

// Duration returns the term's duration in seconds.
func (t *CredentialTerm) Duration() uint64 {
	return uint64(C.zktf_credential_term_duration(t.ptr))
}

// Credential wraps an unsigned zktf_credential handle.
type Credential struct {
	ptr *C.zktf_credential
}

func newCredential(ptr *C.zktf_credential) *Credential {
	if ptr == nil {
		return nil
	}
	c := &Credential{ptr: ptr}
	runtime.AddCleanup(c, func(ptr *C.zktf_credential) {
		C.zktf_credential_destroy(ptr)
	}, c.ptr)
	return c
}

// CredentialBuilder wraps a zktf_credential_builder handle.
type CredentialBuilder struct {
	ptr *C.zktf_credential_builder
}

// NewCredentialBuilder initializes a new credential builder.
func NewCredentialBuilder() *CredentialBuilder {
	ptr := C.zktf_credential_builder_init()
	b := &CredentialBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_credential_builder) {
		C.zktf_credential_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// CredentialType sets the credential's types.
func (b *CredentialBuilder) CredentialType(types *CredentialTypeCollection) *CredentialBuilder {
	C.zktf_credential_builder_credential_type(b.ptr, types.ptr)
	return b
}

// Issuer sets the credential's issuer DID address.
func (b *CredentialBuilder) Issuer(issuer *DIDAddress) *CredentialBuilder {
	C.zktf_credential_builder_issuer(b.ptr, issuer.ptr)
	return b
}

// CredentialSubject sets the credential's subject DID address.
func (b *CredentialBuilder) CredentialSubject(subject *DIDAddress) *CredentialBuilder {
	C.zktf_credential_builder_credential_subject(b.ptr, subject.ptr)
	return b
}

// CredentialSubjectClaim adds a string claim about the subject.
func (b *CredentialBuilder) CredentialSubjectClaim(key, value string) *CredentialBuilder {
	ckey, cval := cstring(key), cstring(value)
	defer free(unsafe.Pointer(ckey))
	defer free(unsafe.Pointer(cval))
	C.zktf_credential_builder_credential_subject_claim(b.ptr, ckey, cval)
	return b
}

// ValidFrom sets the unix timestamp (seconds) the credential is valid from.
func (b *CredentialBuilder) ValidFrom(unix int64) *CredentialBuilder {
	C.zktf_credential_builder_valid_from(b.ptr, C.int64_t(unix))
	return b
}

// ValidUntil sets the unix timestamp (seconds) the credential is valid until.
func (b *CredentialBuilder) ValidUntil(unix int64) *CredentialBuilder {
	C.zktf_credential_builder_valid_until(b.ptr, C.int64_t(unix))
	return b
}

// CredentialSubjectJSON sets the subject claims from a raw JSON document.
func (b *CredentialBuilder) CredentialSubjectJSON(json []byte) *CredentialBuilder {
	buf, length := cbytes(json)
	defer free(unsafe.Pointer(buf))
	C.zktf_credential_builder_credential_subject_json(b.ptr, buf, length)
	return b
}

// SignWith records the signing key and issuance time for the credential.
func (b *CredentialBuilder) SignWith(signer *SigningPublicKey, issuedAtUnix int64) *CredentialBuilder {
	C.zktf_credential_builder_sign_with(b.ptr, signer.ptr, C.int64_t(issuedAtUnix))
	return b
}

// Finish finalizes the unsigned credential.
func (b *CredentialBuilder) Finish() (*Credential, error) {
	var ptr *C.zktf_credential
	if err := status(C.zktf_credential_builder_finish(b.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newCredential(ptr), nil
}

// VerifiableCredential wraps a signed zktf_verifiable_credential handle.
type VerifiableCredential struct {
	ptr *C.zktf_verifiable_credential
}

func newVerifiableCredential(ptr *C.zktf_verifiable_credential) *VerifiableCredential {
	if ptr == nil {
		return nil
	}
	c := &VerifiableCredential{ptr: ptr}
	runtime.AddCleanup(c, func(ptr *C.zktf_verifiable_credential) {
		C.zktf_verifiable_credential_destroy(ptr)
	}, c.ptr)
	return c
}

// VerifiableCredentialDecode decodes a JSON-encoded verifiable credential.
func VerifiableCredentialDecode(data []byte) (*VerifiableCredential, error) {
	buf, length := cbytes(data)
	defer free(unsafe.Pointer(buf))

	var ptr *C.zktf_verifiable_credential
	if err := status(C.zktf_verifiable_credential_decode(&ptr, buf, length)); err != nil {
		return nil, err
	}
	return newVerifiableCredential(ptr), nil
}

// Validate returns an error if the credential is invalid.
func (c *VerifiableCredential) Validate() error {
	return status(C.zktf_verifiable_credential_validate(c.ptr))
}

// TypeOf returns the credential's type strings.
func (c *VerifiableCredential) TypeOf() *CredentialTypeCollection {
	return newCredentialTypeCollection(C.zktf_verifiable_credential_type_of(c.ptr))
}

// Issuer returns the issuer DID address.
func (c *VerifiableCredential) Issuer() *DIDAddress {
	return newDIDAddress(C.zktf_verifiable_credential_issuer(c.ptr))
}

// Subject returns the subject DID address.
func (c *VerifiableCredential) Subject() *DIDAddress {
	return newDIDAddress(C.zktf_verifiable_credential_credential_subject(c.ptr))
}

// SubjectClaim returns a string claim about the subject, or "" if absent.
func (c *VerifiableCredential) SubjectClaim(key string) string {
	ckey := cstring(key)
	defer free(unsafe.Pointer(ckey))

	buf := C.zktf_verifiable_credential_credential_subject_claim(c.ptr, ckey)
	if buf == nil {
		return ""
	}
	defer C.zktf_string_buffer_destroy(buf)
	return C.GoString(C.zktf_string_buffer_ptr(buf))
}

// SubjectJSON returns the subject claims as a raw JSON document, or nil.
func (c *VerifiableCredential) SubjectJSON() []byte {
	return goBytesFromBuffer(C.zktf_verifiable_credential_credential_subject_json(c.ptr))
}

// ValidFrom returns the unix timestamp (seconds) the credential is valid from.
func (c *VerifiableCredential) ValidFrom() int64 {
	return int64(C.zktf_verifiable_credential_valid_from(c.ptr))
}

// ValidUntil returns the unix timestamp (seconds) the credential is valid until.
func (c *VerifiableCredential) ValidUntil() int64 {
	return int64(C.zktf_verifiable_credential_valid_until(c.ptr))
}

// Created returns the unix timestamp (seconds) the credential was created.
func (c *VerifiableCredential) Created() int64 {
	return int64(C.zktf_verifiable_credential_created(c.ptr))
}

// Signer returns the DID address that signed the credential.
func (c *VerifiableCredential) Signer() (*DIDAddress, error) {
	var out *C.zktf_did_address
	if err := status(C.zktf_verifiable_credential_signer(c.ptr, &out)); err != nil {
		return nil, err
	}
	return newDIDAddress(out), nil
}

// SigningKey returns the signing key that signed the credential.
func (c *VerifiableCredential) SigningKey() (*SigningPublicKey, error) {
	var out *C.zktf_signing_public_key
	if err := status(C.zktf_verifiable_credential_signing_key(c.ptr, &out)); err != nil {
		return nil, err
	}
	return newSigningPublicKey(out), nil
}

// RevocationHashes returns the revocation hashes of the credential, one per proof.
func (c *VerifiableCredential) RevocationHashes() ([][]byte, error) {
	var out *C.zktf_collection_bytes_buffer
	if err := status(C.zktf_verifiable_credential_revocation_hashes(c.ptr, &out)); err != nil {
		return nil, err
	}
	return bytesFromBufferCollection(out), nil
}

// Encode returns the JSON-encoded credential.
func (c *VerifiableCredential) Encode() ([]byte, error) {
	var buf *C.zktf_bytes_buffer
	if err := status(C.zktf_verifiable_credential_encode(c.ptr, &buf)); err != nil {
		return nil, err
	}
	defer C.zktf_bytes_buffer_destroy(buf)
	return C.GoBytes(
		unsafe.Pointer(C.zktf_bytes_buffer_buf(buf)),
		C.int(C.zktf_bytes_buffer_len(buf)),
	), nil
}
