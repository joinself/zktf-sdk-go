package status_test

import (
	"errors"
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
		{status.MessageIdInvalid, status.DomainMessage, 20},
		{status.RpcNotFound, status.DomainRpc, 19},
		{status.RpcServiceUnavailable, status.DomainRpc, 28},
		{status.ValueKeyInvalid, status.DomainValue, 1},
		{status.ValueBusy, status.DomainValue, 4},
		{status.ValueCursorInvalid, status.DomainValue, 8},
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

func TestTransient(t *testing.T) {
	for _, c := range []status.Code{status.RpcServiceUnavailable, status.RpcConnectionFailed, status.RpcGatewayTimeout} {
		if !c.Transient() {
			t.Errorf("Code(%#x).Transient() = false, want true", uint32(c))
		}
	}
	for _, c := range []status.Code{status.RpcNotFound, status.RpcGone, status.RpcBadRequest, status.KeyPairNotFound} {
		if c.Transient() {
			t.Errorf("Code(%#x).Transient() = true, want false", uint32(c))
		}
	}
}

func TestValueSentinels(t *testing.T) {
	cases := []struct {
		err  error
		code status.Code
	}{
		{status.ErrValueKeyInvalid, status.ValueKeyInvalid},
		{status.ErrValueValueTooLarge, status.ValueValueTooLarge},
		{status.ErrValueTransactionTooLarge, status.ValueTransactionTooLarge},
		{status.ErrValueBusy, status.ValueBusy},
		{status.ErrValueTransactionClosed, status.ValueTransactionClosed},
		{status.ErrValueTransactionActive, status.ValueTransactionActive},
		{status.ErrValueCheckpointInvalid, status.ValueCheckpointInvalid},
		{status.ErrValueCursorInvalid, status.ValueCursorInvalid},
	}
	for i, tc := range cases {
		if code, ok := status.Of(tc.err); !ok || code != tc.code {
			t.Errorf("status.Of(%v) = %#x, %v, want %#x", tc.err, uint32(code), ok, uint32(tc.code))
		}
		for j, other := range cases {
			if got := errors.Is(tc.err, other.err); got != (i == j) {
				t.Errorf("errors.Is(%v, %v) = %v", tc.err, other.err, got)
			}
		}
	}
}
