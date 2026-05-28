package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

import "runtime"

// CryptoKeyPackage wraps a zktf_crypto_key_package handle (an MLS key package
// used to establish an encrypted session out of band, e.g. via discovery).
type CryptoKeyPackage struct {
	ptr *C.zktf_crypto_key_package
}

func newCryptoKeyPackage(ptr *C.zktf_crypto_key_package) *CryptoKeyPackage {
	if ptr == nil {
		return nil
	}
	k := &CryptoKeyPackage{ptr: ptr}
	runtime.AddCleanup(k, func(ptr *C.zktf_crypto_key_package) {
		C.zktf_crypto_key_package_destroy(ptr)
	}, k.ptr)
	return k
}

// FromAddress returns the signing address the key package is for.
func (k *CryptoKeyPackage) FromAddress() *SigningPublicKey {
	return newSigningPublicKey(C.zktf_crypto_key_package_from_address(k.ptr))
}

// CryptoWelcome wraps a zktf_crypto_welcome handle (an MLS welcome message).
type CryptoWelcome struct {
	ptr *C.zktf_crypto_welcome
}

func newCryptoWelcome(ptr *C.zktf_crypto_welcome) *CryptoWelcome {
	if ptr == nil {
		return nil
	}
	w := &CryptoWelcome{ptr: ptr}
	runtime.AddCleanup(w, func(ptr *C.zktf_crypto_welcome) {
		C.zktf_crypto_welcome_destroy(ptr)
	}, w.ptr)
	return w
}
