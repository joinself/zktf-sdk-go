// Package credential provides credential issuance and verification types for the
// zktf SDK.
package credential

import (
	"time"

	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
)

// Credential type strings, used with Builder.Type.
const (
	TypeEmail                       = "EmailCredential"
	TypePhone                       = "PhoneCredential"
	TypePassport                    = "PassportCredential"
	TypeFacialComparison            = "FacialComparisonCredential"
	TypeLivenessAndFacialComparison = "LivenessAndFacialComparisonCredential"
	TypeBiometricAnchor             = "BiometricAnchorCredential"
	TypeSharingAgreement            = "SharingAgreementCredential"
	TypeOrganisation                = "OrganisationCredential"
	TypeApplication                 = "ApplicationCredential"
)

// Credential field paths (RFC 6901 JSON pointers), used with Claim and predicates.
const (
	FieldType                                              = "/type"
	FieldIssuer                                            = "/issuer"
	FieldValidFrom                                         = "/validFrom"
	FieldValidUntil                                        = "/validUntil"
	FieldSubject                                           = "/credentialSubject/id"
	FieldSubjectClaims                                     = "/credentialSubject"
	FieldSubjectEmailAddress                               = "/credentialSubject/email/emailAddress"
	FieldSubjectPhoneNumber                                = "/credentialSubject/phone/phoneNumber"
	FieldSubjectBiometricAnchorSourceImageHash             = "/credentialSubject/biometricAnchor/sourceImageHash"
	FieldSubjectBiometricAnchorComputedHashes              = "/credentialSubject/biometricAnchor/computedHashes"
	FieldSubjectFacialComparisonSourceImageHash            = "/credentialSubject/facialComparison/sourceImageHash"
	FieldSubjectFacialComparisonTargetImageHash            = "/credentialSubject/facialComparison/targetImageHash"
	FieldSubjectLivenessAndFacialComparisonSourceImageHash = "/credentialSubject/livenessAndFacialComparison/sourceImageHash"
	FieldSubjectLivenessAndFacialComparisonTargetImageHash = "/credentialSubject/livenessAndFacialComparison/targetImageHash"
	FieldSubjectLivenessAndFacialComparisonChallenge       = "/credentialSubject/livenessAndFacialComparison/challenge"
	FieldSubjectLivenessAndFacialComparisonComputedHashes  = "/credentialSubject/livenessAndFacialComparison/computedHashes"
	FieldSubjectPassportDocumentNumber                     = "/credentialSubject/passport/documentNumber"
	FieldSubjectPassportGivenNames                         = "/credentialSubject/passport/givenNames"
	FieldSubjectPassportSurname                            = "/credentialSubject/passport/surname"
	FieldSubjectPassportSex                                = "/credentialSubject/passport/sex"
	FieldSubjectPassportNationality                        = "/credentialSubject/passport/nationality"
	FieldSubjectPassportDateOfBirth                        = "/credentialSubject/passport/dateOfBirth"
	FieldSubjectPassportDateOfExpiration                   = "/credentialSubject/passport/dateOfExpiration"
	FieldSubjectPassportCountryOfIssuance                  = "/credentialSubject/passport/countryOfIssuance"
	FieldSubjectPassportDocumentMrz                        = "/credentialSubject/passport/mrz"
	FieldSubjectPassportImageType                          = "/credentialSubject/passport/imageType"
	FieldSubjectPassportImageHash                          = "/credentialSubject/passport/imageHash"
	FieldSubjectPassportTargetImageHash                    = "/credentialSubject/passport/targetImageHash"
	FieldSubjectOrganisationName                           = "/credentialSubject/organisation/organisationName"
	FieldSubjectApplicationName                            = "/credentialSubject/application/applicationName"
	FieldSubjectApplicationSubsidiaryOf                    = "/credentialSubject/application/subsidiaryOf"
)

// Address is a decentralized identifier (DID) address.
type Address struct {
	h *ffi.DIDAddress
}

// Term describes the duration under which a verifier wishes to access requested
// credentials.
type Term struct {
	h *ffi.CredentialTerm
}

// Credential is an unsigned credential produced by a Builder.
type Credential struct {
	h *ffi.Credential
}

// Verifiable is a signed, verifiable credential.
type Verifiable struct {
	h *ffi.VerifiableCredential
}

func init() {
	ffi.DIDAddressOf = func(o any) *ffi.DIDAddress { return o.(*Address).h }
	ffi.ToDIDAddress = func(h *ffi.DIDAddress) any { return &Address{h: h} }

	ffi.CredentialTermOf = func(o any) *ffi.CredentialTerm { return o.(*Term).h }
	ffi.ToCredentialTerm = func(h *ffi.CredentialTerm) any { return &Term{h: h} }

	ffi.CredentialOf = func(o any) *ffi.Credential { return o.(*Credential).h }
	ffi.ToCredential = func(h *ffi.Credential) any { return &Credential{h: h} }

	ffi.VerifiableCredentialOf = func(o any) *ffi.VerifiableCredential { return o.(*Verifiable).h }
	ffi.ToVerifiableCredential = func(h *ffi.VerifiableCredential) any { return &Verifiable{h: h} }
}

// AddressKey builds a key-method DID address from a signing key.
func AddressKey(key *signing.PublicKey) *Address {
	return &Address{h: ffi.DIDAddressKey(ffi.SigningPublicKeyOf(key))}
}

// ParseAddress decodes a DID string into an address.
func ParseAddress(did string) (*Address, error) {
	a, err := ffi.DIDAddressDecode(did)
	if err != nil {
		return nil, err
	}

	return &Address{h: a}, nil
}

// Key returns the signing public key embedded in the address.
func (a *Address) Key() *signing.PublicKey {
	return ffi.ToSigningPublicKey(a.h.Address()).(*signing.PublicKey)
}

// String returns the encoded DID string.
func (a *Address) String() string { return a.h.String() }

// Preset terms covering the common access durations. Month and Year use the
// average Gregorian second counts.
var (
	TermSingleUse = NewTerm(0)
	TermHour      = NewTerm(time.Hour)
	TermDay       = NewTerm(24 * time.Hour)
	TermWeek      = NewTerm(7 * 24 * time.Hour)
	TermMonth     = NewTerm(2629746 * time.Second)
	TermYear      = NewTerm(31556952 * time.Second)
)

// NewTerm creates a credential term lasting the given duration.
func NewTerm(duration time.Duration) *Term {
	return &Term{h: ffi.NewCredentialTerm(uint64(duration / time.Second))}
}

// Duration returns the term's duration.
func (t *Term) Duration() time.Duration {
	return time.Duration(t.h.Duration()) * time.Second
}

// Builder builds an unsigned credential.
type Builder struct {
	h *ffi.CredentialBuilder
}

// NewBuilder starts building a credential.
func NewBuilder() *Builder { return &Builder{h: ffi.NewCredentialBuilder()} }

// Type sets the credential's types.
func (b *Builder) Type(types ...string) *Builder {
	b.h.CredentialType(ffi.NewCredentialTypes(types))
	return b
}

// Issuer sets the credential's issuer.
func (b *Builder) Issuer(issuer *Address) *Builder {
	b.h.Issuer(issuer.h)
	return b
}

// Subject sets the credential's subject.
func (b *Builder) Subject(subject *Address) *Builder {
	b.h.CredentialSubject(subject.h)
	return b
}

// Claim adds a string claim about the subject.
func (b *Builder) Claim(key, value string) *Builder {
	b.h.CredentialSubjectClaim(key, value)
	return b
}

// ValidFrom sets when the credential becomes valid.
func (b *Builder) ValidFrom(t time.Time) *Builder {
	b.h.ValidFrom(t.Unix())
	return b
}

// ValidUntil sets when the credential stops being valid.
func (b *Builder) ValidUntil(t time.Time) *Builder {
	b.h.ValidUntil(t.Unix())
	return b
}

// ClaimJSON sets the subject claims from a raw JSON document.
func (b *Builder) ClaimJSON(json []byte) *Builder {
	b.h.CredentialSubjectJSON(json)
	return b
}

// SignWith records the signing key and issuance time.
func (b *Builder) SignWith(signer *signing.PublicKey, issuedAt time.Time) *Builder {
	b.h.SignWith(ffi.SigningPublicKeyOf(signer), issuedAt.Unix())
	return b
}

// Finish finalizes the unsigned credential. Sign it via Account.CredentialIssue.
func (b *Builder) Finish() (*Credential, error) {
	c, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &Credential{h: c}, nil
}

// Decode decodes a JSON-encoded verifiable credential.
func Decode(data []byte) (*Verifiable, error) {
	c, err := ffi.VerifiableCredentialDecode(data)
	if err != nil {
		return nil, err
	}

	return &Verifiable{h: c}, nil
}

// Validate returns an error if the credential is invalid.
func (v *Verifiable) Validate() error { return v.h.Validate() }

// Types returns the credential's type strings.
func (v *Verifiable) Types() []string { return v.h.TypeOf().Strings() }

// Issuer returns the issuer address.
func (v *Verifiable) Issuer() *Address { return &Address{h: v.h.Issuer()} }

// Subject returns the subject address.
func (v *Verifiable) Subject() *Address { return &Address{h: v.h.Subject()} }

// Claim returns a string claim about the subject, or "" if absent.
func (v *Verifiable) Claim(key string) string { return v.h.SubjectClaim(key) }

// ClaimJSON returns the subject claims as a raw JSON document, or nil.
func (v *Verifiable) ClaimJSON() []byte { return v.h.SubjectJSON() }

// ValidFrom returns when the credential became valid.
func (v *Verifiable) ValidFrom() time.Time { return time.Unix(v.h.ValidFrom(), 0) }

// ValidUntil returns when the credential stops being valid.
func (v *Verifiable) ValidUntil() time.Time { return time.Unix(v.h.ValidUntil(), 0) }

// Created returns when the credential was created.
func (v *Verifiable) Created() time.Time { return time.Unix(v.h.Created(), 0) }

// Signer returns the DID address that signed the credential.
func (v *Verifiable) Signer() (*Address, error) {
	a, err := v.h.Signer()
	if err != nil {
		return nil, err
	}

	return &Address{h: a}, nil
}

// SigningKey returns the signing key that signed the credential.
func (v *Verifiable) SigningKey() (*signing.PublicKey, error) {
	k, err := v.h.SigningKey()
	if err != nil {
		return nil, err
	}

	return ffi.ToSigningPublicKey(k).(*signing.PublicKey), nil
}

// RevocationHashes returns the credential's revocation hashes, one per proof.
func (v *Verifiable) RevocationHashes() ([][]byte, error) { return v.h.RevocationHashes() }

// Encode returns the JSON-encoded credential.
func (v *Verifiable) Encode() ([]byte, error) { return v.h.Encode() }
