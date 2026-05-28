package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

// Status wraps a non-zero zktf_status code returned by the native library and
// implements the error interface. Public packages return it as a plain `error`,
// so the underlying C enum never appears in any exported signature.
type Status struct {
	code uint32
}

// Code returns the raw zktf_status code. Public packages may compare it against
// the constants re-declared in the public status package.
func (s *Status) Code() uint32 {
	return s.code
}

// Error returns the human readable message for the status code, as provided by
// zktf_status_message.
func (s *Status) Error() string {
	return C.GoString(C.zktf_status_message(C.enum_zktf_status(s.code)))
}

// status converts a raw zktf_status result into a Go error. A zero code (success)
// returns nil; any other code returns a *Status.
func status(code C.enum_zktf_status) error {
	if code == 0 {
		return nil
	}
	return &Status{code: uint32(code)}
}
