package account

import (
	"time"

	"github.com/joinself/zktf-sdk-go/internal/ffi"
)

func NTPNow() time.Time { return time.Unix(ffi.NTPNow(), 0) }

func EnableNTP() { ffi.NTPEnable() }

func DisableNTP() { ffi.NTPDisable() }
