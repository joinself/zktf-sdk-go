package account_test

import (
	"testing"

	"github.com/joinself/zktf-sdk-go/account"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/message"
	"github.com/joinself/zktf-sdk-go/status"
)

func unknownAddress(t *testing.T) *signing.PublicKey {
	t.Helper()

	address, err := offlineAccount(t).KeychainSigningCreate()
	if err != nil {
		t.Fatalf("KeychainSigningCreate: %v", err)
	}

	return address
}

func chat(t *testing.T, body string) *message.Content {
	t.Helper()

	content, err := message.NewChat().Message(body).Finish()
	if err != nil {
		t.Fatalf("build chat: %v", err)
	}

	return content
}

func TestMessageSendBatchEmpty(t *testing.T) {
	acc := offlineAccount(t)

	statuses, err := acc.MessageSendBatch(nil)
	if err != nil {
		t.Fatalf("MessageSendBatch: %v", err)
	}
	if statuses == nil || len(statuses) != 0 {
		t.Fatalf("MessageSendBatch(nil) = %v, want an empty slice", statuses)
	}
}

func TestMessageSendBatchRejectsUnknownRecipients(t *testing.T) {
	acc := offlineAccount(t)

	statuses, err := acc.MessageSendBatch([]account.Outgoing{
		{To: unknownAddress(t), Content: chat(t, "one")},
		{To: unknownAddress(t), Content: chat(t, "two")},
	})
	if err != nil {
		t.Fatalf("MessageSendBatch: %v", err)
	}
	if len(statuses) != 2 {
		t.Fatalf("MessageSendBatch returned %d statuses, want 2", len(statuses))
	}

	for i, s := range statuses {
		if s.Queued() {
			t.Fatalf("status %d queued, want rejected", i)
		}
		code, ok := status.Of(s.Err)
		if !ok || code != status.KeyPairNotFound {
			t.Fatalf("status %d: got %v, want KeyPairNotFound", i, s.Err)
		}
	}
}
