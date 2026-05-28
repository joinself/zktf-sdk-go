//go:build integration

// Package account integration tests exercise the full networked stack: two
// accounts connect to a live zktf network, negotiate an encrypted session, and
// exchange a chat message end to end.
//
// Run against the preview environment (the default):
//
//	CGO_CFLAGS=-I/path/to/zktf-sdk/crates/self-ffi \
//	CGO_LDFLAGS=-L/path/to/zktf-sdk/target/debug \
//	LD_LIBRARY_PATH=/path/to/zktf-sdk/target/debug \
//	go test -tags integration -run TestIntegration -v ./account/
//
// Override the endpoints/network via env vars: ZKTF_RPC / ZKTF_OBJECT /
// ZKTF_MESSAGE / ZKTF_NETWORK (production / sandbox / staging / preview /
// development). With no env vars set the test runs against preview.
package account_test

import (
	"os"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/account"
	"github.com/joinself/zktf-sdk-go/event"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/message"
)

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}

func networkFromEnv() account.Network {
	switch envOr("ZKTF_NETWORK", "preview") {
	case "production":
		return account.NetworkProduction
	case "sandbox":
		return account.NetworkSandbox
	case "staging":
		return account.NetworkStaging
	case "development", "dev":
		return account.NetworkDevelopment
	default:
		return account.NetworkPreview
	}
}

// testAccount configures an in-memory account and returns it plus a channel
// that receives incoming messages.
func testAccount(t *testing.T, name string) (*account.Account, chan *message.Message) {
	t.Helper()

	messages := make(chan *message.Message, 16)
	connected := make(chan struct{}, 1)

	target := &account.Target{
		Network:         networkFromEnv(),
		RPCEndpoint:     envOr("ZKTF_RPC", "https://rpc.preview.joinself.com/"),
		ObjectEndpoint:  envOr("ZKTF_OBJECT", "https://object.preview.joinself.com/"),
		MessageEndpoint: envOr("ZKTF_MESSAGE", "wss://message.preview.joinself.com/"),
	}

	acc, err := account.New(account.Config{
		Target: target,
	}, account.Callbacks{
		OnEvent: func(e *event.Event) {
			if e.Kind() == event.EventConnected {
				select {
				case connected <- struct{}{}:
				default:
				}
			}
		},
		OnMessage: func(msg *message.Message) {
			messages <- msg
		},
	})
	if err != nil {
		t.Fatalf("%s: account.New: %v", name, err)
	}

	select {
	case <-connected:
	case <-time.After(15 * time.Second):
		t.Skipf("%s: did not connect within 15s (network unreachable?)", name)
	}

	return acc, messages
}

func TestIntegrationChatRoundTrip(t *testing.T) {
	alice, _ := testAccount(t, "alice")
	bobby, bobbyMessages := testAccount(t, "bobby")

	aliceInbox, err := alice.InboxOpen(account.WithTimeout(5 * time.Second))
	if err != nil {
		t.Fatalf("alice.InboxOpen: %v", err)
	}

	bobbyInbox, err := bobby.InboxOpen(account.WithTimeout(5 * time.Second))
	if err != nil {
		t.Fatalf("bobby.InboxOpen: %v", err)
	}

	// alice negotiates an encrypted session with bobby; the SDK auto-accepts.
	if err := alice.GroupNegotiate(aliceInbox, bobbyInbox, time.Time{}); err != nil {
		t.Fatalf("alice.GroupNegotiate: %v", err)
	}

	// give the welcome handshake a moment to settle on both sides.
	time.Sleep(2 * time.Second)

	const body = "hello over the wire"

	content, err := message.NewChat().Message(body).Finish()
	if err != nil {
		t.Fatalf("build chat: %v", err)
	}

	if err := alice.MessageSend(bobbyInbox, content); err != nil {
		t.Fatalf("alice.MessageSend: %v", err)
	}

	select {
	case msg := <-bobbyMessages:
		assertChat(t, msg, body, aliceInbox)
	case <-time.After(10 * time.Second):
		t.Fatal("bobby did not receive alice's message within 10s")
	}
}

func assertChat(t *testing.T, msg *message.Message, want string, from *signing.PublicKey) {
	t.Helper()

	if msg.Content().Type() != message.ContentChat {
		t.Fatalf("content type = %d, want chat", msg.Content().Type())
	}

	chat, err := message.ChatDecode(msg.Content())
	if err != nil {
		t.Fatalf("ChatDecode: %v", err)
	}

	if got := chat.Message(); got != want {
		t.Fatalf("chat = %q, want %q", got, want)
	}

	if !msg.From().Matches(from) {
		t.Fatalf("from = %s, want %s", msg.From(), from)
	}
}
