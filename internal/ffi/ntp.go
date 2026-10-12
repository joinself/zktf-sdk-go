package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

func NTPNow() int64 {
	return int64(C.zktf_ntp_now())
}

func NTPEnable() {
	C.zktf_ntp_enable()
}

func NTPDisable() {
	C.zktf_ntp_disable()
}
