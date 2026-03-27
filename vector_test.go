package lucene

import "testing"

type vectorTestDoc struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Embedding []float32 `json:"embedding"`
}

func TestKnnQuery(t *testing.T) {
	b := New[vectorTestDoc]()

	t.Run("basic knn", func(t *testing.T) {
		vector := []float32{0.1, 0.2, 0.3}
		q := b.Knn("embedding", vector)

		if q.Op() != OpKnn {
			t.Errorf("Op() = %v, want OpKnn", q.Op())
		}
		if q.Field() != "embedding" {
			t.Errorf("Field() = %v, want embedding", q.Field())
		}
		if len(q.Vector()) != 3 {
			t.Errorf("len(Vector()) = %d, want 3", len(q.Vector()))
		}
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
	})

	t.Run("with options", func(t *testing.T) {
		vector := []float32{0.1, 0.2, 0.3}
		q := b.Knn("embedding", vector).
			K(10).
			NumCandidates(100).
			Boost(1.5)

		if q.KValue() == nil || *q.KValue() != 10 {
			t.Errorf("KValue() = %v, want 10", q.KValue())
		}
		if q.NumCandidatesValue() == nil || *q.NumCandidatesValue() != 100 {
			t.Errorf("NumCandidatesValue() = %v, want 100", q.NumCandidatesValue())
		}
		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})

	t.Run("with filter", func(t *testing.T) {
		vector := []float32{0.1, 0.2, 0.3}
		filter := b.Term("title", "test")
		q := b.Knn("embedding", vector).K(10).Filter(filter)

		if q.FilterQuery() != filter {
			t.Error("FilterQuery() should return the filter")
		}
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
	})

	t.Run("filter error propagation", func(t *testing.T) {
		vector := []float32{0.1, 0.2, 0.3}
		filter := b.Term("invalid_field", "test")
		q := b.Knn("embedding", vector).Filter(filter)

		if q.Err() == nil {
			t.Error("Err() should propagate filter error")
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		vector := []float32{0.1, 0.2, 0.3}
		q := b.Knn("invalid_field", vector)

		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}
