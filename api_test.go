package lucene

import (
	"errors"
	"testing"
)

func TestOp_Values(t *testing.T) {
	// Verify Op enum values are distinct.
	ops := []Op{
		OpMatch, OpMatchPhrase, OpMatchPhrasePrefix, OpMultiMatch,
		OpQueryString, OpSimpleQueryString,
		OpTerm, OpTerms, OpRange, OpPrefix, OpWildcard, OpRegexp,
		OpFuzzy, OpExists, OpIDs,
		OpBool, OpBoosting, OpConstantScore, OpDisMax,
		OpMatchAll, OpMatchNone, OpNested, OpHasChild, OpHasParent,
		OpKnn,
		OpGeoDistance, OpGeoBoundingBox,
	}

	seen := make(map[Op]bool)
	for _, op := range ops {
		if seen[op] {
			t.Errorf("duplicate Op value: %d", op)
		}
		seen[op] = true
	}
}

func TestQuery_Interface(t *testing.T) {
	q := &query{op: OpTerm, field: "status", value: "active"}

	if q.Op() != OpTerm {
		t.Errorf("Op() = %v, want %v", q.Op(), OpTerm)
	}

	if q.Err() != nil {
		t.Errorf("Err() = %v, want nil", q.Err())
	}
}

func TestQuery_WithError(t *testing.T) {
	err := errors.New("test error")
	q := errQuery(OpTerm, err)

	if q.Op() != OpTerm {
		t.Errorf("Op() = %v, want %v", q.Op(), OpTerm)
	}

	if !errors.Is(q.Err(), err) {
		t.Errorf("Err() = %v, want %v", q.Err(), err)
	}
}
