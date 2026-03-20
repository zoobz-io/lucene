//go:build testing

package benchmarks

import (
	"testing"

	"github.com/zoobz-io/lucene"
	"github.com/zoobz-io/lucene/elasticsearch"
	"github.com/zoobz-io/lucene/opensearch"
)

func BenchmarkRender_Elasticsearch_Simple(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	renderer := elasticsearch.NewRenderer(elasticsearch.V8)
	query := builder.Term("status", "active")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = renderer.RenderQuery(query)
	}
}

func BenchmarkRender_Elasticsearch_Complex(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	renderer := elasticsearch.NewRenderer(elasticsearch.V8)
	query := builder.Bool().
		Must(
			builder.Match("title", "search term"),
			builder.Term("status", "active"),
		).
		Should(
			builder.Match("description", "bonus content"),
		).
		Filter(
			builder.Range("price").Gte(10).Lte(1000),
		).
		MinimumShouldMatch(1)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = renderer.RenderQuery(query)
	}
}

func BenchmarkRender_Elasticsearch_Search(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	renderer := elasticsearch.NewRenderer(elasticsearch.V8)
	search := lucene.NewSearch().
		Query(
			builder.Bool().
				Must(builder.Match("title", "search term")).
				Filter(builder.Term("status", "active")),
		).
		Aggs(
			builder.TermsAgg("by_category", "category").Size(10),
			builder.Avg("avg_price", "price"),
		).
		Size(20).
		From(0)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = renderer.Render(search)
	}
}

func BenchmarkRender_OpenSearch_Simple(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	renderer := opensearch.NewRenderer(opensearch.V2)
	query := builder.Term("status", "active")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = renderer.RenderQuery(query)
	}
}

func BenchmarkRender_OpenSearch_Complex(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	renderer := opensearch.NewRenderer(opensearch.V2)
	query := builder.Bool().
		Must(
			builder.Match("title", "search term"),
			builder.Term("status", "active"),
		).
		Should(
			builder.Match("description", "bonus content"),
		).
		Filter(
			builder.Range("price").Gte(10).Lte(1000),
		).
		MinimumShouldMatch(1)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = renderer.RenderQuery(query)
	}
}

func BenchmarkRender_OpenSearch_Search(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	renderer := opensearch.NewRenderer(opensearch.V2)
	search := lucene.NewSearch().
		Query(
			builder.Bool().
				Must(builder.Match("title", "search term")).
				Filter(builder.Term("status", "active")),
		).
		Aggs(
			builder.TermsAgg("by_category", "category").Size(10),
			builder.Avg("avg_price", "price"),
		).
		Size(20).
		From(0)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = renderer.Render(search)
	}
}

func BenchmarkRender_Aggs(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	renderer := elasticsearch.NewRenderer(elasticsearch.V8)
	aggs := []lucene.Aggregation{
		builder.TermsAgg("by_category", "category").
			Size(10).
			SubAgg(builder.Avg("avg_price", "price")).
			SubAgg(builder.Max("max_price", "price")),
		builder.DateHistogram("by_month", "category").CalendarInterval("month"),
		builder.RangeAgg("price_ranges", "price").
			AddRange(0, 50).
			AddRange(50, 100).
			AddRange(100, nil),
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = renderer.RenderAggs(aggs)
	}
}
