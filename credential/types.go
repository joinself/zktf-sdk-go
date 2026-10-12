package credential

import "github.com/joinself/zktf-sdk-go/internal/ffi"

func CredentialTypesDefault() []string { return ffi.CredentialTypesDefault() }

func CredentialTypesPassport() []string { return ffi.CredentialTypesPassport() }

func CredentialTypesFacialComparison() []string { return ffi.CredentialTypesFacialComparison() }

func CredentialTypesLivenessAndFacialComparison() []string {
	return ffi.CredentialTypesLivenessAndFacialComparison()
}

func CredentialTypesBiometricAnchor() []string { return ffi.CredentialTypesBiometricAnchor() }

func CredentialTypesSharingAgreement() []string { return ffi.CredentialTypesSharingAgreement() }

func CredentialTypesEmail() []string { return ffi.CredentialTypesEmail() }

func CredentialTypesPhone() []string { return ffi.CredentialTypesPhone() }

func CredentialTypesOrganisation() []string { return ffi.CredentialTypesOrganisation() }

func CredentialTypesApplication() []string { return ffi.CredentialTypesApplication() }

func CredentialTypesApplicationInstall() []string { return ffi.CredentialTypesApplicationInstall() }

func PresentationTypesDefault() []string { return ffi.PresentationTypesDefault() }

func PresentationTypesPassport() []string { return ffi.PresentationTypesPassport() }

func PresentationTypesFacialComparison() []string { return ffi.PresentationTypesFacialComparison() }

func PresentationTypesLivenessAndFacialComparison() []string {
	return ffi.PresentationTypesLivenessAndFacialComparison()
}

func PresentationTypesBiometricAnchor() []string { return ffi.PresentationTypesBiometricAnchor() }

func PresentationTypesSharingAgreement() []string { return ffi.PresentationTypesSharingAgreement() }

func PresentationTypesApplicationInstall() []string { return ffi.PresentationTypesApplicationInstall() }
