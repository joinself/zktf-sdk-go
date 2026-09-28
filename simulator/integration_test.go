//go:build integration

// Integration tests exercise the native simulator end to end: a simulated
// device registers against a simulated verifier over an in-process network.
//
//	scripts/fetch-native.sh && set -a && . ./.env && set +a
//	go test -tags integration -v ./...
package simulator_test

import (
	"net"
	"runtime"
	"sync/atomic"
	"testing"

	simulator "github.com/joinself/zktf-sdk-go/simulator"
)

func TestRegister(t *testing.T) {
	network := simulator.NewDefaultNetwork()

	verifier, err := simulator.NewVerifier(network)
	if err != nil {
		t.Fatalf("new verifier: %v", err)
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
}

func freePort(t *testing.T) uint16 {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	return uint16(l.Addr().(*net.TCPAddr).Port)
}

func TestRegisterSurvivesCollection(t *testing.T) {
	network := simulator.NewNetwork(freePort(t), freePort(t), freePort(t), freePort(t))

	verifier, err := simulator.NewVerifier(network)
	if err != nil {
		t.Fatalf("new verifier: %v", err)
	}

	identifier, err := verifier.Identifier()
	if err != nil {
		t.Fatalf("verifier identifier: %v", err)
	}

	var stop atomic.Bool
	collected := make(chan struct{})
	go func() {
		defer close(collected)
		for !stop.Load() {
			runtime.GC()
		}
	}()
	defer func() {
		stop.Store(true)
		<-collected
	}()

	for i := range 5 {
		device := simulator.NewDevice(network)
		device.Expect(simulator.MatchAny(), simulator.Accept())

		if err := device.Register(identifier); err != nil {
			t.Fatalf("register %d: %v", i, err)
		}
	}
}
