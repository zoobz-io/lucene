package lucene

import "testing"

func TestNewHighlight(t *testing.T) {
	h := NewHighlight()
	if h == nil {
		t.Fatal("NewHighlight() returned nil")
	}
}

func TestHighlight_Fields(t *testing.T) {
	h := NewHighlight().Fields("title", "description", "content")

	fields := h.FieldsValue()
	if len(fields) != 3 {
		t.Fatalf("FieldsValue() len = %d, want 3", len(fields))
	}
	if fields[0].Name != "title" {
		t.Errorf("FieldsValue()[0].Name = %s, want title", fields[0].Name)
	}
}

func TestHighlight_PrePostTags(t *testing.T) {
	h := NewHighlight().
		PreTags("<em>", "<strong>").
		PostTags("</em>", "</strong>")

	preTags := h.PreTagsValue()
	if len(preTags) != 2 || preTags[0] != "<em>" {
		t.Errorf("PreTagsValue() = %v, want [<em> <strong>]", preTags)
	}

	postTags := h.PostTagsValue()
	if len(postTags) != 2 || postTags[0] != "</em>" {
		t.Errorf("PostTagsValue() = %v, want [</em> </strong>]", postTags)
	}
}

func TestHighlight_Encoder(t *testing.T) {
	h := NewHighlight().Encoder("html")

	if h.EncoderValue() == nil || *h.EncoderValue() != "html" {
		t.Errorf("EncoderValue() = %v, want html", h.EncoderValue())
	}
}

func TestHighlight_FragmentSize(t *testing.T) {
	h := NewHighlight().FragmentSize(150)

	if h.FragmentSizeValue() == nil || *h.FragmentSizeValue() != 150 {
		t.Errorf("FragmentSizeValue() = %v, want 150", h.FragmentSizeValue())
	}
}

func TestHighlight_NumFragments(t *testing.T) {
	h := NewHighlight().NumFragments(5)

	if h.NumFragmentsValue() == nil || *h.NumFragmentsValue() != 5 {
		t.Errorf("NumFragmentsValue() = %v, want 5", h.NumFragmentsValue())
	}
}

func TestHighlight_Order(t *testing.T) {
	h := NewHighlight().Order("score")

	if h.OrderValue() == nil || *h.OrderValue() != "score" {
		t.Errorf("OrderValue() = %v, want score", h.OrderValue())
	}
}

func TestHighlight_Highlighter(t *testing.T) {
	h := NewHighlight().Highlighter("unified")

	if h.HighlighterValue() == nil || *h.HighlighterValue() != "unified" {
		t.Errorf("HighlighterValue() = %v, want unified", h.HighlighterValue())
	}
}

func TestHighlight_Field(t *testing.T) {
	field := NewHighlightField("title").
		FragmentSize(100).
		NumFragments(3).
		PreTags("<b>").
		PostTags("</b>").
		Build()

	h := NewHighlight().Field(field)

	fields := h.FieldsValue()
	if len(fields) != 1 {
		t.Fatalf("FieldsValue() len = %d, want 1", len(fields))
	}
	if fields[0].Name != "title" {
		t.Error("Field name should be title")
	}
	if fields[0].FragmentSize == nil || *fields[0].FragmentSize != 100 {
		t.Error("Field fragment_size should be 100")
	}
	if fields[0].NumFragments == nil || *fields[0].NumFragments != 3 {
		t.Error("Field number_of_fragments should be 3")
	}
}

func TestHighlightFieldBuilder(t *testing.T) {
	b := New[searchTestDoc]()

	field := NewHighlightField("title").
		FragmentSize(100).
		NumFragments(3).
		PreTags("<b>").
		PostTags("</b>").
		MatchedFields("title", "title.plain").
		FragmentOffset(10).
		NoMatchSize(50).
		RequireFieldMatch(false).
		HighlightQuery(b.Match("title", "search")).
		Build()

	if field.Name != "title" {
		t.Error("Name should be title")
	}
	if field.FragmentSize == nil || *field.FragmentSize != 100 {
		t.Error("FragmentSize should be 100")
	}
	if field.NumFragments == nil || *field.NumFragments != 3 {
		t.Error("NumFragments should be 3")
	}
	if len(field.PreTags) != 1 || field.PreTags[0] != "<b>" {
		t.Error("PreTags should be [<b>]")
	}
	if len(field.PostTags) != 1 || field.PostTags[0] != "</b>" {
		t.Error("PostTags should be [</b>]")
	}
	if len(field.MatchedFields) != 2 {
		t.Error("MatchedFields should have 2 fields")
	}
	if field.FragmentOffset == nil || *field.FragmentOffset != 10 {
		t.Error("FragmentOffset should be 10")
	}
	if field.NoMatchSize == nil || *field.NoMatchSize != 50 {
		t.Error("NoMatchSize should be 50")
	}
	if field.RequireFieldMatch == nil || *field.RequireFieldMatch != false {
		t.Error("RequireFieldMatch should be false")
	}
	if field.HighlightQuery == nil {
		t.Error("HighlightQuery should not be nil")
	}
}

func TestHighlight_Err(t *testing.T) {
	h := NewHighlight().Fields("title")

	if h.Err() != nil {
		t.Errorf("Err() = %v, want nil", h.Err())
	}
}

func TestHighlight_Err_InvalidQuery(t *testing.T) {
	b := New[searchTestDoc]()

	invalidQuery := b.Match("invalid_field", "test")
	field := NewHighlightField("title").
		HighlightQuery(invalidQuery).
		Build()

	h := NewHighlight().Field(field)

	if h.Err() == nil {
		t.Error("Err() should not be nil for invalid highlight query")
	}
}

func TestHighlight_Chained(t *testing.T) {
	h := NewHighlight().
		Fields("title", "description").
		PreTags("<em>").
		PostTags("</em>").
		Encoder("html").
		FragmentSize(150).
		NumFragments(5).
		Order("score").
		Highlighter("unified")

	if h.Err() != nil {
		t.Errorf("Err() = %v, want nil", h.Err())
	}
	if len(h.FieldsValue()) != 2 {
		t.Error("Should have 2 fields")
	}
	if len(h.PreTagsValue()) != 1 {
		t.Error("Should have 1 pre tag")
	}
	if len(h.PostTagsValue()) != 1 {
		t.Error("Should have 1 post tag")
	}
}
