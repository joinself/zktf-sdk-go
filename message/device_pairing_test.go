package message_test

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/message"
	"github.com/joinself/zktf-sdk-go/token"
)

const (
	tokenVersion1     = 1
	tokenKindIdentity = 6
	identityTokenLen  = 205
)

// identityToken assembles a v1 identity token. Only the version and kind bytes
// are inspected on decode, so the signature is left zeroed — this exercises the
// carrying of a token, not its verification.
func identityToken(t *testing.T, issuer, bearer, application []byte) *token.Token {
	t.Helper()

	raw := make([]byte, identityTokenLen)
	raw[0] = tokenVersion1
	raw[1] = tokenKindIdentity

	issued := time.Now().Unix()
	binary.LittleEndian.PutUint64(raw[26:34], uint64(issued))
	binary.LittleEndian.PutUint64(raw[34:42], uint64(issued+3600))

	copy(raw[42:75], application)
	copy(raw[75:108], bearer)
	copy(raw[108:141], issuer)

	decoded, err := token.Decode(raw)
	if err != nil {
		t.Fatalf("token.Decode: %v", err)
	}

	return decoded
}

func address(t *testing.T, fill byte) *signing.PublicKey {
	t.Helper()

	raw := make([]byte, 33)
	for i := range raw[1:] {
		raw[i+1] = fill
	}

	key, err := signing.FromBytes(raw)
	if err != nil {
		t.Fatalf("signing.FromBytes: %v", err)
	}

	return key
}

func pairingOperation(t *testing.T, document, device *signing.PublicKey) *identity.Operation {
	t.Helper()

	op, err := identity.NewOperationBuilder().
		ID(document).
		Sequence(0).
		Timestamp(time.Now()).
		GrantSigning(device, identity.KeyRoleMessaging).
		SignWith(document).
		Finish()
	if err != nil {
		t.Fatalf("operation Finish: %v", err)
	}

	return op
}

func TestDevicePairingResponseCarriesTokens(t *testing.T) {
	document := address(t, 0x11)
	device := address(t, 0x22)
	issuer := address(t, 0x33)

	issued := identityToken(t, issuer.Bytes(), device.Bytes(), document.Bytes())

	response, err := message.NewDevicePairingResponse().
		DocumentAddress(document).
		Operation(pairingOperation(t, document, device)).
		Token(issued).
		Finish()
	if err != nil {
		t.Fatalf("device pairing response Finish: %v", err)
	}

	tokens, err := response.Tokens()
	if err != nil {
		t.Fatalf("Tokens: %v", err)
	}

	if len(tokens) != 1 {
		t.Fatalf("len(Tokens) = %d, want 1", len(tokens))
	}

	if got := tokens[0].Kind(); got != token.KindIdentity {
		t.Fatalf("Kind = %d, want KindIdentity", got)
	}

	if got := tokens[0].Bearer().Bytes(); !bytes.Equal(got, device.Bytes()) {
		t.Fatalf("Bearer = %x, want %x", got, device.Bytes())
	}

	if got := tokens[0].Application().Bytes(); !bytes.Equal(got, document.Bytes()) {
		t.Fatalf("Application = %x, want %x", got, document.Bytes())
	}
}

func TestDevicePairingResponseWithoutTokens(t *testing.T) {
	document := address(t, 0x44)
	device := address(t, 0x55)

	response, err := message.NewDevicePairingResponse().
		DocumentAddress(document).
		Operation(pairingOperation(t, document, device)).
		Finish()
	if err != nil {
		t.Fatalf("device pairing response Finish: %v", err)
	}

	tokens, err := response.Tokens()
	if err != nil {
		t.Fatalf("Tokens: %v", err)
	}

	if len(tokens) != 0 {
		t.Fatalf("len(Tokens) = %d, want 0", len(tokens))
	}
}
