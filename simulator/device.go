package simulator

import (
	"fmt"
	"time"

	"github.com/joinself/zktf-sdk-go/internal/simffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/message"
)

// ContentType selects a message content type for MatchContentType.
type ContentType uint32

const (
	ContentUnknown           ContentType = ContentType(simffi.ContentUnknown)
	ContentCustom            ContentType = ContentType(simffi.ContentCustom)
	ContentChat              ContentType = ContentType(simffi.ContentChat)
	ContentReceipt           ContentType = ContentType(simffi.ContentReceipt)
	ContentCredential        ContentType = ContentType(simffi.ContentCredential)
	ContentIntroduction      ContentType = ContentType(simffi.ContentIntroduction)
	ContentDiscoveryRequest  ContentType = ContentType(simffi.ContentDiscoveryRequest)
	ContentDiscoveryResponse ContentType = ContentType(simffi.ContentDiscoveryResponse)
	ContentExchangeRequest   ContentType = ContentType(simffi.ContentExchangeRequest)
	ContentExchangeResponse  ContentType = ContentType(simffi.ContentExchangeResponse)
)

// Match selects which incoming messages a rule applies to. Build one with
// MatchAny, MatchContentType, or MatchRequestID.
type Match struct {
	kind        simffi.MatchKind
	contentType ContentType
	requestID   []byte
}

// MatchAny matches every incoming message.
func MatchAny() Match { return Match{kind: simffi.MatchAny} }

// MatchContentType matches messages of a given content type.
func MatchContentType(ct ContentType) Match {
	return Match{kind: simffi.MatchContentType, contentType: ct}
}

// MatchRequestID matches the message carrying a specific request id.
func MatchRequestID(id []byte) Match {
	return Match{kind: simffi.MatchRequestID, requestID: id}
}

// Behaviour is how the device reacts to a matched message. Build one with
// Accept, Reject, or Ignore, optionally deferred with After.
type Behaviour struct {
	action simffi.Behaviour
	delay  time.Duration
}

// Accept drives the matched workflow to completion with simulated user consent.
func Accept() Behaviour { return Behaviour{action: simffi.BehaveAccept} }

// Reject responds to the matched workflow with a rejection.
func Reject() Behaviour { return Behaviour{action: simffi.BehaveReject} }

// Ignore drops the matched message without responding.
func Ignore() Behaviour { return Behaviour{action: simffi.BehaveIgnore} }

// After defers the behaviour by d before it is applied.
func (b Behaviour) After(d time.Duration) Behaviour {
	b.delay = d
	return b
}

// Intercepted is a message Intercept diverted to the caller. The device has
// taken no action on it.
type Intercepted struct {
	From    *signing.PublicKey
	To      *signing.PublicKey
	Content *message.Content

	device *Device
}

// Respond answers the intercepted request with outcomes, correlating the reply
// to it: the response goes back to its sender, carries its id in ResponseTo,
// and each outcome is tagged with the action it answers.
func (i *Intercepted) Respond(outcomes ...*message.Outcome) error {
	response := message.NewExchangeResponse().
		ID(i.Content.ID()).
		ResponseTo(i.Content.ID()).
		Status(message.StatusCreated)

	for _, outcome := range outcomes {
		response.Outcome(outcome)
	}

	content, err := response.Finish()
	if err != nil {
		return err
	}

	return i.device.Send(i.From, content)
}

// Actions returns the actions of the intercepted exchange request, so a caller
// can read what is being asked and tag its outcomes with the matching id.
func (i *Intercepted) Actions() ([]*message.Action, error) {
	request, err := message.ExchangeRequestDecode(i.Content)
	if err != nil {
		return nil, err
	}

	return request.Actions()
}

func (c ContentType) message() (message.ContentType, error) {
	switch c {
	case ContentCustom:
		return message.ContentCustom, nil
	case ContentChat:
		return message.ContentChat, nil
	case ContentReceipt:
		return message.ContentReceipt, nil
	case ContentCredential:
		return message.ContentCredential, nil
	case ContentIntroduction:
		return message.ContentIntroduction, nil
	case ContentDiscoveryRequest:
		return message.ContentDiscoveryRequest, nil
	case ContentDiscoveryResponse:
		return message.ContentDiscoveryResponse, nil
	case ContentExchangeRequest:
		return message.ContentExchangeRequest, nil
	case ContentExchangeResponse:
		return message.ContentExchangeResponse, nil
	default:
		return message.ContentUnknown, fmt.Errorf("simulator: cannot decode content type %d", c)
	}
}

// LogLevel selects the native log verbosity for a device's account. Values
// mirror the native zktf_log_level (1..5); the zero value is treated as
// LogError (quiet) by the native layer.
type LogLevel uint32

const (
	LogError LogLevel = iota + 1
	LogWarn
	LogInfo
	LogDebug
	LogTrace
)

// Option configures a device at construction time.
type Option func(*options)

type options struct {
	logLevel LogLevel
}

// WithLogLevel sets the native log verbosity for the device's account. Logs are
// written to stderr, prefixed with the account id, so multiple devices in one
// process can be told apart. Defaults to LogError (effectively quiet).
func WithLogLevel(level LogLevel) Option {
	return func(o *options) { o.logLevel = level }
}

func collectOptions(opts []Option) options {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// Device is a simulated mobile client. Register rules with Expect before
// driving a workflow; the device auto-responds to matching messages on a
// background thread.
type Device struct {
	h *simffi.Device
}

// NewDevice creates a device with its own account connected to the network.
func NewDevice(network *Network, opts ...Option) *Device {
	o := collectOptions(opts)
	return &Device{h: simffi.NewDevice(network.h, uint32(o.logLevel))}
}

// AttachDevice creates a device attached to a real, test-deployed backend at
// the given endpoints. The device always uses test trust anchors, so it can
// only interoperate with test networks — there is no way to target production.
func AttachDevice(rpcEndpoint, objectEndpoint, messagingEndpoint string, opts ...Option) *Device {
	o := collectOptions(opts)
	return &Device{h: simffi.DeviceAttach(rpcEndpoint, objectEndpoint, messagingEndpoint, uint32(o.logLevel))}
}

// Expect registers an auto-response rule. Rules are evaluated in registration
// order; the first match wins.
func (d *Device) Expect(m Match, b Behaviour) {
	d.h.Expect(
		m.kind,
		simffi.ContentType(m.contentType),
		m.requestID,
		b.action,
		uint64(b.delay/time.Millisecond),
	)
}

// Address returns the device's zktf address.
func (d *Device) Address() (*signing.PublicKey, error) {
	b, err := d.h.Address()
	if err != nil {
		return nil, err
	}
	return signing.FromBytes(b)
}

// Inbox returns the device's messaging inbox public key.
func (d *Device) Inbox() (*signing.PublicKey, error) {
	b, err := d.h.Inbox()
	if err != nil {
		return nil, err
	}
	return signing.FromBytes(b)
}

// Register drives the registration workflow against counterparty. It blocks
// until the workflow completes, is rejected, or times out.
func (d *Device) Register(counterparty *signing.PublicKey) error {
	return d.h.Register(counterparty.Bytes())
}

// Connect drives the pairwise connect workflow against counterparty. It blocks
// until the workflow completes, is rejected, or times out.
func (d *Device) Connect(counterparty *signing.PublicKey) error {
	return d.h.Connect(counterparty.Bytes())
}

// Send delivers content to an address, as the application layer would after
// deciding how to answer a request the device left alone.
func (d *Device) Send(to *signing.PublicKey, content *message.Content) error {
	encoded, err := content.Encode()
	if err != nil {
		return err
	}

	return d.h.Send(to.Bytes(), simffi.ContentType(content.Type()), encoded)
}

// Scan consumes an anonymous message, such as the discovery QR a portal shows
// to start a login.
func (d *Device) Scan(anonymousMessage []byte) error {
	return d.h.Scan(anonymousMessage)
}

// InterceptedFuture is a pending diverted message. Resolve it with Wait, or
// discard it with Cancel; either consumes the handle.
type InterceptedFuture struct {
	f      *simffi.InterceptedFuture
	device *Device
}

// Intercept diverts the message carrying requestID to the caller instead of
// driving a workflow for it, and returns the handle it arrives on. Arm it
// before the request is sent, so nothing is missed between arrival and the
// wait. The device does nothing further, so the caller answers it the way the
// host application would.
func (d *Device) Intercept(requestID []byte) *InterceptedFuture {
	f := d.h.Intercept(requestID)
	if f == nil {
		return nil
	}
	return &InterceptedFuture{f: f, device: d}
}

// Wait blocks until the message is diverted or timeout elapses, returning nil
// on timeout. A zero timeout waits indefinitely. Consumes the handle.
func (i *InterceptedFuture) Wait(timeout time.Duration) (*Intercepted, error) {
	raw, err := i.f.Wait(uint64(timeout.Milliseconds()))
	if err != nil || raw == nil {
		return nil, err
	}

	from, err := signing.FromBytes(raw.From)
	if err != nil {
		return nil, err
	}

	to, err := signing.FromBytes(raw.To)
	if err != nil {
		return nil, err
	}

	contentType, err := ContentType(raw.ContentType).message()
	if err != nil {
		return nil, err
	}

	content, err := message.ContentDecode(contentType, raw.Content)
	if err != nil {
		return nil, err
	}

	return &Intercepted{
		From:    from,
		To:      to,
		Content: content,
		device:  i.device,
	}, nil
}

// Cancel discards the handle without taking a message.
func (i *InterceptedFuture) Cancel() { i.f.Cancel() }

// SigningKeyCreate mints a signing keypair the device retains and returns its
// address. Reserving it first is what lets a credential issued in the same
// batch as MintControllerIdentity name the identity as its issuer.
func (d *Device) SigningKeyCreate() (*signing.PublicKey, error) {
	address, err := d.h.SigningKeyCreate()
	if err != nil {
		return nil, err
	}
	return signing.FromBytes(address)
}

// MintControllerIdentity mints a free-standing anchored identity for identifier
// and signs credential as that identity, in one liveness-authorized operation.
// credential is the unsigned credential as JSON; the signed one is returned.
func (d *Device) MintControllerIdentity(identifier *signing.PublicKey, credential []byte) ([]byte, error) {
	return d.h.MintControllerIdentity(identifier.Bytes(), credential)
}

// Close destroys the device's native account
func (d *Device) Close() {
	d.h.Close()
}
