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
	StatusEventAcknowledged StatusEventType = C.STATUS_EVENT_ACKNOWLEDGED
	StatusEventSendFailed   StatusEventType = C.STATUS_EVENT_SEND_FAILED
	StatusEventDropped      StatusEventType = C.STATUS_EVENT_DROPPED
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
	return StatusEventType(C.zktf_status_event_kind(e.ptr))
}

// DisconnectReason returns the reason for a disconnect event as an error, or nil.
func (e *StatusEvent) DisconnectReason() error {
	return status(C.zktf_status_event_disconnect_reason(e.ptr))
}

// SendError returns the error for a send-failed event, or nil.
func (e *StatusEvent) SendError() error {
	return status(C.zktf_status_event_send_error(e.ptr))
}

// ReferenceID returns the message id referenced by an acknowledged/send-failed
// event, or nil.
func (e *StatusEvent) ReferenceID() []byte {
	ref := C.zktf_status_event_reference(e.ptr)
	if ref == nil {
		return nil
	}
	defer C.zktf_reference_destroy(ref)
	return C.GoBytes(unsafe.Pointer(C.zktf_reference_id(ref)), messageIDLen)
}

// Dropped returns the dropped-event details for a STATUS_EVENT_DROPPED, or nil.
func (e *StatusEvent) Dropped() *DroppedEvent {
	return newDroppedEvent(C.zktf_status_event_dropped(e.ptr))
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
	return GroupEventKind(C.zktf_group_event_kind(e.ptr))
}

// Invite extracts the key package wire event (kind == GroupEventInvite).
func (e *GroupEvent) Invite() *KeyPackageEvent {
	return newKeyPackageEvent(C.zktf_group_event_invite(e.ptr))
}

// Welcome extracts the welcome wire event (kind == GroupEventWelcome).
func (e *GroupEvent) Welcome() *WelcomeEvent {
	return newWelcomeEvent(C.zktf_group_event_welcome(e.ptr))
}

// Commit extracts the commit wire event (kind == GroupEventCommit).
func (e *GroupEvent) Commit() *CommitEvent {
	return newCommitEvent(C.zktf_group_event_commit(e.ptr))
}

// Proposal extracts the proposal wire event (kind == GroupEventProposal).
func (e *GroupEvent) Proposal() *ProposalEvent {
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
	return WorkflowEventKind(C.zktf_workflow_event_kind(e.ptr))
}

// WorkflowID returns the workflow id bytes.
func (e *WorkflowEvent) WorkflowID() []byte {
	n := C.zktf_workflow_event_workflow_id_len(e.ptr)
	if n == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(C.zktf_workflow_event_workflow_id_buf(e.ptr)), C.int(n))
}

// TaskID returns the task id bytes (for TaskFailed events), or nil.
func (e *WorkflowEvent) TaskID() []byte {
	n := C.zktf_workflow_event_task_id_len(e.ptr)
	if n == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(C.zktf_workflow_event_task_id_buf(e.ptr)), C.int(n))
}

// Reason returns a human-readable reason for the event, or "" if unset.
func (e *WorkflowEvent) Reason() string {
	n := C.zktf_workflow_event_reason_len(e.ptr)
	if n == 0 {
		return ""
	}
	return string(C.GoBytes(unsafe.Pointer(C.zktf_workflow_event_reason_buf(e.ptr)), C.int(n)))
}

// Attempt returns the attempt number (for TaskFailed events).
func (e *WorkflowEvent) Attempt() uint32 {
	return uint32(C.zktf_workflow_event_attempt(e.ptr))
}

// WillRetry reports whether the failed task will be retried.
func (e *WorkflowEvent) WillRetry() bool {
	return bool(C.zktf_workflow_event_will_retry(e.ptr))
}
