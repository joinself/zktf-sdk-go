package account

import (
	"github.com/joinself/zktf-sdk-go/event"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
)

// Callbacks holds the event handlers for an account. Any handler may be nil.
type Callbacks struct {
	// OnStatus fires when the account connects or disconnects.
	OnStatus func(*event.Status)

	// OnEvent receives batches of stored events: messages, group events,
	// acknowledgements, send failures and drops. Returning event.Done()
	// acknowledges the batch; event.Retry(d) redelivers it after d. When nil,
	// batches are acknowledged without being delivered. Unused when the account
	// is configured with DeliveryPoll.
	OnEvent func(*event.Batch) event.Delivered

	// OnWorkflow fires for background workflow events: task completed / failed.
	OnWorkflow func(*event.Workflow)
}

// adapter implements ffi.AccountCallbacks, re-wrapping the native event handles
// into the public event types before dispatching to user-supplied funcs.
type adapter struct {
	cb Callbacks
}

func (a adapter) OnStatus(e *ffi.StatusEvent) {
	if a.cb.OnStatus != nil {
		a.cb.OnStatus(ffi.ToStatusEvent(e).(*event.Status))
	}
}

func (a adapter) OnEvent(events []ffi.BatchEvent) ffi.Delivered {
	if a.cb.OnEvent == nil {
		return ffi.Delivered{}
	}
	retry, after := a.cb.OnEvent(event.NewBatch(events)).IsRetry()
	return ffi.Delivered{Retry: retry, RetryAfter: after}
}

func (a adapter) OnWorkflow(w *ffi.WorkflowEvent) {
	if a.cb.OnWorkflow != nil {
		a.cb.OnWorkflow(ffi.ToWorkflowEvent(w).(*event.Workflow))
	}
}
