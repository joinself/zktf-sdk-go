package message_test

import (
	"testing"

	"github.com/joinself/zktf-sdk-go/message"
)

// TestChatRoundTrip exercises the full public -> internal/ffi -> C -> back path
// without any network: build chat content, then read its type and message text
// back through the decoded content. This proves the no-linkname / no-exposed-C
// wrapping works end to end at runtime.
func TestChatRoundTrip(t *testing.T) {
	const body = "hello zktf"

	content, err := message.NewChat().Message(body).Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	if got := content.Type(); got != message.ContentChat {
		t.Fatalf("Type = %d, want ContentChat (%d)", got, message.ContentChat)
	}

	chat, err := message.ChatDecode(content)
	if err != nil {
		t.Fatalf("DecodeChat: %v", err)
	}

	if got := chat.Message(); got != body {
		t.Fatalf("Message = %q, want %q", got, body)
	}
}
