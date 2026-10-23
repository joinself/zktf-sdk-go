package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

// Status wraps a non-zero zktf_status code returned by the native library and
// implements the error interface. Public packages return it as a plain `error`,
// so the underlying C type never appears in any exported signature.
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
	return C.GoString(C.zktf_status_message(C.zktf_status(s.code)))
}

// status converts a raw zktf_status result into a Go error. A zero code (success)
// returns nil; any other code returns a *Status.
func status(code C.zktf_status) error {
	if code == 0 {
		return nil
	}
	return &Status{code: uint32(code)}
}

const (
	DomainAccount uint16 = C.ZKTF_DOMAIN_ACCOUNT
	DomainDevice  uint16 = C.ZKTF_DOMAIN_DEVICE
	DomainEvent   uint16 = C.ZKTF_DOMAIN_EVENT
	DomainKeyPair uint16 = C.ZKTF_DOMAIN_KEY_PAIR
)

const (
	StatusOK                       uint32 = C.ZKTF_STATUS_OK
	StatusUnknown                  uint32 = C.ZKTF_STATUS_UNKNOWN
	StatusAccountCallbacksRequired uint32 = C.ZKTF_STATUS_ACCOUNT_CALLBACKS_REQUIRED
	StatusAccountConfigRequired    uint32 = C.ZKTF_STATUS_ACCOUNT_CONFIG_REQUIRED
	StatusAccountAlreadyConfigured uint32 = uint32(DomainAccount)<<16 | C.ZKTF_ACCOUNT_ALREADY_CONFIGURED
	StatusAccountNotConfigured     uint32 = uint32(DomainAccount)<<16 | C.ZKTF_ACCOUNT_NOT_CONFIGURED
	StatusDeviceKeyLost            uint32 = uint32(DomainDevice)<<16 | C.ZKTF_DEVICE_KEY_LOST
	StatusEventPollUnavailable     uint32 = uint32(DomainEvent)<<16 | C.ZKTF_EVENT_POLL_UNAVAILABLE
	StatusKeyPairNotFound          uint32 = uint32(DomainKeyPair)<<16 | C.ZKTF_KEY_PAIR_NOT_FOUND
)

func StatusDomain(code uint32) uint16 {
	return uint16(C.zktf_status_domain(C.zktf_status(code)))
}

func StatusNumber(code uint32) uint16 {
	return uint16(C.zktf_status_code(C.zktf_status(code)))
}
