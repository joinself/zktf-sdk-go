package status_test

import (
	"testing"

	"github.com/joinself/zktf-sdk-go/status"
)

func TestDeviceKeyLostIsDistinct(t *testing.T) {
	for _, c := range []status.Code{status.OK, status.Unknown, status.AccountNotConfigured} {
		if status.DeviceKeyLost == c {
			t.Fatalf("DeviceKeyLost = %d collides with %d", status.DeviceKeyLost, c)
		}
	}
}
