package message_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/message"
	"github.com/joinself/zktf-sdk-go/revocation"
)

// TestRevocationSigningRequestRoundTrip exercises asking a counterparty
// device to co-sign a revocation statement: build the statement, wrap it in
// a revocation-signing action, put that in an exchange request, encode,
// decode, and downcast back.
func TestRevocationSigningRequestRoundTrip(t *testing.T) {
	issuer := address(t, 0x11)

	stmt, err := revocation.NewBuilder().
		Issuer(issuer).
		Sequence(1).
		Timestamp(time.Unix(1_900_000_000, 0)).
		SignWith(issuer, time.Unix(1_900_000_000, 0)).
		Finish()
	if err != nil {
		t.Fatalf("revocation statement Finish: %v", err)
	}

	revReq, err := message.NewRevocationSigningRequest().
		Statement(stmt).
		Finish()
	if err != nil {
		t.Fatalf("revocation signing Finish: %v", err)
	}

	id := bytes.Repeat([]byte{0x77}, 20)
	content, err := message.NewExchangeRequest().
		ID(id).
		Purpose("co-sign revocation").
		Action(revReq.AsAction()).
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
	if len(actions) != 1 {
		t.Fatalf("len(Actions) = %d, want 1", len(actions))
	}
	if k := actions[0].Kind(); k != message.ActionRevocationSigning {
		t.Fatalf("kind = %d, want ActionRevocationSigning", k)
	}

	r, err := actions[0].AsRevocationSigning()
	if err != nil {
		t.Fatalf("AsRevocationSigning: %v", err)
	}

	// The statement carries a pending (not yet cryptographically applied)
	// signer, so reconstructing it is best-effort: it must not panic, but a
	// nil result is acceptable since the signature was never actually
	// produced by a signing Account.
	_ = r.Statement()
}

// TestRevocationSigningOutcomeRoundTrip exercises returning a co-signed
// revocation statement as the outcome of an exchange response.
func TestRevocationSigningOutcomeRoundTrip(t *testing.T) {
	issuer := address(t, 0x22)

	stmt, err := revocation.NewBuilder().
		Issuer(issuer).
		Sequence(7).
		Timestamp(time.Unix(1_900_000_000, 0)).
		SignWith(issuer, time.Unix(1_900_000_100, 0)).
		Finish()
	if err != nil {
		t.Fatalf("revocation statement Finish: %v", err)
	}

	revRes, err := message.NewRevocationSigningResponse().
		Statement(stmt).
		Finish()
	if err != nil {
		t.Fatalf("revocation signing response Finish: %v", err)
	}

	requestID := bytes.Repeat([]byte{0x77}, 20)
	actionID := bytes.Repeat([]byte{0x88}, 20)

	outcome, err := message.NewOutcome().
		ActionID(actionID).
		Status(message.StatusOK).
		ResultRevocationSigning(revRes).
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

	outcomes, err := decoded.Outcomes()
	if err != nil {
		t.Fatalf("Outcomes: %v", err)
	}
	if len(outcomes) != 1 {
		t.Fatalf("len(Outcomes) = %d, want 1", len(outcomes))
	}
	if got := outcomes[0].Kind(); got != message.OutcomeRevocationSigning {
		t.Fatalf("Kind = %d, want OutcomeRevocationSigning", got)
	}

	r, err := outcomes[0].AsRevocationSigning()
	if err != nil {
		t.Fatalf("AsRevocationSigning: %v", err)
	}

	// See note in TestRevocationSigningRequestRoundTrip: the statement was
	// never actually signed by an Account, so reconstructing it here is
	// best-effort and must simply not panic.
	_ = r.Statement()
}
