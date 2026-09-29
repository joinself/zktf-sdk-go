package message_test

import (
	"bytes"
	"runtime"
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

func TestContentDecodeRejectsMalformedBytes(t *testing.T) {
	if _, err := message.ContentDecode(message.ContentChat, []byte("not flatbuffers")); err == nil {
		t.Fatal("ContentDecode: want error, got nil")
	}
}

func TestContentDecodeRejectsEmptyBytes(t *testing.T) {
	if _, err := message.ContentDecode(message.ContentChat, nil); err == nil {
		t.Fatal("ContentDecode: want error, got nil")
	}
}

func TestReceiptRoundTripsManyIDsUnderGC(t *testing.T) {
	ids := make([][]byte, 16)
	builder := message.NewReceipt()
	for i := range ids {
		ids[i] = bytes.Repeat([]byte{byte(i + 1)}, 20)
		builder = builder.Delivered(ids[i])
	}

	content, err := builder.Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	stop := make(chan struct{})
	defer close(stop)
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				runtime.GC()
			}
		}
	}()

	for range 200 {
		r, err := message.ReceiptDecode(content)
		if err != nil {
			t.Fatalf("ReceiptDecode: %v", err)
		}

		delivered := r.Delivered()
		if len(delivered) != len(ids) {
			t.Fatalf("Delivered has %d ids, want %d", len(delivered), len(ids))
		}
		for i := range ids {
			if !bytes.Equal(delivered[i], ids[i]) {
				t.Fatalf("Delivered[%d] = %x, want %x", i, delivered[i], ids[i])
			}
		}
	}
}
