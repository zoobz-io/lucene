//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/zoobz-io/lucene"
	"github.com/zoobz-io/lucene/elasticsearch"
)

type testDoc struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

func elasticsearchEndpoint() string {
	if ep := os.Getenv("ELASTICSEARCH_ENDPOINT"); ep != "" {
		return ep
	}
	return "http://localhost:9200"
}

// postJSON sends a POST request with JSON body.
// #nosec G107 - URL is constructed from test endpoint
func postJSON(t *testing.T, url, body string) (*http.Response, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	return client.Do(req)
}

func skipIfNoElasticsearch(t *testing.T) string {
	t.Helper()
	endpoint := elasticsearchEndpoint()
	client := &http.Client{Timeout: 2 * time.Second}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		t.Skipf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Skipf("Elasticsearch not available at %s: %v", endpoint, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Skipf("Elasticsearch not healthy at %s: status %d", endpoint, resp.StatusCode)
	}
	return endpoint
}

func TestElasticsearch_TermQuery(t *testing.T) {
	endpoint := skipIfNoElasticsearch(t)

	builder := lucene.New[testDoc]()

	query := builder.Term("status", "active")
	renderer := elasticsearch.NewRenderer(elasticsearch.V8)

	body, err := renderer.RenderQuery(query)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	// Validate JSON structure
	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Execute against Elasticsearch (using _validate/query)
	validateURL := fmt.Sprintf("%s/_validate/query", endpoint)
	reqBody := fmt.Sprintf(`{"query": %s}`, string(body))
	resp, err := postJSON(t, validateURL, reqBody)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("validate query failed: %s", string(respBody))
	}
}

func TestElasticsearch_BoolQuery(t *testing.T) {
	endpoint := skipIfNoElasticsearch(t)

	builder := lucene.New[testDoc]()

	query := builder.Bool().
		Must(builder.Match("title", "search term")).
		Filter(builder.Term("status", "active")).
		Should(builder.Range("price").Gte(10))

	renderer := elasticsearch.NewRenderer(elasticsearch.V8)
	body, err := renderer.RenderQuery(query)
	if err != nil {
		t.Fatalf("RenderQuery() error = %v", err)
	}

	validateURL := fmt.Sprintf("%s/_validate/query", endpoint)
	reqBody := fmt.Sprintf(`{"query": %s}`, string(body))
	resp, err := postJSON(t, validateURL, reqBody)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf("validate query failed: %s", string(respBody))
	}
}

func TestElasticsearch_FullSearch(t *testing.T) {
	_ = skipIfNoElasticsearch(t)

	builder := lucene.New[testDoc]()

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

	renderer := elasticsearch.NewRenderer(elasticsearch.V8)
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

func TestElasticsearch_Aggregations(t *testing.T) {
	endpoint := skipIfNoElasticsearch(t)

	builder := lucene.New[testDoc]()

	aggs := []lucene.Aggregation{
		builder.TermsAgg("by_category", "category").Size(10),
		builder.Avg("avg_price", "price"),
		builder.Stats("price_stats", "price"),
	}

	renderer := elasticsearch.NewRenderer(elasticsearch.V8)
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
	searchURL := fmt.Sprintf("%s/_search?size=0", endpoint)
	reqBody := fmt.Sprintf(`{"aggs": %s}`, string(body))
	resp, err := postJSON(t, searchURL, reqBody)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf("aggregation search failed: %s", string(respBody))
	}
}
