package credential_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/account"
	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

var cborldTag = []byte{0xd9, 0xcb, 0x1d}

func signedFixtures(t *testing.T, issuerAddress func(*signing.PublicKey) *identity.Address) (*credential.VerifiableCredential, *credential.VerifiablePresentation) {
	t.Helper()

	acc, err := account.New(account.Config{
		StoragePath: ":memory:",
		Target: &account.Target{
			Network:         account.NetworkPreview,
			RPCEndpoint:     "http://127.0.0.1:1/",
			ObjectEndpoint:  "http://127.0.0.1:1/",
			MessageEndpoint: "ws://127.0.0.1:1/",
		},
	}, account.Callbacks{})
	if err != nil {
		t.Fatalf("account.New: %v", err)
	}
	t.Cleanup(acc.Close)

	issuer, err := acc.KeychainSigningCreate()
	if err != nil {
		t.Fatalf("KeychainSigningCreate: %v", err)
	}
	subject := issuerAddress(issuer)
	holder := identity.AddressKey(issuer)
	at := time.Now().Truncate(time.Second)

	unsigned, err := credential.NewBuilder().
		Type(credential.CredentialTypesEmail()...).
		Issuer(subject).
		Subject(subject).
		ClaimJSON([]byte(`{"email":{"emailAddress":"alice@example.com"}}`)).
		ValidFrom(at).
		SignWith(issuer, at).
		Finish()
	if err != nil {
		t.Fatalf("credential Finish: %v", err)
	}

	vc, err := acc.CredentialIssue(unsigned)
	if err != nil {
		t.Fatalf("CredentialIssue: %v", err)
	}

	vp, err := credential.NewPresentation().
		Type(credential.PresentationTypesDefault()...).
		Holder(holder).
		Credential(vc).
		Finish()
	if err != nil {
		t.Fatalf("presentation Finish: %v", err)
	}
	if err := acc.PresentationSign(vp); err != nil {
		t.Fatalf("PresentationSign: %v", err)
	}

	return vc, vp
}

func checkEncodings(t *testing.T, encode func() ([]byte, error), encodeAs func(credential.Encoding) ([]byte, error), decode func([]byte) ([]byte, error), decodeJSON func([]byte) ([]byte, error)) {
	t.Helper()

	defaultOut, err := encode()
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if !bytes.HasPrefix(defaultOut, cborldTag) {
		t.Fatalf("Encode starts with %x, want %x", defaultOut[:min(3, len(defaultOut))], cborldTag)
	}

	encoded := map[credential.Encoding][]byte{}
	for _, enc := range []credential.Encoding{credential.EncodingJSON, credential.EncodingCBORLD, credential.EncodingCBORLDCompact} {
		out, err := encodeAs(enc)
		if err != nil {
			t.Fatalf("EncodeAs(%d): %v", enc, err)
		}
		encoded[enc] = out

		roundTrip := decode
		if enc == credential.EncodingJSON {
			roundTrip = decodeJSON
		}
		back, err := roundTrip(out)
		if err != nil {
			t.Fatalf("round trip of encoding %d: %v", enc, err)
		}
		if !bytes.Equal(back, defaultOut) {
			t.Fatalf("encoding %d did not round trip:\n got %x\nwant %x", enc, back, defaultOut)
		}
	}

	if !bytes.Equal(encoded[credential.EncodingCBORLDCompact], defaultOut) {
		t.Fatal("EncodeAs(EncodingCBORLDCompact) differs from Encode()")
	}
	if !bytes.HasPrefix(encoded[credential.EncodingCBORLD], cborldTag) {
		t.Fatalf("CBOR-LD starts with %x, want %x", encoded[credential.EncodingCBORLD][:3], cborldTag)
	}
	if len(encoded[credential.EncodingCBORLD]) >= len(encoded[credential.EncodingJSON]) {
		t.Fatalf("CBOR-LD is %d bytes, JSON is %d", len(encoded[credential.EncodingCBORLD]), len(encoded[credential.EncodingJSON]))
	}
	if _, err := decode(encoded[credential.EncodingJSON]); err == nil {
		t.Fatal("Decode accepted JSON")
	}
	if _, err := decodeJSON(defaultOut); err == nil {
		t.Fatal("DecodeJSON accepted CBOR-LD")
	}
}

func TestCredentialEncodeAsRoundTrips(t *testing.T) {
	vc, _ := signedFixtures(t, identity.AddressKey)

	if err := vc.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}

	check := func(decoded *credential.VerifiableCredential, err error) ([]byte, error) {
		if err != nil {
			return nil, err
		}
		if err := decoded.Validate(); err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(decoded.Types(), vc.Types()) {
			t.Fatalf("Types = %v, want %v", decoded.Types(), vc.Types())
		}
		return decoded.Encode()
	}
	checkEncodings(t, vc.Encode, vc.EncodeAs,
		func(data []byte) ([]byte, error) { return check(credential.Decode(data)) },
		func(data []byte) ([]byte, error) { return check(credential.DecodeJSON(data)) })
}

func TestPresentationEncodeAsRoundTrips(t *testing.T) {
	_, vp := signedFixtures(t, zktfAddress)

	if err := vp.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}

	check := func(decoded *credential.VerifiablePresentation, err error) ([]byte, error) {
		if err != nil {
			return nil, err
		}
		if err := decoded.Validate(); err != nil {
			return nil, err
		}
		if len(decoded.Credentials()) != 1 {
			t.Fatalf("len(Credentials) = %d, want 1", len(decoded.Credentials()))
		}
		return decoded.Encode()
	}
	checkEncodings(t, vp.Encode, vp.EncodeAs,
		func(data []byte) ([]byte, error) { return check(credential.DecodePresentation(data)) },
		func(data []byte) ([]byte, error) { return check(credential.DecodePresentationJSON(data)) })
}

func zktfAddress(key *signing.PublicKey) *identity.Address {
	return identity.AddressZktfWithKey(key, key)
}

func TestCompactEncodingIsSmaller(t *testing.T) {
	vc, vp := signedFixtures(t, zktfAddress)

	for name, encodeAs := range map[string]func(credential.Encoding) ([]byte, error){
		"credential":   vc.EncodeAs,
		"presentation": vp.EncodeAs,
	} {
		standard, err := encodeAs(credential.EncodingCBORLD)
		if err != nil {
			t.Fatalf("%s EncodeAs(EncodingCBORLD): %v", name, err)
		}
		compact, err := encodeAs(credential.EncodingCBORLDCompact)
		if err != nil {
			t.Fatalf("%s EncodeAs(EncodingCBORLDCompact): %v", name, err)
		}
		if len(compact) >= len(standard) {
			t.Fatalf("%s compact CBOR-LD is %d bytes, standard is %d", name, len(compact), len(standard))
		}
	}
}
