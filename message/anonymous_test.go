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
