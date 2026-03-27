package lucene

import "testing"

type testProduct struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	InStock     bool      `json:"in_stock"`
	Tags        []string  `json:"tags"`
	Embedding   []float32 `json:"embedding"`
	IgnoreField string    `json:"-"`
	NoTag       string
}

func TestNew(t *testing.T) {
	b := New[testProduct]()

	if b.spec == nil {
		t.Fatal("Builder.spec is nil")
	}

	if len(b.fields) == 0 {
		t.Fatal("Builder.fields is empty")
	}
}

func TestBuilder_Spec(t *testing.T) {
	b := New[testProduct]()

	spec := b.Spec()

	// Should have fields: id, name, price, in_stock, tags, embedding, NoTag.
	// IgnoreField should be excluded (json:"-").
	expectedFields := map[string]FieldKind{
		"id":        KindString,
		"name":      KindString,
		"price":     KindFloat,
		"in_stock":  KindBool,
		"tags":      KindSlice,
		"embedding": KindVector,
		"NoTag":     KindString,
	}

	if len(spec.Fields) != len(expectedFields) {
		t.Errorf("len(spec.Fields) = %d, want %d", len(spec.Fields), len(expectedFields))
	}

	for _, f := range spec.Fields {
		want, ok := expectedFields[f.Name]
		if !ok {
			t.Errorf("unexpected field: %s", f.Name)
			continue
		}
		if f.Kind != want {
			t.Errorf("field %s: Kind = %v, want %v", f.Name, f.Kind, want)
		}
	}
}

func TestBuilder_ResolveField(t *testing.T) {
	b := New[testProduct]()

	t.Run("existing field", func(t *testing.T) {
		spec, err := b.resolveField("name")
		if err != nil {
			t.Fatalf("resolveField(name) error = %v", err)
		}
		if spec.Kind != KindString {
			t.Errorf("spec.Kind = %v, want %v", spec.Kind, KindString)
		}
	})

	t.Run("unknown field", func(t *testing.T) {
		_, err := b.resolveField("nonexistent")
		if err == nil {
			t.Error("resolveField(nonexistent) should return error")
		}
	})
}

func TestBuilder_ValidateField(t *testing.T) {
	b := New[testProduct]()

	t.Run("valid field", func(t *testing.T) {
		spec, errQ := b.validateField(OpTerm, "name")
		if errQ != nil {
			t.Errorf("validateField returned error query: %v", errQ.Err())
		}
		if spec == nil {
			t.Error("validateField returned nil spec for valid field")
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		spec, errQ := b.validateField(OpTerm, "invalid")
		if spec != nil {
			t.Error("validateField should return nil spec for invalid field")
		}
		if errQ == nil {
			t.Error("validateField should return error query for invalid field")
		}
		if errQ.Err() == nil {
			t.Error("error query should have non-nil error")
		}
	})
}

func TestNew_NonStruct(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("New[string]() should panic for non-struct type")
		}
	}()
	_ = New[string]()
}
