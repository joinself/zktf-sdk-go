package account

import "time"

// Option scopes — one marker interface per operation that takes options.
// `With*` constructors return concrete option values that implement every scope
// they apply to. Passing an option to the wrong operation is a compile error,
// not a silent no-op.

// CallOption configures any operation that just waits on the network. Every
// `With*` that semantically applies "to any networked call" implements this.
type CallOption interface {
	applyCallOption(*callOpts)
}

// InboxOpenOption configures Account.InboxOpen.
type InboxOpenOption interface {
	applyInboxOpenOption(*inboxOpenOpts)
}

// ObjectUploadOption configures Account.ObjectUpload.
type ObjectUploadOption interface {
	applyObjectUploadOption(*objectUploadOpts)
}

// Per-scope opts structs. Each operation maintains the fields it cares about
// and applies the options into it; options that don't satisfy an operation's
// scope interface cannot be passed in the first place.

type callOpts struct {
	timeout time.Duration
}

type inboxOpenOpts struct {
	timeout time.Duration
	expires time.Time
}

type objectUploadOpts struct {
	timeout              time.Duration
	objectPersistLocally bool
}

func collectCallOpts(options []CallOption) callOpts {
	var o callOpts

	for _, opt := range options {
		opt.applyCallOption(&o)
	}

	return o
}

func collectInboxOpenOpts(options []InboxOpenOption) inboxOpenOpts {
	var o inboxOpenOpts

	for _, opt := range options {
		opt.applyInboxOpenOption(&o)
	}

	return o
}

func collectObjectUploadOpts(options []ObjectUploadOption) objectUploadOpts {
	var o objectUploadOpts

	for _, opt := range options {
		opt.applyObjectUploadOption(&o)
	}

	return o
}

// WithTimeout overrides the network timeout for an operation. The default is
// ffi.DefaultTimeout (30s). Valid on any operation that takes options.
func WithTimeout(d time.Duration) timeoutOpt { return timeoutOpt{d: d} }

type timeoutOpt struct{ d time.Duration }

func (t timeoutOpt) applyCallOption(o *callOpts)                 { o.timeout = t.d }
func (t timeoutOpt) applyInboxOpenOption(o *inboxOpenOpts)       { o.timeout = t.d }
func (t timeoutOpt) applyObjectUploadOption(o *objectUploadOpts) { o.timeout = t.d }

// WithExpires sets the absolute expiry time for the inbox subscription opened
// by InboxOpen. Zero means no expiry.
func WithExpires(t time.Time) expiresOpt { return expiresOpt{t: t} }

type expiresOpt struct{ t time.Time }

func (e expiresOpt) applyInboxOpenOption(o *inboxOpenOpts) { o.expires = e.t }

// WithObjectPersistLocally controls whether ObjectUpload also writes the object
// to the local store.
func WithObjectPersistLocally(persist bool) persistLocallyOpt {
	return persistLocallyOpt{p: persist}
}

type persistLocallyOpt struct{ p bool }

func (p persistLocallyOpt) applyObjectUploadOption(o *objectUploadOpts) {
	o.objectPersistLocally = p.p
}
