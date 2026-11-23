package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>

extern zktf_status goValueTransaction(zktf_value_txn txn, uintptr_t handle);

static zktf_status c_value_transaction(zktf_value_txn txn, void *user_data) {
	return goValueTransaction(txn, (uintptr_t)user_data);
}

static zktf_status zktf_account_value_transaction_handle(zktf_account *account, int64_t timeout_ms, uintptr_t handle) {
	return zktf_account_value_transaction(account, timeout_ms, c_value_transaction, (void *)handle);
}
*/
import "C"

import (
	"runtime"
	"runtime/cgo"
	"sync/atomic"
	"time"
	"unsafe"
)

// ValueTxn is a value store transaction. It is only usable on the calling
// goroutine and only until the function passed to ValueTransaction returns.
type ValueTxn struct {
	id     C.zktf_value_txn
	closed atomic.Bool
}

// ValueEntry is a stored value with its unix timestamps (seconds). A zero
// Expires means the value never expires.
type ValueEntry struct {
	Key      string
	Value    []byte
	Created  int64
	Modified int64
	Expires  int64
}

type ValueScanKind int

const (
	ValueScanAll ValueScanKind = iota
	ValueScanPrefix
	ValueScanRange
)

// ValueScan describes one page of a scan. Empty Start or End leaves that side
// of a range open, and a zero Limit uses the native default.
type ValueScan struct {
	Kind     ValueScanKind
	Prefix   string
	Start    string
	End      string
	Reverse  bool
	KeysOnly bool
	Limit    int
	After    []byte
}

// ValuePage is one page of scan results. Next is nil on the last page.
type ValuePage struct {
	Entries []ValueEntry
	Next    []byte
}

type valueTransactionCall struct {
	fn         func(*ValueTxn) error
	err        error
	panicked   bool
	panicValue any
}

// ValueTransaction runs fn as one transaction on the account's value store,
// committing when fn returns nil. A timeout <= 0 waits until the lock is free.
func (a *Account) ValueTransaction(timeout time.Duration, fn func(*ValueTxn) error) error {
	defer runtime.KeepAlive(a)

	call := &valueTransactionCall{fn: fn}
	h := cgo.NewHandle(call)
	defer h.Delete()

	code := C.zktf_account_value_transaction_handle(a.ptr, valueTimeoutMS(timeout), C.uintptr_t(h))

	if call.panicked {
		panic(call.panicValue)
	}
	if call.err != nil {
		return call.err
	}

	return status(code)
}

//export goValueTransaction
func goValueTransaction(id C.zktf_value_txn, h C.uintptr_t) (code C.zktf_status) {
	call := cgo.Handle(uintptr(h)).Value().(*valueTransactionCall)
	txn := &ValueTxn{id: id}
	defer txn.closed.Store(true)

	defer func() {
		if r := recover(); r != nil {
			call.panicked = true
			call.panicValue = r
			code = C.ZKTF_STATUS_UNKNOWN
		}
	}()

	if err := call.fn(txn); err != nil {
		call.err = err
		return C.ZKTF_STATUS_UNKNOWN
	}

	return C.ZKTF_STATUS_OK
}

func valueTimeoutMS(d time.Duration) C.int64_t {
	if d <= 0 {
		return -1
	}

	return C.int64_t((d + time.Millisecond - 1) / time.Millisecond)
}

func (t *ValueTxn) open() error {
	if t.closed.Load() {
		return NewStatus(StatusValueTransactionClosed)
	}
	return nil
}

// Lookup returns the value stored under key, including this transaction's
// writes. The boolean is false when no value is stored for that key.
func (t *ValueTxn) Lookup(key string) ([]byte, bool, error) {
	if err := t.open(); err != nil {
		return nil, false, err
	}

	k := cstring(key)
	defer free(unsafe.Pointer(k))

	var buf *C.zktf_bytes_buffer
	code := C.zktf_value_txn_lookup(t.id, k, &buf)
	if code == C.ZKTF_STATUS_VALUE_NOT_FOUND {
		return nil, false, nil
	}
	if err := status(code); err != nil {
		return nil, false, err
	}

	return goBytesFromBuffer(buf), true, nil
}

// Entry returns the value stored under key with its timestamps. The boolean
// is false when no value is stored for that key.
func (t *ValueTxn) Entry(key string) (ValueEntry, bool, error) {
	if err := t.open(); err != nil {
		return ValueEntry{}, false, err
	}

	k := cstring(key)
	defer free(unsafe.Pointer(k))

	var entry *C.zktf_value_entry
	code := C.zktf_value_txn_entry(t.id, k, &entry)
	if code == C.ZKTF_STATUS_VALUE_NOT_FOUND {
		return ValueEntry{}, false, nil
	}
	if err := status(code); err != nil {
		return ValueEntry{}, false, err
	}
	defer C.zktf_value_entry_destroy(entry)

	return valueEntryFrom(entry), true, nil
}

// Store buffers a write of value under key. A zero expiresUnix means the value
// never expires.
func (t *ValueTxn) Store(key string, value []byte, expiresUnix int64) error {
	if err := t.open(); err != nil {
		return err
	}

	k := cstring(key)
	defer free(unsafe.Pointer(k))

	valBuf, valLen := cbytes(value)
	defer free(unsafe.Pointer(valBuf))

	var options *C.zktf_value_store_options
	if expiresUnix != 0 {
		options = C.zktf_value_store_options_init()
		C.zktf_value_store_options_with_expiry(options, C.int64_t(expiresUnix))
		defer C.zktf_value_store_options_destroy(options)
	}

	return status(C.zktf_value_txn_store(t.id, k, valBuf, valLen, options))
}

// Remove buffers removal of the value stored under key.
func (t *ValueTxn) Remove(key string) error {
	if err := t.open(); err != nil {
		return err
	}

	k := cstring(key)
	defer free(unsafe.Pointer(k))

	return status(C.zktf_value_txn_remove(t.id, k))
}

// Expires sets when the value under key expires; zero expiresUnix clears it.
// The boolean is false when no value is stored for that key.
func (t *ValueTxn) Expires(key string, expiresUnix int64) (bool, error) {
	if err := t.open(); err != nil {
		return false, err
	}

	k := cstring(key)
	defer free(unsafe.Pointer(k))

	var exists C.bool
	if err := status(C.zktf_value_txn_expires(t.id, k, C.int64_t(expiresUnix), &exists)); err != nil {
		return false, err
	}

	return bool(exists), nil
}

// Checkpoint marks the transaction's current writes so Restore can return to
// them.
func (t *ValueTxn) Checkpoint() (uint64, error) {
	if err := t.open(); err != nil {
		return 0, err
	}

	var checkpoint C.uint64_t
	if err := status(C.zktf_value_txn_checkpoint(t.id, &checkpoint)); err != nil {
		return 0, err
	}

	return uint64(checkpoint), nil
}

// Restore undoes every write made since checkpoint.
func (t *ValueTxn) Restore(checkpoint uint64) error {
	if err := t.open(); err != nil {
		return err
	}

	return status(C.zktf_value_txn_restore(t.id, C.uint64_t(checkpoint)))
}

// Scan reads one page of scan, including this transaction's writes.
func (t *ValueTxn) Scan(scan ValueScan) (ValuePage, error) {
	if err := t.open(); err != nil {
		return ValuePage{}, err
	}

	s, err := newValueScan(scan)
	if err != nil {
		return ValuePage{}, err
	}
	defer C.zktf_value_scan_destroy(s)

	var page *C.zktf_value_page
	if err := status(C.zktf_value_txn_scan(t.id, s, &page)); err != nil {
		return ValuePage{}, err
	}
	defer C.zktf_value_page_destroy(page)

	n := int(C.zktf_value_page_len(page))
	out := ValuePage{Entries: make([]ValueEntry, n)}
	for i := 0; i < n; i++ {
		out.Entries[i] = valueEntryFrom(C.zktf_value_page_entry_at(page, C.size_t(i)))
	}
	if next := C.zktf_value_page_next(page); next != nil {
		out.Next = goBytesFromBuffer(next)
	}

	return out, nil
}

func newValueScan(scan ValueScan) (*C.zktf_value_scan, error) {
	var s *C.zktf_value_scan
	switch scan.Kind {
	case ValueScanPrefix:
		p := cstring(scan.Prefix)
		defer free(unsafe.Pointer(p))
		s = C.zktf_value_scan_prefix(p)
	case ValueScanRange:
		var start, end *C.char
		if scan.Start != "" {
			start = cstring(scan.Start)
			defer free(unsafe.Pointer(start))
		}
		if scan.End != "" {
			end = cstring(scan.End)
			defer free(unsafe.Pointer(end))
		}
		s = C.zktf_value_scan_range(start, end)
	default:
		s = C.zktf_value_scan_all()
	}
	if s == nil {
		return nil, NewStatus(C.ZKTF_STATUS_PARAMETER_INVALID_UTF8)
	}

	C.zktf_value_scan_reverse(s, C.bool(scan.Reverse))
	C.zktf_value_scan_keys_only(s, C.bool(scan.KeysOnly))
	if scan.Limit != 0 {
		C.zktf_value_scan_limit(s, C.size_t(scan.Limit))
	}
	if scan.After != nil {
		cursorBuf, cursorLen := cbytes(scan.After)
		defer free(unsafe.Pointer(cursorBuf))
		if err := status(C.zktf_value_scan_after(s, cursorBuf, cursorLen)); err != nil {
			C.zktf_value_scan_destroy(s)
			return nil, err
		}
	}

	return s, nil
}

func valueEntryFrom(entry *C.zktf_value_entry) ValueEntry {
	return ValueEntry{
		Key: C.GoString(C.zktf_value_entry_key(entry)),
		Value: C.GoBytes(
			unsafe.Pointer(C.zktf_value_entry_value_buf(entry)),
			C.int(C.zktf_value_entry_value_len(entry)),
		),
		Created:  int64(C.zktf_value_entry_created(entry)),
		Modified: int64(C.zktf_value_entry_modified(entry)),
		Expires:  int64(C.zktf_value_entry_expires(entry)),
	}
}
