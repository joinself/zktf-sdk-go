package message_test

import (
	"bytes"
	"testing"

	"github.com/joinself/zktf-sdk-go/message"
)

func TestAnonymousMessageRoundTrips(t *testing.T) {
	from := address(t, 0x01)

	content, err := message.NewDiscoveryRequest().
		FromAddress(from).
		Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	encoded, err := message.NewAnonymousMessage(content).EncodeAsString()
	if err != nil {
		t.Fatalf("EncodeAsString: %v", err)
	}
	if encoded == "" {
		t.Fatal("expected a non-empty encoded string")
	}

	decoded, err := message.AnonymousMessageDecode(encoded)
	if err != nil {
		t.Fatalf("AnonymousMessageDecode: %v", err)
	}

	request, err := message.DiscoveryRequestDecode(decoded.Content())
	if err != nil {
		t.Fatalf("DiscoveryRequestDecode: %v", err)
	}

	if !bytes.Equal(request.FromAddress().Bytes(), from.Bytes()) {
		t.Fatalf("from address did not round trip")
	}
}

func TestAnonymousMessageEncodesAsQR(t *testing.T) {
	content, err := message.NewDiscoveryRequest().
		FromAddress(address(t, 0x02)).
		Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	svg, err := message.NewAnonymousMessage(content).EncodeAsQR(message.QRSVG)
	if err != nil {
		t.Fatalf("EncodeAsQR: %v", err)
	}
	if !bytes.Contains(svg, []byte("<svg")) {
		t.Fatalf("expected an svg document, got %q", svg[:min(len(svg), 64)])
	}

	unicode, err := message.NewAnonymousMessage(content).EncodeAsQR(message.QRUnicode)
	if err != nil {
		t.Fatalf("EncodeAsQR unicode: %v", err)
	}
	if len(unicode) == 0 {
		t.Fatal("expected a non-empty unicode qr")
	}
}

func TestAnonymousMessageBytesRoundTripWithFlags(t *testing.T) {
	from := address(t, 0x03)

	content, err := message.NewDiscoveryRequest().
		FromAddress(from).
		Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	m := message.NewAnonymousMessage(content)
	if m.HasFlags(message.FlagTargetSandbox) {
		t.Fatal("new message unexpectedly has the sandbox flag")
	}

	m.SetFlags(message.FlagTargetSandbox)
	if !m.HasFlags(message.FlagTargetSandbox) {
		t.Fatal("HasFlags = false after SetFlags")
	}
	if got := m.Flags(); got != message.FlagTargetSandbox {
		t.Fatalf("Flags = %d, want %d", got, message.FlagTargetSandbox)
	}

	encoded, err := m.Encode()
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}

	decoded, err := message.AnonymousMessageDecodeBytes(encoded)
	if err != nil {
		t.Fatalf("AnonymousMessageDecodeBytes: %v", err)
	}

	if !bytes.Equal(decoded.ID(), m.ID()) {
		t.Fatal("id did not round trip")
	}
	if !decoded.HasFlags(message.FlagTargetSandbox) {
		t.Fatal("flags did not round trip")
	}

	request, err := message.DiscoveryRequestDecode(decoded.Content())
	if err != nil {
		t.Fatalf("DiscoveryRequestDecode: %v", err)
	}
	if !bytes.Equal(request.FromAddress().Bytes(), from.Bytes()) {
		t.Fatal("from address did not round trip")
	}
}
