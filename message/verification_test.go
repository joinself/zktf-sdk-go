package message_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/message"
)

func decodeVerification(t *testing.T, req *message.VerificationRequest) *message.VerificationRequest {
	t.Helper()

	content, err := message.NewExchangeRequest().
		Purpose("verify").
		Action(req.AsAction()).
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

	return v
}

func TestVerificationRequestAddressesRoundTrip(t *testing.T) {
	issuer := address(t, 0x21)
	signer := address(t, 0x22)
	subject := address(t, 0x23)
	validUntil := time.Unix(1_900_000_000, 0)

	req, err := message.NewVerificationRequest().
		CredentialType("EmailCredential").
		Issuer(issuer).
		SigningKey(signer).
		Subject(subject).
		ValidUntil(validUntil).
		Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	v := decodeVerification(t, req)

	if got := v.Issuer(); got == nil || !got.Matches(issuer) {
		t.Fatalf("Issuer = %v, want %v", got, issuer)
	}
	if got := v.SigningKey(); got == nil || !got.Matches(signer) {
		t.Fatalf("SigningKey = %v, want %v", got, signer)
	}
	if got := v.Subject(); got == nil || !got.Matches(subject) {
		t.Fatalf("Subject = %v, want %v", got, subject)
	}
	if got := v.ValidUntil(); !got.Equal(validUntil) {
		t.Fatalf("ValidUntil = %v, want %v", got, validUntil)
	}
}

func TestVerificationRequestOptionalAddressesAbsent(t *testing.T) {
	req, err := message.NewVerificationRequest().
		CredentialType("EmailCredential").
		Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	v := decodeVerification(t, req)

	if v.Issuer() != nil || v.SigningKey() != nil || v.Subject() != nil {
		t.Fatal("expected no issuer, signing key or subject")
	}
	if !v.ValidUntil().IsZero() {
		t.Fatalf("ValidUntil = %v, want zero", v.ValidUntil())
	}
}

func TestVerificationParameterStructuredValuesRoundTrip(t *testing.T) {
	object := map[string]any{
		"name":  "alice",
		"age":   int64(30),
		"tags":  []any{"a", "b"},
		"empty": nil,
		"inner": map[string]any{"ok": true},
	}
	array := []any{uint64(7), uint64(8)}

	req, err := message.NewVerificationRequest().
		CredentialType("PassportCredential").
		Parameter("object", object).
		Parameter("array", array).
		Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	v := decodeVerification(t, req)

	got := map[string]any{}
	for _, p := range v.Parameters() {
		got[p.Key()] = p.Value()
	}

	if !reflect.DeepEqual(got["object"], object) {
		t.Fatalf("object = %#v, want %#v", got["object"], object)
	}
	if !reflect.DeepEqual(got["array"], array) {
		t.Fatalf("array = %#v, want %#v", got["array"], array)
	}
}

func TestVerificationParameterRejectsMixedArray(t *testing.T) {
	_, err := message.NewVerificationRequest().
		CredentialType("PassportCredential").
		Parameter("mixed", []any{"a", int64(1)}).
		Finish()
	if err == nil {
		t.Fatal("Finish: want error for a mixed array, got nil")
	}
}

func TestVerificationResponsePresentationsRoundTrip(t *testing.T) {
	vp := presentation(t)

	res, err := message.NewVerificationResponse().
		Presentation(vp).
		Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	outcome, err := message.NewOutcome().
		ActionID(bytes.Repeat([]byte{0x01}, 20)).
		Status(message.StatusOK).
		ResultVerification(res).
		Finish()
	if err != nil {
		t.Fatalf("outcome Finish: %v", err)
	}

	content, err := message.NewExchangeResponse().
		ResponseTo(bytes.Repeat([]byte{0x02}, 20)).
		Status(message.StatusOK).
		Outcome(outcome).
		Finish()
	if err != nil {
		t.Fatalf("exchange response Finish: %v", err)
	}

	decoded, err := message.ExchangeResponseDecode(content)
	if err != nil {
		t.Fatalf("ExchangeResponseDecode: %v", err)
	}
	outcomes, err := decoded.Outcomes()
	if err != nil {
		t.Fatalf("Outcomes: %v", err)
	}
	r, err := outcomes[0].AsVerification()
	if err != nil {
		t.Fatalf("AsVerification: %v", err)
	}

	ps := r.Presentations()
	if len(ps) != 1 {
		t.Fatalf("len(Presentations) = %d, want 1", len(ps))
	}
	if got := ps[0].Types(); !reflect.DeepEqual(got, vp.Types()) {
		t.Fatalf("Types = %v, want %v", got, vp.Types())
	}
}

func presentation(t *testing.T) *credential.VerifiablePresentation {
	t.Helper()

	holder := identity.AddressKey(address(t, 0x31))
	issuer := address(t, 0x32)

	vc, err := credential.NewBuilder().
		Type(credential.CredentialTypesEmail()...).
		Issuer(identity.AddressKey(issuer)).
		SignWith(issuer, time.Unix(1_800_000_000, 0)).
		Subject(holder).
		Claim("emailAddress", "alice@example.com").
		ValidFrom(time.Unix(1_800_000_000, 0)).
		Finish()
	if err != nil {
		t.Fatalf("credential Finish: %v", err)
	}

	vp, err := credential.NewPresentation().
		Type(credential.PresentationTypesDefault()...).
		Holder(holder).
		Credential(vc).
		Finish()
	if err != nil {
		t.Fatalf("presentation Finish: %v", err)
	}

	return vp
}
