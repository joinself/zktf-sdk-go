// Package status exposes the zktf status codes carried by errors returned from
// the SDK, without referencing any cgo types. Errors returned by the SDK are
// plain `error` values; use Code to recover the underlying status code.
package status

import (
	"errors"

	"github.com/joinself/zktf-sdk-go/internal/ffi"
)

// Code is a zktf status code.
type Code uint32

// A subset of the status codes. Callers compare these against the value
// returned by Of. The full set is defined by the native library.
const (
	OK                       Code = 0
	Unknown                  Code = 1
	AccountAlreadyConfigured Code = 2
	AccountCallbacksRequired Code = 3
	AccountConfigRequired    Code = 4
	AccountNotConfigured     Code = 5
)

// Of returns the status code carried by an SDK error and whether the error was
// in fact a status error. Errors that did not originate from a native status
// return false.
func Of(err error) (Code, bool) {
	var s *ffi.Status
	if errors.As(err, &s) {
		return Code(s.Code()), true
	}
	return 0, false
}
