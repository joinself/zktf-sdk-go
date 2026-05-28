package message

import (
	"github.com/joinself/zktf-sdk-go/credential"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/object"
)

// CredentialContent is a decoded credential message payload.
type CredentialContent struct {
	h *ffi.CredentialContent
}

// CredentialContentBuilder builds credential message content.
type CredentialContentBuilder struct {
	h *ffi.CredentialContentBuilder
}

// CredentialDecode decodes message content as a credential payload.
func CredentialDecode(content *Content) (*CredentialContent, error) {
	c, err := ffi.CredentialContentFromContent(content.h)
	if err != nil {
		return nil, err
	}

	return &CredentialContent{h: c}, nil
}

// VerifiablePresentations returns the presentations carried in the content.
func (c *CredentialContent) VerifiablePresentations() []*credential.VerifiablePresentation {
	ps := c.h.VerifiablePresentations()
	out := make([]*credential.VerifiablePresentation, len(ps))

	for i, p := range ps {
		out[i] = ffi.ToVerifiablePresentation(p).(*credential.VerifiablePresentation)
	}

	return out
}

// VerifiableCredentials returns the credentials carried in the content.
func (c *CredentialContent) VerifiableCredentials() []*credential.Verifiable {
	cs := c.h.VerifiableCredentials()
	out := make([]*credential.Verifiable, len(cs))

	for i, vc := range cs {
		out[i] = ffi.ToVerifiableCredential(vc).(*credential.Verifiable)
	}

	return out
}

// Assets returns supporting object assets carried in the content.
func (c *CredentialContent) Assets() []*object.Object {
	os := c.h.Assets()
	out := make([]*object.Object, len(os))

	for i, o := range os {
		out[i] = ffi.ToObject(o).(*object.Object)
	}

	return out
}

// NewCredentialContent starts building credential content.
func NewCredentialContent() *CredentialContentBuilder {
	return &CredentialContentBuilder{h: ffi.NewCredentialContentBuilder()}
}

// VerifiablePresentation adds a presentation.
func (b *CredentialContentBuilder) VerifiablePresentation(p *credential.VerifiablePresentation) *CredentialContentBuilder {
	b.h.VerifiablePresentation(ffi.VerifiablePresentationOf(p))
	return b
}

// VerifiableCredential adds a credential.
func (b *CredentialContentBuilder) VerifiableCredential(c *credential.Verifiable) *CredentialContentBuilder {
	b.h.VerifiableCredential(ffi.VerifiableCredentialOf(c))
	return b
}

// Asset attaches a supporting object asset.
func (b *CredentialContentBuilder) Asset(o *object.Object) *CredentialContentBuilder {
	b.h.Asset(ffi.ObjectOf(o))
	return b
}

// Finish finalizes the credential content, ready to send.
func (b *CredentialContentBuilder) Finish() (*Content, error) {
	c, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Content{h: c}, nil
}
