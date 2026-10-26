package account_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/account"
	"github.com/joinself/zktf-sdk-go/status"
)

func offlineTarget() *account.Target {
	return &account.Target{
		Network:         account.NetworkPreview,
		RPCEndpoint:     "http://127.0.0.1:1/",
		ObjectEndpoint:  "http://127.0.0.1:1/",
		MessageEndpoint: "ws://127.0.0.1:1/",
	}
}

func TestEventPollRequiresPollDelivery(t *testing.T) {
	acc := offlineAccount(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := acc.EventPoll(ctx, 8)
	code, ok := status.Of(err)
	if !ok || code != status.EventPollUnavailable {
		t.Fatalf("EventPoll on a callback account: got %v, want EventPollUnavailable", err)
	}

	code, ok = status.Of(acc.EventCommit())
	if !ok || code != status.EventPollUnavailable {
		t.Fatalf("EventCommit on a callback account: want EventPollUnavailable, got %v", code)
	}
}

func TestEventPollReturnsWhenTheContextEnds(t *testing.T) {
	acc, err := account.New(account.Config{
		StoragePath:   ":memory:",
		Target:        offlineTarget(),
		EventDelivery: account.DeliveryPoll,
	}, account.Callbacks{})
	if err != nil {
		t.Fatalf("account.New: %v", err)
	}
	t.Cleanup(acc.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	started := time.Now()
	batch, err := acc.EventPoll(ctx, 8)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("EventPoll: got batch %v, err %v; want deadline exceeded", batch, err)
	}
	if elapsed := time.Since(started); elapsed > 2*time.Second {
		t.Fatalf("EventPoll took %s to honour a 300ms deadline", elapsed)
	}

	if err := acc.EventCommit(); err != nil {
		t.Fatalf("EventCommit with nothing outstanding: %v", err)
	}
}
