//go:build integration

package simulator_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/account"
	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

type founding struct {
	acc          *account.Account
	controller   *signing.PublicKey
	organisation *signing.PublicKey
	genesis      *identity.Operation
	details      *credential.VerifiableCredential
	presentation *credential.VerifiablePresentation
	registry     *credential.TrustedIssuerRegistry
}

func found(t *testing.T) *founding {
	t.Helper()

	acc, user := applicationUser(t)
	controller, minted := mintController(t, user)

	document, err := acc.IdentityResolve(identity.AddressZktf(controller))
	if err != nil {
		t.Fatalf("resolve controller: %v", err)
	}
	invocation := document.SigningKeys(identity.WithRoles(identity.KeyRoleInvocation))
	assertion := document.SigningKeys(identity.WithRoles(identity.KeyRoleAssertion))
	if len(invocation) == 0 || len(assertion) == 0 {
		t.Fatal("want controller invocation and assertion keys")
	}

	organisation, err := user.SigningKeyCreate()
	if err != nil {
		t.Fatalf("organisation identifier: %v", err)
	}
	method := uint16(identity.AddressMethodZktf)
	operation, err := identity.NewOperationBuilder().
		ID(organisation).
		Sequence(0).
		Timestamp(time.Now()).
		GrantSigningReferenced(method, controller, invocation[0], document.Commitment(), identity.KeyRoleInvocation).
		GrantSigningReferenced(method, controller, assertion[0], document.Commitment(), identity.KeyRoleAssertion).
		Weight(invocation[0], identity.KeyRoleInvocation, 1).
		Threshold(identity.KeyRoleInvocation, 1).
		Finish()
	if err != nil {
		t.Fatalf("build operation: %v", err)
	}

	unsigned, err := credential.NewBuilder().
		Type("VerifiableCredential", "OrganisationDetailsCredential").
		Issuer(identity.AddressZktf(organisation)).
		Claim("organisationName", "Example Ltd").
		ValidFrom(time.Now()).
		SignWith(assertion[0], time.Now()).
		Finish()
	if err != nil {
		t.Fatalf("build details: %v", err)
	}

	signed, err := user.SignIdentity(operation, unsigned)
	if err != nil {
		t.Fatalf("sign identity: %v", err)
	}
	if len(signed.Presentation) == 0 {
		t.Fatal("want a presentation")
	}

	genesis, err := identity.DecodeOperation(organisation, signed.Operation)
	if err != nil {
		t.Fatalf("decode genesis: %v", err)
	}
	details, err := credential.Decode(signed.Credential)
	if err != nil {
		t.Fatalf("decode details: %v", err)
	}
	presentation, err := credential.DecodePresentation(signed.Presentation)
	if err != nil {
		t.Fatalf("decode presentation: %v", err)
	}

	anchor, err := credential.Decode(minted.Anchors[0])
	if err != nil {
		t.Fatalf("decode anchor: %v", err)
	}
	registry := credential.NewTrustedIssuerRegistry()
	registry.AddIssuer(anchor.Issuer())
	for _, kind := range []string{"BiometricAnchorCredential", "LivenessAndFacialComparisonCredential"} {
		if err := registry.GrantAuthority(anchor.Issuer(), kind, time.Now().Add(-time.Hour), time.Time{}); err != nil {
			t.Fatalf("grant authority: %v", err)
		}
	}

	return &founding{
		acc:          acc,
		controller:   controller,
		organisation: organisation,
		genesis:      genesis,
		details:      details,
		presentation: presentation,
		registry:     registry,
	}
}

func TestFoundingIsAuthorizedByTheControllersPresentation(t *testing.T) {
	f := found(t)

	if !bytes.Equal(f.presentation.Holder().Key().Bytes(), f.controller.Bytes()) {
		t.Fatal("want the controller as the presentation holder")
	}

	graph, err := f.acc.CredentialGraphCreate(f.registry, []*credential.VerifiablePresentation{f.presentation})
	if err != nil {
		t.Fatalf("credential graph: %v", err)
	}
	if !graph.ValidDocumentFor(identity.AddressZktf(f.controller)) {
		t.Fatal("want the controller's document to be valid")
	}
	revoked, err := graph.RevokedCredentialsFor(identity.AddressZktf(f.controller))
	if err != nil {
		t.Fatalf("revoked credentials: %v", err)
	}
	if len(revoked) != 0 {
		t.Fatalf("want no revoked credentials, got %d", len(revoked))
	}

	verdict, err := graph.ValidOperationFor(nil, f.genesis)
	if err != nil {
		t.Fatalf("valid operation: %v", err)
	}
	if !verdict.Authorized {
		t.Fatalf("want the genesis authorized, got %+v", verdict.Reasons)
	}
	if verdict.Document == nil {
		t.Fatal("want the organisation document")
	}

	chain, err := graph.AuthorizationFor(f.genesis.Hash())
	if err != nil {
		t.Fatalf("genesis chain: %v", err)
	}
	people := chain.People()
	if len(people) != 1 || !bytes.Equal(people[0].Identity.Key().Bytes(), f.controller.Bytes()) {
		t.Fatalf("want exactly the controller behind the genesis, got %d people", len(people))
	}
	if !chain.AuthorizedBy(identity.AddressZktf(f.controller)) {
		t.Fatal("want the chain authorized by the controller")
	}

	hash, err := f.details.ContentHash()
	if err != nil {
		t.Fatalf("details hash: %v", err)
	}
	details, err := graph.AuthorizationFor(hash)
	if err != nil {
		t.Fatalf("details chain: %v", err)
	}
	assertion, err := details.Verdict(verdict.Document, identity.KeyRoleAssertion)
	if err != nil {
		t.Fatalf("details verdict: %v", err)
	}
	if !assertion.Authorized {
		t.Fatalf("want the details authorized for the organisation, got %+v", assertion.Reasons)
	}
}

func TestFoundingWithoutThePresentationIsUnproven(t *testing.T) {
	f := found(t)

	graph, err := f.acc.CredentialGraphCreate(f.registry, nil)
	if err != nil {
		t.Fatalf("credential graph: %v", err)
	}

	verdict, err := graph.ValidOperationFor(nil, f.genesis)
	if err != nil {
		t.Fatalf("valid operation: %v", err)
	}
	if verdict.Authorized {
		t.Fatal("want the genesis unauthorized without the controller's evidence")
	}
	if len(verdict.Reasons) == 0 || verdict.Reasons[0].Kind != credential.ReasonPersonUnproven {
		t.Fatalf("want an unproven person, got %+v", verdict.Reasons)
	}
	if verdict.Document != nil {
		t.Fatal("want no document for an unauthorized operation")
	}

	var anchors int
	for _, requirement := range verdict.Missing {
		if requirement.Kind == credential.RequirementAnchor {
			anchors++
		}
	}
	if anchors == 0 {
		t.Fatalf("want an anchor requirement, got %+v", verdict.Missing)
	}
}
