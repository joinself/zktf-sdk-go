package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import "time"

// MessageSend sends content to the given recipient address. This call returns
// once the message has been queued locally; delivery is reported via the
// on_status callback (acknowledged / send-failed).
func (a *Account) MessageSend(to *SigningPublicKey, content *Content) error {
	return status(C.zktf_account_message_send(a.ptr, to.ptr, content.ptr))
}

// InboxOpen opens a new messaging inbox, awaiting its address via callback.
func (a *Account) InboxOpen(timeout time.Duration) (*SigningPublicKey, error) {
	fut := C.zktf_account_inbox_open(a.ptr, nil)

	return AwaitSigningPublicKey(fut, timeout)
}

// InboxDefault returns the account's default inbox address synchronously.
func (a *Account) InboxDefault() (*SigningPublicKey, error) {
	var out *C.zktf_signing_public_key

	if err := status(C.zktf_account_inbox_default(a.ptr, &out)); err != nil {
		return nil, err
	}

	return newSigningPublicKey(out), nil
}

// InboxClose closes an open inbox, awaiting completion via callback.
func (a *Account) InboxClose(address *SigningPublicKey, timeout time.Duration) error {
	fut := C.zktf_account_inbox_close(a.ptr, address.ptr)

	return AwaitStatus(fut, timeout)
}

// InboxList returns the addresses of all open inboxes on this account.
func (a *Account) InboxList() ([]*SigningPublicKey, error) {
	var c *C.zktf_collection_signing_public_key

	if err := status(C.zktf_account_inbox_list(a.ptr, &c)); err != nil {
		return nil, err
	}

	return signingPublicKeysFrom(c), nil
}

// GroupNegotiate negotiates an encrypted session between two inbox addresses.
// The SDK auto-accepts the resulting invite/welcome on both sides. expiresUnix
// of 0 means no expiry.
func (a *Account) GroupNegotiate(as, with *SigningPublicKey, expiresUnix int64) error {
	return status(C.zktf_account_group_negotiate(a.ptr, as.ptr, with.ptr, C.int64_t(expiresUnix)))
}
