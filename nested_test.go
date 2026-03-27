package lucene

import "testing"

type nestedTestDoc struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestNestedQuery(t *testing.T) {
	b := New[nestedTestDoc]()

	t.Run("basic nested", func(t *testing.T) {
		inner := b.Match("name", "john")
		q := b.Nested("comments", inner)

		if q.Op() != OpNested {
			t.Errorf("Op() = %v, want OpNested", q.Op())
		}
		if q.Path() != "comments" {
			t.Errorf("Path() = %v, want comments", q.Path())
		}
		if q.InnerQuery() != inner {
			t.Error("InnerQuery() should return the inner query")
		}
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
	})

	t.Run("with options", func(t *testing.T) {
		inner := b.Match("name", "john")
		q := b.Nested("comments", inner).
			ScoreMode("avg").
			IgnoreUnmapped(true)

		if q.ScoreModeValue() == nil || *q.ScoreModeValue() != "avg" {
			t.Errorf("ScoreModeValue() = %v, want avg", q.ScoreModeValue())
		}
		if q.IgnoreUnmappedValue() == nil || *q.IgnoreUnmappedValue() != true {
			t.Errorf("IgnoreUnmappedValue() = %v, want true", q.IgnoreUnmappedValue())
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		inner := b.Match("invalid_field", "john")
		q := b.Nested("comments", inner)

		if q.Err() == nil {
			t.Error("Err() should propagate inner query error")
		}
	})
}

func TestHasChildQuery(t *testing.T) {
	b := New[nestedTestDoc]()

	t.Run("basic has_child", func(t *testing.T) {
		inner := b.Term("id", "child1")
		q := b.HasChild("child_type", inner)

		if q.Op() != OpHasChild {
			t.Errorf("Op() = %v, want OpHasChild", q.Op())
		}
		if q.ChildType() != "child_type" {
			t.Errorf("ChildType() = %v, want child_type", q.ChildType())
		}
		if q.InnerQuery() != inner {
			t.Error("InnerQuery() should return the inner query")
		}
	})

	t.Run("with options", func(t *testing.T) {
		inner := b.MatchAll()
		q := b.HasChild("child_type", inner).
			ScoreMode("max").
			MinChildren(1).
			MaxChildren(10).
			IgnoreUnmapped(true)

		if q.ScoreModeValue() == nil || *q.ScoreModeValue() != "max" {
			t.Errorf("ScoreModeValue() = %v, want max", q.ScoreModeValue())
		}
		if q.MinChildrenValue() == nil || *q.MinChildrenValue() != 1 {
			t.Errorf("MinChildrenValue() = %v, want 1", q.MinChildrenValue())
		}
		if q.MaxChildrenValue() == nil || *q.MaxChildrenValue() != 10 {
			t.Errorf("MaxChildrenValue() = %v, want 10", q.MaxChildrenValue())
		}
		if q.IgnoreUnmappedValue() == nil || *q.IgnoreUnmappedValue() != true {
			t.Errorf("IgnoreUnmappedValue() = %v, want true", q.IgnoreUnmappedValue())
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		inner := b.Match("invalid_field", "test")
		q := b.HasChild("child_type", inner)

		if q.Err() == nil {
			t.Error("Err() should propagate inner query error")
		}
	})
}

func TestHasParentQuery(t *testing.T) {
	b := New[nestedTestDoc]()

	t.Run("basic has_parent", func(t *testing.T) {
		inner := b.Term("id", "parent1")
		q := b.HasParent("parent_type", inner)

		if q.Op() != OpHasParent {
			t.Errorf("Op() = %v, want OpHasParent", q.Op())
		}
		if q.ParentType() != "parent_type" {
			t.Errorf("ParentType() = %v, want parent_type", q.ParentType())
		}
		if q.InnerQuery() != inner {
			t.Error("InnerQuery() should return the inner query")
		}
	})

	t.Run("with options", func(t *testing.T) {
		inner := b.MatchAll()
		q := b.HasParent("parent_type", inner).
			Score(true).
			IgnoreUnmapped(true)

		if q.ScoreValue() == nil || *q.ScoreValue() != true {
			t.Errorf("ScoreValue() = %v, want true", q.ScoreValue())
		}
		if q.IgnoreUnmappedValue() == nil || *q.IgnoreUnmappedValue() != true {
			t.Errorf("IgnoreUnmappedValue() = %v, want true", q.IgnoreUnmappedValue())
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		inner := b.Match("invalid_field", "test")
		q := b.HasParent("parent_type", inner)

		if q.Err() == nil {
			t.Error("Err() should propagate inner query error")
		}
	})
}
