// Async future plumbing.
//
// The native SDK exposes both a blocking `zktf_future_*_wait` and a callback-
// based `zktf_future_*_on_complete` for every future. Blocking the calling
// thread in C is invisible to the Go scheduler — it consumes an OS thread that
// can't be preempted. We avoid that by always using `_on_complete`: the native
// runtime invokes our callback (success, failure, or timeout — the SDK enforces
// the timeout internally) and we deliver the result over a oneshot channel.
//
// The Go side passes a cgo.Handle wrapping the result channel as the C
// user_data, which is shuttled through C as uintptr_t to avoid Go pointers
// crossing into C. Each future kind has its own //export trampoline that
// recovers the channel and sends the typed result.
package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>

extern void goFutureStatusComplete(uintptr_t handle, zktf_status reason);
extern void goFutureSigningPublicKeyComplete(uintptr_t handle, zktf_status reason, zktf_signing_public_key *key);
extern void goFutureIdentityDocumentComplete(uintptr_t handle, zktf_status reason, zktf_identity_document *doc);
extern void goFutureCredentialGraphComplete(uintptr_t handle, zktf_status reason, zktf_credential_graph *graph);
extern void goFutureGroupComplete(uintptr_t handle, zktf_status reason, zktf_group *group);

static void c_future_status_done(void *user_data, zktf_status reason) {
	goFutureStatusComplete((uintptr_t)user_data, reason);
}
static void c_future_signing_public_key_done(void *user_data, zktf_status reason, zktf_signing_public_key *key) {
	goFutureSigningPublicKeyComplete((uintptr_t)user_data, reason, key);
}
static void c_future_identity_document_done(void *user_data, zktf_status reason, zktf_identity_document *doc) {
	goFutureIdentityDocumentComplete((uintptr_t)user_data, reason, doc);
}
static void c_future_credential_graph_done(void *user_data, zktf_status reason, zktf_credential_graph *graph) {
	goFutureCredentialGraphComplete((uintptr_t)user_data, reason, graph);
}
// The header declares zktf_future_group_on_complete as taking zktf_on_group_cb
// (a 2-arg event callback) but the Rust implementation actually invokes the
// callback with 3 args (user_data, status, group). Cast through this 3-arg
// shim so we receive the real signature.
typedef void (*future_group_cb_t)(void *user_data, zktf_status reason, zktf_group *group);
static void c_future_group_done(void *user_data, zktf_status reason, zktf_group *group) {
	goFutureGroupComplete((uintptr_t)user_data, reason, group);
}

static void zktf_future_status_await(struct zktf_future_status *fut, uint32_t timeout_ms, uintptr_t ud) {
	zktf_future_status_on_complete(fut, timeout_ms, c_future_status_done, (void *)ud);
}
static void zktf_future_signing_public_key_await(struct zktf_future_signing_public_key *fut, uint32_t timeout_ms, uintptr_t ud) {
	zktf_future_signing_public_key_on_complete(fut, timeout_ms, c_future_signing_public_key_done, (void *)ud);
}
static void zktf_future_identity_document_await(struct zktf_future_identity_document *fut, uint32_t timeout_ms, uintptr_t ud) {
	zktf_future_identity_document_on_complete(fut, timeout_ms, c_future_identity_document_done, (void *)ud);
}
static void zktf_future_credential_graph_await(struct zktf_future_credential_graph *fut, uint32_t timeout_ms, uintptr_t ud) {
	zktf_future_credential_graph_on_complete(fut, timeout_ms, c_future_credential_graph_done, (void *)ud);
}
static void zktf_future_group_await(struct zktf_future_group *fut, uint32_t timeout_ms, uintptr_t ud) {
	zktf_future_group_on_complete(fut, timeout_ms, (zktf_on_group_cb)c_future_group_done, (void *)ud);
}
*/
import "C"

import (
	"runtime/cgo"
	"time"
)

// DefaultTimeout is the default timeout applied to async operations when the
// caller doesn't pass one.
const DefaultTimeout = 30 * time.Second

// timeoutMS converts a Go timeout into the milliseconds the native SDK expects.
// Zero or negative durations use DefaultTimeout.
func timeoutMS(d time.Duration) C.uint32_t {
	if d <= 0 {
		d = DefaultTimeout
	}

	return C.uint32_t(d / time.Millisecond)
}

// statusResult is the result delivered to a status future's channel.
type statusResult struct{ err error }

// AwaitStatus drives a zktf_future_status to completion via callback and returns
// once the channel has been signalled.
func AwaitStatus(fut *C.zktf_future_status, timeout time.Duration) error {
	ch := make(chan statusResult, 1)
	h := cgo.NewHandle(ch)

	C.zktf_future_status_await(fut, timeoutMS(timeout), C.uintptr_t(h))

	return (<-ch).err
}

//export goFutureStatusComplete
func goFutureStatusComplete(h C.uintptr_t, reason C.enum_zktf_status) {
	handle := cgo.Handle(uintptr(h))
	ch := handle.Value().(chan statusResult)
	handle.Delete()

	ch <- statusResult{err: status(reason)}
}

// signingPublicKeyResult is the result delivered to a signing-public-key future's channel.
type signingPublicKeyResult struct {
	key *SigningPublicKey
	err error
}

// AwaitSigningPublicKey drives a zktf_future_signing_public_key to completion.
func AwaitSigningPublicKey(fut *C.zktf_future_signing_public_key, timeout time.Duration) (*SigningPublicKey, error) {
	ch := make(chan signingPublicKeyResult, 1)
	h := cgo.NewHandle(ch)

	C.zktf_future_signing_public_key_await(fut, timeoutMS(timeout), C.uintptr_t(h))

	r := <-ch
	return r.key, r.err
}

//export goFutureSigningPublicKeyComplete
func goFutureSigningPublicKeyComplete(h C.uintptr_t, reason C.enum_zktf_status, key *C.zktf_signing_public_key) {
	handle := cgo.Handle(uintptr(h))
	ch := handle.Value().(chan signingPublicKeyResult)
	handle.Delete()

	var r signingPublicKeyResult
	if err := status(reason); err != nil {
		r.err = err
	} else {
		r.key = newSigningPublicKey(key)
	}

	ch <- r
}

// identityDocumentResult is the result delivered to an identity-document future's channel.
type identityDocumentResult struct {
	doc *IdentityDocument
	err error
}

// AwaitIdentityDocument drives a zktf_future_identity_document to completion.
func AwaitIdentityDocument(fut *C.zktf_future_identity_document, timeout time.Duration) (*IdentityDocument, error) {
	ch := make(chan identityDocumentResult, 1)
	h := cgo.NewHandle(ch)

	C.zktf_future_identity_document_await(fut, timeoutMS(timeout), C.uintptr_t(h))

	r := <-ch
	return r.doc, r.err
}

//export goFutureIdentityDocumentComplete
func goFutureIdentityDocumentComplete(h C.uintptr_t, reason C.enum_zktf_status, doc *C.zktf_identity_document) {
	handle := cgo.Handle(uintptr(h))
	ch := handle.Value().(chan identityDocumentResult)
	handle.Delete()

	var r identityDocumentResult
	if err := status(reason); err != nil {
		r.err = err
	} else {
		r.doc = newIdentityDocument(doc)
	}

	ch <- r
}

// credentialGraphResult is the result delivered to a credential-graph future's channel.
type credentialGraphResult struct {
	graph *CredentialGraph
	err   error
}

// AwaitCredentialGraph drives a zktf_future_credential_graph to completion.
func AwaitCredentialGraph(fut *C.zktf_future_credential_graph, timeout time.Duration) (*CredentialGraph, error) {
	ch := make(chan credentialGraphResult, 1)
	h := cgo.NewHandle(ch)

	C.zktf_future_credential_graph_await(fut, timeoutMS(timeout), C.uintptr_t(h))

	r := <-ch
	return r.graph, r.err
}

//export goFutureCredentialGraphComplete
func goFutureCredentialGraphComplete(h C.uintptr_t, reason C.enum_zktf_status, graph *C.zktf_credential_graph) {
	handle := cgo.Handle(uintptr(h))
	ch := handle.Value().(chan credentialGraphResult)
	handle.Delete()

	var r credentialGraphResult
	if err := status(reason); err != nil {
		r.err = err
	} else {
		r.graph = newCredentialGraph(graph)
	}

	ch <- r
}

// groupResult is the result delivered to a group future's channel.
type groupResult struct {
	group *Group
	err   error
}

// AwaitGroup drives a zktf_future_group to completion.
func AwaitGroup(fut *C.zktf_future_group, timeout time.Duration) (*Group, error) {
	ch := make(chan groupResult, 1)
	h := cgo.NewHandle(ch)

	C.zktf_future_group_await(fut, timeoutMS(timeout), C.uintptr_t(h))

	r := <-ch
	return r.group, r.err
}

//export goFutureGroupComplete
func goFutureGroupComplete(h C.uintptr_t, reason C.enum_zktf_status, group *C.zktf_group) {
	handle := cgo.Handle(uintptr(h))
	ch := handle.Value().(chan groupResult)
	handle.Delete()

	var r groupResult
	if err := status(reason); err != nil {
		r.err = err
	} else {
		r.group = newGroup(group)
	}

	ch <- r
}
