// Package credential provides credential issuance and verification types for the
// zktf SDK.
package credential

import (
	"time"

	"github.com/joinself/zktf-sdk-go/identity"
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

// Term describes the duration under which a verifier wishes to access requested
// credentials.
type Term struct {
	h *ffi.CredentialTerm
}

// VerifiableCredential is a verifiable credential. A Builder produces one
// unsigned, carrying its pending signers; Account.CredentialIssue signs it.
type VerifiableCredential struct {
	h *ffi.VerifiableCredential
}

func init() {
	ffi.CredentialTermOf = func(o any) *ffi.CredentialTerm { return o.(*Term).h }
	ffi.ToCredentialTerm = func(h *ffi.CredentialTerm) any { return &Term{h: h} }

	ffi.VerifiableCredentialOf = func(o any) *ffi.VerifiableCredential { return o.(*VerifiableCredential).h }
	ffi.ToVerifiableCredential = func(h *ffi.VerifiableCredential) any { return &VerifiableCredential{h: h} }
}

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
func (b *Builder) Issuer(issuer *identity.Address) *Builder {
	b.h.Issuer(ffi.DIDAddressOf(issuer))
	return b
}

// Subject sets the credential's subject.
func (b *Builder) Subject(subject *identity.Address) *Builder {
	b.h.CredentialSubject(ffi.DIDAddressOf(subject))
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

// Finish finalizes the credential, unsigned. Sign it via
// Account.CredentialIssue.
func (b *Builder) Finish() (*VerifiableCredential, error) {
	c, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &VerifiableCredential{h: c}, nil
}

// Decode decodes a JSON-encoded verifiable credential.
func Decode(data []byte) (*VerifiableCredential, error) {
	c, err := ffi.VerifiableCredentialDecode(data)
	if err != nil {
		return nil, err
	}

	return &VerifiableCredential{h: c}, nil
}

// Validate returns an error if the credential is invalid.
func (v *VerifiableCredential) Validate() error { return v.h.Validate() }

// Types returns the credential's type strings.
func (v *VerifiableCredential) Types() []string { return v.h.TypeOf().Strings() }

// Issuer returns the issuer address.
func (v *VerifiableCredential) Issuer() *identity.Address {
	return ffi.ToDIDAddress(v.h.Issuer()).(*identity.Address)
}

// Subject returns the subject address.
func (v *VerifiableCredential) Subject() *identity.Address {
	subject := v.h.Subject()
	if subject == nil {
		return nil
	}
	return ffi.ToDIDAddress(subject).(*identity.Address)
}

// Claim returns a string claim about the subject, or "" if absent.
func (v *VerifiableCredential) Claim(key string) string { return v.h.SubjectClaim(key) }

// ClaimJSON returns the subject claims as a raw JSON document, or nil.
func (v *VerifiableCredential) ClaimJSON() []byte { return v.h.SubjectJSON() }

// ValidFrom returns when the credential became valid.
func (v *VerifiableCredential) ValidFrom() time.Time { return time.Unix(v.h.ValidFrom(), 0) }

// ValidUntil returns when the credential stops being valid.
func (v *VerifiableCredential) ValidUntil() time.Time { return time.Unix(v.h.ValidUntil(), 0) }

// Created returns when the credential was created.
func (v *VerifiableCredential) Created() time.Time { return time.Unix(v.h.Created(), 0) }

// Signer returns the DID address that signed the credential.
func (v *VerifiableCredential) Signer() (*identity.Address, error) {
	a, err := v.h.Signer()
	if err != nil {
		return nil, err
	}

	return ffi.ToDIDAddress(a).(*identity.Address), nil
}

// SigningKey returns the signing key that signed the credential.
func (v *VerifiableCredential) SigningKey() (*signing.PublicKey, error) {
	k, err := v.h.SigningKey()
	if err != nil {
		return nil, err
	}

	return ffi.ToSigningPublicKey(k).(*signing.PublicKey), nil
}

// RevocationHashes returns the credential's revocation hashes, one per proof.
func (v *VerifiableCredential) RevocationHashes() ([][]byte, error) { return v.h.RevocationHashes() }

// Encode returns the JSON-encoded credential.
func (v *VerifiableCredential) Encode() ([]byte, error) { return v.h.Encode() }
