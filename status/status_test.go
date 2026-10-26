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

func TestCodeDomainAndNumber(t *testing.T) {
	cases := []struct {
		code   status.Code
		domain status.Domain
		number uint16
	}{
		{status.OK, 0, 0},
		{status.Unknown, 0, 1},
		{status.AccountCallbacksRequired, 0, 2},
		{status.AccountConfigRequired, 0, 3},
		{status.AccountAlreadyConfigured, status.DomainAccount, 1},
		{status.AccountNotConfigured, status.DomainAccount, 2},
		{status.DeviceKeyLost, status.DomainDevice, 1},
		{status.EventPollUnavailable, status.DomainEvent, 2},
		{status.KeyPairNotFound, status.DomainKeyPair, 5},
	}
	for _, tc := range cases {
		if got := tc.code.Domain(); got != tc.domain {
			t.Errorf("Code(%#x).Domain() = %d, want %d", uint32(tc.code), got, tc.domain)
		}
		if got := tc.code.Number(); got != tc.number {
			t.Errorf("Code(%#x).Number() = %d, want %d", uint32(tc.code), got, tc.number)
		}
		if want := status.Code(uint32(tc.domain)<<16 | uint32(tc.number)); tc.code != want {
			t.Errorf("Code(%#x) != domain %d << 16 | %d", uint32(tc.code), tc.domain, tc.number)
		}
	}
}
