package account_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/joinself/zktf-sdk-go/account"
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sdk-go/message"
	"github.com/joinself/zktf-sdk-go/revocation"
)

func offlineAccount(t *testing.T) *account.Account {
	t.Helper()

	acc, err := account.New(account.Config{
		StoragePath: ":memory:",
		Target: &account.Target{
			Network:         account.NetworkPreview,
			RPCEndpoint:     "http://127.0.0.1:1/",
			ObjectEndpoint:  "http://127.0.0.1:1/",
			MessageEndpoint: "ws://127.0.0.1:1/",
		},
	}, account.Callbacks{})
	if err != nil {
		t.Fatalf("account.New: %v", err)
	}
	t.Cleanup(acc.Close)

	return acc
}

func signedStatement(t *testing.T, acc *account.Account, issuer *signing.PublicKey, hash []byte) *revocation.Statement {
	t.Helper()

	at := time.Unix(1_900_000_000, 0)

	stmt, err := revocation.NewBuilder().
		Issuer(issuer).
		Sequence(1).
		Timestamp(at).
		RevokeBy(hash, at).
		SignWith(issuer, at).
		Finish()
	if err != nil {
		t.Fatalf("revocation Finish: %v", err)
	}

	if err := acc.RevocationSign(stmt); err != nil {
		t.Fatalf("RevocationSign: %v", err)
	}

	return stmt
}

func TestExchangeKeySealAnonymous(t *testing.T) {
	acc := offlineAccount(t)

	key, err := acc.KeychainExchangeCreate()
	if err != nil {
		t.Fatalf("KeychainExchangeCreate: %v", err)
	}

	plaintext := []byte("sealed for the key holder")
	sealed, err := key.SealAnonymous(plaintext)
	if err != nil {
		t.Fatalf("SealAnonymous: %v", err)
	}
	if len(sealed) != len(plaintext)+48 {
		t.Fatalf("len(sealed) = %d, want %d", len(sealed), len(plaintext)+48)
	}
	if bytes.Contains(sealed, plaintext) {
		t.Fatal("sealed output contains the plaintext")
	}

	if _, err := key.SealAnonymous(nil); err != nil {
		t.Fatalf("SealAnonymous empty: %v", err)
	}
}

func TestRevocationStatementMergeAndMerkleRoot(t *testing.T) {
	acc := offlineAccount(t)

	issuer, err := acc.KeychainSigningCreate()
	if err != nil {
		t.Fatalf("KeychainSigningCreate: %v", err)
	}

	hash := bytes.Repeat([]byte{0x33}, 32)
	a := signedStatement(t, acc, issuer, hash)
	b := signedStatement(t, acc, issuer, hash)

	root := a.MerkleRoot()
	if len(root) != 32 {
		t.Fatalf("len(MerkleRoot) = %d, want 32", len(root))
	}
	if !bytes.Equal(root, b.MerkleRoot()) {
		t.Fatal("identical statements have different merkle roots")
	}

	if err := a.Merge(b); err != nil {
		t.Fatalf("Merge: %v", err)
	}
	if !a.SignedBy(issuer) {
		t.Fatal("merged statement not signed by issuer")
	}

	other := signedStatement(t, acc, issuer, bytes.Repeat([]byte{0x44}, 32))
	if bytes.Equal(other.MerkleRoot(), root) {
		t.Fatal("different statements share a merkle root")
	}
	if err := a.Merge(other); err == nil {
		t.Fatal("Merge of a different statement: want error, got nil")
	}
}

func TestRevocationSigningRequestSummaryDescribesIssuer(t *testing.T) {
	acc := offlineAccount(t)

	issuer, err := acc.KeychainSigningCreate()
	if err != nil {
		t.Fatalf("KeychainSigningCreate: %v", err)
	}

	stmt := signedStatement(t, acc, issuer, bytes.Repeat([]byte{0x33}, 32))

	revReq, err := message.NewRevocationSigningRequest().
		Statement(stmt).
		Finish()
	if err != nil {
		t.Fatalf("revocation signing Finish: %v", err)
	}

	content, err := message.NewExchangeRequest().
		Purpose("co-sign revocation").
		Action(revReq.AsAction()).
		Finish()
	if err != nil {
		t.Fatalf("exchange Finish: %v", err)
	}

	summary, err := content.Summary()
	if err != nil {
		t.Fatalf("Summary: %v", err)
	}

	for _, d := range summary.Descriptions() {
		if d.Kind() != message.SummaryRevocation {
			continue
		}
		if got := d.Revocation(); !bytes.Equal(got, issuer.Bytes()) {
			t.Fatalf("Revocation = %x, want %x", got, issuer.Bytes())
		}
		return
	}

	t.Fatal("summary has no revocation description")
}

func TestConfigureWithEachStorageMemorySecurity(t *testing.T) {
	for _, security := range []account.StorageMemorySecurity{
		account.StorageMemorySecurityDefault,
		account.StorageMemorySecurityDisabled,
	} {
		acc, err := account.New(account.Config{
			StoragePath:           ":memory:",
			StorageMemorySecurity: security,
			Target: &account.Target{
				Network:         account.NetworkPreview,
				RPCEndpoint:     "http://127.0.0.1:1/",
				ObjectEndpoint:  "http://127.0.0.1:1/",
				MessageEndpoint: "ws://127.0.0.1:1/",
			},
		}, account.Callbacks{})
		if err != nil {
			t.Fatalf("account.New with storage memory security %d: %v", security, err)
		}
		acc.Close()
	}
}
