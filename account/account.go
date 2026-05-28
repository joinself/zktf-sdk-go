// Package account is the entry point to the zktf SDK: it configures an account,
// drives its event callbacks, and exposes every server-side operation an
// application performs against the network. It is pure Go and never exposes a
// cgo type.
package account

import (
	"time"

	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/credential/predicate"
	"github.com/joinself/zktf-sdk-go/crypto"
	"github.com/joinself/zktf-sdk-go/group"
	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/message"
	"github.com/joinself/zktf-sdk-go/object"
	"github.com/joinself/zktf-sdk-go/revocation"
	"github.com/joinself/zktf-sdk-go/trust"
)

// Network identifies the underlying zktf network. Most callers should use the
// pre-configured TargetProduction or TargetSandbox; the other values exist for
// internal testing.
type Network uint32

const (
	NetworkProduction  Network = Network(ffi.NetworkProduction)
	NetworkSandbox     Network = Network(ffi.NetworkSandbox)
	NetworkStaging     Network = Network(ffi.NetworkStaging)
	NetworkPreview     Network = Network(ffi.NetworkPreview)
	NetworkDevelopment Network = Network(ffi.NetworkDevelopment)
)

// Target carries a zktf network selection plus the endpoints to reach it.
// TargetProduction and TargetSandbox are pre-baked for the two public
// networks; custom deployments are constructed with &Target{Network: ..., ...}.
type Target struct {
	Network         Network
	RPCEndpoint     string
	ObjectEndpoint  string
	MessageEndpoint string
}

// Pre-configured targets for the two public networks.
var (
	TargetProduction = &Target{
		Network:         NetworkProduction,
		RPCEndpoint:     "https://rpc.joinself.com/",
		ObjectEndpoint:  "https://object.joinself.com/",
		MessageEndpoint: "wss://message.joinself.com/",
	}
	TargetSandbox = &Target{
		Network:         NetworkSandbox,
		RPCEndpoint:     "https://rpc.sandbox.joinself.com/",
		ObjectEndpoint:  "https://object.sandbox.joinself.com/",
		MessageEndpoint: "wss://message.sandbox.joinself.com/",
	}
)

// LogLevel controls native library logging verbosity.
type LogLevel uint32

const (
	LogError LogLevel = LogLevel(ffi.LogError)
	LogWarn  LogLevel = LogLevel(ffi.LogWarn)
	LogInfo  LogLevel = LogLevel(ffi.LogInfo)
	LogDebug LogLevel = LogLevel(ffi.LogDebug)
	LogTrace LogLevel = LogLevel(ffi.LogTrace)
)

// Config configures an account.
type Config struct {
	// Target selects the network and its endpoints. Defaults to TargetProduction
	// when nil.
	Target *Target
	// StoragePath is the on-disk path for the account's storage; ":memory:" for
	// ephemeral storage. Defaults to ":memory:" when empty.
	StoragePath string
	// EncryptionKey is a 32-byte key encrypting the account's storage. Defaults
	// to 32 zero bytes when nil.
	EncryptionKey []byte
	// LogLevel controls verbosity of native logging. Defaults to LogError.
	LogLevel LogLevel
}

// Account is a configured zktf account.
type Account struct {
	h *ffi.Account
}

// New allocates, configures, and returns an account.
func New(cfg Config, cb Callbacks) (*Account, error) {
	target := cfg.Target
	if target == nil {
		target = TargetProduction
	}

	storage := cfg.StoragePath
	if storage == "" {
		storage = ":memory:"
	}

	key := cfg.EncryptionKey
	if key == nil {
		key = make([]byte, 32)
	}

	logLevel := cfg.LogLevel
	if logLevel == 0 {
		logLevel = LogError
	}

	a := &Account{h: ffi.NewAccount()}

	err := a.h.Configure(ffi.AccountConfig{
		Network:         ffi.Network(target.Network),
		RPCEndpoint:     target.RPCEndpoint,
		ObjectEndpoint:  target.ObjectEndpoint,
		MessageEndpoint: target.MessageEndpoint,
		StoragePath:     storage,
		EncryptionKey:   key,
		LogLevel:        ffi.LogLevel(logLevel),
	}, adapter{cb: cb})
	if err != nil {
		return nil, err
	}

	return a, nil
}

// MessageSend sends content to the given recipient address. Delivery is reported
// asynchronously via OnEvent (Acknowledged / SendFailed).
func (a *Account) MessageSend(to *signing.PublicKey, content *message.Content) error {
	return a.h.MessageSend(ffi.SigningPublicKeyOf(to), ffi.ContentOf(content))
}

// NotificationSend pushes a notification carrying a content summary to the
// given address.
func (a *Account) NotificationSend(to *signing.PublicKey, summary *message.ContentSummary, options ...CallOption) error {
	o := collectCallOpts(options)

	return a.h.NotificationSend(ffi.SigningPublicKeyOf(to), ffi.MessageContentSummaryOf(summary), o.timeout)
}

// InboxOpen opens a new messaging inbox and returns its address.
func (a *Account) InboxOpen(options ...InboxOpenOption) (*signing.PublicKey, error) {
	o := collectInboxOpenOpts(options)

	k, err := a.h.InboxOpen(o.timeout)
	if err != nil {
		return nil, err
	}

	_ = o.expires // TODO: thread expires through to the FFI when wired

	return ffi.ToSigningPublicKey(k).(*signing.PublicKey), nil
}

// InboxDefault returns the account's default inbox address.
func (a *Account) InboxDefault() (*signing.PublicKey, error) {
	k, err := a.h.InboxDefault()
	if err != nil {
		return nil, err
	}

	return ffi.ToSigningPublicKey(k).(*signing.PublicKey), nil
}

// InboxClose closes an open inbox.
func (a *Account) InboxClose(address *signing.PublicKey, options ...CallOption) error {
	o := collectCallOpts(options)

	return a.h.InboxClose(ffi.SigningPublicKeyOf(address), o.timeout)
}

// InboxList returns the addresses of all open inboxes on this account.
func (a *Account) InboxList() ([]*signing.PublicKey, error) {
	ks, err := a.h.InboxList()
	if err != nil {
		return nil, err
	}

	out := make([]*signing.PublicKey, len(ks))
	for i, k := range ks {
		out[i] = ffi.ToSigningPublicKey(k).(*signing.PublicKey)
	}

	return out, nil
}

// GroupNegotiate establishes an encrypted session between two inbox addresses.
// The SDK auto-accepts the resulting invite/welcome on both sides. A zero
// expires means no expiry.
func (a *Account) GroupNegotiate(as, with *signing.PublicKey, expires time.Time) error {
	var unix int64
	if !expires.IsZero() {
		unix = expires.Unix()
	}

	return a.h.GroupNegotiate(ffi.SigningPublicKeyOf(as), ffi.SigningPublicKeyOf(with), unix)
}

// GroupNegotiateOutOfBand creates an MLS key package the recipient can use to
// establish an encrypted session with this account.
func (a *Account) GroupNegotiateOutOfBand(as *signing.PublicKey, expires time.Time) (*crypto.KeyPackage, error) {
	var unix int64
	if !expires.IsZero() {
		unix = expires.Unix()
	}

	kp, err := a.h.GroupNegotiateOutOfBand(ffi.SigningPublicKeyOf(as), unix)
	if err != nil {
		return nil, err
	}

	return ffi.ToCryptoKeyPackage(kp).(*crypto.KeyPackage), nil
}

// GroupEstablish uses a received key package to establish an encrypted group
// session.
func (a *Account) GroupEstablish(as *signing.PublicKey, kp *crypto.KeyPackage, options ...CallOption) (*group.Group, error) {
	o := collectCallOpts(options)

	g, err := a.h.GroupEstablish(ffi.SigningPublicKeyOf(as), ffi.CryptoKeyPackageOf(kp), o.timeout)
	if err != nil {
		return nil, err
	}

	return ffi.ToGroup(g).(*group.Group), nil
}

// GroupAccept accepts a received welcome to join an encrypted group session.
func (a *Account) GroupAccept(as *signing.PublicKey, welcome *crypto.Welcome, options ...CallOption) (*group.Group, error) {
	o := collectCallOpts(options)

	g, err := a.h.GroupAccept(ffi.SigningPublicKeyOf(as), ffi.CryptoWelcomeOf(welcome), o.timeout)
	if err != nil {
		return nil, err
	}

	return ffi.ToGroup(g).(*group.Group), nil
}

// GroupLookup returns groups matching the lookup query.
func (a *Account) GroupLookup(options ...group.LookupOption) ([]*group.Group, error) {
	l := group.BuildLookup(options...)

	gs, err := a.h.GroupLookup(ffi.GroupLookupOf(l))
	if err != nil {
		return nil, err
	}

	out := make([]*group.Group, len(gs))
	for i, g := range gs {
		out[i] = ffi.ToGroup(g).(*group.Group)
	}

	return out, nil
}

// GroupUpdate publishes a built group update.
func (a *Account) GroupUpdate(r *group.UpdateRequest) error {
	return a.h.GroupUpdate(ffi.GroupUpdateRequestOf(r))
}

// GroupLeave leaves a group.
func (a *Account) GroupLeave(g *group.Group) error {
	return a.h.GroupLeave(ffi.GroupOf(g))
}

// CredentialIssue signs an unsigned credential into a verifiable credential.
func (a *Account) CredentialIssue(c *credential.Credential) (*credential.Verifiable, error) {
	vc, err := a.h.CredentialIssue(ffi.CredentialOf(c))
	if err != nil {
		return nil, err
	}

	return ffi.ToVerifiableCredential(vc).(*credential.Verifiable), nil
}

// CredentialStore stores a verifiable credential in the account's local store.
func (a *Account) CredentialStore(c *credential.Verifiable) error {
	return a.h.CredentialStore(ffi.VerifiableCredentialOf(c))
}

// CredentialLookup returns credentials in the account's local store that
// satisfy the given predicate tree.
func (a *Account) CredentialLookup(tree *predicate.Tree) ([]*credential.Verifiable, error) {
	cs, err := a.h.CredentialLookup(ffi.PredicateTreeOf(tree))
	if err != nil {
		return nil, err
	}

	out := make([]*credential.Verifiable, len(cs))
	for i, c := range cs {
		out[i] = ffi.ToVerifiableCredential(c).(*credential.Verifiable)
	}

	return out, nil
}

// CredentialSharedWith returns the credentials this account has shared with the
// given address that satisfy the predicate tree.
func (a *Account) CredentialSharedWith(with *signing.PublicKey, tree *predicate.Tree) ([]*credential.Verifiable, error) {
	cs, err := a.h.CredentialSharedWith(ffi.SigningPublicKeyOf(with), ffi.PredicateTreeOf(tree))
	if err != nil {
		return nil, err
	}

	out := make([]*credential.Verifiable, len(cs))
	for i, c := range cs {
		out[i] = ffi.ToVerifiableCredential(c).(*credential.Verifiable)
	}

	return out, nil
}

// CredentialGraphCreate validates a set of presentations against the trusted
// issuer registry and returns the resulting verified credential graph.
func (a *Account) CredentialGraphCreate(registry *trust.Registry, presentations []*credential.VerifiablePresentation, options ...CallOption) (*credential.Graph, error) {
	o := collectCallOpts(options)

	ps := make([]*ffi.VerifiablePresentation, len(presentations))
	for i, p := range presentations {
		ps[i] = ffi.VerifiablePresentationOf(p)
	}

	g, err := a.h.CredentialGraphCreate(ffi.TrustedIssuerRegistryOf(registry), ps, o.timeout)
	if err != nil {
		return nil, err
	}

	return ffi.ToCredentialGraph(g).(*credential.Graph), nil
}

// PresentationIssue signs an unsigned presentation into a verifiable presentation.
func (a *Account) PresentationIssue(p *credential.Presentation) (*credential.VerifiablePresentation, error) {
	vp, err := a.h.PresentationIssue(ffi.PresentationOf(p))
	if err != nil {
		return nil, err
	}

	return ffi.ToVerifiablePresentation(vp).(*credential.VerifiablePresentation), nil
}

// IdentityResolve resolves the identity document for an address.
func (a *Account) IdentityResolve(address *credential.Address, options ...CallOption) (*identity.Document, error) {
	o := collectCallOpts(options)

	d, err := a.h.IdentityResolve(ffi.DIDAddressOf(address), o.timeout)
	if err != nil {
		return nil, err
	}

	return ffi.ToIdentityDocument(d).(*identity.Document), nil
}

// IdentityExecute publishes an identity operation.
func (a *Account) IdentityExecute(op *identity.Operation, options ...CallOption) error {
	o := collectCallOpts(options)

	return a.h.IdentityExecute(ffi.IdentityOperationOf(op), o.timeout)
}

// IdentitySign signs an identity operation with this account's keys.
func (a *Account) IdentitySign(op *identity.Operation) error {
	return a.h.IdentitySign(ffi.IdentityOperationOf(op))
}

// IdentityLookup returns DID addresses matching the lookup query.
func (a *Account) IdentityLookup(options ...identity.LookupOption) ([]*credential.Address, error) {
	l := identity.BuildLookup(options...)

	as, err := a.h.IdentityLookup(ffi.IdentityLookupOf(l))
	if err != nil {
		return nil, err
	}

	out := make([]*credential.Address, len(as))
	for i, addr := range as {
		out[i] = ffi.ToDIDAddress(addr).(*credential.Address)
	}

	return out, nil
}

// RevocationSign signs an unsigned revocation statement with this account's keys.
func (a *Account) RevocationSign(statement *revocation.Statement) error {
	return a.h.RevocationSign(ffi.RevocationStatementOf(statement))
}

// RevocationRevoke publishes a signed revocation statement.
func (a *Account) RevocationRevoke(statement *revocation.Statement, options ...CallOption) error {
	o := collectCallOpts(options)

	return a.h.RevocationRevoke(ffi.RevocationStatementOf(statement), o.timeout)
}

// ObjectUpload uploads an object to the object store.
func (a *Account) ObjectUpload(obj *object.Object, options ...ObjectUploadOption) error {
	o := collectObjectUploadOpts(options)

	var uploadOpts *ffi.ObjectUploadOptions
	if o.objectPersistLocally {
		uploadOpts = ffi.NewObjectUploadOptions()
		uploadOpts.PersistLocally(true)
	}

	return a.h.ObjectUpload(ffi.ObjectOf(obj), uploadOpts, o.timeout)
}

// ObjectDownload downloads an object's encrypted bytes from the server.
func (a *Account) ObjectDownload(obj *object.Object, options ...CallOption) error {
	o := collectCallOpts(options)

	return a.h.ObjectDownload(ffi.ObjectOf(obj), o.timeout)
}
