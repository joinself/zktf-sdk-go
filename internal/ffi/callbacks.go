package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>

extern void goOnStatus(void *user_data, struct zktf_status_event *event);
extern void goOnMessage(void *user_data, struct zktf_message *message);
extern void goOnGroup(void *user_data, struct zktf_group_event *event);
extern void goOnWorkflow(void *user_data, struct zktf_workflow_event *event);
extern void goOnLog(struct zktf_log_entry *log_entry);

// c_on_log has no user_data parameter, so the handler is process-global.
static void c_on_log(struct zktf_log_entry *log_entry) {
	goOnLog(log_entry);
}

static void c_on_status(void *user_data, struct zktf_status_event *event) {
	goOnStatus(user_data, event);
}
static void c_on_message(void *user_data, struct zktf_message *message) {
	goOnMessage(user_data, message);
}
static void c_on_group(void *user_data, struct zktf_group_event *event) {
	goOnGroup(user_data, event);
}
static void c_on_workflow(void *user_data, struct zktf_workflow_event *event) {
	goOnWorkflow(user_data, event);
}

static zktf_account_callbacks *zktf_account_callbacks_new(void) {
	zktf_account_callbacks *cb = malloc(sizeof(zktf_account_callbacks));
	cb->on_status = c_on_status;
	cb->on_message = c_on_message;
	cb->on_group = c_on_group;
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
	enum zktf_log_level log_level
) {
	zktf_account_config *c = malloc(sizeof(zktf_account_config));
	c->target = target;
	c->rpc_endpoint = rpc_endpoint;
	c->object_endpoint = object_endpoint;
	c->messaging_endpoint = messaging_endpoint;
	c->storage_path = storage_path;
	c->encryption_key_buf = encryption_key_buf;
	c->encryption_key_len = encryption_key_len;
	c->log_level = log_level;
	c->log_callback = c_on_log;
	c->integrity_callback = NULL;
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
	OnMessage(*Message)
	OnGroup(*GroupEvent)
	OnWorkflow(*WorkflowEvent)
}

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
) *C.zktf_account_config {
	return C.zktf_account_config_new(
		C.enum_zktf_account_target(network),
		rpc, object, messaging, storage,
		keyBuf, keyLen,
		C.enum_zktf_log_level(logLevel),
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

//export goOnMessage
func goOnMessage(userData unsafe.Pointer, message *C.zktf_message) {
	dispatch(userData).OnMessage(newMessage(message))
}

//export goOnGroup
func goOnGroup(userData unsafe.Pointer, event *C.zktf_group_event) {
	dispatch(userData).OnGroup(newGroupEvent(event))
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
