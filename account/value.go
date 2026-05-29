package account

import "time"

// ValueKeys lists the keys of values stored on this account. An empty prefix
// lists every key; otherwise only keys with the given prefix are returned.
func (a *Account) ValueKeys(prefix string) ([]string, error) {
	return a.h.ValueKeys(prefix)
}

// ValueLookup returns the value stored under key. The boolean is false when no
// value is stored for that key.
func (a *Account) ValueLookup(key string) ([]byte, bool, error) {
	return a.h.ValueLookup(key)
}

// ValueStore stores a key/value pair. A zero expires means the value never
// expires; otherwise it is removed at that time.
func (a *Account) ValueStore(key string, value []byte, expires time.Time) error {
	var unix int64
	if !expires.IsZero() {
		unix = expires.Unix()
	}

	return a.h.ValueStore(key, value, unix)
}

// ValueRemove deletes the value stored under key.
func (a *Account) ValueRemove(key string) error {
	return a.h.ValueRemove(key)
}
