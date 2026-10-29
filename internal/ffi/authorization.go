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

type AuthorizationReasonKind int

const (
	AuthorizationReasonRejected       AuthorizationReasonKind = C.AUTHORIZATION_REASON_REJECTED
	AuthorizationReasonQuorumNotMet   AuthorizationReasonKind = C.AUTHORIZATION_REASON_QUORUM_NOT_MET
	AuthorizationReasonPersonUnproven AuthorizationReasonKind = C.AUTHORIZATION_REASON_PERSON_UNPROVEN
)

type AuthorizationGap int

const (
	AuthorizationGapNone            AuthorizationGap = C.AUTHORIZATION_GAP_NONE
	AuthorizationGapAnchorMissing   AuthorizationGap = C.AUTHORIZATION_GAP_ANCHOR_MISSING
	AuthorizationGapAnchorMismatch  AuthorizationGap = C.AUTHORIZATION_GAP_ANCHOR_MISMATCH
	AuthorizationGapLivenessMissing AuthorizationGap = C.AUTHORIZATION_GAP_LIVENESS_MISSING
	AuthorizationGapEvidenceRevoked AuthorizationGap = C.AUTHORIZATION_GAP_EVIDENCE_REVOKED
)

type AuthorizationRequirementKind int

const (
	AuthorizationRequirementSignatures AuthorizationRequirementKind = C.AUTHORIZATION_REQUIREMENT_SIGNATURES
	AuthorizationRequirementAnchor     AuthorizationRequirementKind = C.AUTHORIZATION_REQUIREMENT_ANCHOR
	AuthorizationRequirementLiveness   AuthorizationRequirementKind = C.AUTHORIZATION_REQUIREMENT_LIVENESS
)

type AuthorizationReason struct {
	Kind      AuthorizationReasonKind
	Error     error
	Document  *DIDAddress
	Role      uint64
	Weight    uint64
	Threshold uint64
	Signer    *SigningPublicKey
	Identity  *DIDAddress
	Gap       AuthorizationGap
}

type AuthorizationRequirement struct {
	Kind       AuthorizationRequirementKind
	Document   *DIDAddress
	Role       uint64
	Weight     uint64
	Candidates []*SigningPublicKey
	Identity   *DIDAddress
	Commitment []byte
	Hash       []byte
}

type AuthorizationPerson struct {
	Identity   *DIDAddress
	Commitment []byte
	Anchor     *VerifiableCredential
	Evidence   *VerifiableCredential
}

type AuthorizationLink struct {
	From       *DIDAddress
	To         *DIDAddress
	Key        *SigningPublicKey
	Commitment []byte
	Roles      uint64
}

type AuthorizationVerdict struct {
	Authorized bool
	Reasons    []AuthorizationReason
}

type OperationVerdict struct {
	Verdict  AuthorizationVerdict
	Document *IdentityDocument
	Missing  []AuthorizationRequirement
}

type AuthorizationChain struct {
	ptr *C.zktf_credential_authorization_chain
}

func newAuthorizationChain(ptr *C.zktf_credential_authorization_chain) *AuthorizationChain {
	if ptr == nil {
		return nil
	}
	c := &AuthorizationChain{ptr: ptr}
	runtime.AddCleanup(c, func(ptr *C.zktf_credential_authorization_chain) {
		C.zktf_credential_authorization_chain_destroy(ptr)
	}, c.ptr)
	return c
}

func (g *CredentialGraph) ValidOperationFor(document *DIDAddress, operation *IdentityOperation) (*OperationVerdict, error) {
	defer runtime.KeepAlive(g)
	defer runtime.KeepAlive(document)
	defer runtime.KeepAlive(operation)

	var address *C.zktf_did_address
	if document != nil {
		address = document.ptr
	}

	var ptr *C.zktf_credential_operation_verdict
	if err := status(C.zktf_credential_graph_valid_operation_for(g.ptr, address, operation.ptr, &ptr)); err != nil {
		return nil, err
	}
	defer C.zktf_credential_operation_verdict_destroy(ptr)

	verdict := C.zktf_credential_operation_verdict_verdict(ptr)
	defer C.zktf_credential_verdict_destroy(verdict)

	missing := C.zktf_credential_operation_verdict_missing(ptr)
	defer C.zktf_credential_authorization_requirements_destroy(missing)

	return &OperationVerdict{
		Verdict:  authorizationVerdictFrom(verdict),
		Document: newIdentityDocument(C.zktf_credential_operation_verdict_document(ptr)),
		Missing:  authorizationRequirementsFrom(missing),
	}, nil
}

func (g *CredentialGraph) AuthorizationFor(hash []byte) (*AuthorizationChain, error) {
	defer runtime.KeepAlive(g)
	buf, length := cbytes(hash)
	defer free(unsafe.Pointer(buf))

	var ptr *C.zktf_credential_authorization_chain
	if err := status(C.zktf_credential_graph_authorization_for(g.ptr, buf, C.uintptr_t(length), &ptr)); err != nil {
		return nil, err
	}
	return newAuthorizationChain(ptr), nil
}

func (c *AuthorizationChain) People() []AuthorizationPerson {
	defer runtime.KeepAlive(c)
	count := int(C.zktf_credential_authorization_chain_people_len(c.ptr))
	out := make([]AuthorizationPerson, 0, count)
	for i := 0; i < count; i++ {
		p := C.zktf_credential_authorization_chain_person_at(c.ptr, C.size_t(i))
		out = append(out, AuthorizationPerson{
			Identity:   newDIDAddress(C.zktf_credential_authorization_person_identity(p)),
			Commitment: goBytesFromBuffer(C.zktf_credential_authorization_person_commitment(p)),
			Anchor:     newVerifiableCredential(C.zktf_credential_authorization_person_anchor(p)),
			Evidence:   newVerifiableCredential(C.zktf_credential_authorization_person_evidence(p)),
		})
		C.zktf_credential_authorization_person_destroy(p)
	}
	return out
}

func (c *AuthorizationChain) Links() []AuthorizationLink {
	defer runtime.KeepAlive(c)
	count := int(C.zktf_credential_authorization_chain_links_len(c.ptr))
	out := make([]AuthorizationLink, 0, count)
	for i := 0; i < count; i++ {
		l := C.zktf_credential_authorization_chain_link_at(c.ptr, C.size_t(i))
		out = append(out, AuthorizationLink{
			From:       newDIDAddress(C.zktf_credential_authorization_link_from(l)),
			To:         newDIDAddress(C.zktf_credential_authorization_link_to(l)),
			Key:        newSigningPublicKey(C.zktf_credential_authorization_link_key(l)),
			Commitment: goBytesFromBuffer(C.zktf_credential_authorization_link_commitment(l)),
			Roles:      uint64(C.zktf_credential_authorization_link_roles(l)),
		})
		C.zktf_credential_authorization_link_destroy(l)
	}
	return out
}

func (c *AuthorizationChain) AuthorizedBy(party *DIDAddress) bool {
	defer runtime.KeepAlive(c)
	defer runtime.KeepAlive(party)
	return bool(C.zktf_credential_authorization_chain_authorized_by(c.ptr, party.ptr))
}

func (c *AuthorizationChain) Verdict(document *IdentityDocument, role uint64) (AuthorizationVerdict, error) {
	defer runtime.KeepAlive(c)
	defer runtime.KeepAlive(document)

	var ptr *C.zktf_credential_verdict
	if err := status(C.zktf_credential_authorization_chain_verdict(c.ptr, document.ptr, C.uint64_t(role), &ptr)); err != nil {
		return AuthorizationVerdict{}, err
	}
	defer C.zktf_credential_verdict_destroy(ptr)
	return authorizationVerdictFrom(ptr), nil
}

func (c *AuthorizationChain) Missing(document *IdentityDocument, role uint64) ([]AuthorizationRequirement, error) {
	defer runtime.KeepAlive(c)
	defer runtime.KeepAlive(document)

	var ptr *C.zktf_credential_authorization_requirements
	if err := status(C.zktf_credential_authorization_chain_missing(c.ptr, document.ptr, C.uint64_t(role), &ptr)); err != nil {
		return nil, err
	}
	defer C.zktf_credential_authorization_requirements_destroy(ptr)
	return authorizationRequirementsFrom(ptr), nil
}

func authorizationVerdictFrom(ptr *C.zktf_credential_verdict) AuthorizationVerdict {
	count := int(C.zktf_credential_verdict_reasons_len(ptr))
	reasons := make([]AuthorizationReason, 0, count)
	for i := 0; i < count; i++ {
		r := C.zktf_credential_verdict_reason_at(ptr, C.size_t(i))
		reasons = append(reasons, AuthorizationReason{
			Kind:      AuthorizationReasonKind(C.zktf_credential_authorization_reason_kind_of(r)),
			Error:     status(C.zktf_credential_authorization_reason_error(r)),
			Document:  newDIDAddress(C.zktf_credential_authorization_reason_document(r)),
			Role:      uint64(C.zktf_credential_authorization_reason_role(r)),
			Weight:    uint64(C.zktf_credential_authorization_reason_weight(r)),
			Threshold: uint64(C.zktf_credential_authorization_reason_threshold(r)),
			Signer:    newSigningPublicKey(C.zktf_credential_authorization_reason_signer(r)),
			Identity:  newDIDAddress(C.zktf_credential_authorization_reason_identity(r)),
			Gap:       AuthorizationGap(C.zktf_credential_authorization_reason_gap(r)),
		})
		C.zktf_credential_authorization_reason_destroy(r)
	}
	return AuthorizationVerdict{
		Authorized: bool(C.zktf_credential_verdict_authorized(ptr)),
		Reasons:    reasons,
	}
}

func authorizationRequirementsFrom(ptr *C.zktf_credential_authorization_requirements) []AuthorizationRequirement {
	count := int(C.zktf_credential_authorization_requirements_len(ptr))
	out := make([]AuthorizationRequirement, 0, count)
	for i := 0; i < count; i++ {
		r := C.zktf_credential_authorization_requirements_at(ptr, C.size_t(i))
		candidates := make([]*SigningPublicKey, int(C.zktf_credential_authorization_requirement_candidates_len(r)))
		for j := range candidates {
			candidates[j] = newSigningPublicKey(C.zktf_credential_authorization_requirement_candidate_at(r, C.size_t(j)))
		}
		out = append(out, AuthorizationRequirement{
			Kind:       AuthorizationRequirementKind(C.zktf_credential_authorization_requirement_kind_of(r)),
			Document:   newDIDAddress(C.zktf_credential_authorization_requirement_document(r)),
			Role:       uint64(C.zktf_credential_authorization_requirement_role(r)),
			Weight:     uint64(C.zktf_credential_authorization_requirement_weight(r)),
			Candidates: candidates,
			Identity:   newDIDAddress(C.zktf_credential_authorization_requirement_identity(r)),
			Commitment: goBytesFromBuffer(C.zktf_credential_authorization_requirement_commitment(r)),
			Hash:       goBytesFromBuffer(C.zktf_credential_authorization_requirement_hash(r)),
		})
		C.zktf_credential_authorization_requirement_destroy(r)
	}
	return out
}
