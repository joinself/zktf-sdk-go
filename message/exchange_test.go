package message_test

import (
	"bytes"
	"crypto/rand"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/message"
)

// TestVerificationRequestRoundTrip exercises the verifier flow: build a
// verification action, wrap it in an exchange request, encode, decode, and
// downcast back — all without a network. This proves the polymorphic
// Action/Outcome wrappers and the exchange request builder/decoder work.
func TestVerificationRequestRoundTrip(t *testing.T) {
	verReq, err := message.NewVerificationRequest().
		CredentialType("VerifiableCredential", "EmailCredential").
		Finish()
	if err != nil {
		t.Fatalf("verification Finish: %v", err)
	}

	id := bytes.Repeat([]byte{0x42}, 16)
	content, err := message.NewExchangeRequest().
		ID(id).
		Purpose("verify email").
		Expires(time.Unix(1_900_000_000, 0)).
		Action(verReq.AsAction()).
		Finish()
	if err != nil {
		t.Fatalf("exchange Finish: %v", err)
	}
	if content.Type() != message.ContentExchangeRequest {
		t.Fatalf("Type = %d, want ContentExchangeRequest", content.Type())
	}

	decoded, err := message.ExchangeRequestDecode(content)
	if err != nil {
		t.Fatalf("DecodeExchangeRequest: %v", err)
	}
	if got := decoded.Purpose(); got != "verify email" {
		t.Fatalf("Purpose = %q, want %q", got, "verify email")
	}
	if got := decoded.Expires().Unix(); got != 1_900_000_000 {
		t.Fatalf("Expires = %d, want 1900000000", got)
	}

	actions, err := decoded.Actions()
	if err != nil {
		t.Fatalf("Actions: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("len(Actions) = %d, want 1", len(actions))
	}
	if k := actions[0].Kind(); k != message.ActionCredentialVerification {
		t.Fatalf("kind = %d, want CredentialVerification", k)
	}
	v, err := actions[0].AsVerification()
	if err != nil {
		t.Fatalf("AsVerification: %v", err)
	}
	if got := v.CredentialTypes(); len(got) != 2 ||
		got[0] != "VerifiableCredential" || got[1] != "EmailCredential" {
		t.Fatalf("CredentialTypes = %v, want [VerifiableCredential EmailCredential]", got)
	}
}

// TestPresentationRequestRoundTrip exercises building a presentation request
// (challenge + types) and downcasting back via an exchange request.
func TestPresentationRequestRoundTrip(t *testing.T) {
	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		t.Fatalf("rand: %v", err)
	}

	presReq, err := message.NewPresentationRequest().
		PresentationType("VerifiablePresentation").
		Challenge(challenge).
		Finish()
	if err != nil {
		t.Fatalf("presentation Finish: %v", err)
	}

	content, err := message.NewExchangeRequest().
		Purpose("share presentation").
		Action(presReq.AsAction()).
		Finish()
	if err != nil {
		t.Fatalf("exchange Finish: %v", err)
	}

	decoded, err := message.ExchangeRequestDecode(content)
	if err != nil {
		t.Fatalf("DecodeExchangeRequest: %v", err)
	}
	actions, err := decoded.Actions()
	if err != nil {
		t.Fatalf("Actions: %v", err)
	}
	if len(actions) != 1 || actions[0].Kind() != message.ActionCredentialPresentation {
		t.Fatalf("actions = %+v, want one credential-presentation action", actions)
	}
	p, err := actions[0].AsPresentation()
	if err != nil {
		t.Fatalf("AsPresentation: %v", err)
	}
	if got := p.PresentationTypes(); len(got) != 1 || got[0] != "VerifiablePresentation" {
		t.Fatalf("PresentationTypes = %v, want [VerifiablePresentation]", got)
	}
	if got := p.Challenge(); !bytes.Equal(got, challenge) {
		t.Fatalf("Challenge = %x, want %x", got, challenge)
	}
}

// TestExchangeResponseRoundTrip exercises wrapping an outcome containing a
// verification result and decoding it back.
func TestExchangeResponseRoundTrip(t *testing.T) {
	verRes, err := message.NewVerificationResponse().Finish()
	if err != nil {
		t.Fatalf("verification response Finish: %v", err)
	}

	requestID := bytes.Repeat([]byte{0x42}, 16)
	actionID := bytes.Repeat([]byte{0x99}, 16)

	outcome, err := message.NewOutcome().
		ActionID(actionID).
		Status(message.StatusOK).
		ResultVerification(verRes).
		Finish()
	if err != nil {
		t.Fatalf("outcome Finish: %v", err)
	}

	content, err := message.NewExchangeResponse().
		ResponseTo(requestID).
		Status(message.StatusOK).
		Outcome(outcome).
		Finish()
	if err != nil {
		t.Fatalf("exchange response Finish: %v", err)
	}

	decoded, err := message.ExchangeResponseDecode(content)
	if err != nil {
		t.Fatalf("DecodeExchangeResponse: %v", err)
	}
	if got := decoded.ResponseTo(); !bytes.Equal(got, requestID) {
		t.Fatalf("ResponseTo = %x, want %x", got, requestID)
	}
	if got := decoded.Status(); got != message.StatusOK {
		t.Fatalf("Status = %d, want StatusOK", got)
	}

	outcomes, err := decoded.Outcomes()
	if err != nil {
		t.Fatalf("Outcomes: %v", err)
	}
	if len(outcomes) != 1 {
		t.Fatalf("len(Outcomes) = %d, want 1", len(outcomes))
	}
	if got := outcomes[0].ActionID(); !bytes.Equal(got, actionID) {
		t.Fatalf("ActionID = %x, want %x", got, actionID)
	}
	if got := outcomes[0].Kind(); got != message.OutcomeCredentialVerification {
		t.Fatalf("Kind = %d, want CredentialVerification", got)
	}
}
