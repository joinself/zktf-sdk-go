package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// StatusEventType mirrors zktf_status_event_type.
type StatusEventType uint32

const (
	StatusEventConnected    StatusEventType = C.STATUS_EVENT_CONNECTED
	StatusEventDisconnected StatusEventType = C.STATUS_EVENT_DISCONNECTED
)

// StatusEvent wraps a zktf_status_event delivered to the on_status callback.
type StatusEvent struct {
	ptr *C.zktf_status_event
}

func newStatusEvent(ptr *C.zktf_status_event) *StatusEvent {
	if ptr == nil {
		return nil
	}
	e := &StatusEvent{ptr: ptr}
	runtime.AddCleanup(e, func(ptr *C.zktf_status_event) {
		C.zktf_status_event_destroy(ptr)
	}, e.ptr)
	return e
}

// Kind returns the kind of status event.
func (e *StatusEvent) Kind() StatusEventType {
	defer runtime.KeepAlive(e)
	return StatusEventType(C.zktf_status_event_kind(e.ptr))
}

// DisconnectReason returns the reason for a disconnect event as an error, or nil.
func (e *StatusEvent) DisconnectReason() error {
	defer runtime.KeepAlive(e)
	return status(C.zktf_status_event_disconnect_reason(e.ptr))
}

type EventKind uint32

const (
	EventMessage      EventKind = C.EVENT_MESSAGE
	EventGroup        EventKind = C.EVENT_GROUP
	EventAcknowledged EventKind = C.EVENT_ACKNOWLEDGED
	EventFailed       EventKind = C.EVENT_FAILED
	EventDropped      EventKind = C.EVENT_DROPPED
	EventPaired       EventKind = C.EVENT_PAIRED
)

type BatchEvent struct {
	Kind      EventKind
	Message   *Message
	Group     *GroupEvent
	Reference *Reference
	Err       error
	Dropped   *DroppedEvent
	Paired    *SigningPublicKey
}

func takeEventBatch(batch *C.zktf_event_batch) []BatchEvent {
	if batch == nil {
		return nil
	}
	defer C.zktf_event_batch_destroy(batch)

	n := int(C.zktf_event_batch_len(batch))
	events := make([]BatchEvent, 0, n)
	for i := 0; i < n; i++ {
		idx := C.size_t(i)
		ev := BatchEvent{Kind: EventKind(C.zktf_event_batch_kind(batch, idx))}
		switch ev.Kind {
		case EventMessage:
			ev.Message = newMessage(C.zktf_event_batch_message(batch, idx))
		case EventGroup:
			ev.Group = newGroupEvent(C.zktf_event_batch_group(batch, idx))
		case EventAcknowledged:
			ev.Reference = takeReference(C.zktf_event_batch_acknowledged(batch, idx))
		case EventFailed:
			failure := C.zktf_event_batch_failed(batch, idx)
			if failure != nil {
				ev.Reference = takeReference(C.zktf_failure_reference(failure))
				ev.Err = status(C.zktf_failure_error(failure))
				C.zktf_failure_destroy(failure)
			}
		case EventDropped:
			ev.Dropped = newDroppedEvent(C.zktf_event_batch_dropped(batch, idx))
		case EventPaired:
			ev.Paired = newSigningPublicKey(C.zktf_event_batch_paired(batch, idx))
		}
		events = append(events, ev)
	}

	return events
}

type ReferenceType uint32

const (
	ReferenceCommit     ReferenceType = C.REFERENCE_COMMIT
	ReferenceKeyPackage ReferenceType = C.REFERENCE_KEY_PACKAGE
	ReferenceMessage    ReferenceType = C.REFERENCE_MESSAGE
	ReferenceProposal   ReferenceType = C.REFERENCE_PROPOSAL
	ReferenceWelcome    ReferenceType = C.REFERENCE_WELCOME
)

type Reference struct {
	ID   []byte
	Type ReferenceType
	From *SigningPublicKey
	To   *SigningPublicKey
}

func takeReference(ref *C.zktf_reference) *Reference {
	if ref == nil {
		return nil
	}
	defer C.zktf_reference_destroy(ref)

	r := &Reference{
		Type: ReferenceType(C.zktf_reference_type_of(ref)),
		From: newSigningPublicKey(C.zktf_reference_from_address(ref)),
		To:   newSigningPublicKey(C.zktf_reference_to_address(ref)),
	}
	if id := C.zktf_reference_id(ref); id != nil {
		r.ID = C.GoBytes(unsafe.Pointer(id), messageIDLen)
	}
	return r
}

// GroupEventKind mirrors zktf_group_event_type.
type GroupEventKind uint32

const (
	GroupEventInvite   GroupEventKind = C.GROUP_EVENT_INVITE
	GroupEventWelcome  GroupEventKind = C.GROUP_EVENT_WELCOME
	GroupEventCommit   GroupEventKind = C.GROUP_EVENT_COMMIT
	GroupEventProposal GroupEventKind = C.GROUP_EVENT_PROPOSAL
)

// GroupEvent wraps a zktf_group_event delivered to on_group. Use Kind to
// dispatch and the per-kind accessors to extract the carried wire event.
type GroupEvent struct {
	ptr *C.zktf_group_event
}

func newGroupEvent(ptr *C.zktf_group_event) *GroupEvent {
	if ptr == nil {
		return nil
	}
	e := &GroupEvent{ptr: ptr}
	runtime.AddCleanup(e, func(ptr *C.zktf_group_event) {
		C.zktf_group_event_destroy(ptr)
	}, e.ptr)
	return e
}

// Kind returns the kind of group event.
func (e *GroupEvent) Kind() GroupEventKind {
	defer runtime.KeepAlive(e)
	return GroupEventKind(C.zktf_group_event_kind(e.ptr))
}

// Invite extracts the key package wire event (kind == GroupEventInvite).
func (e *GroupEvent) Invite() *KeyPackageEvent {
	defer runtime.KeepAlive(e)
	return newKeyPackageEvent(C.zktf_group_event_invite(e.ptr))
}

// Welcome extracts the welcome wire event (kind == GroupEventWelcome).
func (e *GroupEvent) Welcome() *WelcomeEvent {
	defer runtime.KeepAlive(e)
	return newWelcomeEvent(C.zktf_group_event_welcome(e.ptr))
}

// Commit extracts the commit wire event (kind == GroupEventCommit).
func (e *GroupEvent) Commit() *CommitEvent {
	defer runtime.KeepAlive(e)
	return newCommitEvent(C.zktf_group_event_commit(e.ptr))
}

// Proposal extracts the proposal wire event (kind == GroupEventProposal).
func (e *GroupEvent) Proposal() *ProposalEvent {
	defer runtime.KeepAlive(e)
	return newProposalEvent(C.zktf_group_event_proposal(e.ptr))
}

// WorkflowEventKind mirrors zktf_workflow_event_type.
type WorkflowEventKind uint32

const (
	WorkflowEventCompleted  WorkflowEventKind = C.WORKFLOW_EVENT_COMPLETED
	WorkflowEventTaskFailed WorkflowEventKind = C.WORKFLOW_EVENT_TASK_FAILED
)

// WorkflowEvent wraps a zktf_workflow_event delivered to on_workflow.
type WorkflowEvent struct {
	ptr *C.zktf_workflow_event
}

func newWorkflowEvent(ptr *C.zktf_workflow_event) *WorkflowEvent {
	if ptr == nil {
		return nil
	}
	e := &WorkflowEvent{ptr: ptr}
	runtime.AddCleanup(e, func(ptr *C.zktf_workflow_event) {
		C.zktf_workflow_event_destroy(ptr)
	}, e.ptr)
	return e
}

// Kind returns the kind of workflow event.
func (e *WorkflowEvent) Kind() WorkflowEventKind {
	defer runtime.KeepAlive(e)
	return WorkflowEventKind(C.zktf_workflow_event_kind(e.ptr))
}

// WorkflowID returns the workflow id bytes.
func (e *WorkflowEvent) WorkflowID() []byte {
	defer runtime.KeepAlive(e)
	n := C.zktf_workflow_event_workflow_id_len(e.ptr)
	if n == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(C.zktf_workflow_event_workflow_id_buf(e.ptr)), C.int(n))
}

// TaskID returns the task id bytes (for TaskFailed events), or nil.
func (e *WorkflowEvent) TaskID() []byte {
	defer runtime.KeepAlive(e)
	n := C.zktf_workflow_event_task_id_len(e.ptr)
	if n == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(C.zktf_workflow_event_task_id_buf(e.ptr)), C.int(n))
}

// Reason returns a human-readable reason for the event, or "" if unset.
func (e *WorkflowEvent) Reason() string {
	defer runtime.KeepAlive(e)
	n := C.zktf_workflow_event_reason_len(e.ptr)
	if n == 0 {
		return ""
	}
	return string(C.GoBytes(unsafe.Pointer(C.zktf_workflow_event_reason_buf(e.ptr)), C.int(n)))
}

// Attempt returns the attempt number (for TaskFailed events).
func (e *WorkflowEvent) Attempt() uint32 {
	defer runtime.KeepAlive(e)
	return uint32(C.zktf_workflow_event_attempt(e.ptr))
}

// WillRetry reports whether the failed task will be retried.
func (e *WorkflowEvent) WillRetry() bool {
	defer runtime.KeepAlive(e)
	return bool(C.zktf_workflow_event_will_retry(e.ptr))
}

func (e *WorkflowEvent) RootWorkflowID() []byte {
	defer runtime.KeepAlive(e)
	n := C.zktf_workflow_event_root_workflow_id_len(e.ptr)
	if n == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(C.zktf_workflow_event_root_workflow_id_buf(e.ptr)), C.int(n))
}
