package elasticsearch

import (
	"encoding/json"
	"testing"

	"github.com/zoobz-io/lucene"
)

type testDoc struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Status string  `json:"status"`
	Price  float64 `json:"price"`
}

func mustBuilder(t *testing.T) *lucene.Builder[testDoc] {
	t.Helper()
	return lucene.New[testDoc]()
}

func TestRenderer_RenderQuery_Term(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V8)

	t.Run("basic term", func(t *testing.T) {
		q := b.Term("status", "active")
		out, err := r.RenderQuery(q)
		if err != nil {
			t.Fatalf("RenderQuery() error = %v", err)
		}

		var result map[string]any
		if err := json.Unmarshal(out, &result); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		term := result["term"].(map[string]any)
		status := term["status"].(map[string]any)
		if status["value"] != "active" {
			t.Errorf("value = %v, want active", status["value"])
		}
	})

	t.Run("term with boost", func(t *testing.T) {
		q := b.Term("status", "active").Boost(1.5)
		out, err := r.RenderQuery(q)
		if err != nil {
			t.Fatalf("RenderQuery() error = %v", err)
		}

		var result map[string]any
		if err := json.Unmarshal(out, &result); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		term := result["term"].(map[string]any)
		status := term["status"].(map[string]any)
		if status["boost"] != 1.5 {
			t.Errorf("boost = %v, want 1.5", status["boost"])
		}
	})
}

func TestRenderer_RenderQuery_Range(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V8)

	q := b.Range("price").Gte(10).Lt(100)
	out, err := r.RenderQuery(q)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	rng := result["range"].(map[string]any)
	price := rng["price"].(map[string]any)
	if price["gte"] != float64(10) {
		t.Errorf("gte = %v, want 10", price["gte"])
	}
	if price["lt"] != float64(100) {
		t.Errorf("lt = %v, want 100", price["lt"])
	}
}

func TestRenderer_RenderQuery_Match(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V8)

	t.Run("basic match", func(t *testing.T) {
		q := b.Match("name", "test query")
		out, err := r.RenderQuery(q)
		if err != nil {
			t.Fatalf("RenderQuery() error = %v", err)
		}

		var result map[string]any
		if err := json.Unmarshal(out, &result); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		match := result["match"].(map[string]any)
		name := match["name"].(map[string]any)
		if name["query"] != "test query" {
			t.Errorf("query = %v, want test query", name["query"])
		}
	})

	t.Run("match with options", func(t *testing.T) {
		q := b.Match("name", "test").Fuzziness("AUTO").Operator("and")
		out, err := r.RenderQuery(q)
		if err != nil {
			t.Fatalf("RenderQuery() error = %v", err)
		}

		var result map[string]any
		if err := json.Unmarshal(out, &result); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		match := result["match"].(map[string]any)
		name := match["name"].(map[string]any)
		if name["fuzziness"] != "AUTO" {
			t.Errorf("fuzziness = %v, want AUTO", name["fuzziness"])
		}
		if name["operator"] != "and" {
			t.Errorf("operator = %v, want and", name["operator"])
		}
	})
}

func TestRenderer_RenderQuery_Bool(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V8)

	q := b.Bool().
		Must(b.Term("status", "active")).
		Should(b.Match("name", "test")).
		MinimumShouldMatch(1)

	out, err := r.RenderQuery(q)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	boolQ := result["bool"].(map[string]any)
	must := boolQ["must"].([]any)
	if len(must) != 1 {
		t.Errorf("len(must) = %d, want 1", len(must))
	}
	should := boolQ["should"].([]any)
	if len(should) != 1 {
		t.Errorf("len(should) = %d, want 1", len(should))
	}
	if boolQ["minimum_should_match"] != float64(1) {
		t.Errorf("minimum_should_match = %v, want 1", boolQ["minimum_should_match"])
	}
}

func TestRenderer_RenderQuery_MatchAll(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V8)

	q := b.MatchAll()
	out, err := r.RenderQuery(q)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if _, ok := result["match_all"]; !ok {
		t.Error("expected match_all key")
	}
}

func TestRenderer_RenderQuery_Error(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V8)

	q := b.Term("invalid", "value")
	_, err := r.RenderQuery(q)
	if err == nil {
		t.Error("RenderQuery() should return error for query with error")
	}
}

func TestRenderer_RenderQuery_MultiMatch(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V8)

	q := b.MultiMatch("search term", "name", "status").
		Type("best_fields").
		TieBreaker(0.3)

	out, err := r.RenderQuery(q)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	mm := result["multi_match"].(map[string]any)
	if mm["query"] != "search term" {
		t.Errorf("query = %v, want search term", mm["query"])
	}
	fields := mm["fields"].([]any)
	if len(fields) != 2 {
		t.Errorf("len(fields) = %d, want 2", len(fields))
	}
	if mm["type"] != "best_fields" {
		t.Errorf("type = %v, want best_fields", mm["type"])
	}
	if mm["tie_breaker"] != 0.3 {
		t.Errorf("tie_breaker = %v, want 0.3", mm["tie_breaker"])
	}
}

func TestRenderer_RenderQuery_Exists(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V8)

	q := b.Exists("name")
	out, err := r.RenderQuery(q)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	exists := result["exists"].(map[string]any)
	if exists["field"] != "name" {
		t.Errorf("field = %v, want name", exists["field"])
	}
}

func TestRenderer_RenderQuery_IDs(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V8)

	q := b.IDs("doc1", "doc2")
	out, err := r.RenderQuery(q)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	ids := result["ids"].(map[string]any)
	values := ids["values"].([]any)
	if len(values) != 2 {
		t.Errorf("len(values) = %d, want 2", len(values))
	}
}

func TestRenderer_Render_FullSearch(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V8)

	h := lucene.NewHighlight().
		Fields("name").
		PreTags("<em>").
		PostTags("</em>")

	s := lucene.NewSearch().
		Query(b.Match("name", "test")).
		Aggs(b.TermsAgg("by_status", "status").Size(10)).
		Size(20).
		From(5).
		Sort(lucene.SortField{Field: "price", Order: "asc"}).
		Source("name", "price").
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

	// Verify query
	if _, ok := result["query"]; !ok {
		t.Error("expected query key")
	}

	// Verify aggregations
	if _, ok := result["aggs"]; !ok {
		t.Error("expected aggs key")
	}

	// Verify size and from
	if result["size"] != float64(20) {
		t.Errorf("size = %v, want 20", result["size"])
	}
	if result["from"] != float64(5) {
		t.Errorf("from = %v, want 5", result["from"])
	}

	// Verify sort
	sort := result["sort"].([]any)
	if len(sort) != 1 {
		t.Errorf("len(sort) = %d, want 1", len(sort))
	}

	// Verify _source
	source := result["_source"].([]any)
	if len(source) != 2 {
		t.Errorf("len(_source) = %d, want 2", len(source))
	}

	// Verify highlight
	if _, ok := result["highlight"]; !ok {
		t.Error("expected highlight key")
	}

	// Verify track_total_hits
	if result["track_total_hits"] != true {
		t.Errorf("track_total_hits = %v, want true", result["track_total_hits"])
	}

	// Verify min_score
	if result["min_score"] != 0.5 {
		t.Errorf("min_score = %v, want 0.5", result["min_score"])
	}

	// Verify timeout
	if result["timeout"] != "10s" {
		t.Errorf("timeout = %v, want 10s", result["timeout"])
	}
}

func TestRenderer_Render_SourceFiltering(t *testing.T) {
	r := NewRenderer(V8)

	t.Run("includes only", func(t *testing.T) {
		s := lucene.NewSearch().Source("name", "price")
		out, err := r.Render(s)
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		var result map[string]any
		if err := json.Unmarshal(out, &result); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		source := result["_source"].([]any)
		if len(source) != 2 {
			t.Errorf("len(_source) = %d, want 2", len(source))
		}
	})

	t.Run("includes and excludes", func(t *testing.T) {
		s := lucene.NewSearch().
			SourceIncludes("name", "price").
			SourceExcludes("internal")
		out, err := r.Render(s)
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		var result map[string]any
		if err := json.Unmarshal(out, &result); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		source := result["_source"].(map[string]any)
		if _, ok := source["includes"]; !ok {
			t.Error("expected includes key")
		}
		if _, ok := source["excludes"]; !ok {
			t.Error("expected excludes key")
		}
	})

	t.Run("excludes only", func(t *testing.T) {
		s := lucene.NewSearch().SourceExcludes("internal")
		out, err := r.Render(s)
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		var result map[string]any
		if err := json.Unmarshal(out, &result); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		source := result["_source"].(map[string]any)
		if _, ok := source["excludes"]; !ok {
			t.Error("expected excludes key")
		}
	})
}

func TestRenderer_Render_Highlight(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V8)

	field := lucene.NewHighlightField("name").
		FragmentSize(100).
		NumFragments(3).
		NoMatchSize(50).
		Build()

	h := lucene.NewHighlight().
		Field(field).
		PreTags("<b>").
		PostTags("</b>").
		Encoder("html").
		Order("score").
		Highlighter("unified")

	s := lucene.NewSearch().
		Query(b.MatchAll()).
		Highlight(h)

	out, err := r.Render(s)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	highlight := result["highlight"].(map[string]any)

	// Global settings
	preTags := highlight["pre_tags"].([]any)
	if len(preTags) != 1 || preTags[0] != "<b>" {
		t.Errorf("pre_tags = %v, want [<b>]", preTags)
	}
	postTags := highlight["post_tags"].([]any)
	if len(postTags) != 1 || postTags[0] != "</b>" {
		t.Errorf("post_tags = %v, want [</b>]", postTags)
	}
	if highlight["encoder"] != "html" {
		t.Errorf("encoder = %v, want html", highlight["encoder"])
	}
	if highlight["order"] != "score" {
		t.Errorf("order = %v, want score", highlight["order"])
	}
	if highlight["type"] != "unified" {
		t.Errorf("type = %v, want unified", highlight["type"])
	}

	// Field settings
	fields := highlight["fields"].(map[string]any)
	nameField := fields["name"].(map[string]any)
	if nameField["fragment_size"] != float64(100) {
		t.Errorf("fragment_size = %v, want 100", nameField["fragment_size"])
	}
	if nameField["number_of_fragments"] != float64(3) {
		t.Errorf("number_of_fragments = %v, want 3", nameField["number_of_fragments"])
	}
	if nameField["no_match_size"] != float64(50) {
		t.Errorf("no_match_size = %v, want 50", nameField["no_match_size"])
	}
}

func TestRenderer_Render_NilSearch(t *testing.T) {
	r := NewRenderer(V8)
	_, err := r.Render(nil)
	if err == nil {
		t.Error("Render(nil) should return error")
	}
}

func TestRenderer_Render_SearchWithError(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V8)

	s := lucene.NewSearch().Query(b.Match("invalid", "test"))
	_, err := r.Render(s)
	if err == nil {
		t.Error("Render() should return error for search with invalid query")
	}
}

func TestRenderer_RenderQuery_HybridUnsupported(t *testing.T) {
	b := mustBuilder(t)
	r := NewRenderer(V8)

	q := b.Hybrid(b.Match("name", "test"))
	_, err := r.RenderQuery(q)
	if err == nil {
		t.Error("RenderQuery() should return error for hybrid query on Elasticsearch")
	}
}
