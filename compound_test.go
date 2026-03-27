package lucene

import "testing"

type compoundTestDoc struct {
	Name   string  `json:"name"`
	Status string  `json:"status"`
	Price  float64 `json:"price"`
}

func TestBoolQuery(t *testing.T) {
	b := New[compoundTestDoc]()

	t.Run("empty bool", func(t *testing.T) {
		q := b.Bool()
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
		if q.Op() != OpBool {
			t.Errorf("Op() = %v, want %v", q.Op(), OpBool)
		}
	})

	t.Run("with must clauses", func(t *testing.T) {
		q := b.Bool().
			Must(b.Term("status", "active")).
			Must(b.Range("price").Gte(10))

		if len(q.MustClauses()) != 2 {
			t.Errorf("len(MustClauses()) = %d, want 2", len(q.MustClauses()))
		}
	})

	t.Run("with all clause types", func(t *testing.T) {
		q := b.Bool().
			Must(b.Term("status", "active")).
			Should(b.Match("name", "test")).
			MustNot(b.Term("status", "deleted")).
			Filter(b.Exists("price")).
			MinimumShouldMatch(1)

		if len(q.MustClauses()) != 1 {
			t.Errorf("len(MustClauses()) = %d, want 1", len(q.MustClauses()))
		}
		if len(q.ShouldClauses()) != 1 {
			t.Errorf("len(ShouldClauses()) = %d, want 1", len(q.ShouldClauses()))
		}
		if len(q.MustNotClauses()) != 1 {
			t.Errorf("len(MustNotClauses()) = %d, want 1", len(q.MustNotClauses()))
		}
		if len(q.FilterClauses()) != 1 {
			t.Errorf("len(FilterClauses()) = %d, want 1", len(q.FilterClauses()))
		}
		if q.MinimumShouldMatchValue() == nil || *q.MinimumShouldMatchValue() != 1 {
			t.Errorf("MinimumShouldMatchValue() = %v, want 1", q.MinimumShouldMatchValue())
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		q := b.Bool().Must(b.Term("invalid", "value"))
		if q.Err() == nil {
			t.Error("Err() should propagate child errors")
		}
	})
}

func TestMatchAllQuery(t *testing.T) {
	b := New[compoundTestDoc]()

	t.Run("basic", func(t *testing.T) {
		q := b.MatchAll()
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
		if q.Op() != OpMatchAll {
			t.Errorf("Op() = %v, want %v", q.Op(), OpMatchAll)
		}
	})

	t.Run("with boost", func(t *testing.T) {
		q := b.MatchAll().Boost(1.5)
		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})
}

func TestMatchNoneQuery(t *testing.T) {
	b := New[compoundTestDoc]()

	q := b.MatchNone()
	if q.Err() != nil {
		t.Errorf("Err() = %v, want nil", q.Err())
	}
	if q.Op() != OpMatchNone {
		t.Errorf("Op() = %v, want %v", q.Op(), OpMatchNone)
	}
}

func TestConvenienceMethods(t *testing.T) {
	b := New[compoundTestDoc]()

	t.Run("And", func(t *testing.T) {
		q := b.And(
			b.Term("status", "active"),
			b.Exists("price"),
		)
		if len(q.MustClauses()) != 2 {
			t.Errorf("len(MustClauses()) = %d, want 2", len(q.MustClauses()))
		}
	})

	t.Run("Or", func(t *testing.T) {
		q := b.Or(
			b.Term("status", "active"),
			b.Term("status", "pending"),
		)
		if len(q.ShouldClauses()) != 2 {
			t.Errorf("len(ShouldClauses()) = %d, want 2", len(q.ShouldClauses()))
		}
		if q.MinimumShouldMatchValue() == nil || *q.MinimumShouldMatchValue() != 1 {
			t.Errorf("MinimumShouldMatchValue() = %v, want 1", q.MinimumShouldMatchValue())
		}
	})

	t.Run("Not", func(t *testing.T) {
		q := b.Not(b.Term("status", "deleted"))
		if len(q.MustNotClauses()) != 1 {
			t.Errorf("len(MustNotClauses()) = %d, want 1", len(q.MustNotClauses()))
		}
	})

	t.Run("Not nil", func(t *testing.T) {
		q := b.Not(nil)
		if q.Err() == nil {
			t.Error("Not(nil) should return error")
		}
	})
}

func TestBoostingQuery(t *testing.T) {
	b := New[compoundTestDoc]()

	t.Run("basic boosting", func(t *testing.T) {
		positive := b.Match("name", "search")
		negative := b.Term("status", "deprecated")
		q := b.Boosting().
			Positive(positive).
			Negative(negative).
			NegativeBoost(0.5)

		if q.Op() != OpBoosting {
			t.Errorf("Op() = %v, want OpBoosting", q.Op())
		}
		if q.PositiveQuery() != positive {
			t.Error("PositiveQuery() should return the positive query")
		}
		if q.NegativeQuery() != negative {
			t.Error("NegativeQuery() should return the negative query")
		}
		if q.NegativeBoostValue() == nil || *q.NegativeBoostValue() != 0.5 {
			t.Errorf("NegativeBoostValue() = %v, want 0.5", q.NegativeBoostValue())
		}
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
	})

	t.Run("error propagation from positive", func(t *testing.T) {
		q := b.Boosting().
			Positive(b.Term("invalid", "value")).
			Negative(b.MatchAll())

		if q.Err() == nil {
			t.Error("Err() should propagate positive query error")
		}
	})

	t.Run("error propagation from negative", func(t *testing.T) {
		q := b.Boosting().
			Positive(b.MatchAll()).
			Negative(b.Term("invalid", "value"))

		if q.Err() == nil {
			t.Error("Err() should propagate negative query error")
		}
	})
}

func TestDisMaxQuery(t *testing.T) {
	b := New[compoundTestDoc]()

	t.Run("basic dis_max", func(t *testing.T) {
		q := b.DisMax(
			b.Match("name", "search"),
			b.Term("status", "active"),
		)

		if q.Op() != OpDisMax {
			t.Errorf("Op() = %v, want OpDisMax", q.Op())
		}
		if len(q.Queries()) != 2 {
			t.Errorf("len(Queries()) = %d, want 2", len(q.Queries()))
		}
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
	})

	t.Run("with options", func(t *testing.T) {
		q := b.DisMax(
			b.Match("name", "search"),
			b.Term("status", "active"),
		).TieBreaker(0.3).Boost(1.5)

		if q.TieBreakerValue() == nil || *q.TieBreakerValue() != 0.3 {
			t.Errorf("TieBreakerValue() = %v, want 0.3", q.TieBreakerValue())
		}
		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		q := b.DisMax(
			b.Match("name", "search"),
			b.Term("invalid", "value"),
		)

		if q.Err() == nil {
			t.Error("Err() should propagate child query error")
		}
	})
}

func TestConstantScoreQuery(t *testing.T) {
	b := New[compoundTestDoc]()

	t.Run("basic constant_score", func(t *testing.T) {
		filter := b.Term("status", "active")
		q := b.ConstantScore(filter)

		if q.Op() != OpConstantScore {
			t.Errorf("Op() = %v, want OpConstantScore", q.Op())
		}
		if q.FilterQuery() != filter {
			t.Error("FilterQuery() should return the filter query")
		}
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
	})

	t.Run("with boost", func(t *testing.T) {
		q := b.ConstantScore(b.Term("status", "active")).Boost(1.5)

		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		q := b.ConstantScore(b.Term("invalid", "value"))

		if q.Err() == nil {
			t.Error("Err() should propagate filter query error")
		}
	})
}
