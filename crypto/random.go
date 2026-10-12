package crypto

import "github.com/joinself/zktf-sdk-go/internal/ffi"

func SHA3256(data []byte) []byte { return ffi.SHA3256(data) }

func RandomID() []byte { return ffi.RandomID() }

func RandomBytes(n int) []byte { return ffi.RandomBytes(n) }

func FileKeyLength() int { return ffi.FileKeyLength() }

func EncryptFile(source, destination string) ([]byte, error) {
	return ffi.EncryptFile(source, destination)
}

func DecryptFile(source, destination string, key []byte) error {
	return ffi.DecryptFile(source, destination, key)
}
