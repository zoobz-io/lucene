package lucene

// Highlight represents highlight configuration for search results.
type Highlight struct {
	fields       []HighlightField
	preTags      []string
	postTags     []string
	encoder      *string
	fragmentSize *int
	numFragments *int
	order        *string
	highlighter  *string
	err          error
}

// HighlightField represents a single field highlight configuration.
type HighlightField struct {
	Name             string
	FragmentSize     *int
	NumFragments     *int
	PreTags          []string
	PostTags         []string
	HighlightQuery   Query
	MatchedFields    []string
	FragmentOffset   *int
	NoMatchSize      *int
	RequireFieldMatch *bool
}

// NewHighlight creates a new highlight configuration.
func NewHighlight() *Highlight {
	return &Highlight{}
}

// Fields adds fields to highlight.
func (h *Highlight) Fields(names ...string) *Highlight {
	for _, name := range names {
		h.fields = append(h.fields, HighlightField{Name: name})
	}
	return h
}

// Field adds a single field with custom configuration.
func (h *Highlight) Field(f HighlightField) *Highlight {
	h.fields = append(h.fields, f)
	return h
}

// PreTags sets the pre-tags for highlighting.
func (h *Highlight) PreTags(tags ...string) *Highlight {
	h.preTags = tags
	return h
}

// PostTags sets the post-tags for highlighting.
func (h *Highlight) PostTags(tags ...string) *Highlight {
	h.postTags = tags
	return h
}

// Encoder sets the encoder (default or html).
func (h *Highlight) Encoder(e string) *Highlight {
	h.encoder = &e
	return h
}

// FragmentSize sets the size of fragments.
func (h *Highlight) FragmentSize(n int) *Highlight {
	h.fragmentSize = &n
	return h
}

// NumFragments sets the number of fragments.
func (h *Highlight) NumFragments(n int) *Highlight {
	h.numFragments = &n
	return h
}

// Order sets the fragment order (score or none).
func (h *Highlight) Order(o string) *Highlight {
	h.order = &o
	return h
}

// Highlighter sets the highlighter type (unified, plain, or fvh).
func (h *Highlight) Highlighter(t string) *Highlight {
	h.highlighter = &t
	return h
}

// FieldsValue returns the highlight fields.
func (h *Highlight) FieldsValue() []HighlightField { return h.fields }

// PreTagsValue returns the pre-tags.
func (h *Highlight) PreTagsValue() []string { return h.preTags }

// PostTagsValue returns the post-tags.
func (h *Highlight) PostTagsValue() []string { return h.postTags }

// EncoderValue returns the encoder if set.
func (h *Highlight) EncoderValue() *string { return h.encoder }

// FragmentSizeValue returns the fragment size if set.
func (h *Highlight) FragmentSizeValue() *int { return h.fragmentSize }

// NumFragmentsValue returns the number of fragments if set.
func (h *Highlight) NumFragmentsValue() *int { return h.numFragments }

// OrderValue returns the order if set.
func (h *Highlight) OrderValue() *string { return h.order }

// HighlighterValue returns the highlighter type if set.
func (h *Highlight) HighlighterValue() *string { return h.highlighter }

// Err returns any error in the highlight configuration.
func (h *Highlight) Err() error {
	if h.err != nil {
		return h.err
	}
	for _, f := range h.fields {
		if f.HighlightQuery != nil {
			if err := f.HighlightQuery.Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

// NewHighlightField creates a new highlight field configuration.
func NewHighlightField(name string) *HighlightFieldBuilder {
	return &HighlightFieldBuilder{field: HighlightField{Name: name}}
}

// HighlightFieldBuilder builds a highlight field configuration.
type HighlightFieldBuilder struct {
	field HighlightField
}

// FragmentSize sets the fragment size for this field.
func (b *HighlightFieldBuilder) FragmentSize(n int) *HighlightFieldBuilder {
	b.field.FragmentSize = &n
	return b
}

// NumFragments sets the number of fragments for this field.
func (b *HighlightFieldBuilder) NumFragments(n int) *HighlightFieldBuilder {
	b.field.NumFragments = &n
	return b
}

// PreTags sets the pre-tags for this field.
func (b *HighlightFieldBuilder) PreTags(tags ...string) *HighlightFieldBuilder {
	b.field.PreTags = tags
	return b
}

// PostTags sets the post-tags for this field.
func (b *HighlightFieldBuilder) PostTags(tags ...string) *HighlightFieldBuilder {
	b.field.PostTags = tags
	return b
}

// HighlightQuery sets a custom query for highlighting.
func (b *HighlightFieldBuilder) HighlightQuery(q Query) *HighlightFieldBuilder {
	b.field.HighlightQuery = q
	return b
}

// MatchedFields sets fields to combine for highlighting.
func (b *HighlightFieldBuilder) MatchedFields(fields ...string) *HighlightFieldBuilder {
	b.field.MatchedFields = fields
	return b
}

// FragmentOffset sets the offset for fragments.
func (b *HighlightFieldBuilder) FragmentOffset(n int) *HighlightFieldBuilder {
	b.field.FragmentOffset = &n
	return b
}

// NoMatchSize sets the text size to show when no match.
func (b *HighlightFieldBuilder) NoMatchSize(n int) *HighlightFieldBuilder {
	b.field.NoMatchSize = &n
	return b
}

// RequireFieldMatch sets whether to require field match.
func (b *HighlightFieldBuilder) RequireFieldMatch(v bool) *HighlightFieldBuilder {
	b.field.RequireFieldMatch = &v
	return b
}

// Build returns the configured HighlightField.
func (b *HighlightFieldBuilder) Build() HighlightField {
	return b.field
}
