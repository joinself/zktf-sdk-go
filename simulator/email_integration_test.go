//go:build integration

package simulator_test

import (
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/credential/predicate"
	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/message"
	simulator "github.com/joinself/zktf-sdk-go/simulator"
)

func TestAcquireEmailIsPresentedUnderThePairwiseIdentity(t *testing.T) {
	network := simulator.NewEphemeralNetwork()
	app, err := network.ApplicationCreate("Acme Ltd", "Acme Test")
	if err != nil {
		t.Fatalf("application create: %v", err)
	}
	verifier, err := simulator.NewVerifier(network)
	if err != nil {
		t.Fatalf("verifier: %v", err)
	}
	issuer, err := verifier.Identifier()
	if err != nil {
		t.Fatalf("verifier identifier: %v", err)
	}

	paired := make(chan *signing.PublicKey, 1)
	messages := make(chan *message.Message, 16)
	acc := messagingServer(t, network, paired, messages)
	offer, err := acc.SetupPairingOffer()
	if err != nil {
		t.Fatalf("setup pairing offer: %v", err)
	}
	if err := network.ApplicationPair("Acme Test", offer); err != nil {
		t.Fatalf("application pair: %v", err)
	}
	select {
	case <-paired:
	case <-time.After(30 * time.Second):
		t.Fatal("server did not report paired")
	}

	user := simulator.NewDevice(network)
	t.Cleanup(user.Close)
	if err := user.AcquireEmail("alice@acme.test"); err != nil {
		t.Fatalf("acquire email: %v", err)
	}
	if err := user.Expect(simulator.MatchAny(), simulator.Accept()); err != nil {
		t.Fatalf("expect: %v", err)
	}
	if err := user.Register(app); err != nil {
		t.Fatalf("register: %v", err)
	}

	var from *signing.PublicKey
	var introduction *message.Introduction
	for introduction == nil {
		select {
		case m := <-messages:
			if m.Content().Type() != message.ContentIntroduction {
				continue
			}
			from = m.From()
			introduction, err = message.IntroductionDecode(m.Content())
			if err != nil {
				t.Fatalf("decode introduction: %v", err)
			}
		case <-time.After(10 * time.Second):
			t.Fatal("server did not receive an introduction")
		}
	}
	pairwiseIntroduction, err := introduction.PairwiseIntroduction()
	if err != nil {
		t.Fatalf("pairwise introduction: %v", err)
	}
	relationship, err := acc.PairwiseValidateIntroduction(from, pairwiseIntroduction)
	if err != nil {
		t.Fatalf("validate introduction: %v", err)
	}

	challenge := make([]byte, 32)
	challenge[0] = 7
	presentation := message.NewPresentationRequest().
		Authenticate(relationship, challenge).
		PresentationType(credential.TypeEmail).
		Predicates(predicate.NewTree(predicate.And(
			predicate.Contains("/type", credential.TypeEmail),
			predicate.NotEmpty(credential.FieldSubjectEmailAddress),
		)))
	request, err := presentation.Finish()
	if err != nil {
		t.Fatalf("presentation request: %v", err)
	}
	exchange, err := message.NewExchangeRequest().Action(request.AsAction()).Finish()
	if err != nil {
		t.Fatalf("exchange request: %v", err)
	}
	if err := acc.MessageSend(from, exchange); err != nil {
		t.Fatalf("send: %v", err)
	}

	var response *message.ExchangeResponse
	for response == nil {
		select {
		case m := <-messages:
			if m.Content().Type() != message.ContentExchangeResponse {
				continue
			}
			response, err = message.ExchangeResponseDecode(m.Content())
			if err != nil {
				t.Fatalf("decode response: %v", err)
			}
		case <-time.After(30 * time.Second):
			t.Fatal("server did not receive a presentation response")
		}
	}
	if response.Status() != message.StatusCreated {
		t.Fatalf("want created, got %v: %s", response.Status(), response.ErrorMessage())
	}
	outcomes, err := response.Outcomes()
	if err != nil || len(outcomes) == 0 {
		t.Fatalf("outcomes: %v", err)
	}
	result, err := outcomes[0].AsPresentation()
	if err != nil {
		t.Fatalf("presentation outcome: %v", err)
	}

	registry := credential.NewTrustedIssuerRegistry()
	registry.AddIssuer(identity.AddressZktf(issuer))
	for _, kind := range []string{credential.TypeEmail, "BiometricAnchorCredential", "LivenessAndFacialComparisonCredential"} {
		if err := registry.GrantAuthority(identity.AddressZktf(issuer), kind, time.Now().Add(-time.Hour), time.Time{}); err != nil {
			t.Fatalf("grant authority: %v", err)
		}
	}
	graph, err := acc.CredentialGraphCreate(registry, result.Presentations())
	if err != nil {
		t.Fatalf("credential graph: %v", err)
	}
	if !graph.ValidAuthenticationFor(relationship, challenge) {
		t.Fatal("want the response to authenticate the pairwise identity")
	}
	valid, err := graph.ValidCredentialsFor(relationship.DocumentAddress())
	if err != nil {
		t.Fatalf("valid credentials: %v", err)
	}
	var email string
	for _, c := range valid {
		if hasType(c.Types(), credential.TypeEmail) {
			email = c.Claim("/email/emailAddress")
		}
	}
	if email != "alice@acme.test" {
		t.Fatalf("want the acquired address presented, got %q", email)
	}
}
