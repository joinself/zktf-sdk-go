package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import "unsafe"

// ValueKeys lists stored value keys, optionally filtered by prefix. An empty
// prefix lists every key.
func (a *Account) ValueKeys(prefix string) ([]string, error) {
	var p *C.char
	if prefix != "" {
		p = cstring(prefix)
		defer free(unsafe.Pointer(p))
	}

	var c *C.zktf_collection_value_key
	if err := status(C.zktf_account_value_keys(a.ptr, p, &c)); err != nil {
		return nil, err
	}

	return valueKeysFrom(c), nil
}

// ValueLookup returns the value stored under key. The boolean is false when no
// value is stored for that key.
func (a *Account) ValueLookup(key string) ([]byte, bool, error) {
	k := cstring(key)
	defer free(unsafe.Pointer(k))

	var buf *C.zktf_bytes_buffer
	code := C.zktf_account_value_lookup(a.ptr, k, &buf)
	if code == C.STATUS_VALUE_NOT_FOUND {
		return nil, false, nil
	}
	if err := status(code); err != nil {
		return nil, false, err
	}

	return goBytesFromBuffer(buf), true, nil
}

// ValueStore stores a key/value pair. A zero expiresUnix means the value never
// expires; otherwise it is removed at that absolute unix timestamp (seconds).
func (a *Account) ValueStore(key string, value []byte, expiresUnix int64) error {
	k := cstring(key)
	defer free(unsafe.Pointer(k))

	valBuf, valLen := cbytes(value)
	defer free(unsafe.Pointer(valBuf))

	var options *C.zktf_value_store_options
	if expiresUnix != 0 {
		options = C.zktf_value_store_options_init()
		C.zktf_value_store_options_with_expiry(options, C.int64_t(expiresUnix))
		defer C.zktf_value_store_options_destroy(options)
	}

	return status(C.zktf_account_value_store(a.ptr, k, valBuf, valLen, options))
}

// ValueRemove deletes the value stored under key.
func (a *Account) ValueRemove(key string) error {
	k := cstring(key)
	defer free(unsafe.Pointer(k))

	return status(C.zktf_account_value_remove(a.ptr, k))
}

// valueKeysFrom copies a caller-owned zktf_collection_value_key into Go strings
// and destroys the collection.
func valueKeysFrom(c *C.zktf_collection_value_key) []string {
	if c == nil {
		return nil
	}
	defer C.zktf_collection_value_key_destroy(c)

	n := int(C.zktf_collection_value_key_len(c))
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = C.GoString(C.zktf_collection_value_key_at(c, C.size_t(i)))
	}

	return out
}
