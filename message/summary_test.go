package message_test

import (
	"testing"

	"github.com/joinself/zktf-sdk-go/message"
)

// TestContentSummaryOf builds chat content and confirms its summary reports
// the correct type. Exercises the new MessageContentSummary wrapper.
func TestContentSummaryOf(t *testing.T) {
	content, err := message.NewChat().Message("hi").Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}
	summary, err := content.Summary()
	if err != nil {
		t.Fatalf("Summary: %v", err)
	}
	if got := summary.Type(); got != message.ContentChat {
		t.Fatalf("Type = %d, want ContentChat", got)
	}
	if len(summary.ID()) != 20 {
		t.Fatalf("ID len = %d, want 20", len(summary.ID()))
	}
}
