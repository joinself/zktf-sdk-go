// Package ffi is the single cgo boundary for the zktf Go SDK.
//
// It is the ONLY package in this module that does `import "C"`. Every other
// package is pure Go and talks to the native zktf library exclusively through
// the Go types and functions exported here. Keeping all cgo in one package
// avoids two problems that the previous-generation self-go-sdk suffered from:
//
//  1. cgo gives every package that imports "C" its own *distinct* set of C
//     types, so sharing a C pointer across packages required `//go:linkname`.
//     With a single cgo package there is no cross-package C type and no linkname
//     is ever needed.
//
//  2. The C types never escape: each wrapper struct holds its `*C.zktf_*`
//     pointer in an UNEXPORTED field, so no exported signature — here or in any
//     public package — ever mentions a C type.
//
// Build prerequisites: the native header `zktf-sdk.h` must be on the C include
// path and `libzktf_sdk` on the linker path. For local development point cgo at
// the zktf-sdk checkout, e.g.:
//
//	CGO_CFLAGS=-I/path/to/zktf-sdk/crates/zktf-ffi \
//	CGO_LDFLAGS=-L/path/to/zktf-sdk/target/debug \
//	LD_LIBRARY_PATH=/path/to/zktf-sdk/target/debug \
//	go build ./...
package ffi

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo linux LDFLAGS: -lzktf_sdk
#cgo darwin LDFLAGS: -lzktf_sdk -framework CoreFoundation -framework SystemConfiguration -framework Security
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import "unsafe"

// cbytes copies a Go byte slice into C-allocated memory. The returned pointer
// must be released with free. A nil/empty slice yields a nil pointer.
func cbytes(b []byte) (*C.uint8_t, C.size_t) {
	if len(b) == 0 {
		return nil, 0
	}
	return (*C.uint8_t)(C.CBytes(b)), C.size_t(len(b))
}

// cstring copies a Go string into a C-allocated NUL-terminated string. The
// returned pointer must be released with free.
func cstring(s string) *C.char {
	return C.CString(s)
}

// free releases memory allocated by cbytes / cstring.
func free(p unsafe.Pointer) {
	C.free(p)
}
