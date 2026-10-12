package crypto_test

import (
	"bytes"
	"testing"

	"github.com/joinself/zktf-sdk-go/crypto"
)

func TestRandomID(t *testing.T) {
	a, b := crypto.RandomID(), crypto.RandomID()
	if len(a) != 20 || len(b) != 20 {
		t.Fatalf("RandomID lengths = %d, %d, want 20", len(a), len(b))
	}
	if bytes.Equal(a, b) {
		t.Fatal("RandomID returned the same id twice")
	}
}
