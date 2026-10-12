package credential_test

import (
	"bytes"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

func key(t *testing.T, fill byte) *signing.PublicKey {
	t.Helper()

	raw := make([]byte, 33)
	for i := range raw[1:] {
		raw[i+1] = fill
	}

	k, err := signing.FromBytes(raw)
	if err != nil {
		t.Fatalf("FromBytes: %v", err)
	}

	return k
}

func unsignedCredential(t *testing.T, claims []byte) *credential.VerifiableCredential {
	t.Helper()

	issuer := key(t, 0x01)
	at := time.Unix(1_800_000_000, 0)

	vc, err := credential.NewBuilder().
		Type(credential.CredentialTypesEmail()...).
		Issuer(identity.AddressKey(issuer)).
		Subject(identity.AddressKey(key(t, 0x02))).
		ClaimJSON(claims).
		ValidFrom(at).
		SignWith(issuer, at).
		Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	return vc
}

func TestClaimValueDecodesTypedClaims(t *testing.T) {
	vc := unsignedCredential(t, []byte(`{"email":{"emailAddress":"alice@example.com","verified":true,"score":-3,"tags":["a","b"],"nested":{"n":null}}}`))

	got, ok := vc.ClaimValue("/credentialSubject/email")
	if !ok {
		t.Fatal("ClaimValue(/credentialSubject/email) not found")
	}

	m, isMap := got.(map[string]any)
	if !isMap {
		t.Fatalf("ClaimValue = %#v, want an object", got)
	}
	if m["emailAddress"] != "alice@example.com" {
		t.Fatalf("emailAddress = %#v", m["emailAddress"])
	}
	if m["verified"] != true {
		t.Fatalf("verified = %#v", m["verified"])
	}
	if m["score"] != int64(-3) {
		t.Fatalf("score = %#v (%T)", m["score"], m["score"])
	}
	if !reflect.DeepEqual(m["tags"], []any{"a", "b"}) {
		t.Fatalf("tags = %#v", m["tags"])
	}
	if !reflect.DeepEqual(m["nested"], map[string]any{"n": nil}) {
		t.Fatalf("nested = %#v", m["nested"])
	}

	if got, ok := vc.ClaimValue(credential.FieldSubjectEmailAddress); !ok || got != "alice@example.com" {
		t.Fatalf("ClaimValue(%s) = %#v, %v", credential.FieldSubjectEmailAddress, got, ok)
	}

	if _, ok := vc.ClaimValue("/credentialSubject/missing"); ok {
		t.Fatal("ClaimValue(missing) reported found")
	}
}

func TestContentHashIsStable(t *testing.T) {
	vc := unsignedCredential(t, []byte(`{"email":{"emailAddress":"alice@example.com"}}`))

	a, err := vc.ContentHash()
	if err != nil {
		t.Fatalf("ContentHash: %v", err)
	}
	if len(a) != 32 {
		t.Fatalf("len(ContentHash) = %d, want 32", len(a))
	}

	encoded, err := vc.Encode()
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	decoded, err := credential.Decode(encoded)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	b, err := decoded.ContentHash()
	if err != nil {
		t.Fatalf("decoded ContentHash: %v", err)
	}
	if !bytes.Equal(a, b) {
		t.Fatalf("ContentHash changed across encode/decode: %x vs %x", a, b)
	}

	other := unsignedCredential(t, []byte(`{"email":{"emailAddress":"bob@example.com"}}`))
	c, err := other.ContentHash()
	if err != nil {
		t.Fatalf("other ContentHash: %v", err)
	}
	if bytes.Equal(a, c) {
		t.Fatal("different credentials share a content hash")
	}
}

func TestTypeSets(t *testing.T) {
	sets := map[string][]string{
		"CredentialTypesDefault":                       credential.CredentialTypesDefault(),
		"CredentialTypesPassport":                      credential.CredentialTypesPassport(),
		"CredentialTypesFacialComparison":              credential.CredentialTypesFacialComparison(),
		"CredentialTypesLivenessAndFacialComparison":   credential.CredentialTypesLivenessAndFacialComparison(),
		"CredentialTypesBiometricAnchor":               credential.CredentialTypesBiometricAnchor(),
		"CredentialTypesSharingAgreement":              credential.CredentialTypesSharingAgreement(),
		"CredentialTypesEmail":                         credential.CredentialTypesEmail(),
		"CredentialTypesPhone":                         credential.CredentialTypesPhone(),
		"CredentialTypesOrganisation":                  credential.CredentialTypesOrganisation(),
		"CredentialTypesApplication":                   credential.CredentialTypesApplication(),
		"CredentialTypesApplicationInstall":            credential.CredentialTypesApplicationInstall(),
		"PresentationTypesDefault":                     credential.PresentationTypesDefault(),
		"PresentationTypesPassport":                    credential.PresentationTypesPassport(),
		"PresentationTypesFacialComparison":            credential.PresentationTypesFacialComparison(),
		"PresentationTypesLivenessAndFacialComparison": credential.PresentationTypesLivenessAndFacialComparison(),
		"PresentationTypesBiometricAnchor":             credential.PresentationTypesBiometricAnchor(),
		"PresentationTypesSharingAgreement":            credential.PresentationTypesSharingAgreement(),
		"PresentationTypesApplicationInstall":          credential.PresentationTypesApplicationInstall(),
	}

	for name, types := range sets {
		if len(types) == 0 {
			t.Fatalf("%s is empty", name)
		}
	}

	if !slices.Contains(credential.CredentialTypesPassport(), credential.TypePassport) {
		t.Fatalf("CredentialTypesPassport = %v, want it to contain %s", credential.CredentialTypesPassport(), credential.TypePassport)
	}
	if !slices.Contains(credential.PresentationTypesPassport(), credential.PresentationTypePassport) {
		t.Fatalf("PresentationTypesPassport = %v, want it to contain %s", credential.PresentationTypesPassport(), credential.PresentationTypePassport)
	}
}
