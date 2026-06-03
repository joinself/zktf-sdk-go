package message_test

import (
	"bytes"
	"testing"

	"github.com/joinself/zktf-sdk-go/message"
)

func TestCustomRoundTrip(t *testing.T) {
	payload := []byte(`{"k":"v"}`)

	content, err := message.NewCustom().Payload(payload).Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	if content.Type() != message.ContentCustom {
		t.Fatalf("Type = %d, want ContentCustom", content.Type())
	}

	c, err := message.CustomDecode(content)
	if err != nil {
		t.Fatalf("CustomDecode: %v", err)
	}

	if got := c.Payload(); !bytes.Equal(got, payload) {
		t.Fatalf("Payload = %q, want %q", got, payload)
	}
}

func TestReceiptRoundTrip(t *testing.T) {
	id := bytes.Repeat([]byte{0xab}, 20)

	content, err := message.NewReceipt().Delivered(id).Read(id).Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	if content.Type() != message.ContentReceipt {
		t.Fatalf("Type = %d, want ContentReceipt", content.Type())
	}

	r, err := message.ReceiptDecode(content)
	if err != nil {
		t.Fatalf("ReceiptDecode: %v", err)
	}

	if d := r.Delivered(); len(d) != 1 || !bytes.Equal(d[0], id) {
		t.Fatalf("Delivered = %x, want [%x]", d, id)
	}

	if rd := r.Read(); len(rd) != 1 || !bytes.Equal(rd[0], id) {
		t.Fatalf("Read = %x, want [%x]", rd, id)
	}
}
