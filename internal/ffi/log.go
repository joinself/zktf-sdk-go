package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

import (
	"sync"
	"time"
)

// LogField is a single structured key/value pair attached to a log entry.
type LogField struct {
	Key   string
	Value string
}

// LogEntry is a structured log record emitted by the native library. It is a
// pure-Go snapshot: all strings are copied out before the native entry is freed.
type LogEntry struct {
	Level     LogLevel
	AccountID string
	Target    string
	Message   string
	Timestamp time.Time
	Fields    []LogField
}

// LogHandler receives log entries from the native library.
type LogHandler func(LogEntry)

// The native log callback carries no user_data, so the handler is process-global.
var (
	logMu      sync.RWMutex
	logHandler LogHandler
)

// SetLogHandler registers the process-global log handler. A nil handler is a
// safe no-op.
func SetLogHandler(h LogHandler) {
	logMu.Lock()
	logHandler = h
	logMu.Unlock()
}

func dispatchLog(e LogEntry) {
	logMu.RLock()
	h := logHandler
	logMu.RUnlock()
	if h != nil {
		h(e)
	}
}

// newLogEntry copies every accessor off the native entry into a Go LogEntry.
// The caller frees the native entry; the returned value holds no C pointers.
func newLogEntry(ptr *C.zktf_log_entry) LogEntry {
	e := LogEntry{
		Level:     LogLevel(C.zktf_log_entry_level(ptr)),
		Target:    C.GoString(C.zktf_log_entry_target(ptr)),
		Message:   C.GoString(C.zktf_log_entry_message(ptr)),
		Timestamp: time.UnixMilli(int64(C.zktf_log_entry_timestamp_ms(ptr))),
	}

	if id := C.zktf_log_entry_account_id(ptr); id != nil {
		e.AccountID = C.GoString(id)
	}

	if n := int(C.zktf_log_entry_fields_count(ptr)); n > 0 {
		e.Fields = make([]LogField, n)
		for i := 0; i < n; i++ {
			e.Fields[i] = LogField{
				Key:   C.GoString(C.zktf_log_entry_field_key(ptr, C.size_t(i))),
				Value: C.GoString(C.zktf_log_entry_field_value(ptr, C.size_t(i))),
			}
		}
	}

	return e
}
