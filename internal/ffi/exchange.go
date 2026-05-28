package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

const (
	exchangePublicKeyEncodedLen = 66
	exchangePublicKeyBytesLen   = 33
)

// ExchangePublicKey wraps a zktf_exchange_public_key handle.
type ExchangePublicKey struct {
	ptr *C.zktf_exchange_public_key
}

func newExchangePublicKey(ptr *C.zktf_exchange_public_key) *ExchangePublicKey {
	if ptr == nil {
		return nil
	}
	k := &ExchangePublicKey{ptr: ptr}
	runtime.AddCleanup(k, func(ptr *C.zktf_exchange_public_key) {
		C.zktf_exchange_public_key_destroy(ptr)
	}, k.ptr)
	return k
}

// ExchangePublicKeyFromAddress decodes a hex address into an exchange key.
func ExchangePublicKeyFromAddress(hex string) (*ExchangePublicKey, error) {
	buf, length := cbytes([]byte(hex))
	defer free(unsafe.Pointer(buf))

	var ptr *C.zktf_exchange_public_key
	if err := status(C.zktf_exchange_public_key_decode(&ptr, buf, length)); err != nil {
		return nil, err
	}
	return newExchangePublicKey(ptr), nil
}

// ExchangePublicKeyFromBytes constructs an exchange key from raw bytes.
func ExchangePublicKeyFromBytes(data []byte) (*ExchangePublicKey, error) {
	buf, length := cbytes(data)
	defer free(unsafe.Pointer(buf))

	var ptr *C.zktf_exchange_public_key
	if err := status(C.zktf_exchange_public_key_from_bytes(&ptr, buf, length)); err != nil {
		return nil, err
	}
	return newExchangePublicKey(ptr), nil
}

// String returns the hex encoded address.
func (k *ExchangePublicKey) String() string {
	buf := C.malloc(exchangePublicKeyEncodedLen)
	defer C.free(buf)
	if err := status(C.zktf_exchange_public_key_encode(
		k.ptr, (*C.uint8_t)(buf), exchangePublicKeyEncodedLen,
	)); err != nil {
		return ""
	}
	return string(C.GoBytes(buf, exchangePublicKeyEncodedLen))
}

// Bytes returns the raw bytes of the key.
func (k *ExchangePublicKey) Bytes() []byte {
	buf := C.malloc(exchangePublicKeyBytesLen)
	defer C.free(buf)
	if err := status(C.zktf_exchange_public_key_as_bytes(
		k.ptr, (*C.uint8_t)(buf), exchangePublicKeyBytesLen,
	)); err != nil {
		return nil
	}
	return C.GoBytes(buf, exchangePublicKeyBytesLen)
}
