package ffi

/*
#include <zktf-sdk.h>
*/
import "C"

func CredentialTypesDefault() []string {
	return stringsFromBufferCollection(C.zktf_credential_type_default())
}

func CredentialTypesPassport() []string {
	return stringsFromBufferCollection(C.zktf_credential_type_passport())
}

func CredentialTypesFacialComparison() []string {
	return stringsFromBufferCollection(C.zktf_credential_type_facial_comparison())
}

func CredentialTypesLivenessAndFacialComparison() []string {
	return stringsFromBufferCollection(C.zktf_credential_type_liveness_and_facial_comparison())
}

func CredentialTypesBiometricAnchor() []string {
	return stringsFromBufferCollection(C.zktf_credential_type_biometric_anchor())
}

func CredentialTypesSharingAgreement() []string {
	return stringsFromBufferCollection(C.zktf_credential_type_sharing_agreement())
}

func CredentialTypesEmail() []string {
	return stringsFromBufferCollection(C.zktf_credential_type_email())
}

func CredentialTypesPhone() []string {
	return stringsFromBufferCollection(C.zktf_credential_type_phone())
}

func CredentialTypesOrganisation() []string {
	return stringsFromBufferCollection(C.zktf_credential_type_organisation())
}

func CredentialTypesApplication() []string {
	return stringsFromBufferCollection(C.zktf_credential_type_application())
}

func CredentialTypesApplicationInstall() []string {
	return stringsFromBufferCollection(C.zktf_credential_type_application_install())
}

func PresentationTypesDefault() []string {
	return stringsFromBufferCollection(C.zktf_presentation_type_default())
}

func PresentationTypesPassport() []string {
	return stringsFromBufferCollection(C.zktf_presentation_type_passport())
}

func PresentationTypesFacialComparison() []string {
	return stringsFromBufferCollection(C.zktf_presentation_type_facial_comparison())
}

func PresentationTypesLivenessAndFacialComparison() []string {
	return stringsFromBufferCollection(C.zktf_presentation_type_liveness_and_facial_comparison())
}

func PresentationTypesBiometricAnchor() []string {
	return stringsFromBufferCollection(C.zktf_presentation_type_biometric_anchor())
}

func PresentationTypesSharingAgreement() []string {
	return stringsFromBufferCollection(C.zktf_presentation_type_sharing_agreement())
}

func PresentationTypesApplicationInstall() []string {
	return stringsFromBufferCollection(C.zktf_presentation_type_application_install())
}
