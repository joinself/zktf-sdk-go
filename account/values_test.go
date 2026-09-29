package account_test

import (
	"runtime"
	"sort"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/account"
)

func TestValueKeysRoundTripUnderGC(t *testing.T) {
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
	defer acc.Close()

	want := []string{"k/alpha", "k/bravo", "k/charlie", "k/delta"}
	for _, key := range want {
		if err := acc.ValueStore(key, []byte(key), time.Time{}); err != nil {
			t.Fatalf("ValueStore(%s): %v", key, err)
		}
	}

	stop := make(chan struct{})
	defer close(stop)
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				runtime.GC()
			}
		}
	}()

	for range 200 {
		keys, err := acc.ValueKeys("k/")
		if err != nil {
			t.Fatalf("ValueKeys: %v", err)
		}
		sort.Strings(keys)
		if len(keys) != len(want) {
			t.Fatalf("ValueKeys = %v, want %v", keys, want)
		}
		for i := range want {
			if keys[i] != want[i] {
				t.Fatalf("ValueKeys = %v, want %v", keys, want)
			}
		}
	}
}
