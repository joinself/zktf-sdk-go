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

type Domain uint16

const (
	DomainAccount Domain = Domain(ffi.DomainAccount)
	DomainDevice  Domain = Domain(ffi.DomainDevice)
	DomainEvent   Domain = Domain(ffi.DomainEvent)
	DomainKeyPair Domain = Domain(ffi.DomainKeyPair)
	DomainMessage Domain = Domain(ffi.DomainMessage)
	DomainPairing Domain = Domain(ffi.DomainPairing)
	DomainRpc     Domain = Domain(ffi.DomainRpc)
)

// A subset of the status codes. Callers compare these against the value
// returned by Of. The full set is defined by the native library.
const (
	OK                       Code = Code(ffi.StatusOK)
	Unknown                  Code = Code(ffi.StatusUnknown)
	AccountAlreadyConfigured Code = Code(ffi.StatusAccountAlreadyConfigured)
	AccountCallbacksRequired Code = Code(ffi.StatusAccountCallbacksRequired)
	AccountConfigRequired    Code = Code(ffi.StatusAccountConfigRequired)
	AccountNotConfigured     Code = Code(ffi.StatusAccountNotConfigured)
	DeviceKeyLost            Code = Code(ffi.StatusDeviceKeyLost)
	EventPollUnavailable     Code = Code(ffi.StatusEventPollUnavailable)
	KeyPairNotFound          Code = Code(ffi.StatusKeyPairNotFound)
	PairingCodeInvalid       Code = Code(ffi.StatusPairingCodeInvalid)
	PairingGrantInvalid      Code = Code(ffi.StatusPairingGrantInvalid)
	PairingRolesInvalid      Code = Code(ffi.StatusPairingRolesInvalid)
	MessageIdInvalid         Code = Code(ffi.StatusMessageIdInvalid)
	RpcAborted               Code = Code(ffi.StatusRpcAborted)
	RpcBadGateway            Code = Code(ffi.StatusRpcBadGateway)
	RpcBadRequest            Code = Code(ffi.StatusRpcBadRequest)
	RpcConnectionFailed      Code = Code(ffi.StatusRpcConnectionFailed)
	RpcConnectionTimeout     Code = Code(ffi.StatusRpcConnectionTimeout)
	RpcForbidden             Code = Code(ffi.StatusRpcForbidden)
	RpcGatewayTimeout        Code = Code(ffi.StatusRpcGatewayTimeout)
	RpcGone                  Code = Code(ffi.StatusRpcGone)
	RpcInternalServerError   Code = Code(ffi.StatusRpcInternalServerError)
	RpcNotFound              Code = Code(ffi.StatusRpcNotFound)
	RpcRequestTimeout        Code = Code(ffi.StatusRpcRequestTimeout)
	RpcResourceExhausted     Code = Code(ffi.StatusRpcResourceExhausted)
	RpcServiceUnavailable    Code = Code(ffi.StatusRpcServiceUnavailable)
	RpcUnauthorized          Code = Code(ffi.StatusRpcUnauthorized)
)

func (c Code) Domain() Domain {
	return Domain(ffi.StatusDomain(uint32(c)))
}

func (c Code) Number() uint16 {
	return ffi.StatusNumber(uint32(c))
}

// Transient reports whether the code is a network or server availability
// failure that may succeed when retried.
func (c Code) Transient() bool {
	switch c {
	case RpcAborted, RpcBadGateway, RpcConnectionFailed, RpcConnectionTimeout,
		RpcGatewayTimeout, RpcRequestTimeout, RpcResourceExhausted, RpcServiceUnavailable:
		return true
	}
	return false
}

// Transient reports whether err carries a transient status.
func Transient(err error) bool {
	c, ok := Of(err)
	return ok && c.Transient()
}

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
