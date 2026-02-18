//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/zoobzio/lucene"
	"github.com/zoobzio/lucene/opensearch"
)

func opensearchEndpoint() string {
	if ep := os.Getenv("OPENSEARCH_ENDPOINT"); ep != "" {
		return ep
	}
	return "http://localhost:9200"
}

func skipIfNoOpenSearch(t *testing.T) string {
	t.Helper()
	endpoint := opensearchEndpoint()
	client := &http.Client{Timeout: 2 * time.Second}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		t.Skipf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Skipf("OpenSearch not available at %s: %v", endpoint, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Skipf("OpenSearch not healthy at %s: status %d", endpoint, resp.StatusCode)
	}
	return endpoint
}

// postJSONOpenSearch sends a POST request with JSON body.
// #nosec G107 - URL is constructed from test endpoint
func postJSONOpenSearch(t *testing.T, endpoint, path, body string) (*http.Response, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := endpoint + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	return client.Do(req)
}

func TestOpenSearch_TermQuery(t *testing.T) {
	endpoint := skipIfNoOpenSearch(t)

	builder, err := lucene.New[testDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	query := builder.Term("status", "active")
	renderer := opensearch.NewRenderer(opensearch.V2)

	body, err := renderer.RenderQuery(query)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	// Validate JSON structure
	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Execute against OpenSearch (using _validate/query)
	reqBody := `{"query": ` + string(body) + `}`
	resp, err := postJSONOpenSearch(t, endpoint, "/_validate/query", reqBody)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("validate query failed: %s", string(respBody))
	}
}

func TestOpenSearch_BoolQuery(t *testing.T) {
	endpoint := skipIfNoOpenSearch(t)

	builder, err := lucene.New[testDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	query := builder.Bool().
		Must(builder.Match("title", "search term")).
		Filter(builder.Term("status", "active")).
		Should(builder.Range("price").Gte(10))

	renderer := opensearch.NewRenderer(opensearch.V2)
	body, err := renderer.RenderQuery(query)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	reqBody := `{"query": ` + string(body) + `}`
	resp, err := postJSONOpenSearch(t, endpoint, "/_validate/query", reqBody)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf("validate query failed: %s", string(respBody))
	}
}

func TestOpenSearch_FullSearch(t *testing.T) {
	_ = skipIfNoOpenSearch(t)

	builder, err := lucene.New[testDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	search := lucene.NewSearch().
		Query(
			builder.Bool().
				Must(builder.Match("title", "test")).
				Filter(builder.Term("status", "active")),
		).
		Aggs(
			builder.TermsAgg("by_category", "category").Size(10),
		).
		Size(20).
		From(0)

	renderer := opensearch.NewRenderer(opensearch.V2)
	body, err := renderer.Render(search)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	// Validate JSON structure
	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Check expected top-level keys
	if _, ok := parsed["query"]; !ok {
		t.Error("missing 'query' key in search body")
	}
	if _, ok := parsed["aggs"]; !ok {
		t.Error("missing 'aggs' key in search body")
	}
	if _, ok := parsed["size"]; !ok {
		t.Error("missing 'size' key in search body")
	}
}

func TestOpenSearch_Aggregations(t *testing.T) {
	endpoint := skipIfNoOpenSearch(t)

	builder, err := lucene.New[testDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	aggs := []lucene.Aggregation{
		builder.TermsAgg("by_category", "category").Size(10),
		builder.Avg("avg_price", "price"),
		builder.Stats("price_stats", "price"),
	}

	renderer := opensearch.NewRenderer(opensearch.V2)
	body, err := renderer.RenderAggs(aggs)
	if err != nil {
		t.Fatalf("RenderAggs() error = %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if _, ok := parsed["by_category"]; !ok {
		t.Error("missing 'by_category' aggregation")
	}
	if _, ok := parsed["avg_price"]; !ok {
		t.Error("missing 'avg_price' aggregation")
	}

	// Execute search with aggregations
	reqBody := `{"aggs": ` + string(body) + `}`
	resp, err := postJSONOpenSearch(t, endpoint, "/_search?size=0", reqBody)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf("aggregation search failed: %s", string(respBody))
	}
}

func TestOpenSearch_KnnQuery(t *testing.T) {
	_ = skipIfNoOpenSearch(t)

	type vectorDoc struct {
		Title     string    `json:"title"`
		Embedding []float32 `json:"embedding"`
	}

	builder, err := lucene.New[vectorDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	vector := []float32{0.1, 0.2, 0.3, 0.4, 0.5}
	query := builder.Knn("embedding", vector).K(10).NumCandidates(100)

	renderer := opensearch.NewRenderer(opensearch.V2)
	body, err := renderer.RenderQuery(query)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Note: KNN query validation requires an index with vector field mapping
	// This test just validates the JSON structure
	t.Logf("KNN query JSON: %s", string(body))
}
