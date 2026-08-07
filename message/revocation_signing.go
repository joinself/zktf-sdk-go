package message

import (
	"github.com/joinself/zktf-sdk-go/internal/ffi"
	"github.com/joinself/zktf-sdk-go/revocation"
)

// RevocationSigningRequest asks a counterparty device to co-sign a
// revocation statement.
type RevocationSigningRequest struct {
	h *ffi.RevocationSigningAction
}

// RevocationSigningRequestBuilder builds a revocation-signing request.
type RevocationSigningRequestBuilder struct {
	h *ffi.RevocationSigningActionBuilder
}

// RevocationSigningResponse is the response to a revocation-signing request.
type RevocationSigningResponse struct {
	h *ffi.RevocationSigningResult
}

// RevocationSigningResponseBuilder builds a revocation-signing response.
type RevocationSigningResponseBuilder struct {
	h *ffi.RevocationSigningResultBuilder
}

func init() {
	ffi.RevocationSigningActionOf = func(o any) *ffi.RevocationSigningAction {
		return o.(*RevocationSigningRequest).h
	}
	ffi.ToRevocationSigningAction = func(h *ffi.RevocationSigningAction) any {
		return &RevocationSigningRequest{h: h}
	}

	ffi.RevocationSigningResultOf = func(o any) *ffi.RevocationSigningResult {
		return o.(*RevocationSigningResponse).h
	}
	ffi.ToRevocationSigningResult = func(h *ffi.RevocationSigningResult) any {
		return &RevocationSigningResponse{h: h}
	}
}

// Statement returns the revocation statement this device is asked to co-sign,
// or nil if the statement could not be reconstructed.
func (r *RevocationSigningRequest) Statement() *revocation.Statement {
	s := r.h.Statement()
	if s == nil {
		return nil
	}

	return ffi.ToRevocationStatement(s).(*revocation.Statement)
}

// AsAction wraps this request into a generic Action.
func (r *RevocationSigningRequest) AsAction() *Action { return &Action{h: r.h.AsAction()} }

// NewRevocationSigningRequest starts building a revocation-signing request.
func NewRevocationSigningRequest() *RevocationSigningRequestBuilder {
	return &RevocationSigningRequestBuilder{h: ffi.NewRevocationSigningActionBuilder()}
}

// Statement sets the revocation statement to be co-signed.
func (b *RevocationSigningRequestBuilder) Statement(statement *revocation.Statement) *RevocationSigningRequestBuilder {
	b.h.Statement(ffi.RevocationStatementOf(statement))
	return b
}

// Finish finalizes the request.
func (b *RevocationSigningRequestBuilder) Finish() (*RevocationSigningRequest, error) {
	a, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &RevocationSigningRequest{h: a}, nil
}

// Statement returns the co-signed revocation statement, or nil if the
// statement could not be reconstructed.
func (r *RevocationSigningResponse) Statement() *revocation.Statement {
	s := r.h.Statement()
	if s == nil {
		return nil
	}

	return ffi.ToRevocationStatement(s).(*revocation.Statement)
}

// NewRevocationSigningResponse starts building a revocation-signing response.
func NewRevocationSigningResponse() *RevocationSigningResponseBuilder {
	return &RevocationSigningResponseBuilder{h: ffi.NewRevocationSigningResultBuilder()}
}

// Statement sets the co-signed revocation statement.
func (b *RevocationSigningResponseBuilder) Statement(statement *revocation.Statement) *RevocationSigningResponseBuilder {
	b.h.Statement(ffi.RevocationStatementOf(statement))
	return b
}

// Finish finalizes the response.
func (b *RevocationSigningResponseBuilder) Finish() (*RevocationSigningResponse, error) {
	r, err := b.h.Finish()
	if err != nil {
		return nil, err
	}

	return &RevocationSigningResponse{h: r}, nil
}
