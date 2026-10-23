package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>

extern void goOnStatus(void *user_data, struct zktf_status_event *event);
extern enum zktf_delivered goOnEvent(void *user_data, struct zktf_event_batch *batch);
extern void goOnWorkflow(void *user_data, struct zktf_workflow_event *event);
extern void goOnLog(struct zktf_log_entry *log_entry);

// c_on_log has no user_data parameter, so the handler is process-global.
static void c_on_log(struct zktf_log_entry *log_entry) {
	goOnLog(log_entry);
}

static void c_on_status(void *user_data, struct zktf_status_event *event) {
	goOnStatus(user_data, event);
}
static enum zktf_delivered c_on_event(void *user_data, struct zktf_event_batch *batch) {
	return goOnEvent(user_data, batch);
}
static void c_on_workflow(void *user_data, struct zktf_workflow_event *event) {
	goOnWorkflow(user_data, event);
}

static zktf_account_callbacks *zktf_account_callbacks_new(void) {
	zktf_account_callbacks *cb = calloc(1, sizeof(zktf_account_callbacks));
	cb->on_status = c_on_status;
	cb->on_event = c_on_event;
	cb->on_workflow = c_on_workflow;
	return cb;
}

static void zktf_account_callbacks_destroy(zktf_account_callbacks *cb) {
	free(cb);
}

static zktf_account_config *zktf_account_config_new(
	enum zktf_account_target target,
	char *rpc_endpoint,
	char *object_endpoint,
	char *messaging_endpoint,
	char *storage_path,
	uint8_t *encryption_key_buf,
	size_t encryption_key_len,
	enum zktf_log_level log_level,
	bool disable_process_hardening,
	bool disable_storage_memory_security,
	enum zktf_event_delivery event_delivery,
	uint32_t event_max_batch,
	uint32_t event_retry_delay_ms
) {
	zktf_account_config *c = calloc(1, sizeof(zktf_account_config));
	c->target = target;
	c->rpc_endpoint = rpc_endpoint;
	c->object_endpoint = object_endpoint;
	c->messaging_endpoint = messaging_endpoint;
	c->storage_path = storage_path;
	c->encryption_key_buf = encryption_key_buf;
	c->encryption_key_len = encryption_key_len;
	c->log_level = log_level;
	c->log_callback = c_on_log;
	c->disable_process_hardening = disable_process_hardening;
	c->disable_storage_memory_security = disable_storage_memory_security;
	c->event_delivery = event_delivery;
	c->event_max_batch = event_max_batch;
	c->event_retry_delay_ms = event_retry_delay_ms;
	return c;
}

static void zktf_account_config_destroy(zktf_account_config *c) {
	free(c);
}

// zktf_account_configure_handle wraps zktf_account_configure so we can pass the
// cgo.Handle as a uintptr_t (its underlying representation) rather than as a
// raw void* — go vet flags uintptr→unsafe.Pointer→void* conversions but
// uintptr_t→void* in C is benign.
static enum zktf_status zktf_account_configure_handle(
	zktf_account *account,
	const zktf_account_config *config,
	const zktf_account_callbacks *callbacks,
	uintptr_t user_data
) {
	return zktf_account_configure(account, config, callbacks, (void *)user_data);
}
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

// AccountCallbacks is the Go-side dispatch interface for the native account
// callbacks. The public account package implements it with an adapter that
// re-wraps the ffi event types into public types.
type AccountCallbacks interface {
	OnStatus(*StatusEvent)
	OnEvent([]BatchEvent) Delivered
	OnWorkflow(*WorkflowEvent)
}

type Delivered struct {
	Retry bool
}

type EventDelivery uint32

const (
	EventDeliveryCallback EventDelivery = C.EVENT_DELIVERY_CALLBACK
	EventDeliveryPoll     EventDelivery = C.EVENT_DELIVERY_POLL
)

func dispatch(userData unsafe.Pointer) AccountCallbacks {
	return cgo.Handle(uintptr(userData)).Value().(AccountCallbacks)
}

// The C helper functions below are `static`, so in cgo each importing file gets
// its own translation unit and cannot see them. These Go wrappers live in the
// same file as the helpers and are called by account.go (same package, so the
// C *types* are shared even though the static C *functions* are not).

func newAccountConfig(
	network Network,
	rpc, object, messaging, storage *C.char,
	keyBuf *C.uint8_t, keyLen C.size_t,
	logLevel LogLevel,
	disableProcessHardening, disableStorageMemorySecurity bool,
	eventDelivery EventDelivery, eventMaxBatch, eventRetryDelayMs uint32,
) *C.zktf_account_config {
	return C.zktf_account_config_new(
		C.enum_zktf_account_target(network),
		rpc, object, messaging, storage,
		keyBuf, keyLen,
		C.enum_zktf_log_level(logLevel),
		C.bool(disableProcessHardening),
		C.bool(disableStorageMemorySecurity),
		C.enum_zktf_event_delivery(eventDelivery),
		C.uint32_t(eventMaxBatch),
		C.uint32_t(eventRetryDelayMs),
	)
}

func destroyAccountConfig(c *C.zktf_account_config) { C.zktf_account_config_destroy(c) }

func newAccountCallbacks() *C.zktf_account_callbacks { return C.zktf_account_callbacks_new() }

func destroyAccountCallbacks(cb *C.zktf_account_callbacks) { C.zktf_account_callbacks_destroy(cb) }

func accountConfigure(
	acc *C.zktf_account,
	config *C.zktf_account_config,
	callbacks *C.zktf_account_callbacks,
	h cgo.Handle,
) error {
	// h (a cgo.Handle, underlying uintptr) is passed as uintptr_t and cast to
	// void* on the C side — see zktf_account_configure_handle.
	return status(C.zktf_account_configure_handle(acc, config, callbacks, C.uintptr_t(h)))
}

//export goOnStatus
func goOnStatus(userData unsafe.Pointer, event *C.zktf_status_event) {
	dispatch(userData).OnStatus(newStatusEvent(event))
}

//export goOnEvent
func goOnEvent(userData unsafe.Pointer, batch *C.zktf_event_batch) C.enum_zktf_delivered {
	if dispatch(userData).OnEvent(takeEventBatch(batch)).Retry {
		return C.DELIVERED_RETRY
	}
	return C.DELIVERED_DONE
}

//export goOnWorkflow
func goOnWorkflow(userData unsafe.Pointer, event *C.zktf_workflow_event) {
	dispatch(userData).OnWorkflow(newWorkflowEvent(event))
}

//export goOnLog
func goOnLog(entry *C.zktf_log_entry) {
	dispatchLog(newLogEntry(entry))
	C.zktf_log_entry_destroy(entry)
}
