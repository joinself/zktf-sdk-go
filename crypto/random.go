package crypto

import "github.com/joinself/zktf-sdk-go/internal/ffi"

func RandomID() []byte { return ffi.RandomID() }
