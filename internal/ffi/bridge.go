// Bridge trampolines for the public packages.
//
// Each public type that needs to cross a package boundary has a pair of
// function-typed vars here. The owning public package populates both in its
// init():
//
//	<Type>Of(any) *<Type>      // extract the ffi handle from a public wrapper
//	To<Type>(*<Type>) any      // build a public wrapper from an ffi handle
//
// Sibling packages call these (with a type assertion on the create side) so
// public types can keep their handle field unexported AND not expose Wrap/
// Unwrap helpers in their godoc. Misuse panics at runtime — the trampolines
// are only ever wired up correctly from within this SDK.
package ffi

// keypair
var (
	SigningPublicKeyOf  func(any) *SigningPublicKey
	ToSigningPublicKey  func(*SigningPublicKey) any
	ExchangePublicKeyOf func(any) *ExchangePublicKey
	ToExchangePublicKey func(*ExchangePublicKey) any
)

// credential
var (
	DIDAddressOf               func(any) *DIDAddress
	ToDIDAddress               func(*DIDAddress) any
	CredentialTermOf           func(any) *CredentialTerm
	ToCredentialTerm           func(*CredentialTerm) any
	CredentialOf               func(any) *Credential
	ToCredential               func(*Credential) any
	VerifiableCredentialOf     func(any) *VerifiableCredential
	ToVerifiableCredential     func(*VerifiableCredential) any
	PresentationOf             func(any) *Presentation
	ToPresentation             func(*Presentation) any
	VerifiablePresentationOf   func(any) *VerifiablePresentation
	ToVerifiablePresentation   func(*VerifiablePresentation) any
	CredentialGraphOf          func(any) *CredentialGraph
	ToCredentialGraph          func(*CredentialGraph) any
	RevocationProofOf          func(any) *RevocationProof
	ToRevocationProof          func(*RevocationProof) any
	CredentialExchangeOf       func(any) *CredentialExchange
	ToCredentialExchange       func(*CredentialExchange) any
)

// credential/predicate
var (
	PredicateTreeOf func(any) *PredicateTree
	ToPredicateTree func(*PredicateTree) any
)

// identity
var (
	IdentityDocumentOf  func(any) *IdentityDocument
	ToIdentityDocument  func(*IdentityDocument) any
	IdentityOperationOf func(any) *IdentityOperation
	ToIdentityOperation func(*IdentityOperation) any
	IdentityLookupOf    func(any) *IdentityLookup
	ToIdentityLookup    func(*IdentityLookup) any
)

// group
var (
	GroupOf              func(any) *Group
	ToGroup              func(*Group) any
	GroupLookupOf        func(any) *GroupLookup
	ToGroupLookup        func(*GroupLookup) any
	GroupUpdateRequestOf func(any) *GroupUpdateRequest
	ToGroupUpdateRequest func(*GroupUpdateRequest) any
)

// crypto
var (
	CryptoKeyPackageOf func(any) *CryptoKeyPackage
	ToCryptoKeyPackage func(*CryptoKeyPackage) any
	CryptoWelcomeOf    func(any) *CryptoWelcome
	ToCryptoWelcome    func(*CryptoWelcome) any
)

// pairwise
var (
	PairwiseIdentityOf     func(any) *PairwiseIdentity
	ToPairwiseIdentity     func(*PairwiseIdentity) any
	PairwiseIntroductionOf func(any) *PairwiseIntroduction
	ToPairwiseIntroduction func(*PairwiseIntroduction) any
)

// revocation
var (
	RevocationStatementOf func(any) *RevocationStatement
	ToRevocationStatement func(*RevocationStatement) any
)

// token / object / trust
var (
	TokenOf                 func(any) *Token
	ToToken                 func(*Token) any
	TokenRequestOf          func(any) *TokenRequest
	ToTokenRequest          func(*TokenRequest) any
	ObjectOf                func(any) *Object
	ToObject                func(*Object) any
	TrustedIssuerRegistryOf func(any) *TrustedIssuerRegistry
	ToTrustedIssuerRegistry func(*TrustedIssuerRegistry) any
)

// message — content + meta
var (
	ContentOf               func(any) *Content
	ToContent               func(*Content) any
	MessageOf               func(any) *Message
	ToMessage               func(*Message) any
	MessageContentSummaryOf func(any) *MessageContentSummary
	ToMessageContentSummary func(*MessageContentSummary) any
)

// message — exchange action/outcome polymorphism
var (
	ActionOf         func(any) *Action
	ToAction         func(*Action) any
	OutcomeOf        func(any) *Outcome
	ToOutcome        func(*Outcome) any
	ExchangeRequestOf  func(any) *ExchangeRequest
	ToExchangeRequest  func(*ExchangeRequest) any
	ExchangeResponseOf func(any) *ExchangeResponse
	ToExchangeResponse func(*ExchangeResponse) any
)

// message — per-kind request/response bodies
var (
	PresentationActionOf       func(any) *PresentationAction
	ToPresentationAction       func(*PresentationAction) any
	PresentationResultOf       func(any) *PresentationResult
	ToPresentationResult       func(*PresentationResult) any
	VerificationActionOf       func(any) *VerificationAction
	ToVerificationAction       func(*VerificationAction) any
	VerificationResultOf       func(any) *VerificationResult
	ToVerificationResult       func(*VerificationResult) any
	IdentitySigningActionOf    func(any) *IdentitySigningAction
	ToIdentitySigningAction    func(*IdentitySigningAction) any
	IdentitySigningResultOf    func(any) *IdentitySigningResult
	ToIdentitySigningResult    func(*IdentitySigningResult) any
	DevicePairingActionOf      func(any) *DevicePairingAction
	ToDevicePairingAction      func(*DevicePairingAction) any
	DevicePairingResultOf      func(any) *DevicePairingResult
	ToDevicePairingResult      func(*DevicePairingResult) any
)

// event — status + group + workflow + the wire events nested under group
var (
	StatusEventOf    func(any) *StatusEvent
	ToStatusEvent    func(*StatusEvent) any
	GroupEventOf     func(any) *GroupEvent
	ToGroupEvent     func(*GroupEvent) any
	WorkflowEventOf  func(any) *WorkflowEvent
	ToWorkflowEvent  func(*WorkflowEvent) any
	KeyPackageEventOf func(any) *KeyPackageEvent
	ToKeyPackageEvent func(*KeyPackageEvent) any
	WelcomeEventOf   func(any) *WelcomeEvent
	ToWelcomeEvent   func(*WelcomeEvent) any
	CommitEventOf    func(any) *CommitEvent
	ToCommitEvent    func(*CommitEvent) any
	ProposalEventOf  func(any) *ProposalEvent
	ToProposalEvent  func(*ProposalEvent) any
	DroppedEventOf   func(any) *DroppedEvent
	ToDroppedEvent   func(*DroppedEvent) any
)
