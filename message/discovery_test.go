package message_test

import (
	"bytes"
	"testing"

	"github.com/joinself/zktf-sdk-go/message"
	"github.com/joinself/zktf-sdk-go/token"
)

func discoveryResponse(t *testing.T, tokens ...*token.Token) *message.DiscoveryResponse {
	t.Helper()

	builder := message.NewDiscoveryResponse().
		ResponseTo([]byte("request-id")).
		Status(message.StatusAccepted)

	for _, tk := range tokens {
		builder = builder.Token(tk)
	}

	content, err := builder.Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	decoded, err := message.DiscoveryResponseDecode(content)
	if err != nil {
		t.Fatalf("DiscoveryResponseDecode: %v", err)
	}

	return decoded
}

func TestDiscoveryResponseCarriesNoTokensByDefault(t *testing.T) {
	tokens, err := discoveryResponse(t).Tokens()
	if err != nil {
		t.Fatalf("Tokens: %v", err)
	}

	if len(tokens) != 0 {
		t.Fatalf("expected no tokens, got %d", len(tokens))
	}
}

func TestDiscoveryResponseCarriesAnIdentityToken(t *testing.T) {
	issuer := address(t, 1)
	bearer := address(t, 2)
	application := address(t, 3)

	issued := identityToken(t, issuer.Bytes(), bearer.Bytes(), application.Bytes())

	tokens, err := discoveryResponse(t, issued).Tokens()
	if err != nil {
		t.Fatalf("Tokens: %v", err)
	}

	if len(tokens) != 1 {
		t.Fatalf("expected one token, got %d", len(tokens))
	}

	if got := tokens[0].Bearer().Bytes(); !bytes.Equal(got, bearer.Bytes()) {
		t.Fatalf("bearer mismatch: got %x want %x", got, bearer.Bytes())
	}

	if got := tokens[0].Application().Bytes(); !bytes.Equal(got, application.Bytes()) {
		t.Fatalf("application mismatch: got %x want %x", got, application.Bytes())
	}
}

func TestDiscoveryResponseCarriesSeveralTokensInOrder(t *testing.T) {
	first := identityToken(t, address(t, 1).Bytes(), address(t, 2).Bytes(), address(t, 3).Bytes())
	second := identityToken(t, address(t, 4).Bytes(), address(t, 5).Bytes(), address(t, 6).Bytes())

	tokens, err := discoveryResponse(t, first, second).Tokens()
	if err != nil {
		t.Fatalf("Tokens: %v", err)
	}

	if len(tokens) != 2 {
		t.Fatalf("expected two tokens, got %d", len(tokens))
	}

	for i, want := range []*token.Token{first, second} {
		expected, err := want.Encode()
		if err != nil {
			t.Fatalf("Encode: %v", err)
		}

		got, err := tokens[i].Encode()
		if err != nil {
			t.Fatalf("Encode: %v", err)
		}

		if !bytes.Equal(got, expected) {
			t.Fatalf("token %d did not survive the round trip in order", i)
		}
	}
}
