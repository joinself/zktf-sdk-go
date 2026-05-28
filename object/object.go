// Package object provides encrypted objects (message attachments / blobs).
package object

import "github.com/joinself/zktf-sdk-go/internal/ffi"

// Object is an encrypted object that can be attached to messages.
type Object struct {
	h *ffi.Object
}

func init() {
	ffi.ObjectOf = func(o any) *ffi.Object {
		return o.(*Object).h
	}
	ffi.ToObject = func(h *ffi.Object) any {
		return &Object{h: h}
	}
}

// New builds an object from raw data and a mime type.
func New(mime string, data []byte) (*Object, error) {
	o, err := ffi.ObjectCreate(mime, data)
	if err != nil {
		return nil, err
	}

	return &Object{h: o}, nil
}

// ID returns the hash of the encrypted data, or nil if the object has not yet
// been uploaded.
func (o *Object) ID() []byte { return o.h.ID() }

// Hash returns the hash of the unencrypted data, or nil if unavailable.
func (o *Object) Hash() []byte { return o.h.Hash() }

// MimeType returns the object's mime type.
func (o *Object) MimeType() string { return o.h.MimeType() }

// Key returns the object's encryption key, or nil if not present.
func (o *Object) Key() []byte { return o.h.Key() }

// Data returns the object's data buffer.
func (o *Object) Data() []byte { return o.h.Data() }
