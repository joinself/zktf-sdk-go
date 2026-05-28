// Package event provides events delivered to an account's callbacks: account
// status events, group events (invite/welcome/commit/proposal), and background
// workflow events.
package event

import (
	"time"

	"github.com/joinself/zktf-sdk-go/crypto"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// EventKind identifies the kind of a status Event.
type EventKind uint32

const (
	EventConnected    EventKind = EventKind(ffi.StatusEventConnected)
	EventDisconnected EventKind = EventKind(ffi.StatusEventDisconnected)
	EventAcknowledged EventKind = EventKind(ffi.StatusEventAcknowledged)
	EventSendFailed   EventKind = EventKind(ffi.StatusEventSendFailed)
	EventDropped      EventKind = EventKind(ffi.StatusEventDropped)
)

// Event is an account-level status event. Dispatch on Kind() and use the
// matching accessor.
type Event struct {
	h *ffi.StatusEvent
}

// Group is a polymorphic group event. Dispatch on Kind() and call the
// matching extractor (Invite/Welcome/Commit/Proposal).
type Group struct {
	h *ffi.GroupEvent
}

// Workflow is a workflow event.
type Workflow struct {
	h *ffi.WorkflowEvent
}

// KeyPackage is a wire key-package invite event.
type KeyPackage struct {
	h *ffi.KeyPackageEvent
}

// Welcome is a wire welcome event.
type Welcome struct {
	h *ffi.WelcomeEvent
}

// Commit is a wire commit event.
type Commit struct {
	h *ffi.CommitEvent
}

// Proposal is a wire proposal event.
type Proposal struct {
	h *ffi.ProposalEvent
}

// Dropped describes a range of messages dropped between two addresses.
type Dropped struct {
	h *ffi.DroppedEvent
}

func init() {
	ffi.StatusEventOf = func(o any) *ffi.StatusEvent { return o.(*Event).h }
	ffi.ToStatusEvent = func(h *ffi.StatusEvent) any { return &Event{h: h} }

	ffi.GroupEventOf = func(o any) *ffi.GroupEvent { return o.(*Group).h }
	ffi.ToGroupEvent = func(h *ffi.GroupEvent) any { return &Group{h: h} }

	ffi.WorkflowEventOf = func(o any) *ffi.WorkflowEvent { return o.(*Workflow).h }
	ffi.ToWorkflowEvent = func(h *ffi.WorkflowEvent) any { return &Workflow{h: h} }

	ffi.KeyPackageEventOf = func(o any) *ffi.KeyPackageEvent { return o.(*KeyPackage).h }
	ffi.ToKeyPackageEvent = func(h *ffi.KeyPackageEvent) any { return &KeyPackage{h: h} }

	ffi.WelcomeEventOf = func(o any) *ffi.WelcomeEvent { return o.(*Welcome).h }
	ffi.ToWelcomeEvent = func(h *ffi.WelcomeEvent) any { return &Welcome{h: h} }

	ffi.CommitEventOf = func(o any) *ffi.CommitEvent { return o.(*Commit).h }
	ffi.ToCommitEvent = func(h *ffi.CommitEvent) any { return &Commit{h: h} }

	ffi.ProposalEventOf = func(o any) *ffi.ProposalEvent { return o.(*Proposal).h }
	ffi.ToProposalEvent = func(h *ffi.ProposalEvent) any { return &Proposal{h: h} }

	ffi.DroppedEventOf = func(o any) *ffi.DroppedEvent { return o.(*Dropped).h }
	ffi.ToDroppedEvent = func(h *ffi.DroppedEvent) any { return &Dropped{h: h} }
}

// Kind returns the kind of event.
func (e *Event) Kind() EventKind { return EventKind(e.h.Kind()) }

// DisconnectReason returns the reason for a disconnected event, or nil.
func (e *Event) DisconnectReason() error { return e.h.DisconnectReason() }

// SendError returns the error for a send-failed event, or nil.
func (e *Event) SendError() error { return e.h.SendError() }

// ReferenceID returns the message id referenced by acknowledged/send-failed.
func (e *Event) ReferenceID() []byte { return e.h.ReferenceID() }

// Dropped returns the dropped-event details for a dropped event, or nil.
func (e *Event) Dropped() *Dropped {
	d := e.h.Dropped()
	if d == nil {
		return nil
	}

	return &Dropped{h: d}
}

// GroupKind identifies the kind of a Group event.
type GroupKind uint32

const (
	GroupInvite   GroupKind = GroupKind(ffi.GroupEventInvite)
	GroupWelcome  GroupKind = GroupKind(ffi.GroupEventWelcome)
	GroupCommit   GroupKind = GroupKind(ffi.GroupEventCommit)
	GroupProposal GroupKind = GroupKind(ffi.GroupEventProposal)
)

// Kind returns the kind of group event.
func (g *Group) Kind() GroupKind { return GroupKind(g.h.Kind()) }

// Invite extracts the wire key package for an Invite event.
func (g *Group) Invite() *KeyPackage { return &KeyPackage{h: g.h.Invite()} }

// Welcome extracts the wire welcome for a Welcome event.
func (g *Group) Welcome() *Welcome { return &Welcome{h: g.h.Welcome()} }

// Commit extracts the wire commit for a Commit event.
func (g *Group) Commit() *Commit { return &Commit{h: g.h.Commit()} }

// Proposal extracts the wire proposal for a Proposal event.
func (g *Group) Proposal() *Proposal { return &Proposal{h: g.h.Proposal()} }

// From returns the sender's address.
func (e *KeyPackage) From() *signing.PublicKey {
	return ffi.ToSigningPublicKey(e.h.FromAddress()).(*signing.PublicKey)
}

// To returns the recipient's address.
func (e *KeyPackage) To() *signing.PublicKey {
	return ffi.ToSigningPublicKey(e.h.ToAddress()).(*signing.PublicKey)
}

// Sequence returns the event sequence number.
func (e *KeyPackage) Sequence() uint64 { return e.h.Sequence() }

// Timestamp returns the event time.
func (e *KeyPackage) Timestamp() time.Time { return time.Unix(e.h.Timestamp(), 0) }

// KeyPackage extracts the MLS key package suitable for Account.GroupEstablish.
func (e *KeyPackage) KeyPackage() *crypto.KeyPackage {
	return ffi.ToCryptoKeyPackage(e.h.CryptoKeyPackage()).(*crypto.KeyPackage)
}

// From returns the sender's address.
func (e *Welcome) From() *signing.PublicKey {
	return ffi.ToSigningPublicKey(e.h.FromAddress()).(*signing.PublicKey)
}

// To returns the recipient's address.
func (e *Welcome) To() *signing.PublicKey {
	return ffi.ToSigningPublicKey(e.h.ToAddress()).(*signing.PublicKey)
}

// Sequence returns the event sequence number.
func (e *Welcome) Sequence() uint64 { return e.h.Sequence() }

// Timestamp returns the event time.
func (e *Welcome) Timestamp() time.Time { return time.Unix(e.h.Timestamp(), 0) }

// Welcome extracts the MLS welcome suitable for Account.GroupAccept.
func (e *Welcome) Welcome() *crypto.Welcome {
	return ffi.ToCryptoWelcome(e.h.CryptoWelcome()).(*crypto.Welcome)
}

// From returns the sender's address.
func (e *Commit) From() *signing.PublicKey {
	return ffi.ToSigningPublicKey(e.h.FromAddress()).(*signing.PublicKey)
}

// To returns the recipient's address.
func (e *Commit) To() *signing.PublicKey {
	return ffi.ToSigningPublicKey(e.h.ToAddress()).(*signing.PublicKey)
}

// Sequence returns the event sequence number.
func (e *Commit) Sequence() uint64 { return e.h.Sequence() }

// Timestamp returns the event time.
func (e *Commit) Timestamp() time.Time { return time.Unix(e.h.Timestamp(), 0) }

// From returns the sender's address.
func (e *Proposal) From() *signing.PublicKey {
	return ffi.ToSigningPublicKey(e.h.FromAddress()).(*signing.PublicKey)
}

// To returns the recipient's address.
func (e *Proposal) To() *signing.PublicKey {
	return ffi.ToSigningPublicKey(e.h.ToAddress()).(*signing.PublicKey)
}

// Sequence returns the event sequence number.
func (e *Proposal) Sequence() uint64 { return e.h.Sequence() }

// Timestamp returns the event time.
func (e *Proposal) Timestamp() time.Time { return time.Unix(e.h.Timestamp(), 0) }

// From returns the sender's address.
func (e *Dropped) From() *signing.PublicKey {
	return ffi.ToSigningPublicKey(e.h.FromAddress()).(*signing.PublicKey)
}

// To returns the recipient address.
func (e *Dropped) To() *signing.PublicKey {
	return ffi.ToSigningPublicKey(e.h.ToAddress()).(*signing.PublicKey)
}

// FromSequence returns the first sequence number in the dropped range.
func (e *Dropped) FromSequence() uint64 { return e.h.FromSequence() }

// ToSequence returns the last sequence number in the dropped range.
func (e *Dropped) ToSequence() uint64 { return e.h.ToSequence() }

// Reason returns the reason the messages were dropped.
func (e *Dropped) Reason() error { return e.h.Reason() }

// WorkflowKind identifies a workflow event.
type WorkflowKind uint32

const (
	WorkflowCompleted  WorkflowKind = WorkflowKind(ffi.WorkflowEventCompleted)
	WorkflowTaskFailed WorkflowKind = WorkflowKind(ffi.WorkflowEventTaskFailed)
)

// Kind returns the kind of workflow event.
func (e *Workflow) Kind() WorkflowKind { return WorkflowKind(e.h.Kind()) }

// WorkflowID returns the workflow id.
func (e *Workflow) WorkflowID() []byte { return e.h.WorkflowID() }

// TaskID returns the failed task id (for TaskFailed events), or nil.
func (e *Workflow) TaskID() []byte { return e.h.TaskID() }

// Reason returns a human-readable reason, or "".
func (e *Workflow) Reason() string { return e.h.Reason() }

// Attempt returns the attempt number.
func (e *Workflow) Attempt() uint32 { return e.h.Attempt() }

// WillRetry reports whether the failed task will be retried.
func (e *Workflow) WillRetry() bool { return e.h.WillRetry() }
