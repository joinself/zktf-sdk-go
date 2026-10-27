package simffi

/*
#include <stdlib.h>
#include <zktf-sim.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type CriterionKind uint32

const (
	CriterionRequestID CriterionKind = iota
	CriterionContentType
	CriterionFrom
	CriterionFromIdentity
	CriterionTo
)

type Criterion struct {
	Kind        CriterionKind
	ContentType ContentType
	Bytes       []byte
}

func newMatcher(criteria []Criterion) (*C.zktf_sim_match, error) {
	matcher := C.zktf_sim_match_new()

	for _, criterion := range criteria {
		if err := addCriterion(matcher, criterion); err != nil {
			C.zktf_sim_match_destroy(matcher)
			return nil, err
		}
	}

	return matcher, nil
}

func addCriterion(matcher *C.zktf_sim_match, criterion Criterion) error {
	if criterion.Kind == CriterionContentType {
		return status(C.zktf_sim_match_content_type(matcher, C.enum_zktf_sim_content_type(criterion.ContentType)))
	}

	buf, length := cbytes(criterion.Bytes)
	defer free(unsafe.Pointer(buf))

	switch criterion.Kind {
	case CriterionRequestID:
		return status(C.zktf_sim_match_request_id(matcher, buf, length))
	case CriterionFrom:
		return status(C.zktf_sim_match_from(matcher, buf, length))
	case CriterionFromIdentity:
		return status(C.zktf_sim_match_from_identity(matcher, buf, length))
	case CriterionTo:
		return status(C.zktf_sim_match_to(matcher, buf, length))
	default:
		return &Status{code: uint32(C.ZKTF_SIM_INVALID_INPUT)}
	}
}

func (d *Device) ExpectMatch(criteria []Criterion, behaviour Behaviour, delayMs uint64) error {
	defer runtime.KeepAlive(d)

	matcher, err := newMatcher(criteria)
	if err != nil {
		return err
	}
	defer C.zktf_sim_match_destroy(matcher)

	C.zktf_sim_device_expect_match(d.ptr, matcher, C.enum_zktf_sim_behaviour(behaviour), C.uint64_t(delayMs))

	return nil
}

func (d *Device) InterceptMatch(criteria []Criterion) (*InterceptedFuture, error) {
	defer runtime.KeepAlive(d)

	matcher, err := newMatcher(criteria)
	if err != nil {
		return nil, err
	}
	defer C.zktf_sim_match_destroy(matcher)

	ptr := C.zktf_sim_device_intercept_match(d.ptr, matcher)
	if ptr == nil {
		return nil, &Status{code: uint32(C.ZKTF_SIM_NULL_ARGUMENT)}
	}

	return &InterceptedFuture{ptr: ptr}, nil
}

func (d *Device) PairwiseAddress(counterparty []byte) ([]byte, error) {
	defer runtime.KeepAlive(d)
	buf, length := cbytes(counterparty)
	defer free(unsafe.Pointer(buf))

	return keyBytes(func(out *C.uint8_t) C.enum_zktf_sim_status {
		return C.zktf_sim_device_pairwise_address(d.ptr, buf, length, out, signingKeyBytesLen)
	})
}
