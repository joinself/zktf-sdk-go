// Package predicate provides credential-predicate trees used by verifiers.
// Field paths are RFC 6901 JSON pointers.
package predicate

import "github.com/joinself/zktf-sdk-go/internal/ffi"

// PredicatorKind is the operator a predicator applies.
type PredicatorKind uint32

const (
	PredicatorEquals             PredicatorKind = PredicatorKind(ffi.PredicatorEquals)
	PredicatorNotEquals          PredicatorKind = PredicatorKind(ffi.PredicatorNotEquals)
	PredicatorGreaterThan        PredicatorKind = PredicatorKind(ffi.PredicatorGreaterThan)
	PredicatorGreaterThanOrEqual PredicatorKind = PredicatorKind(ffi.PredicatorGreaterThanOrEqual)
	PredicatorLessThan           PredicatorKind = PredicatorKind(ffi.PredicatorLessThan)
	PredicatorLessThanOrEqual    PredicatorKind = PredicatorKind(ffi.PredicatorLessThanOrEqual)
	PredicatorContains           PredicatorKind = PredicatorKind(ffi.PredicatorContains)
	PredicatorNotContains        PredicatorKind = PredicatorKind(ffi.PredicatorNotContains)
	PredicatorOneOf              PredicatorKind = PredicatorKind(ffi.PredicatorOneOf)
	PredicatorNotOneOf           PredicatorKind = PredicatorKind(ffi.PredicatorNotOneOf)
	PredicatorEmpty              PredicatorKind = PredicatorKind(ffi.PredicatorEmpty)
	PredicatorNotEmpty           PredicatorKind = PredicatorKind(ffi.PredicatorNotEmpty)
)

// Tree is a built predicate tree ready to evaluate against credentials.
type Tree struct {
	h *ffi.PredicateTree
}

func init() {
	ffi.PredicateTreeOf = func(o any) *ffi.PredicateTree { return o.(*Tree).h }
	ffi.ToPredicateTree = func(h *ffi.PredicateTree) any { return &Tree{h: h} }
}

// Predicate is a leaf or composite predicate node.
type Predicate struct {
	h *ffi.Predicate
}

// Equals checks if a credential field equals the given value.
func Equals(field, value string) *Predicate {
	return &Predicate{h: ffi.PredicateEquals(field, value)}
}

// NotEquals is the negation of Equals.
func NotEquals(field, value string) *Predicate {
	return &Predicate{h: ffi.PredicateNotEquals(field, value)}
}

// GreaterThan checks if a credential field is greater than value.
func GreaterThan(field, value string) *Predicate {
	return &Predicate{h: ffi.PredicateGreaterThan(field, value)}
}

// GreaterThanOrEquals checks if a credential field is >= value.
func GreaterThanOrEquals(field, value string) *Predicate {
	return &Predicate{h: ffi.PredicateGreaterThanOrEquals(field, value)}
}

// LessThan checks if a credential field is less than value.
func LessThan(field, value string) *Predicate {
	return &Predicate{h: ffi.PredicateLessThan(field, value)}
}

// LessThanOrEquals checks if a credential field is <= value.
func LessThanOrEquals(field, value string) *Predicate {
	return &Predicate{h: ffi.PredicateLessThanOrEquals(field, value)}
}

// Contains checks if a credential field contains value.
func Contains(field, value string) *Predicate {
	return &Predicate{h: ffi.PredicateContains(field, value)}
}

// NotContains is the negation of Contains.
func NotContains(field, value string) *Predicate {
	return &Predicate{h: ffi.PredicateNotContains(field, value)}
}

// OneOf checks if a credential field is one of the given values.
func OneOf(field string, values []string) *Predicate {
	return &Predicate{h: ffi.PredicateOneOf(field, values)}
}

// NotOneOf is the negation of OneOf.
func NotOneOf(field string, values []string) *Predicate {
	return &Predicate{h: ffi.PredicateNotOneOf(field, values)}
}

// Empty checks if a credential field is empty.
func Empty(field string) *Predicate {
	return &Predicate{h: ffi.PredicateEmpty(field)}
}

// NotEmpty checks if a credential field is not empty.
func NotEmpty(field string) *Predicate {
	return &Predicate{h: ffi.PredicateNotEmpty(field)}
}

// And combines two predicates with logical AND.
func And(a, b *Predicate) *Predicate {
	return &Predicate{h: ffi.PredicateAnd(a.h, b.h)}
}

// Or combines two predicates with logical OR.
func Or(a, b *Predicate) *Predicate {
	return &Predicate{h: ffi.PredicateOr(a.h, b.h)}
}

// NewTree builds a tree rooted at the given predicate.
func NewTree(root *Predicate) *Tree {
	return &Tree{h: ffi.NewPredicateTree(root.h)}
}

// DecodeTree decodes an encoded predicate tree.
func DecodeTree(data []byte) (*Tree, error) {
	t, err := ffi.PredicateTreeDecode(data)
	if err != nil {
		return nil, err
	}

	return &Tree{h: t}, nil
}

// Encode returns the encoded bytes of the tree.
func (t *Tree) Encode() []byte { return t.h.Encode() }

// Graphviz renders the tree in graphviz dot format.
func (t *Tree) Graphviz() string { return t.h.Graphviz() }

// Report describes which requirements remain unsatisfied by a set of credentials.
type Report struct {
	h *ffi.PredicateReport
}

// Requirements returns the per-requirement solutions in the report.
func (r *Report) Requirements() []*Solution {
	ss := r.h.Requirements()
	out := make([]*Solution, len(ss))

	for i, s := range ss {
		out[i] = &Solution{h: s}
	}

	return out
}

// Solution is one requirement's set of predicators that would satisfy it.
type Solution struct {
	h *ffi.PredicateSolution
}

// Predicators returns the predicators that satisfy this solution.
func (s *Solution) Predicators() []*Predicator {
	ps := s.h.Predicators()
	out := make([]*Predicator, len(ps))

	for i, p := range ps {
		out[i] = &Predicator{h: p}
	}

	return out
}

// Predicator is a field/op/values triple describing a needed predicate.
type Predicator struct {
	h *ffi.Predicator
}

// Kind returns the predicator's operator.
func (p *Predicator) Kind() PredicatorKind { return PredicatorKind(p.h.Kind()) }

// Field returns the credential field (JSON pointer) the predicator operates on.
func (p *Predicator) Field() string { return p.h.Field() }

// Values returns the predicator's value(s).
func (p *Predicator) Values() []string { return p.h.Values() }
