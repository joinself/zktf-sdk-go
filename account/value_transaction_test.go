package account_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/account"
	"github.com/joinself/zktf-sdk-go/status"
)

func newValueAccount(t *testing.T) *account.Account {
	t.Helper()
	acc, err := account.New(account.Config{
		StoragePath: ":memory:",
		Target: &account.Target{
			Network:         account.NetworkPreview,
			RPCEndpoint:     "http://127.0.0.1:1/",
			ObjectEndpoint:  "http://127.0.0.1:1/",
			MessageEndpoint: "ws://127.0.0.1:1/",
		},
	}, account.Callbacks{})
	if err != nil {
		t.Fatalf("account.New: %v", err)
	}
	t.Cleanup(acc.Close)
	return acc
}

func lookupValue(t *testing.T, acc *account.Account, key string) ([]byte, bool) {
	t.Helper()
	value, ok, err := acc.ValueLookup(key)
	if err != nil {
		t.Fatalf("ValueLookup(%s): %v", key, err)
	}
	return value, ok
}

func TestValueTransactionCommits(t *testing.T) {
	acc := newValueAccount(t)
	if err := acc.ValueStore("gone", []byte("x"), time.Time{}); err != nil {
		t.Fatalf("ValueStore: %v", err)
	}

	err := acc.ValueTransaction(0, func(txn *account.ValueTransaction) error {
		if err := txn.Store("kept", []byte("value"), time.Time{}); err != nil {
			return err
		}
		return txn.Remove("gone")
	})
	if err != nil {
		t.Fatalf("ValueTransaction: %v", err)
	}

	if value, ok := lookupValue(t, acc, "kept"); !ok || string(value) != "value" {
		t.Fatalf("kept = %q, %v", value, ok)
	}
	if _, ok := lookupValue(t, acc, "gone"); ok {
		t.Fatal("gone still stored after commit")
	}
}

func TestValueTransactionRollsBackOnError(t *testing.T) {
	acc := newValueAccount(t)
	cause := errors.New("abort")

	err := acc.ValueTransaction(time.Second, func(txn *account.ValueTransaction) error {
		if err := txn.Store("discarded", []byte("x"), time.Time{}); err != nil {
			return err
		}
		return cause
	})
	if err != cause {
		t.Fatalf("ValueTransaction = %v, want %v", err, cause)
	}
	if _, ok := lookupValue(t, acc, "discarded"); ok {
		t.Fatal("write survived rollback")
	}
}

func TestValueTransactionRollsBackOnPanic(t *testing.T) {
	acc := newValueAccount(t)
	type marker struct{ n int }
	cause := &marker{n: 7}

	func() {
		defer func() {
			if r := recover(); r != cause {
				t.Fatalf("recovered %v, want %v", r, cause)
			}
		}()
		_ = acc.ValueTransaction(0, func(txn *account.ValueTransaction) error {
			if err := txn.Store("discarded", []byte("x"), time.Time{}); err != nil {
				return err
			}
			panic(cause)
		})
		t.Fatal("ValueTransaction returned instead of panicking")
	}()

	if _, ok := lookupValue(t, acc, "discarded"); ok {
		t.Fatal("write survived panic")
	}
	if err := acc.ValueTransaction(0, func(*account.ValueTransaction) error { return nil }); err != nil {
		t.Fatalf("ValueTransaction after panic: %v", err)
	}
}

func TestValueTransactionReadsOwnWrites(t *testing.T) {
	acc := newValueAccount(t)
	expires := time.Now().Add(time.Hour).Truncate(time.Second)

	err := acc.ValueTransaction(0, func(txn *account.ValueTransaction) error {
		if err := txn.Store("a/1", []byte("one"), expires); err != nil {
			return err
		}
		if err := txn.Store("a/2", []byte("two"), time.Time{}); err != nil {
			return err
		}

		value, ok, err := txn.Lookup("a/1")
		if err != nil || !ok || string(value) != "one" {
			t.Errorf("Lookup(a/1) = %q, %v, %v", value, ok, err)
		}
		if _, ok, err := txn.Lookup("missing"); err != nil || ok {
			t.Errorf("Lookup(missing) = %v, %v", ok, err)
		}

		entry, ok, err := txn.Entry("a/1")
		if err != nil || !ok {
			t.Errorf("Entry(a/1) = %v, %v", ok, err)
		}
		if entry.Key != "a/1" || string(entry.Value) != "one" || !entry.Expires.Equal(expires) {
			t.Errorf("Entry(a/1) = %+v", entry)
		}
		if _, ok, err := txn.Entry("missing"); err != nil || ok {
			t.Errorf("Entry(missing) = %v, %v", ok, err)
		}

		page, err := txn.Scan(account.ValueScanPrefix("a/"))
		if err != nil {
			return err
		}
		if len(page.Entries) != 2 || page.Entries[0].Key != "a/1" || page.Entries[1].Key != "a/2" || page.Next != nil {
			return fmt.Errorf("Scan = %+v", page)
		}
		if !page.Entries[0].Expires.Equal(expires) || !page.Entries[1].Expires.IsZero() {
			t.Errorf("scan expiry = %v, %v", page.Entries[0].Expires, page.Entries[1].Expires)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("ValueTransaction: %v", err)
	}
}

func TestValueTransactionScanPages(t *testing.T) {
	acc := newValueAccount(t)
	want := []string{"p/1", "p/2", "p/3", "p/4", "p/5"}
	for _, key := range append([]string{"o/1", "q/1"}, want...) {
		if err := acc.ValueStore(key, []byte(key), time.Time{}); err != nil {
			t.Fatalf("ValueStore(%s): %v", key, err)
		}
	}

	scanKeys := func(scan account.ValueScan) []string {
		var keys []string
		err := acc.ValueTransaction(0, func(txn *account.ValueTransaction) error {
			for {
				page, err := txn.Scan(scan)
				if err != nil {
					return err
				}
				if len(page.Entries) > 2 {
					t.Errorf("page has %d entries, limit 2", len(page.Entries))
				}
				for _, e := range page.Entries {
					keys = append(keys, e.Key)
				}
				if page.Next == nil {
					return nil
				}
				scan = scan.After(*page.Next)
			}
		})
		if err != nil {
			t.Fatalf("ValueTransaction: %v", err)
		}
		return keys
	}

	if got := scanKeys(account.ValueScanPrefix("p/").Limit(2)); !equalKeys(got, want) {
		t.Errorf("prefix scan = %v, want %v", got, want)
	}
	if got := scanKeys(account.ValueScanRange("p/2", "p/5").Limit(2).Reverse()); !equalKeys(got, []string{"p/4", "p/3", "p/2"}) {
		t.Errorf("reverse range scan = %v", got)
	}
	if got := scanKeys(account.ValueScanAll().Limit(2)); len(got) != 7 {
		t.Errorf("all scan = %v", got)
	}

	err := acc.ValueTransaction(0, func(txn *account.ValueTransaction) error {
		page, err := txn.Scan(account.ValueScanPrefix("p/").KeysOnly())
		if err != nil {
			return err
		}
		for _, e := range page.Entries {
			if len(e.Value) != 0 {
				t.Errorf("keys-only entry %s has value %q", e.Key, e.Value)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("ValueTransaction: %v", err)
	}
}

func equalKeys(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range want {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func TestValueTransactionExpires(t *testing.T) {
	acc := newValueAccount(t)
	at := time.Now().Add(24 * time.Hour).Truncate(time.Second)

	err := acc.ValueTransaction(0, func(txn *account.ValueTransaction) error {
		if err := txn.Store("k", []byte("v"), time.Time{}); err != nil {
			return err
		}
		if ok, err := txn.Expires("k", at); err != nil || !ok {
			t.Errorf("Expires(k, at) = %v, %v", ok, err)
		}
		if ok, err := txn.Expires("missing", at); err != nil || ok {
			t.Errorf("Expires(missing) = %v, %v", ok, err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("ValueTransaction: %v", err)
	}

	err = acc.ValueTransaction(0, func(txn *account.ValueTransaction) error {
		entry, _, err := txn.Entry("k")
		if err != nil {
			return err
		}
		if !entry.Expires.Equal(at) {
			t.Errorf("expires = %v, want %v", entry.Expires, at)
		}
		if _, err := txn.Expires("k", time.Time{}); err != nil {
			return err
		}
		entry, _, err = txn.Entry("k")
		if err != nil {
			return err
		}
		if !entry.Expires.IsZero() {
			t.Errorf("expires after clear = %v, want zero", entry.Expires)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("ValueTransaction: %v", err)
	}
}

func TestValueTransactionCheckpointRestore(t *testing.T) {
	acc := newValueAccount(t)

	err := acc.ValueTransaction(0, func(txn *account.ValueTransaction) error {
		if err := txn.Store("first", []byte("1"), time.Time{}); err != nil {
			return err
		}
		checkpoint, err := txn.Checkpoint()
		if err != nil {
			return err
		}
		if err := txn.Store("second", []byte("2"), time.Time{}); err != nil {
			return err
		}
		later, err := txn.Checkpoint()
		if err != nil {
			return err
		}
		if err := txn.Restore(checkpoint); err != nil {
			return err
		}
		if _, ok, _ := txn.Lookup("second"); ok {
			t.Error("second visible after restore")
		}
		if err := txn.Restore(later); !errors.Is(err, status.ErrValueCheckpointInvalid) {
			t.Errorf("Restore(later) = %v, want ErrValueCheckpointInvalid", err)
		}
		return txn.Restore(checkpoint)
	})
	if err != nil {
		t.Fatalf("ValueTransaction: %v", err)
	}

	if _, ok := lookupValue(t, acc, "first"); !ok {
		t.Error("first not committed")
	}
	if _, ok := lookupValue(t, acc, "second"); ok {
		t.Error("restored write committed")
	}
}

func TestValueTransactionClosedAfterReturn(t *testing.T) {
	acc := newValueAccount(t)

	var escaped *account.ValueTransaction
	err := acc.ValueTransaction(0, func(txn *account.ValueTransaction) error {
		escaped = txn
		return nil
	})
	if err != nil {
		t.Fatalf("ValueTransaction: %v", err)
	}

	err = escaped.Store("late", []byte("x"), time.Time{})
	if !errors.Is(err, status.ErrValueTransactionClosed) {
		t.Fatalf("Store after return = %v, want ErrValueTransactionClosed", err)
	}
	if code, ok := status.Of(err); !ok || code != status.ValueTransactionClosed {
		t.Fatalf("status.Of = %v, %v", code, ok)
	}
	if _, _, err := escaped.Lookup("late"); !errors.Is(err, status.ErrValueTransactionClosed) {
		t.Fatalf("Lookup after return = %v", err)
	}
	if _, err := escaped.Scan(account.ValueScanAll()); !errors.Is(err, status.ErrValueTransactionClosed) {
		t.Fatalf("Scan after return = %v", err)
	}
}

func TestValueTransactionFromOtherGoroutine(t *testing.T) {
	acc := newValueAccount(t)

	err := acc.ValueTransaction(0, func(txn *account.ValueTransaction) error {
		done := make(chan error)
		go func() { done <- txn.Store("k", []byte("v"), time.Time{}) }()
		if err := <-done; !errors.Is(err, status.ErrValueTransactionClosed) {
			t.Errorf("Store from other goroutine = %v, want ErrValueTransactionClosed", err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("ValueTransaction: %v", err)
	}
}

func TestValueTransactionNested(t *testing.T) {
	acc := newValueAccount(t)

	var inner, store error
	err := acc.ValueTransaction(0, func(*account.ValueTransaction) error {
		inner = acc.ValueTransaction(0, func(*account.ValueTransaction) error { return nil })
		store = acc.ValueStore("k", []byte("v"), time.Time{})
		return nil
	})
	if err != nil {
		t.Fatalf("ValueTransaction: %v", err)
	}
	if !errors.Is(inner, status.ErrValueTransactionActive) {
		t.Errorf("nested ValueTransaction = %v, want ErrValueTransactionActive", inner)
	}
	if !errors.Is(store, status.ErrValueTransactionActive) {
		t.Errorf("ValueStore inside transaction = %v, want ErrValueTransactionActive", store)
	}
}

func TestValueTransactionBusy(t *testing.T) {
	acc := newValueAccount(t)

	holding := make(chan struct{})
	release := make(chan struct{})
	held := make(chan error, 1)
	go func() {
		held <- acc.ValueTransaction(0, func(*account.ValueTransaction) error {
			close(holding)
			<-release
			return nil
		})
	}()
	<-holding

	err := acc.ValueTransaction(50*time.Millisecond, func(*account.ValueTransaction) error {
		t.Error("ran while the lock was held")
		return nil
	})
	close(release)

	if !errors.Is(err, status.ErrValueBusy) {
		t.Fatalf("ValueTransaction while held = %v, want ErrValueBusy", err)
	}
	if err := <-held; err != nil {
		t.Fatalf("holding ValueTransaction: %v", err)
	}
}

func TestValueTransactionKeyInvalid(t *testing.T) {
	acc := newValueAccount(t)

	err := acc.ValueTransaction(0, func(txn *account.ValueTransaction) error {
		return txn.Store("", []byte("x"), time.Time{})
	})
	if !errors.Is(err, status.ErrValueKeyInvalid) {
		t.Fatalf("Store(\"\") = %v, want ErrValueKeyInvalid", err)
	}
	if errors.Is(err, status.ErrValueBusy) {
		t.Fatal("ErrValueKeyInvalid matched ErrValueBusy")
	}
}
