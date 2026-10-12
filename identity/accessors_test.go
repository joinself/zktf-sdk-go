package identity_test

import (
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

func otherKey(t *testing.T) *signing.PublicKey {
	t.Helper()

	raw := make([]byte, 33)
	for i := range raw[1:] {
		raw[i+1] = 0x42
	}

	key, err := signing.FromBytes(raw)
	if err != nil {
		t.Fatalf("FromBytes: %v", err)
	}

	return key
}

func TestAddressMethodAndSigningKey(t *testing.T) {
	address := testKey(t)
	key := otherKey(t)

	if got := identity.AddressKey(address).Method(); got != identity.AddressMethodKey {
		t.Fatalf("AddressKey Method = %d, want AddressMethodKey", got)
	}
	if got := identity.AddressZktf(address).Method(); got != identity.AddressMethodZktf {
		t.Fatalf("AddressZktf Method = %d, want AddressMethodZktf", got)
	}

	if k := identity.AddressZktf(address).SigningKey(); k != nil {
		t.Fatalf("AddressZktf SigningKey = %s, want nil", k)
	}

	withKey := identity.AddressZktfWithKey(address, key)
	k := withKey.SigningKey()
	if k == nil || !k.Matches(key) {
		t.Fatalf("AddressZktfWithKey SigningKey = %v, want %s", k, key)
	}

	parsed, err := identity.ParseAddress(withKey.String())
	if err != nil {
		t.Fatalf("ParseAddress: %v", err)
	}
	if k := parsed.SigningKey(); k == nil || !k.Matches(key) {
		t.Fatalf("parsed SigningKey = %v, want %s", k, key)
	}
}

func TestOperationThresholdsAndWeights(t *testing.T) {
	document := testKey(t)
	key := otherKey(t)

	op, err := identity.NewOperationBuilder().
		ID(document).
		Sequence(0).
		Timestamp(time.Unix(1_900_000_000, 0)).
		GrantSigning(key, identity.KeyRoleVerification).
		Threshold(identity.KeyRoleVerification, 2).
		Weight(key, identity.KeyRoleVerification, 3).
		Finish()
	if err != nil {
		t.Fatalf("Finish: %v", err)
	}

	thresholds := op.Thresholds()
	if len(thresholds) != 1 {
		t.Fatalf("len(Thresholds) = %d, want 1", len(thresholds))
	}
	if thresholds[0].Role != identity.KeyRoleVerification || thresholds[0].Value != 2 {
		t.Fatalf("Thresholds[0] = %+v, want verification/2", thresholds[0])
	}

	weights := op.Weights()
	if len(weights) != 1 {
		t.Fatalf("len(Weights) = %d, want 1", len(weights))
	}
	w := weights[0]
	if w.Role != identity.KeyRoleVerification || w.Value != 3 || w.Key == nil || !w.Key.Matches(key) {
		t.Fatalf("Weights[0] = %+v, want %s verification/3", w, key)
	}
}
