package account

import (
	"github.com/joinself/zktf-sdk-go/event"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/message"
)

// Callbacks holds the four event handlers for an account, mirroring the native
// SDK's callback structure. Any handler may be nil.
type Callbacks struct {
	// OnEvent fires for any account-level status event: connected, disconnected,
	// acknowledged, send-failed, dropped. Dispatch on Event.Kind().
	OnEvent func(*event.Event)

	// OnMessage fires for incoming messages.
	OnMessage func(*message.Message)

	// OnGroup fires for incoming group events: invite, welcome, commit,
	// proposal. Dispatch on event.Group.Kind() and pick the matching extractor.
	OnGroup func(*event.Group)

	// OnWorkflow fires for background workflow events: task completed / failed.
	OnWorkflow func(*event.Workflow)
}

// adapter implements ffi.AccountCallbacks, re-wrapping the native event handles
// into the public event types before dispatching to user-supplied funcs.
type adapter struct {
	cb Callbacks
}

func (a adapter) OnStatus(e *ffi.StatusEvent) {
	if a.cb.OnEvent != nil {
		a.cb.OnEvent(ffi.ToStatusEvent(e).(*event.Event))
	}
}

func (a adapter) OnMessage(m *ffi.Message) {
	if a.cb.OnMessage != nil {
		a.cb.OnMessage(ffi.ToMessage(m).(*message.Message))
	}
}

func (a adapter) OnGroup(g *ffi.GroupEvent) {
	if a.cb.OnGroup != nil {
		a.cb.OnGroup(ffi.ToGroupEvent(g).(*event.Group))
	}
}

func (a adapter) OnWorkflow(w *ffi.WorkflowEvent) {
	if a.cb.OnWorkflow != nil {
		a.cb.OnWorkflow(ffi.ToWorkflowEvent(w).(*event.Workflow))
	}
}
