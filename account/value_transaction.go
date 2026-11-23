package account

import (
	"time"

	"github.com/joinself/zktf-sdk-go/internal/ffi"
)

// ValueTransaction reads and writes the value store within a single
// transaction. It is only valid inside the function passed to
// Account.ValueTransaction and must not be shared with other goroutines; once
// that function returns, every call fails with status.ErrValueTransactionClosed.
type ValueTransaction struct {
	t *ffi.ValueTxn
}

// ValueEntry is a stored value with its timestamps. A zero Expires means the
// value never expires. Value is empty for keys-only scans.
type ValueEntry struct {
	Key      string
	Value    []byte
	Created  time.Time
	Modified time.Time
	Expires  time.Time
}

// ValueCursor marks where a scan stopped. Pass it to ValueScan.After to read
// the next page.
type ValueCursor struct {
	b []byte
}

// ValuePage is one page of scan results. Next is nil on the last page.
type ValuePage struct {
	Entries []ValueEntry
	Next    *ValueCursor
}

// ValueCheckpoint marks a point in a transaction's writes to restore to.
type ValueCheckpoint struct {
	token uint64
}

// ValueScan describes a scan over the value store. Scans read 256 keys per
// page in ascending order unless changed with Limit and Reverse.
type ValueScan struct {
	s ffi.ValueScan
}

// ValueScanAll scans every key.
func ValueScanAll() ValueScan {
	return ValueScan{s: ffi.ValueScan{Kind: ffi.ValueScanAll}}
}

// ValueScanPrefix scans keys starting with prefix.
func ValueScanPrefix(prefix string) ValueScan {
	return ValueScan{s: ffi.ValueScan{Kind: ffi.ValueScanPrefix, Prefix: prefix}}
}

// ValueScanRange scans keys from start (inclusive) to end (exclusive). An
// empty start or end leaves that side open.
func ValueScanRange(start, end string) ValueScan {
	return ValueScan{s: ffi.ValueScan{Kind: ffi.ValueScanRange, Start: start, End: end}}
}

// Reverse reads keys in descending order.
func (s ValueScan) Reverse() ValueScan {
	s.s.Reverse = true
	return s
}

// Limit sets the page size, clamped to 1..4096.
func (s ValueScan) Limit(n int) ValueScan {
	s.s.Limit = min(max(n, 1), 4096)
	return s
}

// KeysOnly returns keys and timestamps without loading values.
func (s ValueScan) KeysOnly() ValueScan {
	s.s.KeysOnly = true
	return s
}

// After continues the scan from cursor.
func (s ValueScan) After(cursor ValueCursor) ValueScan {
	s.s.After = cursor.b
	return s
}

// ValueTransaction runs fn as one transaction on the account's value store.
// The writes fn makes are committed when it returns nil and discarded when it
// returns an error or panics; the error is returned unchanged and the panic is
// re-raised. A timeout <= 0 waits until the lock is free; otherwise the call
// fails with status.ErrValueBusy once it elapses. fn runs on the calling
// goroutine and must not hand txn to other goroutines.
func (a *Account) ValueTransaction(timeout time.Duration, fn func(txn *ValueTransaction) error) error {
	return a.h.ValueTransaction(timeout, func(t *ffi.ValueTxn) error {
		return fn(&ValueTransaction{t: t})
	})
}

// Lookup returns the value stored under key, including this transaction's
// writes. The boolean is false when no value is stored for that key.
func (txn *ValueTransaction) Lookup(key string) ([]byte, bool, error) {
	return txn.t.Lookup(key)
}

// Entry returns the value stored under key with its timestamps. The boolean
// is false when no value is stored for that key.
func (txn *ValueTransaction) Entry(key string) (ValueEntry, bool, error) {
	e, ok, err := txn.t.Entry(key)
	if err != nil || !ok {
		return ValueEntry{}, ok, err
	}

	return valueEntryFrom(e), true, nil
}

// Store writes value under key. A zero expires means the value never expires.
func (txn *ValueTransaction) Store(key string, value []byte, expires time.Time) error {
	return txn.t.Store(key, value, unixOrZero(expires))
}

// Remove deletes the value stored under key.
func (txn *ValueTransaction) Remove(key string) error {
	return txn.t.Remove(key)
}

// Expires sets when the value under key expires without rewriting it. A zero
// at makes it never expire. The boolean is false when no value is stored for
// that key.
func (txn *ValueTransaction) Expires(key string, at time.Time) (bool, error) {
	return txn.t.Expires(key, unixOrZero(at))
}

// Scan reads one page of scan, including this transaction's writes.
func (txn *ValueTransaction) Scan(scan ValueScan) (ValuePage, error) {
	p, err := txn.t.Scan(scan.s)
	if err != nil {
		return ValuePage{}, err
	}

	page := ValuePage{Entries: make([]ValueEntry, len(p.Entries))}
	for i, e := range p.Entries {
		page.Entries[i] = valueEntryFrom(e)
	}
	if p.Next != nil {
		page.Next = &ValueCursor{b: p.Next}
	}

	return page, nil
}

// Checkpoint marks the transaction's current writes so Restore can return to
// them.
func (txn *ValueTransaction) Checkpoint() (ValueCheckpoint, error) {
	token, err := txn.t.Checkpoint()
	if err != nil {
		return ValueCheckpoint{}, err
	}

	return ValueCheckpoint{token: token}, nil
}

// Restore undoes every write made since checkpoint. Checkpoints taken after it
// become invalid.
func (txn *ValueTransaction) Restore(checkpoint ValueCheckpoint) error {
	return txn.t.Restore(checkpoint.token)
}

func valueEntryFrom(e ffi.ValueEntry) ValueEntry {
	entry := ValueEntry{
		Key:      e.Key,
		Value:    e.Value,
		Created:  time.Unix(e.Created, 0),
		Modified: time.Unix(e.Modified, 0),
	}
	if e.Expires != 0 {
		entry.Expires = time.Unix(e.Expires, 0)
	}

	return entry
}

func unixOrZero(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}

	return t.Unix()
}
