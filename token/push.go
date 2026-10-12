package token

import (
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

type Push struct {
	h *ffi.PlatformPush
}

func init() {
	ffi.PlatformPushOf = func(o any) *ffi.PlatformPush { return o.(*Push).h }
	ffi.ToPlatformPush = func(h *ffi.PlatformPush) any { return &Push{h: h} }
}

func NewFCMPush(application *signing.PublicKey, credential string) *Push {
	return &Push{h: ffi.NewPlatformPushFCM(ffi.SigningPublicKeyOf(application), credential)}
}

func NewAPNSPush(application *signing.PublicKey, credential string) *Push {
	return &Push{h: ffi.NewPlatformPushAPNS(ffi.SigningPublicKeyOf(application), credential)}
}

func (b *PushBuilder) PushCredential(push *Push) *PushBuilder {
	b.h.PushCredential(push.h)
	return b
}
