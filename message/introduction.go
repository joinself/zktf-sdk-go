package message

import (
	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/object"
	"github.com/joinself/zktf-sdk-go/pairwise"
	"github.com/joinself/zktf-sdk-go/token"
)

// Introduction is a decoded introduction message.
type Introduction struct {
	h *ffi.Introduction
}

// IntroductionBuilder builds introduction message content.
type IntroductionBuilder struct {
	h *ffi.IntroductionBuilder
}

// IntroductionDecode decodes message content as an introduction.
func IntroductionDecode(content *Content) (*Introduction, error) {
	i, err := ffi.IntroductionFromContent(content.h)
	if err != nil {
		return nil, err
	}

	return &Introduction{h: i}, nil
}

// DocumentAddress returns the sender's document DID address.
func (i *Introduction) DocumentAddress() *credential.Address {
	return ffi.ToDIDAddress(i.h.DocumentAddress()).(*credential.Address)
}

// Presentations returns the verified presentations shared by the sender.
func (i *Introduction) Presentations() []*credential.VerifiablePresentation {
	ps := i.h.Presentations()
	out := make([]*credential.VerifiablePresentation, len(ps))

	for j, p := range ps {
		out[j] = ffi.ToVerifiablePresentation(p).(*credential.VerifiablePresentation)
	}

	return out
}

// Tokens returns the tokens issued by the sender.
func (i *Introduction) Tokens() ([]*token.Token, error) {
	ts, err := i.h.Tokens()
	if err != nil {
		return nil, err
	}

	out := make([]*token.Token, len(ts))
	for j, t := range ts {
		out[j] = ffi.ToToken(t).(*token.Token)
	}

	return out, nil
}

// Assets returns supporting object assets.
func (i *Introduction) Assets() []*object.Object {
	os := i.h.Assets()
	out := make([]*object.Object, len(os))

	for j, o := range os {
		out[j] = ffi.ToObject(o).(*object.Object)
	}

	return out
}

// PairwiseIntroduction extracts the embedded pairwise introduction.
func (i *Introduction) PairwiseIntroduction() (*pairwise.Introduction, error) {
	pi, err := i.h.PairwiseIntroduction()
	if err != nil {
		return nil, err
	}

	return ffi.ToPairwiseIntroduction(pi).(*pairwise.Introduction), nil
}

// NewIntroduction starts building an introduction.
func NewIntroduction() *IntroductionBuilder {
	return &IntroductionBuilder{h: ffi.NewIntroductionBuilder()}
}

// DocumentAddress sets the document address the sender identifies as.
func (b *IntroductionBuilder) DocumentAddress(address *credential.Address) *IntroductionBuilder {
	b.h.DocumentAddress(ffi.DIDAddressOf(address))
	return b
}

// Presentation adds a verifiable presentation.
func (b *IntroductionBuilder) Presentation(p *credential.VerifiablePresentation) *IntroductionBuilder {
	b.h.Presentation(ffi.VerifiablePresentationOf(p))
	return b
}

// Token attaches a token.
func (b *IntroductionBuilder) Token(t *token.Token) *IntroductionBuilder {
	b.h.Token(ffi.TokenOf(t))
	return b
}

// Asset attaches a supporting object asset.
func (b *IntroductionBuilder) Asset(o *object.Object) *IntroductionBuilder {
	b.h.Asset(ffi.ObjectOf(o))
	return b
}

// Finish finalizes the introduction content, ready to send.
func (b *IntroductionBuilder) Finish() (*Content, error) {
	c, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Content{h: c}, nil
}
