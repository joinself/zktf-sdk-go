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

func credentialValueToGo(v *C.zktf_credential_value) any {
	if v == nil {
		return nil
	}
	defer C.zktf_credential_value_destroy(v)

	switch C.zktf_credential_value_type_of(v) {
	case C.CREDENTIAL_VALUE_BYTES:
		return goBytesFromBuffer(C.zktf_credential_value_as_bytes(v))
	case C.CREDENTIAL_VALUE_STRING:
		return goStringFromBuffer(C.zktf_credential_value_as_string(v))
	case C.CREDENTIAL_VALUE_BOOL:
		return bool(C.zktf_credential_value_as_bool(v))
	case C.CREDENTIAL_VALUE_INTEGER:
		return int64(C.zktf_credential_value_as_integer(v))
	case C.CREDENTIAL_VALUE_UNSIGNED:
		return uint64(C.zktf_credential_value_as_unsigned(v))
	case C.CREDENTIAL_VALUE_FLOAT:
		return float64(C.zktf_credential_value_as_float(v))
	case C.CREDENTIAL_VALUE_ARRAY:
		n := int(C.zktf_credential_value_len(v))
		out := make([]any, n)
		for i := 0; i < n; i++ {
			out[i] = credentialValueToGo(C.zktf_credential_value_array_at(v, C.uintptr_t(i)))
		}
		return out
	case C.CREDENTIAL_VALUE_OBJECT:
		n := int(C.zktf_credential_value_len(v))
		out := make(map[string]any, n)
		for i := 0; i < n; i++ {
			key := goStringFromBuffer(C.zktf_credential_value_object_key_at(v, C.uintptr_t(i)))
			out[key] = credentialValueToGo(C.zktf_credential_value_object_value_at(v, C.uintptr_t(i)))
		}
		return out
	default:
		return nil
	}
}

func (c *VerifiableCredential) SubjectValue(path string) (any, bool) {
	defer runtime.KeepAlive(c)
	cpath := cstring(path)
	defer free(unsafe.Pointer(cpath))

	v := C.zktf_verifiable_credential_credential_subject_value(c.ptr, cpath)
	if v == nil {
		return nil, false
	}
	return credentialValueToGo(v), true
}
