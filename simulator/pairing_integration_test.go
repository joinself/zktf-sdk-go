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
	"github.com/joinself/zktf-sdk-go/message"
	simulator "github.com/joinself/zktf-sdk-go/simulator"
)

func TestSetupPairingOffer(t *testing.T) {
	network := simulator.NewNetwork(freePort(t), freePort(t), freePort(t), freePort(t))
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
		OnEvent: func(e *event.Event) {
			if e.Kind() == event.EventConnected {
				select {
				case connected <- struct{}{}:
				default:
				}
			}
		},
	})
	if err != nil {
		t.Fatalf("account: %v", err)
	}
	defer acc.Close()

	select {
	case <-connected:
	case <-time.After(10 * time.Second):
		t.Fatal("account did not connect")
	}

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
