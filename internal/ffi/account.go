package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import (
	"runtime"
	"runtime/cgo"
	"time"
	"unsafe"
)

// Network mirrors zktf_account_target — the network whose trust roots the
// account is anchored to.
type Network uint32

const (
	NetworkProduction  Network = C.TARGET_PRODUCTION
	NetworkSandbox     Network = C.TARGET_SANDBOX
	NetworkStaging     Network = C.TARGET_STAGING
	NetworkPreview     Network = C.TARGET_PREVIEW
	NetworkDevelopment Network = C.TARGET_DEVELOPMENT
)

// LogLevel mirrors zktf_log_level.
type LogLevel uint32

const (
	LogError LogLevel = C.LOG_ERROR
	LogWarn  LogLevel = C.LOG_WARN
	LogInfo  LogLevel = C.LOG_INFO
	LogDebug LogLevel = C.LOG_DEBUG
	LogTrace LogLevel = C.LOG_TRACE
)

// AccountConfig holds the configuration for an account.
type AccountConfig struct {
	Network         Network
	RPCEndpoint     string
	ObjectEndpoint  string
	MessageEndpoint string
	StoragePath     string
	EncryptionKey   []byte
	LogLevel        LogLevel

	DisableProcessHardening bool
	StorageMemorySecurity   StorageMemorySecurity

	EventDelivery     EventDelivery
	EventMaxBatch     uint32
	EventRetryDelayMs uint32

	PairingRoles IdentityKeyRole
	PairingLabel string
}

// Account wraps a zktf_account handle.
type Account struct {
	ptr     *C.zktf_account
	handle  cgo.Handle
	cleanup runtime.Cleanup
}

// NewAccount allocates an unconfigured account.
func NewAccount() *Account {
	a := &Account{ptr: C.zktf_account_init()}
	defer runtime.KeepAlive(a)

	a.cleanup = runtime.AddCleanup(a, func(ptr *C.zktf_account) {
		C.zktf_account_destroy(ptr)
	}, a.ptr)

	return a
}

// Close destroys the native account and releases the cgo.Handle pinning its
// callbacks, stopping the GC cleanup to avoid a double free. No-op if already
// closed.
func (a *Account) Close() {
	defer runtime.KeepAlive(a)
	if a.ptr == nil {
		return
	}
	a.cleanup.Stop()
	C.zktf_account_destroy(a.ptr)
	a.ptr = nil
	if a.handle != 0 {
		a.handle.Delete()
		a.handle = 0
	}
}

// Configure configures the account and registers its callbacks. The callbacks
// are kept alive via a cgo.Handle for the lifetime of the account.
func (a *Account) Configure(cfg AccountConfig, cb AccountCallbacks) error {
	rpc := cstring(cfg.RPCEndpoint)
	object := cstring(cfg.ObjectEndpoint)
	messaging := cstring(cfg.MessageEndpoint)
	storage := cstring(cfg.StoragePath)
	keyBuf, keyLen := cbytes(cfg.EncryptionKey)

	var label *C.char
	if cfg.PairingLabel != "" {
		label = cstring(cfg.PairingLabel)
	}

	defer func() {
		free(unsafe.Pointer(label))
		free(unsafe.Pointer(rpc))
		free(unsafe.Pointer(object))
		free(unsafe.Pointer(messaging))
		free(unsafe.Pointer(storage))
		free(unsafe.Pointer(keyBuf))
	}()

	config := newAccountConfig(
		cfg.Network,
		rpc, object, messaging, storage,
		keyBuf, keyLen,
		cfg.LogLevel,
		cfg.DisableProcessHardening,
		cfg.StorageMemorySecurity,
		cfg.EventDelivery,
		cfg.EventMaxBatch,
		cfg.EventRetryDelayMs,
		cfg.PairingRoles, label,
	)
	defer destroyAccountConfig(config)

	callbacks := newAccountCallbacks()
	defer destroyAccountCallbacks(callbacks)

	a.handle = cgo.NewHandle(cb)

	if err := accountConfigure(a.ptr, config, callbacks, a.handle); err != nil {
		a.handle.Delete()
		return err
	}

	return nil
}

func (a *Account) EventPoll(max int, timeout time.Duration) ([]BatchEvent, error) {
	defer runtime.KeepAlive(a)
	var batch *C.zktf_event_batch
	err := status(C.zktf_account_event_poll(a.ptr, C.size_t(max), C.uint64_t(timeout.Milliseconds()), &batch))
	if err != nil {
		return nil, err
	}
	return takeEventBatch(batch), nil
}

func (a *Account) EventCommit() error {
	defer runtime.KeepAlive(a)
	return status(C.zktf_account_event_commit(a.ptr))
}
