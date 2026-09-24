//go:build integration

// Integration tests exercise the native simulator end to end: a simulated
// device registers against a simulated verifier over an in-process network.
//
//	scripts/fetch-native.sh && set -a && . ./.env && set +a
//	go test -tags integration -v ./...
package simulator_test

import (
	"sync"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/identity"
	simulator "github.com/joinself/zktf-sdk-go/simulator"
)

func TestRegister(t *testing.T) {
	registeredDevice(t)
}

var (
	shared    sync.Once
	network   *simulator.Network
	verifier  *simulator.Verifier
	sharedErr error
)

func registeredDevice(t *testing.T) *simulator.Device {
	t.Helper()

	shared.Do(func() {
		network = simulator.NewDefaultNetwork()
		verifier, sharedErr = simulator.NewVerifier(network)
	})
	if sharedErr != nil {
		t.Fatalf("new verifier: %v", sharedErr)
	}

	device := simulator.NewDevice(network)
	device.Expect(simulator.MatchAny(), simulator.Accept())

	identifier, err := verifier.Identifier()
	if err != nil {
		t.Fatalf("verifier identifier: %v", err)
	}

	if err := device.Register(identifier); err != nil {
		t.Fatalf("register: %v", err)
	}

	return device
}

func TestControllerAnchor(t *testing.T) {
	device := registeredDevice(t)

	identifier, err := device.SigningKeyCreate()
	if err != nil {
		t.Fatalf("signing key create: %v", err)
	}

	subject := identity.AddressZktf(identifier)
	unsigned, err := credential.NewBuilder().
		Type("VerifiableCredential", "OrganisationCredential").
		Issuer(subject).
		Subject(subject).
		Claim("organisationName", "Example Ltd").
		ValidFrom(time.Now()).
		SignWith(identifier, time.Now()).
		Finish()
	if err != nil {
		t.Fatalf("build credential: %v", err)
	}
	encoded, err := unsigned.Encode()
	if err != nil {
		t.Fatalf("encode credential: %v", err)
	}

	if _, err := device.MintControllerIdentity(identifier, encoded); err != nil {
		t.Fatalf("mint controller identity: %v", err)
	}

	anchor, err := device.ControllerAnchor(identifier, 10*time.Second)
	if err != nil {
		t.Fatalf("controller anchor: %v", err)
	}

	decoded, err := credential.Decode(anchor)
	if err != nil {
		t.Fatalf("decode anchor: %v", err)
	}

	found := false
	for _, kind := range decoded.Types() {
		if kind == "BiometricAnchorCredential" {
			found = true
		}
	}
	if !found {
		t.Fatalf("want a biometric anchor credential, got %v", decoded.Types())
	}
}

func TestSignIdentity(t *testing.T) {
	device := registeredDevice(t)

	identifier, err := device.SigningKeyCreate()
	if err != nil {
		t.Fatalf("identifier: %v", err)
	}
	invocation, err := device.SigningKeyCreate()
	if err != nil {
		t.Fatalf("invocation key: %v", err)
	}

	operation, err := identity.NewOperationBuilder().
		ID(identifier).
		Sequence(0).
		Timestamp(time.Now()).
		GrantSigning(invocation, identity.KeyRoleInvocation).
		Finish()
	if err != nil {
		t.Fatalf("build operation: %v", err)
	}

	signed, details, err := device.SignIdentity(operation, nil)
	if err != nil {
		t.Fatalf("sign identity: %v", err)
	}
	if details != nil {
		t.Fatalf("want no credential without one to sign, got %d bytes", len(details))
	}

	decoded, err := identity.DecodeOperation(identifier, signed)
	if err != nil {
		t.Fatalf("decode signed operation: %v", err)
	}
	if !decoded.SignedBy(identifier) || !decoded.SignedBy(invocation) {
		t.Fatal("want the operation signed by the identifier and the granted key")
	}
}
