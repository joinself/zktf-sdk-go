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
	"github.com/joinself/zktf-sdk-go/keypair/exchange"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/message"
	"github.com/joinself/zktf-sdk-go/object"
	"github.com/joinself/zktf-sdk-go/revocation"
	"github.com/joinself/zktf-sdk-go/token"
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

// LogField is a single structured key/value pair attached to a LogEntry.
type LogField struct {
	Key   string
	Value string
}

// LogEntry is a structured log record emitted by the native library.
type LogEntry struct {
	Level     LogLevel
	AccountID string
	Target    string
	Message   string
	Timestamp time.Time
	Fields    []LogField
}

// LogHandler receives structured log entries from the native library.
type LogHandler func(LogEntry)

// SetLogHandler registers a process-global handler for native log entries. The
// native log callback carries no per-account context, so the handler is shared
// across every account in the process. A nil handler disables delivery. New
// also registers Config.LogHandler when set.
func SetLogHandler(h LogHandler) {
	if h == nil {
		ffi.SetLogHandler(nil)
		return
	}

	ffi.SetLogHandler(func(e ffi.LogEntry) {
		h(toLogEntry(e))
	})
}

func toLogEntry(e ffi.LogEntry) LogEntry {
	fields := make([]LogField, len(e.Fields))
	for i, f := range e.Fields {
		fields[i] = LogField{Key: f.Key, Value: f.Value}
	}

	return LogEntry{
		Level:     LogLevel(e.Level),
		AccountID: e.AccountID,
		Target:    e.Target,
		Message:   e.Message,
		Timestamp: e.Timestamp,
		Fields:    fields,
	}
}

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
	// LogHandler, when set, is registered as the process-global log handler via
	// SetLogHandler before the account is configured. Because the native log
	// callback is process-global, the most recently configured account wins.
	LogHandler LogHandler
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

	if cfg.LogHandler != nil {
		SetLogHandler(cfg.LogHandler)
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

// KeychainSigningCreate generates a new signing key in the keychain and returns
// its public address.
func (a *Account) KeychainSigningCreate() (*signing.PublicKey, error) {
	k, err := a.h.KeychainSigningCreate()
	if err != nil {
		return nil, err
	}

	return ffi.ToSigningPublicKey(k).(*signing.PublicKey), nil
}

// KeychainExchangeCreate generates a new exchange key in the keychain and
// returns its public address.
func (a *Account) KeychainExchangeCreate() (*exchange.PublicKey, error) {
	k, err := a.h.KeychainExchangeCreate()
	if err != nil {
		return nil, err
	}

	return ffi.ToExchangePublicKey(k).(*exchange.PublicKey), nil
}

// KeychainSign signs payload with the keychain key identified by address.
func (a *Account) KeychainSign(address *signing.PublicKey, payload []byte) ([]byte, error) {
	return a.h.KeychainSign(ffi.SigningPublicKeyOf(address), payload)
}

// KeychainLookupOption is one of the variadic filters accepted by
// KeychainLookup.
type KeychainLookupOption func(*keychainLookupOpts)

type keychainLookupOpts struct {
	identity *signing.PublicKey
	roles    identity.KeyRole
	hasRoles bool
}

// ByIdentity restricts the lookup to keys associated with the given identity.
func ByIdentity(key *signing.PublicKey) KeychainLookupOption {
	return func(o *keychainLookupOpts) {
		o.identity = key
	}
}

// WithRoles restricts the lookup to keys carrying every role in roles. Only
// applies in combination with ByIdentity.
func WithRoles(roles identity.KeyRole) KeychainLookupOption {
	return func(o *keychainLookupOpts) {
		o.roles = roles
		o.hasRoles = true
	}
}

// KeychainLookup returns the signing keys held in the keychain that satisfy the
// given filters. With no filters it returns every signing key.
func (a *Account) KeychainLookup(options ...KeychainLookupOption) ([]*signing.PublicKey, error) {
	var o keychainLookupOpts
	for _, opt := range options {
		opt(&o)
	}

	l := ffi.NewKeychainLookup()
	if o.identity != nil {
		l.ByIdentity(ffi.SigningPublicKeyOf(o.identity))
	}
	if o.hasRoles {
		l.WithRoles(ffi.IdentityKeyRole(o.roles))
	}

	ks, err := a.h.KeychainLookup(l)
	if err != nil {
		return nil, err
	}

	out := make([]*signing.PublicKey, len(ks))
	for i, k := range ks {
		out[i] = ffi.ToSigningPublicKey(k).(*signing.PublicKey)
	}

	return out, nil
}

// SetupPairingCode sets the account up for pairing with an application identity
// and returns the pairing code. It fails if the account is already paired.
func (a *Account) SetupPairingCode() (string, error) {
	return a.h.SetupPairingCode()
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

// CredentialExchangeTrack records that a credential was exchanged with an address.
func (a *Account) CredentialExchangeTrack(with *signing.PublicKey, c *credential.Verifiable) error {
	return a.h.CredentialExchangeTrack(ffi.SigningPublicKeyOf(with), ffi.VerifiableCredentialOf(c))
}

// CredentialExchangeLog returns the credential exchange log, optionally
// restricted to exchanges with an address and to credentials satisfying a
// predicate tree. Either filter may be nil.
func (a *Account) CredentialExchangeLog(with *signing.PublicKey, tree *predicate.Tree) ([]*credential.Exchange, error) {
	var w *ffi.SigningPublicKey
	if with != nil {
		w = ffi.SigningPublicKeyOf(with)
	}

	var t *ffi.PredicateTree
	if tree != nil {
		t = ffi.PredicateTreeOf(tree)
	}

	es, err := a.h.CredentialExchangeLog(w, t)
	if err != nil {
		return nil, err
	}

	out := make([]*credential.Exchange, len(es))
	for i, e := range es {
		out[i] = ffi.ToCredentialExchange(e).(*credential.Exchange)
	}

	return out, nil
}

// CredentialGraphCreate validates a set of presentations against the trusted
// issuer registry and returns the resulting verified credential graph.
func (a *Account) CredentialGraphCreate(registry *credential.TrustedIssuerRegistry, presentations []*credential.VerifiablePresentation, options ...CallOption) (*credential.Graph, error) {
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

// PresentationSign signs a presentation with any account keys it requires.
func (a *Account) PresentationSign(p *credential.VerifiablePresentation) error {
	return a.h.PresentationSign(ffi.VerifiablePresentationOf(p))
}

// PresentationStore stores a presentation on the account for later retrieval.
func (a *Account) PresentationStore(p *credential.VerifiablePresentation) error {
	return a.h.PresentationStore(ffi.VerifiablePresentationOf(p))
}

// PresentationLookup returns presentations stored on the account that satisfy
// the predicate tree. A nil tree returns every stored presentation.
func (a *Account) PresentationLookup(tree *predicate.Tree) ([]*credential.VerifiablePresentation, error) {
	var t *ffi.PredicateTree
	if tree != nil {
		t = ffi.PredicateTreeOf(tree)
	}

	vps, err := a.h.PresentationLookup(t)
	if err != nil {
		return nil, err
	}

	out := make([]*credential.VerifiablePresentation, len(vps))
	for i, vp := range vps {
		out[i] = ffi.ToVerifiablePresentation(vp).(*credential.VerifiablePresentation)
	}

	return out, nil
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

// ObjectStore stores an object in the account's local data store.
func (a *Account) ObjectStore(obj *object.Object) error {
	return a.h.ObjectStore(ffi.ObjectOf(obj))
}

// ObjectRetrieve loads a locally stored object by its id.
func (a *Account) ObjectRetrieve(objectID []byte) (*object.Object, error) {
	o, err := a.h.ObjectRetrieve(objectID)
	if err != nil {
		return nil, err
	}

	return ffi.ToObject(o).(*object.Object), nil
}

// ValueKeys lists the keys of values stored on this account. An empty prefix
// lists every key; otherwise only keys with the given prefix are returned.
func (a *Account) ValueKeys(prefix string) ([]string, error) {
	return a.h.ValueKeys(prefix)
}

// ValueLookup returns the value stored under key. The boolean is false when no
// value is stored for that key.
func (a *Account) ValueLookup(key string) ([]byte, bool, error) {
	return a.h.ValueLookup(key)
}

// ValueStore stores a key/value pair. A zero expires means the value never
// expires; otherwise it is removed at that time.
func (a *Account) ValueStore(key string, value []byte, expires time.Time) error {
	var unix int64
	if !expires.IsZero() {
		unix = expires.Unix()
	}

	return a.h.ValueStore(key, value, unix)
}

// ValueRemove deletes the value stored under key.
func (a *Account) ValueRemove(key string) error {
	return a.h.ValueRemove(key)
}

// TokenIssue issues a fresh token from a validated request.
func (a *Account) TokenIssue(req *token.Request) (*token.Token, error) {
	tk, err := a.h.TokenIssue(ffi.TokenRequestOf(req))
	if err != nil {
		return nil, err
	}

	return ffi.ToToken(tk).(*token.Token), nil
}

// TokenStore stores a token. Its issuer, bearer and local owner are derived
// from the token itself.
func (a *Account) TokenStore(tk *token.Token) error {
	return a.h.TokenStore(ffi.TokenOf(tk))
}
