package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

import "runtime"

// KeyPackageEvent wraps a zktf_key_package wire event delivered to OnGroup as
// an invite. It carries the routing fields plus a conversion to the MLS-level
// crypto key package usable with Account.Establish.
type KeyPackageEvent struct {
	ptr *C.zktf_key_package
}

func newKeyPackageEvent(ptr *C.zktf_key_package) *KeyPackageEvent {
	if ptr == nil {
		return nil
	}
	e := &KeyPackageEvent{ptr: ptr}
	runtime.AddCleanup(e, func(ptr *C.zktf_key_package) {
		C.zktf_key_package_destroy(ptr)
	}, e.ptr)
	return e
}

// FromAddress returns the sender's address.
func (e *KeyPackageEvent) FromAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_key_package_from_address(e.ptr))
}

// ToAddress returns the recipient address.
func (e *KeyPackageEvent) ToAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_key_package_to_address(e.ptr))
}

// Sequence returns the event's sequence number.
func (e *KeyPackageEvent) Sequence() uint64 { return uint64(C.zktf_key_package_sequence(e.ptr)) }

// Timestamp returns the event's unix timestamp.
func (e *KeyPackageEvent) Timestamp() int64 { return int64(C.zktf_key_package_timestamp(e.ptr)) }

// CryptoKeyPackage extracts the MLS key package suitable for Account.Establish.
func (e *KeyPackageEvent) CryptoKeyPackage() *CryptoKeyPackage {
	return newCryptoKeyPackage(C.zktf_key_package_crypto_key_package(e.ptr))
}

// WelcomeEvent wraps a zktf_welcome wire event delivered to OnGroup.
type WelcomeEvent struct {
	ptr *C.zktf_welcome
}

func newWelcomeEvent(ptr *C.zktf_welcome) *WelcomeEvent {
	if ptr == nil {
		return nil
	}
	e := &WelcomeEvent{ptr: ptr}
	runtime.AddCleanup(e, func(ptr *C.zktf_welcome) {
		C.zktf_welcome_destroy(ptr)
	}, e.ptr)
	return e
}

// FromAddress returns the sender's address.
func (e *WelcomeEvent) FromAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_welcome_from_address(e.ptr))
}

// ToAddress returns the recipient address.
func (e *WelcomeEvent) ToAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_welcome_to_address(e.ptr))
}

// Sequence returns the event's sequence number.
func (e *WelcomeEvent) Sequence() uint64 { return uint64(C.zktf_welcome_sequence(e.ptr)) }

// Timestamp returns the event's unix timestamp.
func (e *WelcomeEvent) Timestamp() int64 { return int64(C.zktf_welcome_timestamp(e.ptr)) }

// CryptoWelcome extracts the MLS welcome suitable for Account.Accept.
func (e *WelcomeEvent) CryptoWelcome() *CryptoWelcome {
	return newCryptoWelcome(C.zktf_welcome_crypto_welcome(e.ptr))
}

// CommitEvent wraps a zktf_commit wire event.
type CommitEvent struct {
	ptr *C.zktf_commit
}

func newCommitEvent(ptr *C.zktf_commit) *CommitEvent {
	if ptr == nil {
		return nil
	}
	e := &CommitEvent{ptr: ptr}
	runtime.AddCleanup(e, func(ptr *C.zktf_commit) {
		C.zktf_commit_destroy(ptr)
	}, e.ptr)
	return e
}

// FromAddress returns the sender's address.
func (e *CommitEvent) FromAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_commit_from_address(e.ptr))
}

// ToAddress returns the recipient address.
func (e *CommitEvent) ToAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_commit_to_address(e.ptr))
}

// Sequence returns the event's sequence number.
func (e *CommitEvent) Sequence() uint64 { return uint64(C.zktf_commit_sequence(e.ptr)) }

// Timestamp returns the event's unix timestamp.
func (e *CommitEvent) Timestamp() int64 { return int64(C.zktf_commit_timestamp(e.ptr)) }

// ProposalEvent wraps a zktf_proposal wire event.
type ProposalEvent struct {
	ptr *C.zktf_proposal
}

func newProposalEvent(ptr *C.zktf_proposal) *ProposalEvent {
	if ptr == nil {
		return nil
	}
	e := &ProposalEvent{ptr: ptr}
	runtime.AddCleanup(e, func(ptr *C.zktf_proposal) {
		C.zktf_proposal_destroy(ptr)
	}, e.ptr)
	return e
}

// FromAddress returns the sender's address.
func (e *ProposalEvent) FromAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_proposal_from_address(e.ptr))
}

// ToAddress returns the recipient address.
func (e *ProposalEvent) ToAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_proposal_to_address(e.ptr))
}

// Sequence returns the event's sequence number.
func (e *ProposalEvent) Sequence() uint64 { return uint64(C.zktf_proposal_sequence(e.ptr)) }

// Timestamp returns the event's unix timestamp.
func (e *ProposalEvent) Timestamp() int64 { return int64(C.zktf_proposal_timestamp(e.ptr)) }

// DroppedEvent wraps a zktf_dropped_event handle carried by a STATUS_EVENT_DROPPED.
type DroppedEvent struct {
	ptr *C.zktf_dropped_event
}

func newDroppedEvent(ptr *C.zktf_dropped_event) *DroppedEvent {
	if ptr == nil {
		return nil
	}
	e := &DroppedEvent{ptr: ptr}
	runtime.AddCleanup(e, func(ptr *C.zktf_dropped_event) {
		C.zktf_dropped_event_destroy(ptr)
	}, e.ptr)
	return e
}

// FromAddress returns the sender's address.
func (e *DroppedEvent) FromAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_dropped_event_from_address(e.ptr))
}

// ToAddress returns the recipient address.
func (e *DroppedEvent) ToAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_dropped_event_to_address(e.ptr))
}

// FromSequence returns the starting sequence number of the dropped range.
func (e *DroppedEvent) FromSequence() uint64 {
	return uint64(C.zktf_dropped_event_from_sequence(e.ptr))
}

// ToSequence returns the ending sequence number of the dropped range.
func (e *DroppedEvent) ToSequence() uint64 {
	return uint64(C.zktf_dropped_event_to_sequence(e.ptr))
}

// Reason returns the reason the messages were dropped.
func (e *DroppedEvent) Reason() error {
	return status(C.zktf_dropped_event_reason(e.ptr))
}
