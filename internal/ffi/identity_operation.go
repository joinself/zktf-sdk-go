package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import (
	"runtime"
	"time"
	"unsafe"
)

// IdentityKeyRole is a bitmask of roles a key may hold in an identity document.
type IdentityKeyRole uint64

const (
	KeyRoleVerification   IdentityKeyRole = C.KEY_ROLE_VERIFICATION
	KeyRoleAssertion      IdentityKeyRole = C.KEY_ROLE_ASSERTION
	KeyRoleAuthentication IdentityKeyRole = C.KEY_ROLE_AUTHENTICATION
	KeyRoleDelegation     IdentityKeyRole = C.KEY_ROLE_DELEGATION
	KeyRoleInvocation     IdentityKeyRole = C.KEY_ROLE_INVOCATION
	KeyRoleKeyAgreement   IdentityKeyRole = C.KEY_ROLE_KEYAGREEMENT
	KeyRoleMessaging      IdentityKeyRole = C.KEY_ROLE_MESSAGING
)

// KeypairType mirrors zktf_keypair_type.
type KeypairType uint32

const (
	KeypairSigning  KeypairType = C.KEYPAIR_SIGNING
	KeypairExchange KeypairType = C.KEYPAIR_EXCHANGE
)

// AddressMethod mirrors zktf_address_method.
type AddressMethod uint32

const (
	AddressMethodZktf AddressMethod = C.METHOD_ZKTF
	AddressMethodKey  AddressMethod = C.METHOD_KEY
)

// OperationActionKind mirrors zktf_identity_operation_action_type.
type OperationActionKind uint32

const (
	OperationActionGrant      OperationActionKind = C.OPERATION_ACTION_GRANT
	OperationActionModify     OperationActionKind = C.OPERATION_ACTION_MODIFY
	OperationActionRevoke     OperationActionKind = C.OPERATION_ACTION_REVOKE
	OperationActionRecover    OperationActionKind = C.OPERATION_ACTION_RECOVER
	OperationActionDeactivate OperationActionKind = C.OPERATION_ACTION_DEACTIVATE
)

// OperationDescriptionKind mirrors zktf_identity_operation_description_type.
type OperationDescriptionKind uint32

const (
	DescriptionKindNone      OperationDescriptionKind = C.OPERATION_DESCRIPTION_NONE
	DescriptionKindEmbedded  OperationDescriptionKind = C.OPERATION_DESCRIPTION_EMBEDDED
	DescriptionKindReference OperationDescriptionKind = C.OPERATION_DESCRIPTION_REFERENCE
)

const operationHashLen = 32

// IdentityOperationBuilder builds an identity-document operation (key grants,
// modifications, revocations, recovery, deactivation).
type IdentityOperationBuilder struct {
	ptr *C.zktf_identity_operation_builder
}

// NewIdentityOperationBuilder initializes an operation builder.
func NewIdentityOperationBuilder() *IdentityOperationBuilder {
	ptr := C.zktf_identity_operation_builder_init()
	b := &IdentityOperationBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_identity_operation_builder) {
		C.zktf_identity_operation_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// ID sets the document address the operation targets.
func (b *IdentityOperationBuilder) ID(id *SigningPublicKey) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_id(b.ptr, id.ptr)
	return b
}

// Previous sets the hash of the previous operation in the sequence.
func (b *IdentityOperationBuilder) Previous(hash []byte) *IdentityOperationBuilder {
	buf, length := cbytes(hash)
	defer free(unsafe.Pointer(buf))
	C.zktf_identity_operation_builder_previous(b.ptr, buf, length)
	return b
}

// Sequence sets the operation sequence number.
func (b *IdentityOperationBuilder) Sequence(seq uint32) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_sequence(b.ptr, C.uint32_t(seq))
	return b
}

// Timestamp sets the operation timestamp.
func (b *IdentityOperationBuilder) Timestamp(unix int64) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_timestamp(b.ptr, C.int64_t(unix))
	return b
}

// Anchor attaches a biometric anchor and nonce.
func (b *IdentityOperationBuilder) Anchor(anchor, nonce []byte) *IdentityOperationBuilder {
	aBuf, _ := cbytes(anchor)
	defer free(unsafe.Pointer(aBuf))
	nBuf, _ := cbytes(nonce)
	defer free(unsafe.Pointer(nBuf))
	C.zktf_identity_operation_builder_anchor(b.ptr, aBuf, nBuf)
	return b
}

// Commitment sets a commitment value attached to the operation.
func (b *IdentityOperationBuilder) Commitment(commitment []byte) *IdentityOperationBuilder {
	buf, _ := cbytes(commitment)
	defer free(unsafe.Pointer(buf))
	C.zktf_identity_operation_builder_commitment(b.ptr, buf)
	return b
}

// Threshold sets the threshold required to satisfy a given role.
func (b *IdentityOperationBuilder) Threshold(role IdentityKeyRole, threshold uint64) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_threshold(b.ptr, C.zktf_identity_key_role(role), C.uint64_t(threshold))
	return b
}

// Weight assigns a weight to a signing key for a given role.
func (b *IdentityOperationBuilder) Weight(key *SigningPublicKey, role IdentityKeyRole, weight uint64) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_weight(b.ptr, key.ptr, C.zktf_identity_key_role(role), C.uint64_t(weight))
	return b
}

// SigningGrantEmbedded grants a signing key the given roles, with the key embedded.
func (b *IdentityOperationBuilder) SigningGrantEmbedded(key *SigningPublicKey, roles IdentityKeyRole) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_signing_grant_embedded(b.ptr, key.ptr, C.zktf_identity_key_role(roles))
	return b
}

// SigningGrantReferenced grants a signing key the given roles via a referenced description.
func (b *IdentityOperationBuilder) SigningGrantReferenced(method uint16, controller, key *SigningPublicKey, commitment []byte, roles IdentityKeyRole) *IdentityOperationBuilder {
	buf, _ := cbytes(commitment)
	defer free(unsafe.Pointer(buf))
	C.zktf_identity_operation_builder_signing_grant_referenced(
		b.ptr,
		C.uint16_t(method),
		controller.ptr,
		key.ptr,
		buf,
		C.zktf_identity_key_role(roles),
	)
	return b
}

// SigningModify modifies the roles of an existing signing key.
func (b *IdentityOperationBuilder) SigningModify(key *SigningPublicKey, roles IdentityKeyRole) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_signing_modify(b.ptr, key.ptr, C.zktf_identity_key_role(roles))
	return b
}

// SigningRevoke revokes a signing key, effective from the given timestamp.
func (b *IdentityOperationBuilder) SigningRevoke(key *SigningPublicKey, effectiveFromUnix int64) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_signing_revoke(b.ptr, key.ptr, C.int64_t(effectiveFromUnix))
	return b
}

// ExchangeGrantEmbedded grants an exchange key the given roles, embedded.
func (b *IdentityOperationBuilder) ExchangeGrantEmbedded(key *ExchangePublicKey, roles IdentityKeyRole) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_exchange_grant_embedded(b.ptr, key.ptr, C.zktf_identity_key_role(roles))
	return b
}

// ExchangeModify modifies the roles of an existing exchange key.
func (b *IdentityOperationBuilder) ExchangeModify(key *ExchangePublicKey, roles IdentityKeyRole) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_exchange_modify(b.ptr, key.ptr, C.zktf_identity_key_role(roles))
	return b
}

// ExchangeRevoke revokes an exchange key, effective from the given timestamp.
func (b *IdentityOperationBuilder) ExchangeRevoke(key *ExchangePublicKey, effectiveFromUnix int64) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_exchange_revoke(b.ptr, key.ptr, C.int64_t(effectiveFromUnix))
	return b
}

// Recover stages a recovery operation effective from the given timestamp.
func (b *IdentityOperationBuilder) Recover(effectiveFromUnix int64) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_recover(b.ptr, C.int64_t(effectiveFromUnix))
	return b
}

// Deactivate stages a deactivation effective from the given timestamp.
func (b *IdentityOperationBuilder) Deactivate(effectiveFromUnix int64) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_deactivate(b.ptr, C.int64_t(effectiveFromUnix))
	return b
}

// SignWith records the signing key for the operation.
func (b *IdentityOperationBuilder) SignWith(signer *SigningPublicKey) *IdentityOperationBuilder {
	C.zktf_identity_operation_builder_sign_with(b.ptr, signer.ptr)
	return b
}

// Finish finalizes the operation.
func (b *IdentityOperationBuilder) Finish() (*IdentityOperation, error) {
	var out *C.zktf_identity_operation
	if err := status(C.zktf_identity_operation_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newIdentityOperation(out), nil
}

// IdentityOperationDecode decodes an encoded identity operation for the given document address.
func IdentityOperationDecode(documentAddress *SigningPublicKey, data []byte) (*IdentityOperation, error) {
	buf, length := cbytes(data)
	defer free(unsafe.Pointer(buf))
	var out *C.zktf_identity_operation
	if err := status(C.zktf_identity_operation_decode(documentAddress.ptr, buf, length, &out)); err != nil {
		return nil, err
	}
	return newIdentityOperation(out), nil
}

// Encode returns the encoded bytes of the operation.
func (o *IdentityOperation) Encode() ([]byte, error) {
	var buf *C.zktf_bytes_buffer
	if err := status(C.zktf_identity_operation_encode(o.ptr, &buf)); err != nil {
		return nil, err
	}
	return goBytesFromBuffer(buf), nil
}

// Sequence returns the operation sequence number.
func (o *IdentityOperation) Sequence() uint32 {
	return uint32(C.zktf_identity_operation_sequence(o.ptr))
}

// Hash returns the 32-byte operation hash.
func (o *IdentityOperation) Hash() []byte {
	return C.GoBytes(unsafe.Pointer(C.zktf_identity_operation_hash(o.ptr)), operationHashLen)
}

// SignedBy reports whether the operation has been signed by the given key.
func (o *IdentityOperation) SignedBy(signer *SigningPublicKey) bool {
	return bool(C.zktf_identity_operation_signed_by(o.ptr, signer.ptr))
}

// Merge merges signatures from another operation into this one.
func (o *IdentityOperation) Merge(other *IdentityOperation) error {
	return status(C.zktf_identity_operation_merge(o.ptr, other.ptr))
}

// Actions returns the actions described by the operation.
func (o *IdentityOperation) Actions() []*OperationAction {
	c := C.zktf_identity_operation_actions(o.ptr)
	if c == nil {
		return nil
	}
	defer C.zktf_collection_identity_operation_action_destroy(c)
	n := int(C.zktf_collection_identity_operation_action_len(c))
	out := make([]*OperationAction, n)
	for i := 0; i < n; i++ {
		out[i] = newOperationAction(C.zktf_collection_identity_operation_action_at(c, C.size_t(i)))
	}
	return out
}

// OperationAction is one of the actions inside an identity operation.
type OperationAction struct {
	ptr *C.zktf_identity_operation_action
}

func newOperationAction(ptr *C.zktf_identity_operation_action) *OperationAction {
	if ptr == nil {
		return nil
	}
	a := &OperationAction{ptr: ptr}
	runtime.AddCleanup(a, func(ptr *C.zktf_identity_operation_action) {
		C.zktf_identity_operation_action_destroy(ptr)
	}, a.ptr)
	return a
}

// Kind returns the kind of action.
func (a *OperationAction) Kind() OperationActionKind {
	return OperationActionKind(C.zktf_identity_operation_action_type_of(a.ptr))
}

// Roles returns the roles assigned by a grant or modify action.
func (a *OperationAction) Roles() IdentityKeyRole {
	return IdentityKeyRole(C.zktf_identity_operation_action_roles(a.ptr))
}

// EffectiveFrom returns the unix timestamp the action takes effect.
func (a *OperationAction) EffectiveFrom() int64 {
	return int64(C.zktf_identity_operation_action_from(a.ptr))
}

// DescriptionKind returns the kind of description on the action.
func (a *OperationAction) DescriptionKind() OperationDescriptionKind {
	return OperationDescriptionKind(C.zktf_identity_operation_action_description_type(a.ptr))
}

// DescriptionEmbedded returns the embedded description, or nil.
func (a *OperationAction) DescriptionEmbedded() *OperationDescriptionEmbedded {
	return newOperationDescriptionEmbedded(C.zktf_identity_operation_action_description_as_embedded(a.ptr))
}

// DescriptionReference returns the reference description, or nil.
func (a *OperationAction) DescriptionReference() *OperationDescriptionReference {
	return newOperationDescriptionReference(C.zktf_identity_operation_action_description_as_reference(a.ptr))
}

// OperationDescriptionEmbedded describes a key embedded in an action.
type OperationDescriptionEmbedded struct {
	ptr *C.zktf_identity_operation_description_embedded
}

func newOperationDescriptionEmbedded(ptr *C.zktf_identity_operation_description_embedded) *OperationDescriptionEmbedded {
	if ptr == nil {
		return nil
	}
	d := &OperationDescriptionEmbedded{ptr: ptr}
	runtime.AddCleanup(d, func(ptr *C.zktf_identity_operation_description_embedded) {
		C.zktf_identity_operation_description_embedded_destroy(ptr)
	}, d.ptr)
	return d
}

// AddressType returns the type of address (signing or exchange).
func (d *OperationDescriptionEmbedded) AddressType() KeypairType {
	return KeypairType(C.zktf_identity_operation_description_embedded_address_type(d.ptr))
}

// AddressAsSigning returns the address as a signing key, or nil if exchange.
func (d *OperationDescriptionEmbedded) AddressAsSigning() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_identity_operation_description_embedded_address_as_signing(d.ptr))
}

// AddressAsExchange returns the address as an exchange key, or nil if signing.
func (d *OperationDescriptionEmbedded) AddressAsExchange() *ExchangePublicKey {
	return newExchangePublicKey(C.zktf_identity_operation_description_embedded_address_as_exchange(d.ptr))
}

// Controller returns the controller address, or nil.
func (d *OperationDescriptionEmbedded) Controller() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_identity_operation_description_embedded_controller(d.ptr))
}

// OperationDescriptionReference describes a key referenced by another method.
type OperationDescriptionReference struct {
	ptr *C.zktf_identity_operation_description_reference
}

func newOperationDescriptionReference(ptr *C.zktf_identity_operation_description_reference) *OperationDescriptionReference {
	if ptr == nil {
		return nil
	}
	d := &OperationDescriptionReference{ptr: ptr}
	runtime.AddCleanup(d, func(ptr *C.zktf_identity_operation_description_reference) {
		C.zktf_identity_operation_description_reference_destroy(ptr)
	}, d.ptr)
	return d
}

// Method returns the DID method used by the reference.
func (d *OperationDescriptionReference) Method() AddressMethod {
	return AddressMethod(C.zktf_identity_operation_description_reference_method(d.ptr))
}

// AddressType returns the type of address (signing or exchange).
func (d *OperationDescriptionReference) AddressType() KeypairType {
	return KeypairType(C.zktf_identity_operation_description_reference_address_type(d.ptr))
}

// AddressAsSigning returns the address as a signing key, or nil if exchange.
func (d *OperationDescriptionReference) AddressAsSigning() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_identity_operation_description_reference_address_as_signing(d.ptr))
}

// AddressAsExchange returns the address as an exchange key, or nil if signing.
func (d *OperationDescriptionReference) AddressAsExchange() *ExchangePublicKey {
	return newExchangePublicKey(C.zktf_identity_operation_description_reference_address_as_exchange(d.ptr))
}

// Controller returns the controller address, or nil.
func (d *OperationDescriptionReference) Controller() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_identity_operation_description_reference_controller(d.ptr))
}

// IdentityLookup builds a query for identities.
type IdentityLookup struct {
	ptr *C.zktf_identity_lookup
}

// NewIdentityLookup initializes an identity lookup query.
func NewIdentityLookup() *IdentityLookup {
	ptr := C.zktf_identity_lookup_init()
	l := &IdentityLookup{ptr: ptr}
	runtime.AddCleanup(l, func(ptr *C.zktf_identity_lookup) {
		C.zktf_identity_lookup_destroy(ptr)
	}, l.ptr)
	return l
}

// ByKey restricts the lookup to identities the given key is associated with.
func (l *IdentityLookup) ByKey(key *SigningPublicKey) *IdentityLookup {
	C.zktf_identity_lookup_by_key(l.ptr, key.ptr)
	return l
}

// IdentityExecute publishes an identity operation via callback, returning once
// the result has been delivered.
func (a *Account) IdentityExecute(operation *IdentityOperation, timeout time.Duration) error {
	fut := C.zktf_account_identity_execute(a.ptr, operation.ptr)

	return AwaitStatus(fut, timeout)
}

// IdentitySign signs an identity operation with this account's keys.
func (a *Account) IdentitySign(operation *IdentityOperation) error {
	return status(C.zktf_account_identity_sign(a.ptr, operation.ptr))
}

// IdentityLookup returns DID addresses matching the lookup query.
func (a *Account) IdentityLookup(lookup *IdentityLookup) ([]*DIDAddress, error) {
	var c *C.zktf_collection_did_address
	if err := status(C.zktf_account_identity_lookup(a.ptr, lookup.ptr, &c)); err != nil {
		return nil, err
	}
	return didAddressesFrom(c), nil
}
