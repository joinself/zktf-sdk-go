package credential

import (
	"github.com/joinself/zktf-sdk-go/identity"
	"github.com/joinself/zktf-sdk-go/internal/ffi"
)

// Presentation type strings, used with PresentationBuilder.Type.
const (
	PresentationTypePassport                    = "PassportPresentation"
	PresentationTypeFacialComparison            = "FacialComparisonPresentation"
	PresentationTypeLivenessAndFacialComparison = "LivenessAndFacialComparisonPresentation"
	PresentationTypeBiometricAnchor             = "BiometricAnchorPresentation"
	PresentationTypeSharingAgreement            = "SharingAgreementPresentation"
	PresentationTypeProfile                     = "ProfilePresentation"
	PresentationTypeContactDetails              = "ContactDetailsPresentation"
	PresentationTypeApplicationPublisher        = "ApplicationPublisherPresentation"
)

// Presentation is an unsigned presentation produced by a PresentationBuilder.
type Presentation struct {
	h *ffi.Presentation
}

// VerifiablePresentation is a signed, verifiable presentation.
type VerifiablePresentation struct {
	h *ffi.VerifiablePresentation
}

func init() {
	ffi.PresentationOf = func(o any) *ffi.Presentation { return o.(*Presentation).h }
	ffi.ToPresentation = func(h *ffi.Presentation) any { return &Presentation{h: h} }

	ffi.VerifiablePresentationOf = func(o any) *ffi.VerifiablePresentation {
		return o.(*VerifiablePresentation).h
	}
	ffi.ToVerifiablePresentation = func(h *ffi.VerifiablePresentation) any {
		return &VerifiablePresentation{h: h}
	}
}

// PresentationBuilder builds an unsigned presentation.
type PresentationBuilder struct {
	h *ffi.PresentationBuilder
}

// NewPresentation starts building a presentation.
func NewPresentation() *PresentationBuilder {
	return &PresentationBuilder{h: ffi.NewPresentationBuilder()}
}

// Type sets the presentation's types.
func (b *PresentationBuilder) Type(types ...string) *PresentationBuilder {
	b.h.PresentationType(ffi.NewPresentationTypes(types))
	return b
}

// Holder sets the holder/bearer address.
func (b *PresentationBuilder) Holder(holder *identity.Address) *PresentationBuilder {
	b.h.Holder(ffi.DIDAddressOf(holder))
	return b
}

// Credential adds a verifiable credential to the presentation.
func (b *PresentationBuilder) Credential(credentials ...*Verifiable) *PresentationBuilder {
	for _, c := range credentials {
		b.h.CredentialAdd(c.h)
	}

	return b
}

// Finish finalizes the unsigned presentation.
func (b *PresentationBuilder) Finish() (*Presentation, error) {
	p, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Presentation{h: p}, nil
}

// DecodePresentation decodes a JSON-encoded verifiable presentation.
func DecodePresentation(data []byte) (*VerifiablePresentation, error) {
	p, err := ffi.VerifiablePresentationDecode(data)
	if err != nil {
		return nil, err
	}

	return &VerifiablePresentation{h: p}, nil
}

// Validate returns an error if the presentation is invalid.
func (v *VerifiablePresentation) Validate() error { return v.h.Validate() }

// Types returns the presentation's type strings.
func (v *VerifiablePresentation) Types() []string { return v.h.Types() }

// Holder returns the holder address, or nil.
func (v *VerifiablePresentation) Holder() *identity.Address {
	h := v.h.Holder()
	if h == nil {
		return nil
	}

	return ffi.ToDIDAddress(h).(*identity.Address)
}

// Credentials returns the credentials contained in the presentation.
func (v *VerifiablePresentation) Credentials() []*Verifiable {
	cs := v.h.Credentials()
	out := make([]*Verifiable, len(cs))

	for i, c := range cs {
		out[i] = &Verifiable{h: c}
	}

	return out
}

// Encode returns the JSON-encoded presentation.
func (v *VerifiablePresentation) Encode() ([]byte, error) { return v.h.Encode() }
