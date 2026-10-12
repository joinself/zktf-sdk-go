package pairwise_test

import (
	"bytes"
	"testing"

	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/pairwise"
)

func documentAddress(t *testing.T) *identity.Address {
	t.Helper()

	raw := make([]byte, 33)
	for i := range raw[1:] {
		raw[i+1] = 0x17
	}

	key, err := signing.FromBytes(raw)
	if err != nil {
		t.Fatalf("FromBytes: %v", err)
	}

	return identity.AddressZktf(key)
}

func TestNewIdentityRoundTrips(t *testing.T) {
	address := documentAddress(t)
	anchor := bytes.Repeat([]byte{0xab}, 32)

	id := pairwise.NewIdentity(address, anchor)
	if got := id.BiometricAnchorHash(); !bytes.Equal(got, anchor) {
		t.Fatalf("BiometricAnchorHash = %x, want %x", got, anchor)
	}

	decoded, err := pairwise.DecodeIdentity(id.Encode())
	if err != nil {
		t.Fatalf("DecodeIdentity: %v", err)
	}
	if got := decoded.DocumentAddress().String(); got != address.String() {
		t.Fatalf("DocumentAddress = %s, want %s", got, address)
	}
	if got := decoded.BiometricAnchorHash(); !bytes.Equal(got, anchor) {
		t.Fatalf("decoded BiometricAnchorHash = %x, want %x", got, anchor)
	}
}

func TestNewIdentityWithoutAnchor(t *testing.T) {
	id := pairwise.NewIdentity(documentAddress(t), nil)
	if got := id.BiometricAnchorHash(); got != nil {
		t.Fatalf("BiometricAnchorHash = %x, want nil", got)
	}
}
