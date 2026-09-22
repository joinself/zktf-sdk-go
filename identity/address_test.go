package identity_test

import (
	"strings"
	"testing"

	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

func testKey(t *testing.T) *signing.PublicKey {
	t.Helper()

	key, err := signing.FromAddress("00ab1d970ba97d16c9339546dcd079066cc66406cfa7d7e421fe6ec8eaaf0fb5ea")
	if err != nil {
		t.Fatalf("FromAddress: %v", err)
	}

	return key
}

func TestAddressZktfIsADocumentDID(t *testing.T) {
	key := testKey(t)

	zktf := identity.AddressZktf(key).String()
	if !strings.HasPrefix(zktf, "did:zktf:") {
		t.Fatalf("AddressZktf = %q, want a did:zktf: address", zktf)
	}

	// the distinction matters: a did:key: issuer requires that exact key to
	// sign a credential, a document may use any key it grants the role.
	if k := identity.AddressKey(key).String(); !strings.HasPrefix(k, "did:key:") {
		t.Fatalf("AddressKey = %q, want a did:key: address", k)
	}

	if zktf == identity.AddressKey(key).String() {
		t.Fatal("AddressZktf and AddressKey produced the same address")
	}
}

func TestAddressZktfWithKeyNamesTheKey(t *testing.T) {
	key := testKey(t)

	withKey := identity.AddressZktfWithKey(key, key).String()
	if !strings.Contains(withKey, "#") {
		t.Fatalf("AddressZktfWithKey = %q, want a fragment naming the key", withKey)
	}
}

func TestAddressZktfRoundTrips(t *testing.T) {
	zktf := identity.AddressZktf(testKey(t))

	parsed, err := identity.ParseAddress(zktf.String())
	if err != nil {
		t.Fatalf("ParseAddress: %v", err)
	}

	if parsed.String() != zktf.String() {
		t.Fatalf("round trip = %q, want %q", parsed.String(), zktf.String())
	}
}
