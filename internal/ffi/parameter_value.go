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

// ParameterValue is a typed value carried by a verification parameter. Values
// map to and from native Go types via NewParameterValue and Value.
type ParameterValue struct {
	ptr *C.zktf_message_content_parameter_value
}

func newParameterValue(ptr *C.zktf_message_content_parameter_value) *ParameterValue {
	if ptr == nil {
		return nil
	}
	v := &ParameterValue{ptr: ptr}
	runtime.AddCleanup(v, func(ptr *C.zktf_message_content_parameter_value) {
		C.zktf_message_content_parameter_value_destroy(ptr)
	}, v.ptr)
	return v
}

// NewParameterValue builds a parameter value from a supported Go type. Accepted
// types are []byte, string, bool, the signed integer types (int, int8, int16,
// int32, int64), the unsigned integer types (uint, uint8, uint16, uint32,
// uint64), float32, float64, [][]byte and []string. Any other type yields nil.
func NewParameterValue(v any) *ParameterValue {
	var ptr *C.zktf_message_content_parameter_value

	switch val := v.(type) {
	case []byte:
		// from_raw_parts on the Rust side requires a non-null pointer even for a
		// zero-length slice, so allocate via CBytes rather than the cbytes helper.
		buf := C.CBytes(val)
		ptr = C.zktf_message_content_parameter_value_bytes_create((*C.uint8_t)(buf), C.uintptr_t(len(val)))
		C.free(buf)
	case string:
		cs := cstring(val)
		ptr = C.zktf_message_content_parameter_value_string_create(cs)
		free(unsafe.Pointer(cs))
	case bool:
		ptr = C.zktf_message_content_parameter_value_bool_create(C.bool(val))
	case int:
		ptr = C.zktf_message_content_parameter_value_integer_create(C.int64_t(val))
	case int8:
		ptr = C.zktf_message_content_parameter_value_integer_create(C.int64_t(val))
	case int16:
		ptr = C.zktf_message_content_parameter_value_integer_create(C.int64_t(val))
	case int32:
		ptr = C.zktf_message_content_parameter_value_integer_create(C.int64_t(val))
	case int64:
		ptr = C.zktf_message_content_parameter_value_integer_create(C.int64_t(val))
	case uint:
		ptr = C.zktf_message_content_parameter_value_unsigned_create(C.uint64_t(val))
	case uint8:
		ptr = C.zktf_message_content_parameter_value_unsigned_create(C.uint64_t(val))
	case uint16:
		ptr = C.zktf_message_content_parameter_value_unsigned_create(C.uint64_t(val))
	case uint32:
		ptr = C.zktf_message_content_parameter_value_unsigned_create(C.uint64_t(val))
	case uint64:
		ptr = C.zktf_message_content_parameter_value_unsigned_create(C.uint64_t(val))
	case float32:
		ptr = C.zktf_message_content_parameter_value_float_create(C.double(val))
	case float64:
		ptr = C.zktf_message_content_parameter_value_float_create(C.double(val))
	case [][]byte:
		col := C.zktf_collection_bytes_buffer_init()
		for _, b := range val {
			buf := C.CBytes(b)
			C.zktf_collection_bytes_buffer_append(col, (*C.uint8_t)(buf), C.size_t(len(b)))
			C.free(buf)
		}
		ptr = C.zktf_message_content_parameter_value_array_bytes_create(col)
		C.zktf_collection_bytes_buffer_destroy(col)
	case []string:
		col := C.zktf_collection_string_buffer_init()
		for _, s := range val {
			cs := cstring(s)
			C.zktf_collection_string_buffer_append(col, cs)
			free(unsafe.Pointer(cs))
		}
		ptr = C.zktf_message_content_parameter_value_array_string_create(col)
		C.zktf_collection_string_buffer_destroy(col)
	default:
		return nil
	}

	return newParameterValue(ptr)
}

// Value decodes the parameter value into a native Go type: []byte, string,
// bool, int64, uint64, float64, [][]byte or []string. Null and object values,
// for which the ABI exposes no accessor, decode to nil.
func (v *ParameterValue) Value() any {
	switch C.zktf_message_content_parameter_value_type_of(v.ptr) {
	case C.PARAMETER_VALUE_BYTES:
		return goBytesFromBuffer(C.zktf_message_content_parameter_value_as_bytes(v.ptr))
	case C.PARAMETER_VALUE_STRING:
		return goStringFromBuffer(C.zktf_message_content_parameter_value_as_string(v.ptr))
	case C.PARAMETER_VALUE_BOOL:
		return bool(C.zktf_message_content_parameter_value_as_bool(v.ptr))
	case C.PARAMETER_VALUE_INTEGER:
		return int64(C.zktf_message_content_parameter_value_as_integer(v.ptr))
	case C.PARAMETER_VALUE_UNSIGNED:
		return uint64(C.zktf_message_content_parameter_value_as_unsigned(v.ptr))
	case C.PARAMETER_VALUE_FLOAT:
		return float64(C.zktf_message_content_parameter_value_as_float(v.ptr))
	case C.PARAMETER_VALUE_ARRAY:
		// Arrays are homogeneous: a string array decodes only via as_array_string
		// and a byte array only via as_array_bytes (each returns null otherwise).
		if s := C.zktf_message_content_parameter_value_as_array_string(v.ptr); s != nil {
			return stringsFromBufferCollection(s)
		}
		if b := C.zktf_message_content_parameter_value_as_array_bytes(v.ptr); b != nil {
			return bytesFromBufferCollection(b)
		}
		return nil
	default:
		return nil
	}
}
