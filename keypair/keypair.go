// Package keypair provides shared types for the signing and exchange keypair
// subpackages.
package keypair

// PublicKey is a public key of any keypair type (signing or exchange). It is
// implemented by *signing.PublicKey and *exchange.PublicKey, allowing document
// queries to accept either kind of key.
type PublicKey interface {
	// Bytes returns the raw bytes of the public key.
	Bytes() []byte
	// String returns the hex encoded address.
	String() string
}
