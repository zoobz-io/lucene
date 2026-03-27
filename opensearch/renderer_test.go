package opensearch

import (
	"encoding/json"
	"testing"

	"github.com/zoobz-io/lucene"
)

type testDoc struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Embedding []float32 `json:"embedding"`
}

func mustBuilder(t *testing.T) *lucene.Builder[testDoc] {
	t.Helper()
	return lucene.New[testDoc]()
}

func TestRenderer_RenderQuery_Term(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V2)

	q := b.Term("name", "test")
	out, err := r.RenderQuery(q)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if _, ok := result["term"]; !ok {
		t.Error("expected term key")
	}
}

func TestRenderer_RenderQuery_Knn(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V2)

	vector := []float32{0.1, 0.2, 0.3}
	q := b.Knn("embedding", vector).K(10)
	out, err := r.RenderQuery(q)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	knn, ok := result["knn"].(map[string]any)
	if !ok {
		t.Fatal("expected knn key")
	}

	// OpenSearch puts the field name as a key under knn.
	embedding, ok := knn["embedding"].(map[string]any)
	if !ok {
		t.Fatal("expected embedding key under knn")
	}

	if embedding["k"] != float64(10) {
		t.Errorf("k = %v, want 10", embedding["k"])
	}
}

func TestRenderer_RenderQuery_Hybrid(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V2)

	vector := []float32{0.1, 0.2, 0.3}
	q := b.Hybrid(
		b.Match("name", "search term"),
		b.Knn("embedding", vector).K(10),
	)
	out, err := r.RenderQuery(q)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	hybrid, ok := result["hybrid"].(map[string]any)
	if !ok {
		t.Fatal("expected hybrid key")
	}

	queries, ok := hybrid["queries"].([]any)
	if !ok {
		t.Fatal("expected queries array under hybrid")
	}

	if len(queries) != 2 {
		t.Errorf("len(queries) = %d, want 2", len(queries))
	}
}

func TestRenderer_RenderQuery_Bool(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V2)

	q := b.Bool().Must(b.Term("name", "test"))
	out, err := r.RenderQuery(q)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	boolQ, ok := result["bool"].(map[string]any)
	if !ok {
		t.Fatal("expected bool key")
	}

	must, ok := boolQ["must"].([]any)
	if !ok {
		t.Fatal("expected must array")
	}

	if len(must) != 1 {
		t.Errorf("len(must) = %d, want 1", len(must))
	}
}

func TestRenderer_Render_FullSearch(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V2)

	h := lucene.NewHighlight().
		Fields("name").
		PreTags("<em>").
		PostTags("</em>")

	s := lucene.NewSearch().
		Query(b.Match("name", "test")).
		Size(20).
		From(5).
		Sort(lucene.SortField{Field: "name", Order: "asc"}).
		Source("name").
		Highlight(h).
		TrackTotalHits(true).
		MinScore(0.5).
		Timeout("10s")

	out, err := r.Render(s)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if _, ok := result["query"]; !ok {
		t.Error("expected query key")
	}
	if result["size"] != float64(20) {
		t.Errorf("size = %v, want 20", result["size"])
	}
	if result["from"] != float64(5) {
		t.Errorf("from = %v, want 5", result["from"])
	}
	if _, ok := result["sort"]; !ok {
		t.Error("expected sort key")
	}
	if _, ok := result["_source"]; !ok {
		t.Error("expected _source key")
	}
	if _, ok := result["highlight"]; !ok {
		t.Error("expected highlight key")
	}
	if result["track_total_hits"] != true {
		t.Errorf("track_total_hits = %v, want true", result["track_total_hits"])
	}
}

func TestRenderer_Render_NilSearch(t *testing.T) {
	r := NewRenderer(V2)
	_, err := r.Render(nil)
	if err == nil {
		t.Error("Render(nil) should return error")
	}
}

func TestRenderer_Render_SearchWithError(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V2)

	s := lucene.NewSearch().Query(b.Match("invalid", "test"))
	_, err := r.Render(s)
	if err == nil {
		t.Error("Render() should return error for search with invalid query")
	}
}

func TestRenderer_RenderAggs(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V2)

	aggs := []lucene.Aggregation{
		b.TermsAgg("by_name", "name").Size(10),
	}

	out, err := r.RenderAggs(aggs)
	if err != nil {
		t.Fatalf("RenderAggs() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if _, ok := result["by_name"]; !ok {
		t.Error("expected by_name key")
	}
}

func TestRenderer_RenderQuery_Error(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V2)

	q := b.Term("invalid", "value")
	_, err := r.RenderQuery(q)
	if err == nil {
		t.Error("RenderQuery() should return error for query with error")
	}
}
