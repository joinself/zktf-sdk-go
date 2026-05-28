package ffi

/*
#include <zktf-sdk.h>
#include <stdlib.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// PredicatorKind mirrors zktf_credential_predicator_type.
type PredicatorKind uint32

const (
	PredicatorEquals             PredicatorKind = C.PREDICATOR_EQUALS
	PredicatorNotEquals          PredicatorKind = C.PREDICATOR_NOT_EQUALS
	PredicatorGreaterThan        PredicatorKind = C.PREDICATOR_GREATER_THAN
	PredicatorGreaterThanOrEqual PredicatorKind = C.PREDICATOR_GREATER_THAN_OR_EQUALS
	PredicatorLessThan           PredicatorKind = C.PREDICATOR_LESS_THAN
	PredicatorLessThanOrEqual    PredicatorKind = C.PREDICATOR_LESS_THAN_OR_EQUALS
	PredicatorContains           PredicatorKind = C.PREDICATOR_CONTAINS
	PredicatorNotContains        PredicatorKind = C.PREDICATOR_NOT_CONTAINS
	PredicatorOneOf              PredicatorKind = C.PREDICATOR_ONE_OF
	PredicatorNotOneOf           PredicatorKind = C.PREDICATOR_NOT_ONE_OF
	PredicatorEmpty              PredicatorKind = C.PREDICATOR_EMPTY
	PredicatorNotEmpty           PredicatorKind = C.PREDICATOR_NOT_EMPTY
)

// Predicate wraps a zktf_credential_predicate handle. Predicates are intermediate
// values combined via And/Or and ultimately rooted in a PredicateTree.
type Predicate struct {
	ptr *C.zktf_credential_predicate
}

func newPredicate(ptr *C.zktf_credential_predicate) *Predicate {
	if ptr == nil {
		return nil
	}
	p := &Predicate{ptr: ptr}
	runtime.AddCleanup(p, func(ptr *C.zktf_credential_predicate) {
		C.zktf_credential_predicate_destroy(ptr)
	}, p.ptr)
	return p
}

// fieldValuePredicate is a helper for the field/value binary predicates.
func fieldValuePredicate(
	c func(field, value *C.char) *C.zktf_credential_predicate,
	field, value string,
) *Predicate {
	cf, cv := cstring(field), cstring(value)
	defer free(unsafe.Pointer(cf))
	defer free(unsafe.Pointer(cv))
	return newPredicate(c(cf, cv))
}

// PredicateEquals checks if the field equals value. Field is an RFC 6901 JSON pointer.
func PredicateEquals(field, value string) *Predicate {
	return fieldValuePredicate(func(f, v *C.char) *C.zktf_credential_predicate {
		return C.zktf_credential_predicate_equals(f, v)
	}, field, value)
}

// PredicateNotEquals is the negation of PredicateEquals.
func PredicateNotEquals(field, value string) *Predicate {
	return fieldValuePredicate(func(f, v *C.char) *C.zktf_credential_predicate {
		return C.zktf_credential_predicate_not_equals(f, v)
	}, field, value)
}

// PredicateGreaterThan checks if field > value.
func PredicateGreaterThan(field, value string) *Predicate {
	return fieldValuePredicate(func(f, v *C.char) *C.zktf_credential_predicate {
		return C.zktf_credential_predicate_greater_than(f, v)
	}, field, value)
}

// PredicateGreaterThanOrEquals checks if field >= value.
func PredicateGreaterThanOrEquals(field, value string) *Predicate {
	return fieldValuePredicate(func(f, v *C.char) *C.zktf_credential_predicate {
		return C.zktf_credential_predicate_greater_than_or_equals(f, v)
	}, field, value)
}

// PredicateLessThan checks if field < value.
func PredicateLessThan(field, value string) *Predicate {
	return fieldValuePredicate(func(f, v *C.char) *C.zktf_credential_predicate {
		return C.zktf_credential_predicate_less_than(f, v)
	}, field, value)
}

// PredicateLessThanOrEquals checks if field <= value.
func PredicateLessThanOrEquals(field, value string) *Predicate {
	return fieldValuePredicate(func(f, v *C.char) *C.zktf_credential_predicate {
		return C.zktf_credential_predicate_less_than_or_equals(f, v)
	}, field, value)
}

// PredicateContains checks if field contains value.
func PredicateContains(field, value string) *Predicate {
	return fieldValuePredicate(func(f, v *C.char) *C.zktf_credential_predicate {
		return C.zktf_credential_predicate_contains(f, v)
	}, field, value)
}

// PredicateNotContains is the negation of PredicateContains.
func PredicateNotContains(field, value string) *Predicate {
	return fieldValuePredicate(func(f, v *C.char) *C.zktf_credential_predicate {
		return C.zktf_credential_predicate_not_contains(f, v)
	}, field, value)
}

// PredicateOneOf checks if field is one of the given values.
func PredicateOneOf(field string, values []string) *Predicate {
	cf := cstring(field)
	defer free(unsafe.Pointer(cf))

	collection := stringBufferCollection(values)
	defer C.zktf_collection_string_buffer_destroy(collection)

	return newPredicate(C.zktf_credential_predicate_one_of(cf, collection))
}

// PredicateNotOneOf is the negation of PredicateOneOf.
func PredicateNotOneOf(field string, values []string) *Predicate {
	cf := cstring(field)
	defer free(unsafe.Pointer(cf))

	collection := stringBufferCollection(values)
	defer C.zktf_collection_string_buffer_destroy(collection)

	return newPredicate(C.zktf_credential_predicate_not_one_of(cf, collection))
}

// PredicateEmpty checks if field is empty.
func PredicateEmpty(field string) *Predicate {
	cf := cstring(field)
	defer free(unsafe.Pointer(cf))
	return newPredicate(C.zktf_credential_predicate_empty(cf))
}

// PredicateNotEmpty checks if field is not empty.
func PredicateNotEmpty(field string) *Predicate {
	cf := cstring(field)
	defer free(unsafe.Pointer(cf))
	return newPredicate(C.zktf_credential_predicate_not_empty(cf))
}

// PredicateAnd combines two predicates with logical AND.
func PredicateAnd(a, b *Predicate) *Predicate {
	return newPredicate(C.zktf_credential_predicate_and(a.ptr, b.ptr))
}

// PredicateOr combines two predicates with logical OR.
func PredicateOr(a, b *Predicate) *Predicate {
	return newPredicate(C.zktf_credential_predicate_or(a.ptr, b.ptr))
}

// stringBufferCollection builds a zktf_collection_string_buffer from Go strings.
// Caller destroys the collection.
func stringBufferCollection(values []string) *C.zktf_collection_string_buffer {
	collection := C.zktf_collection_string_buffer_init()
	for _, v := range values {
		cv := cstring(v)
		C.zktf_collection_string_buffer_append(collection, cv)
		free(unsafe.Pointer(cv))
	}
	return collection
}

// PredicateTree is a built tree of predicates ready to evaluate against credentials.
type PredicateTree struct {
	ptr *C.zktf_credential_predicate_tree
}

func newPredicateTree(ptr *C.zktf_credential_predicate_tree) *PredicateTree {
	if ptr == nil {
		return nil
	}
	t := &PredicateTree{ptr: ptr}
	runtime.AddCleanup(t, func(ptr *C.zktf_credential_predicate_tree) {
		C.zktf_credential_predicate_tree_destroy(ptr)
	}, t.ptr)
	return t
}

// NewPredicateTree builds a tree rooted at the given predicate.
func NewPredicateTree(root *Predicate) *PredicateTree {
	return newPredicateTree(C.zktf_credential_predicate_tree_init(root.ptr))
}

// PredicateTreeDecode decodes an encoded predicate tree.
func PredicateTreeDecode(data []byte) (*PredicateTree, error) {
	buf, length := cbytes(data)
	defer free(unsafe.Pointer(buf))
	var out *C.zktf_credential_predicate_tree
	if err := status(C.zktf_credential_predicate_tree_decode(&out, buf, length)); err != nil {
		return nil, err
	}
	return newPredicateTree(out), nil
}

// Encode returns the encoded bytes of the tree.
func (t *PredicateTree) Encode() []byte {
	return goBytesFromBuffer(C.zktf_credential_predicate_tree_encode(t.ptr))
}

// Graphviz renders the tree in graphviz dot format.
func (t *PredicateTree) Graphviz() string {
	buf := C.zktf_credential_predicate_tree_graphviz(t.ptr)
	if buf == nil {
		return ""
	}
	defer C.zktf_string_buffer_destroy(buf)
	return C.GoString(C.zktf_string_buffer_ptr(buf))
}

// FindOptimalMatch selects the optimal set of credentials matching the tree,
// or returns nil if no match is possible.
func (t *PredicateTree) FindOptimalMatch(credentials []*VerifiableCredential) []*VerifiableCredential {
	in := verifiableCredentialCollection(credentials)
	defer C.zktf_collection_verifiable_credential_destroy(in)

	out := C.zktf_credential_predicate_tree_find_optimal_match(t.ptr, in)
	return verifiableCredentialsFrom(out)
}

// FindMissingPredicates returns a report of predicates the credentials do not satisfy.
func (t *PredicateTree) FindMissingPredicates(credentials []*VerifiableCredential) *PredicateReport {
	in := verifiableCredentialCollection(credentials)
	defer C.zktf_collection_verifiable_credential_destroy(in)
	return newPredicateReport(C.zktf_credential_predicate_tree_find_missing_predicates(t.ptr, in))
}

// verifiableCredentialCollection builds a zktf_collection_verifiable_credential
// from Go credentials. Caller destroys.
func verifiableCredentialCollection(credentials []*VerifiableCredential) *C.zktf_collection_verifiable_credential {
	c := C.zktf_collection_verifiable_credential_init()
	for _, vc := range credentials {
		C.zktf_collection_verifiable_credential_append(c, vc.ptr)
	}
	return c
}

// PredicateReport describes which requirements (predicate solutions) remain unsatisfied.
type PredicateReport struct {
	ptr *C.zktf_credential_predicate_report
}

func newPredicateReport(ptr *C.zktf_credential_predicate_report) *PredicateReport {
	if ptr == nil {
		return nil
	}
	r := &PredicateReport{ptr: ptr}
	runtime.AddCleanup(r, func(ptr *C.zktf_credential_predicate_report) {
		C.zktf_credential_predicate_report_destroy(ptr)
	}, r.ptr)
	return r
}

// Requirements returns the per-requirement solutions in the report.
func (r *PredicateReport) Requirements() []*PredicateSolution {
	n := int(C.zktf_credential_predicate_report_requirements_len(r.ptr))
	out := make([]*PredicateSolution, n)
	for i := 0; i < n; i++ {
		out[i] = newPredicateSolution(C.zktf_credential_predicate_report_requirements_at(r.ptr, C.size_t(i)))
	}
	return out
}

// PredicateSolution is the per-requirement set of predicators that would satisfy it.
type PredicateSolution struct {
	ptr *C.zktf_credential_predicate_solution
}

func newPredicateSolution(ptr *C.zktf_credential_predicate_solution) *PredicateSolution {
	if ptr == nil {
		return nil
	}
	s := &PredicateSolution{ptr: ptr}
	runtime.AddCleanup(s, func(ptr *C.zktf_credential_predicate_solution) {
		C.zktf_credential_predicate_solution_destroy(ptr)
	}, s.ptr)
	return s
}

// Predicators returns the predicators required to satisfy this solution.
func (s *PredicateSolution) Predicators() []*Predicator {
	n := int(C.zktf_credential_predicate_solution_len(s.ptr))
	out := make([]*Predicator, n)
	for i := 0; i < n; i++ {
		c := C.zktf_credential_predicate_solution_at(s.ptr, C.size_t(i))
		if c == nil {
			continue
		}
		// Each at() returns a per-credential collection of predicators.
		predicators := predicatorsFrom(c)
		out = append(out[:0], predicators...)
	}
	return out
}

// predicatorsFrom copies a zktf_collection_credential_predicator into Go wrappers.
func predicatorsFrom(c *C.zktf_collection_credential_predicator) []*Predicator {
	defer C.zktf_collection_credential_predicator_destroy(c)
	n := int(C.zktf_collection_credential_predicator_len(c))
	out := make([]*Predicator, n)
	for i := 0; i < n; i++ {
		out[i] = newPredicator(C.zktf_collection_credential_predicator_at(c, C.size_t(i)))
	}
	return out
}

// Predicator is a single field/op/values triple describing a needed predicate.
type Predicator struct {
	ptr *C.zktf_credential_predicator
}

func newPredicator(ptr *C.zktf_credential_predicator) *Predicator {
	if ptr == nil {
		return nil
	}
	p := &Predicator{ptr: ptr}
	runtime.AddCleanup(p, func(ptr *C.zktf_credential_predicator) {
		C.zktf_credential_predicator_destroy(ptr)
	}, p.ptr)
	return p
}

// Kind returns the predicator's operator.
func (p *Predicator) Kind() PredicatorKind {
	return PredicatorKind(C.zktf_credential_predicator_predicator_type(p.ptr))
}

// Field returns the credential field (JSON pointer) the predicator operates on.
func (p *Predicator) Field() string {
	buf := C.zktf_credential_predicator_field(p.ptr)
	if buf == nil {
		return ""
	}
	defer C.zktf_string_buffer_destroy(buf)
	return C.GoString(C.zktf_string_buffer_ptr(buf))
}

// Values returns the predicator's value(s).
func (p *Predicator) Values() []string {
	c := C.zktf_credential_predicator_values(p.ptr)
	if c == nil {
		return nil
	}
	defer C.zktf_collection_string_buffer_destroy(c)
	n := int(C.zktf_collection_string_buffer_len(c))
	out := make([]string, n)
	for i := 0; i < n; i++ {
		buf := C.zktf_collection_string_buffer_at(c, C.size_t(i))
		if buf == nil {
			continue
		}
		out[i] = C.GoString(C.zktf_string_buffer_ptr(buf))
	}
	return out
}
