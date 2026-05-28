package object_test

import (
	"bytes"
	"testing"

	"github.com/joinself/zktf-sdk-go/object"
)

func TestObjectRoundTrip(t *testing.T) {
	data := []byte("attachment bytes")

	o, err := object.New("text/plain", data)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if got := o.MimeType(); got != "text/plain" {
		t.Fatalf("MimeType = %q, want text/plain", got)
	}
	if got := o.Data(); !bytes.Equal(got, data) {
		t.Fatalf("Data = %q, want %q", got, data)
	}
	// id/key are only populated once the object is uploaded, so they are nil
	// for a freshly created object.
	if o.ID() != nil {
		t.Fatalf("ID = %x, want nil before upload", o.ID())
	}
}
