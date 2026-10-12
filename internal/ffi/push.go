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

type PlatformPush struct {
	ptr *C.zktf_platform_push
}

func newPlatformPush(ptr *C.zktf_platform_push) *PlatformPush {
	if ptr == nil {
		return nil
	}
	p := &PlatformPush{ptr: ptr}
	runtime.AddCleanup(p, func(ptr *C.zktf_platform_push) {
		C.zktf_platform_push_destroy(ptr)
	}, p.ptr)
	return p
}

func NewPlatformPushFCM(application *SigningPublicKey, credential string) *PlatformPush {
	defer runtime.KeepAlive(application)
	c := cstring(credential)
	defer free(unsafe.Pointer(c))
	return newPlatformPush(C.zktf_platform_push_fcm(application.ptr, c))
}

func NewPlatformPushAPNS(application *SigningPublicKey, credential string) *PlatformPush {
	defer runtime.KeepAlive(application)
	c := cstring(credential)
	defer free(unsafe.Pointer(c))
	return newPlatformPush(C.zktf_platform_push_apns(application.ptr, c))
}

func (a *Account) PushCredentialSet(push *PlatformPush) error {
	defer runtime.KeepAlive(a)
	defer runtime.KeepAlive(push)
	return status(C.zktf_account_push_credential_set(a.ptr, push.ptr))
}

func (a *Account) PushCredentialGet() (*PlatformPush, error) {
	defer runtime.KeepAlive(a)
	var out *C.zktf_platform_push
	if err := status(C.zktf_account_push_credential_get(a.ptr, &out)); err != nil {
		return nil, err
	}
	return newPlatformPush(out), nil
}

func (b *PushTokenBuilder) PushCredential(push *PlatformPush) *PushTokenBuilder {
	defer runtime.KeepAlive(b)
	defer runtime.KeepAlive(push)
	C.zktf_push_token_builder_push_credential(b.ptr, push.ptr)
	return b
}

type NetworkState uint32

const (
	NetworkStateUnavailable NetworkState = C.NETWORK_STATE_UNAVAILABLE
	NetworkStateAvailable   NetworkState = C.NETWORK_STATE_AVAILABLE
)

func (a *Account) NetworkStateSet(state NetworkState) error {
	defer runtime.KeepAlive(a)
	return status(C.zktf_account_network_state_set(a.ptr, C.enum_zktf_network_state(state)))
}
