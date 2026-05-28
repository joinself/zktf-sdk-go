package ffi

import (
	"reflect"
	"testing"
)

func TestCredentialTypeCollectionRoundTrip(t *testing.T) {
	in := []string{"VerifiableCredential", "EmailCredential"}
	got := NewCredentialTypes(in).Strings()
	if !reflect.DeepEqual(got, in) {
		t.Fatalf("CredentialTypes = %v, want %v", got, in)
	}
}

func TestPresentationTypeCollectionRoundTrip(t *testing.T) {
	in := []string{"VerifiablePresentation", "PassportPresentation"}
	got := NewPresentationTypes(in).Strings()
	if !reflect.DeepEqual(got, in) {
		t.Fatalf("PresentationTypes = %v, want %v", got, in)
	}
}
