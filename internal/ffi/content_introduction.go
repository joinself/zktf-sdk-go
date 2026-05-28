package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import "runtime"

// Introduction wraps a zktf_message_content_introduction handle.
type Introduction struct {
	ptr *C.zktf_message_content_introduction
}

func newIntroduction(ptr *C.zktf_message_content_introduction) *Introduction {
	if ptr == nil {
		return nil
	}
	i := &Introduction{ptr: ptr}
	runtime.AddCleanup(i, func(ptr *C.zktf_message_content_introduction) {
		C.zktf_message_content_introduction_destroy(ptr)
	}, i.ptr)
	return i
}

// IntroductionFromContent decodes message content as an introduction.
func IntroductionFromContent(content *Content) (*Introduction, error) {
	var ptr *C.zktf_message_content_introduction
	if err := status(C.zktf_message_content_as_introduction(content.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newIntroduction(ptr), nil
}

// DocumentAddress returns the sender's document DID address.
func (i *Introduction) DocumentAddress() *DIDAddress {
	return newDIDAddress(C.zktf_message_content_introduction_document_address(i.ptr))
}

// Presentations returns the verified presentations shared by the sender.
func (i *Introduction) Presentations() []*VerifiablePresentation {
	return verifiablePresentationsFrom(C.zktf_message_content_introduction_presentations(i.ptr))
}

// Tokens returns the tokens issued by the sender.
func (i *Introduction) Tokens() ([]*Token, error) {
	var c *C.zktf_collection_token
	if err := status(C.zktf_message_content_introduction_tokens(i.ptr, &c)); err != nil {
		return nil, err
	}
	return tokensFrom(c), nil
}

// Assets returns supporting object assets attached to the introduction.
func (i *Introduction) Assets() []*Object {
	return objectsFrom(C.zktf_message_content_introduction_assets(i.ptr))
}

// PairwiseIntroduction extracts the pairwise introduction (suitable for
// validating with Account.PairwiseValidateIntroduction).
func (i *Introduction) PairwiseIntroduction() (*PairwiseIntroduction, error) {
	var out *C.zktf_pairwise_introduction
	if err := status(C.zktf_message_content_introduction_introduction(i.ptr, &out)); err != nil {
		return nil, err
	}
	return newPairwiseIntroduction(out), nil
}

// IntroductionBuilder builds an introduction message content.
type IntroductionBuilder struct {
	ptr *C.zktf_message_content_introduction_builder
}

// NewIntroductionBuilder initializes an introduction builder.
func NewIntroductionBuilder() *IntroductionBuilder {
	ptr := C.zktf_message_content_introduction_builder_init()
	b := &IntroductionBuilder{ptr: ptr}
	runtime.AddCleanup(b, func(ptr *C.zktf_message_content_introduction_builder) {
		C.zktf_message_content_introduction_builder_destroy(ptr)
	}, b.ptr)
	return b
}

// DocumentAddress sets the document address the sender wants to identify as.
func (b *IntroductionBuilder) DocumentAddress(address *DIDAddress) *IntroductionBuilder {
	C.zktf_message_content_introduction_builder_document_address(b.ptr, address.ptr)
	return b
}

// Presentation adds a verifiable presentation to the introduction.
func (b *IntroductionBuilder) Presentation(p *VerifiablePresentation) *IntroductionBuilder {
	C.zktf_message_content_introduction_builder_presentation(b.ptr, p.ptr)
	return b
}

// Token attaches a token (e.g. a delegation/send token) to the introduction.
func (b *IntroductionBuilder) Token(t *Token) *IntroductionBuilder {
	C.zktf_message_content_introduction_builder_token(b.ptr, t.ptr)
	return b
}

// Asset attaches a supporting object asset.
func (b *IntroductionBuilder) Asset(o *Object) *IntroductionBuilder {
	C.zktf_message_content_introduction_builder_asset(b.ptr, o.ptr)
	return b
}

// Finish finalizes the introduction content, ready to send.
func (b *IntroductionBuilder) Finish() (*Content, error) {
	var out *C.zktf_message_content
	if err := status(C.zktf_message_content_introduction_builder_finish(b.ptr, &out)); err != nil {
		return nil, err
	}
	return newContent(out), nil
}
