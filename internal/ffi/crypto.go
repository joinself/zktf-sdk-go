package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

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
	defer runtime.KeepAlive(k)
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

const (
	sha3256Len     = 32
	randomIDLen    = 20
	sealedBoxExtra = 48
)

func SHA3256(data []byte) []byte {
	buf, length := cbytes(data)
	defer free(unsafe.Pointer(buf))
	out := make([]byte, sha3256Len)
	C.zktf_crypto_sha3_256(buf, length, (*C.uint8_t)(unsafe.Pointer(&out[0])))
	return out
}

func RandomID() []byte {
	out := make([]byte, randomIDLen)
	C.zktf_crypto_random_id((*C.uint8_t)(unsafe.Pointer(&out[0])))
	return out
}

func RandomBytes(n int) []byte {
	if n <= 0 {
		return nil
	}
	out := make([]byte, n)
	C.zktf_crypto_random_bytes((*C.uint8_t)(unsafe.Pointer(&out[0])), C.size_t(n))
	return out
}

func FileKeyLength() int {
	return int(C.zktf_crypto_aes_gcm_encrypt_file_key_length())
}

func EncryptFile(source, destination string) ([]byte, error) {
	src, dst := cstring(source), cstring(destination)
	defer free(unsafe.Pointer(src))
	defer free(unsafe.Pointer(dst))
	key := make([]byte, FileKeyLength())
	if err := status(C.zktf_crypto_aes_gcm_encrypt_file(src, dst, (*C.uint8_t)(unsafe.Pointer(&key[0])))); err != nil {
		return nil, err
	}
	return key, nil
}

func DecryptFile(source, destination string, key []byte) error {
	if len(key) != FileKeyLength() {
		return status(C.STATUS_INPUT_BUFFER_INSUFFICIENT)
	}
	src, dst := cstring(source), cstring(destination)
	defer free(unsafe.Pointer(src))
	defer free(unsafe.Pointer(dst))
	k, _ := cbytes(key)
	defer free(unsafe.Pointer(k))
	return status(C.zktf_crypto_aes_gcm_decrypt_file(src, dst, k))
}
