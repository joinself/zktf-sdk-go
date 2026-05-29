//go:build integration

// Package account integration tests exercise the full networked stack: two
// accounts connect to a live zktf network, negotiate an encrypted session, and
// exchange a chat message end to end.
//
// Run against the preview environment (the default):
//
//	CGO_CFLAGS=-I/path/to/zktf-sdk/crates/zktf-ffi \
//	CGO_LDFLAGS=-L/path/to/zktf-sdk/target/debug \
//	LD_LIBRARY_PATH=/path/to/zktf-sdk/target/debug \
//	go test -tags integration -run TestIntegration -v ./account/
//
// Override the endpoints/network via env vars: ZKTF_RPC / ZKTF_OBJECT /
// ZKTF_MESSAGE / ZKTF_NETWORK (production / sandbox / staging / preview /
// development). With no env vars set the test runs against preview.
package account_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/account"
	"github.com/joinself/zktf-sdk-go/event"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/message"
	"github.com/joinself/zktf-sdk-go/object"
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

// TestIntegrationVerificationEvidence uploads an object so it is assigned an
// id, attaches it as evidence (alongside a parameter) on a verification
// request, then round-trips the request through encode/decode and reads the
// evidence and parameter back. Evidence requires an uploaded object, so this
// can only run against a live object store.
func TestIntegrationVerificationEvidence(t *testing.T) {
	alice, _ := testAccount(t, "alice")

	if _, err := alice.InboxOpen(account.WithTimeout(5 * time.Second)); err != nil {
		t.Fatalf("alice.InboxOpen: %v", err)
	}

	const (
		evidenceType = "passport"
		mime         = "image/png"
	)
	data := []byte("fake-passport-scan")

	obj, err := object.New(mime, data)
	if err != nil {
		t.Fatalf("object.New: %v", err)
	}

	if err := alice.ObjectUpload(obj, account.WithTimeout(10*time.Second)); err != nil {
		// Uploading to the object store requires an authorized account; an
		// ephemeral test account may be rejected depending on the environment.
		// Evidence needs an uploaded (id-bearing) object, so skip if we cannot
		// obtain one here.
		t.Skipf("alice.ObjectUpload: %v (object store unavailable/unauthorized in this environment)", err)
	}
	if len(obj.ID()) == 0 {
		t.Fatal("object has no id after upload")
	}

	verReq, err := message.NewVerificationRequest().
		CredentialType("PassportCredential").
		Evidence(evidenceType, obj).
		Parameter("country", "GB").
		Finish()
	if err != nil {
		t.Fatalf("verification Finish: %v", err)
	}

	content, err := message.NewExchangeRequest().
		Purpose("verify passport").
		Action(verReq.AsAction()).
		Finish()
	if err != nil {
		t.Fatalf("exchange Finish: %v", err)
	}

	decoded, err := message.ExchangeRequestDecode(content)
	if err != nil {
		t.Fatalf("ExchangeRequestDecode: %v", err)
	}
	actions, err := decoded.Actions()
	if err != nil {
		t.Fatalf("Actions: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("len(Actions) = %d, want 1", len(actions))
	}
	v, err := actions[0].AsVerification()
	if err != nil {
		t.Fatalf("AsVerification: %v", err)
	}

	ev := v.Evidence()
	if len(ev) != 1 {
		t.Fatalf("len(Evidence) = %d, want 1", len(ev))
	}
	if got := ev[0].Type(); got != evidenceType {
		t.Fatalf("evidence Type = %q, want %q", got, evidenceType)
	}
	// the decoded evidence references the object by id/key/mime; the payload is
	// fetched separately via download, so only the reference is asserted here.
	if got := ev[0].Object().ID(); !bytes.Equal(got, obj.ID()) {
		t.Fatalf("evidence Object ID = %x, want %x", got, obj.ID())
	}
	if got := ev[0].Object().MimeType(); got != mime {
		t.Fatalf("evidence Object MimeType = %q, want %q", got, mime)
	}

	params := v.Parameters()
	if len(params) != 1 {
		t.Fatalf("len(Parameters) = %d, want 1", len(params))
	}
	if params[0].Key() != "country" || params[0].Value() != "GB" {
		t.Fatalf("parameter = %q=%v, want country=GB", params[0].Key(), params[0].Value())
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

	// a normal chat carries no internal metadata payload; this exercises the
	// Metadata accessor over a real wire message and confirms the envelope change
	// didn't break ordinary messaging.
	if _, ok := msg.Metadata(); ok {
		t.Fatal("unexpected metadata on a normal chat message")
	}
}
