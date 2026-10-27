//go:build integration

package simulator_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/account"
	"github.com/joinself/zktf-sdk-go/event"
	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/message"
	simulator "github.com/joinself/zktf-sdk-go/simulator"
)

func pairingServer(t *testing.T, network *simulator.Network, paired chan<- *signing.PublicKey) *account.Account {
	return messagingServer(t, network, paired, nil)
}

func messagingServer(t *testing.T, network *simulator.Network, paired chan<- *signing.PublicKey, messages chan<- *message.Message) *account.Account {
	t.Helper()

	connected := make(chan struct{}, 1)
	acc, err := account.New(account.Config{
		Target: &account.Target{
			Network:         account.NetworkDevelopment,
			RPCEndpoint:     network.RPCEndpoint(),
			ObjectEndpoint:  network.ObjectEndpoint(),
			MessageEndpoint: network.MessagingEndpoint(),
		},
		PairingRoles: identity.KeyRoleMessaging | identity.KeyRoleAssertion,
		PairingLabel: "pairing-test",
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
				if e.Kind() == event.KindPaired && paired != nil {
					paired <- e.PairedApplication()
				}
				if e.Kind() == event.KindMessage && messages != nil {
					select {
					case messages <- e.Message():
					default:
					}
				}
			}
			return event.Done()
		},
	})
	if err != nil {
		t.Fatalf("account: %v", err)
	}
	t.Cleanup(acc.Close)

	select {
	case <-connected:
	case <-time.After(10 * time.Second):
		t.Fatal("account did not connect")
	}

	return acc
}

func TestSetupPairingOffer(t *testing.T) {
	network := simulator.NewEphemeralNetwork()
	acc := pairingServer(t, network, nil)

	offer, err := acc.SetupPairingOffer()
	if err != nil {
		t.Fatalf("setup pairing offer: %v", err)
	}
	code, ok := strings.CutPrefix(offer, "ZKTFPAIR1:")
	if !ok {
		t.Fatalf("want a ZKTFPAIR1: offer, got %q", offer)
	}
	anonymous, err := message.AnonymousMessageDecode(code)
	if err != nil {
		t.Fatalf("decode offer: %v", err)
	}
	request, err := message.DiscoveryRequestDecode(anonymous.Content())
	if err != nil {
		t.Fatalf("decode discovery request: %v", err)
	}
	if request.FromAddress() == nil {
		t.Fatal("want the offer to carry the server's address")
	}
	if request.KeyPackage() != nil {
		t.Fatal("want no key package in a server offer")
	}

	again, err := acc.SetupPairingOffer()
	if err != nil {
		t.Fatalf("setup pairing offer again: %v", err)
	}
	if again != offer {
		t.Fatal("want the pending offer reused")
	}

	unicode, err := acc.SetupPairingQR(message.QRUnicode)
	if err != nil {
		t.Fatalf("setup pairing qr: %v", err)
	}
	if len(unicode) == 0 {
		t.Fatal("want a unicode qr rendering")
	}

	svg, err := acc.SetupPairingQR(message.QRSVG)
	if err != nil {
		t.Fatalf("setup pairing svg: %v", err)
	}
	if !bytes.Contains(svg, []byte("<svg")) {
		t.Fatal("want an svg qr rendering")
	}
}

func TestApplicationPair(t *testing.T) {
	network := simulator.NewEphemeralNetwork()

	app, err := network.ApplicationCreate("Acme Ltd", "Acme Test")
	if err != nil {
		t.Fatalf("application create: %v", err)
	}
	sibling, err := network.ApplicationCreate("Acme Ltd", "Acme Admin")
	if err != nil {
		t.Fatalf("second application create: %v", err)
	}
	if bytes.Equal(app.Bytes(), sibling.Bytes()) {
		t.Fatal("want distinct applications")
	}
	if _, err := network.ApplicationCreate("Other Ltd", "Acme Test"); err == nil {
		t.Fatal("want a duplicate application name rejected")
	}

	paired := make(chan *signing.PublicKey, 1)
	acc := pairingServer(t, network, paired)

	offer, err := acc.SetupPairingOffer()
	if err != nil {
		t.Fatalf("setup pairing offer: %v", err)
	}
	if err := network.ApplicationPair("Missing", offer); err == nil {
		t.Fatal("want an unknown application rejected")
	}
	if err := network.ApplicationPair("Acme Test", offer); err != nil {
		t.Fatalf("application pair: %v", err)
	}

	select {
	case application := <-paired:
		if !bytes.Equal(application.Bytes(), app.Bytes()) {
			t.Fatal("want the server paired to the application")
		}
	case <-time.After(30 * time.Second):
		t.Fatal("server did not report paired")
	}
}

func TestInterceptFromPairedApplication(t *testing.T) {
	network := simulator.NewEphemeralNetwork()

	app, err := network.ApplicationCreate("Acme Ltd", "Acme Test")
	if err != nil {
		t.Fatalf("application create: %v", err)
	}

	paired := make(chan *signing.PublicKey, 1)
	messages := make(chan *message.Message, 16)
	acc := messagingServer(t, network, paired, messages)

	offer, err := acc.SetupPairingOffer()
	if err != nil {
		t.Fatalf("setup pairing offer: %v", err)
	}
	if err := network.ApplicationPair("Acme Test", offer); err != nil {
		t.Fatalf("application pair: %v", err)
	}
	select {
	case <-paired:
	case <-time.After(30 * time.Second):
		t.Fatal("server did not report paired")
	}

	user := simulator.NewDevice(network)
	if err := user.Register(app); err != nil {
		t.Fatalf("register: %v", err)
	}

	pairwise, err := user.PairwiseAddress(app)
	if err != nil {
		t.Fatalf("pairwise address: %v", err)
	}
	address, err := user.Address()
	if err != nil {
		t.Fatalf("address: %v", err)
	}
	if bytes.Equal(pairwise.Bytes(), address.Bytes()) {
		t.Fatal("want a pairwise address distinct from the device address")
	}
	if _, err := user.PairwiseAddress(address); err == nil {
		t.Fatal("want no pairwise address for an unknown counterparty")
	}

	var from *signing.PublicKey
	select {
	case m := <-messages:
		from = m.From()
	case <-time.After(10 * time.Second):
		t.Fatal("server did not hear from the registered device")
	}

	pending, err := user.InterceptWhen(simulator.MatchAll(
		simulator.MatchFromIdentity(app),
		simulator.MatchContentType(simulator.ContentChat),
	))
	if err != nil {
		t.Fatalf("intercept when: %v", err)
	}

	chat, err := message.NewChat().Message("hello").Finish()
	if err != nil {
		t.Fatalf("chat: %v", err)
	}
	if err := acc.MessageSend(from, chat); err != nil {
		t.Fatalf("send: %v", err)
	}

	intercepted, err := pending.Wait(10 * time.Second)
	if err != nil {
		t.Fatalf("wait: %v", err)
	}
	if intercepted == nil {
		t.Fatal("want the chat from the paired application intercepted")
	}
}
