//go:build integration

package simulator_test

import (
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/account"
	"github.com/joinself/zktf-sdk-go/event"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/message"
	simulator "github.com/joinself/zktf-sdk-go/simulator"
	"github.com/joinself/zktf-sdk-go/status"
)

func simulatedAccount(t *testing.T, network *simulator.Network, name string) (*account.Account, chan *message.Message) {
	t.Helper()

	messages := make(chan *message.Message, 16)
	connected := make(chan struct{}, 1)

	acc, err := account.New(account.Config{
		StoragePath: ":memory:",
		Target: &account.Target{
			Network:         account.NetworkDevelopment,
			RPCEndpoint:     network.RPCEndpoint(),
			ObjectEndpoint:  network.ObjectEndpoint(),
			MessageEndpoint: network.MessagingEndpoint(),
		},
	}, account.Callbacks{
		OnStatus: func(s *event.Status) {
			if s.Kind() == event.StatusConnected {
				select {
				case connected <- struct{}{}:
				default:
				}
			}
		},
		OnEvent: func(b *event.Batch) event.Delivered {
			for _, e := range b.Events() {
				if e.Kind() == event.KindMessage {
					messages <- e.Message()
				}
			}
			return event.Done()
		},
	})
	if err != nil {
		t.Fatalf("%s: account.New: %v", name, err)
	}
	t.Cleanup(acc.Close)

	select {
	case <-connected:
	case <-time.After(15 * time.Second):
		t.Fatalf("%s: did not connect within 15s", name)
	}

	return acc, messages
}

func TestMessageSendBatch(t *testing.T) {
	network := simulator.NewNetwork(freePort(t), freePort(t), freePort(t), freePort(t))

	alice, _ := simulatedAccount(t, network, "alice")
	bobby, bobbyMessages := simulatedAccount(t, network, "bobby")

	aliceInbox, err := alice.InboxOpen(account.WithTimeout(5 * time.Second))
	if err != nil {
		t.Fatalf("alice.InboxOpen: %v", err)
	}

	bobbyInbox, err := bobby.InboxOpen(account.WithTimeout(5 * time.Second))
	if err != nil {
		t.Fatalf("bobby.InboxOpen: %v", err)
	}

	if err := alice.GroupNegotiate(aliceInbox, bobbyInbox, time.Time{}); err != nil {
		t.Fatalf("alice.GroupNegotiate: %v", err)
	}

	time.Sleep(2 * time.Second)

	stranger, err := bobby.KeychainSigningCreate()
	if err != nil {
		t.Fatalf("bobby.KeychainSigningCreate: %v", err)
	}

	bodies := []string{"one", "two", "three"}
	recipients := []*signing.PublicKey{bobbyInbox, stranger, bobbyInbox}

	outgoing := make([]account.Outgoing, len(bodies))
	for i, body := range bodies {
		content, err := message.NewChat().Message(body).Finish()
		if err != nil {
			t.Fatalf("build chat %q: %v", body, err)
		}
		outgoing[i] = account.Outgoing{To: recipients[i], Content: content}
	}

	statuses, err := alice.MessageSendBatch(outgoing)
	if err != nil {
		t.Fatalf("alice.MessageSendBatch: %v", err)
	}
	if len(statuses) != 3 {
		t.Fatalf("MessageSendBatch returned %d statuses, want 3", len(statuses))
	}
	if !statuses[0].Queued() || !statuses[2].Queued() {
		t.Fatalf("statuses = %v, want the first and last queued", statuses)
	}
	if statuses[1].Queued() {
		t.Fatal("status 1 queued, want rejected")
	}
	if code, ok := status.Of(statuses[1].Err); !ok || code != status.KeyPairNotFound {
		t.Fatalf("status 1: got %v, want KeyPairNotFound", statuses[1].Err)
	}

	for _, want := range []string{"one", "three"} {
		select {
		case msg := <-bobbyMessages:
			chat, err := message.ChatDecode(msg.Content())
			if err != nil {
				t.Fatalf("ChatDecode: %v", err)
			}
			if got := chat.Message(); got != want {
				t.Fatalf("bobby received %q, want %q", got, want)
			}
			if !msg.From().Matches(aliceInbox) {
				t.Fatalf("message from %s, want %s", msg.From(), aliceInbox)
			}
		case <-time.After(10 * time.Second):
			t.Fatalf("bobby did not receive %q within 10s", want)
		}
	}

	select {
	case msg := <-bobbyMessages:
		t.Fatalf("bobby received an unexpected message: %v", msg)
	case <-time.After(500 * time.Millisecond):
	}
}
