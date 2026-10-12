package account

import (
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/token"
)

type NetworkState uint32

const (
	NetworkUnavailable NetworkState = NetworkState(ffi.NetworkStateUnavailable)
	NetworkAvailable   NetworkState = NetworkState(ffi.NetworkStateAvailable)
)

func (a *Account) NetworkStateSet(state NetworkState) error {
	return a.h.NetworkStateSet(ffi.NetworkState(state))
}

func (a *Account) PushCredentialSet(push *token.Push) error {
	return a.h.PushCredentialSet(ffi.PlatformPushOf(push))
}

func (a *Account) PushCredentialGet() (*token.Push, error) {
	p, err := a.h.PushCredentialGet()
	if err != nil || p == nil {
		return nil, err
	}

	return ffi.ToPlatformPush(p).(*token.Push), nil
}
