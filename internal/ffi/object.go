package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import (
	"runtime"
	"time"
	"unsafe"
)

const (
	objectIDLen   = 32
	objectHashLen = 32
	objectKeyLen  = 44
)

// Object wraps a zktf_object handle (an encrypted attachment / blob).
type Object struct {
	ptr *C.zktf_object
}

func newObject(ptr *C.zktf_object) *Object {
	if ptr == nil {
		return nil
	}
	o := &Object{ptr: ptr}
	runtime.AddCleanup(o, func(ptr *C.zktf_object) {
		C.zktf_object_destroy(ptr)
	}, o.ptr)
	return o
}

// ObjectCreate builds an object from raw data and a mime type.
func ObjectCreate(mime string, data []byte) (*Object, error) {
	cmime := cstring(mime)
	defer free(unsafe.Pointer(cmime))
	buf, length := cbytes(data)
	defer free(unsafe.Pointer(buf))

	var ptr *C.zktf_object
	if err := status(C.zktf_object_create(&ptr, cmime, buf, C.uintptr_t(length))); err != nil {
		return nil, err
	}
	return newObject(ptr), nil
}

// ID returns the hash of the encrypted data, or nil if the object has not yet
// been uploaded (the id is only available once the encrypted data is hashed).
func (o *Object) ID() []byte {
	id := C.zktf_object_id(o.ptr)
	if id == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(id), objectIDLen)
}

// Hash returns the hash of the unencrypted data, or nil if unavailable.
func (o *Object) Hash() []byte {
	h := C.zktf_object_hash(o.ptr)
	if h == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(h), objectHashLen)
}

// MimeType returns the object's mime type.
func (o *Object) MimeType() string {
	return C.GoString(C.zktf_object_mime(o.ptr))
}

// Key returns the object's 44-byte encryption key, or nil if not present.
func (o *Object) Key() []byte {
	k := C.zktf_object_key(o.ptr)
	if k == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(k), objectKeyLen)
}

// Data returns the object's data buffer.
func (o *Object) Data() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.zktf_object_data_buf(o.ptr)),
		C.int(C.zktf_object_data_len(o.ptr)),
	)
}

// ObjectUploadOptions configures Account.UploadObject.
type ObjectUploadOptions struct {
	ptr *C.zktf_object_upload_options
}

// NewObjectUploadOptions initializes upload options.
func NewObjectUploadOptions() *ObjectUploadOptions {
	ptr := C.zktf_object_upload_options_init()
	o := &ObjectUploadOptions{ptr: ptr}
	runtime.AddCleanup(o, func(ptr *C.zktf_object_upload_options) {
		C.zktf_object_upload_options_destroy(ptr)
	}, o.ptr)
	return o
}

// PersistLocally controls whether the uploaded object is also written to the
// local object store.
func (o *ObjectUploadOptions) PersistLocally(persist bool) *ObjectUploadOptions {
	C.zktf_object_upload_options_persist_locally(o.ptr, C.bool(persist))
	return o
}

// ObjectUpload uploads an object to the object store via callback. Pass nil
// options for defaults.
func (a *Account) ObjectUpload(obj *Object, options *ObjectUploadOptions, timeout time.Duration) error {
	var optsPtr *C.zktf_object_upload_options
	if options != nil {
		optsPtr = options.ptr
	}

	fut := C.zktf_account_object_upload(a.ptr, obj.ptr, optsPtr)

	return AwaitStatus(fut, timeout)
}

// ObjectDownload downloads an object's encrypted bytes and key from the server,
// via callback.
func (a *Account) ObjectDownload(obj *Object, timeout time.Duration) error {
	fut := C.zktf_account_object_download(a.ptr, obj.ptr)

	return AwaitStatus(fut, timeout)
}
