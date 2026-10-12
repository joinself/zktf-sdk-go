package crypto_test

import (
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/joinself/zktf-sdk-go/crypto"
)

func TestSHA3256(t *testing.T) {
	want := "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"
	if got := hex.EncodeToString(crypto.SHA3256(nil)); got != want {
		t.Fatalf("SHA3256(empty) = %s, want %s", got, want)
	}

	want = "3a985da74fe225b2045c172d6bd390bd855f086e3e9d525b46bfe24511431532"
	if got := hex.EncodeToString(crypto.SHA3256([]byte("abc"))); got != want {
		t.Fatalf("SHA3256(abc) = %s, want %s", got, want)
	}
}

func TestRandom(t *testing.T) {
	a, b := crypto.RandomID(), crypto.RandomID()
	if len(a) != 20 || len(b) != 20 {
		t.Fatalf("RandomID lengths = %d, %d, want 20", len(a), len(b))
	}
	if bytes.Equal(a, b) {
		t.Fatal("RandomID returned the same id twice")
	}

	if got := crypto.RandomBytes(0); got != nil {
		t.Fatalf("RandomBytes(0) = %x, want nil", got)
	}
	r := crypto.RandomBytes(64)
	if len(r) != 64 {
		t.Fatalf("len(RandomBytes(64)) = %d", len(r))
	}
	if bytes.Equal(r, make([]byte, 64)) {
		t.Fatal("RandomBytes returned all zeroes")
	}
}

func TestEncryptDecryptFile(t *testing.T) {
	dir := t.TempDir()
	plain := filepath.Join(dir, "plain")
	sealed := filepath.Join(dir, "sealed")
	opened := filepath.Join(dir, "opened")

	data := bytes.Repeat([]byte("zktf file encryption "), 4096)
	if err := os.WriteFile(plain, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	key, err := crypto.EncryptFile(plain, sealed)
	if err != nil {
		t.Fatalf("EncryptFile: %v", err)
	}
	if len(key) != crypto.FileKeyLength() {
		t.Fatalf("len(key) = %d, want %d", len(key), crypto.FileKeyLength())
	}

	ciphertext, err := os.ReadFile(sealed)
	if err != nil {
		t.Fatalf("ReadFile sealed: %v", err)
	}
	if bytes.Contains(ciphertext, data[:64]) {
		t.Fatal("encrypted file contains plaintext")
	}

	if err := crypto.DecryptFile(sealed, opened, key); err != nil {
		t.Fatalf("DecryptFile: %v", err)
	}
	got, err := os.ReadFile(opened)
	if err != nil {
		t.Fatalf("ReadFile opened: %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Fatal("decrypted file does not match the original")
	}

	if err := crypto.DecryptFile(sealed, opened, key[:10]); err == nil {
		t.Fatal("DecryptFile with a short key: want error, got nil")
	}
}
