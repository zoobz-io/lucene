//go:build testing

package benchmarks

import (
	"testing"

	"github.com/zoobz-io/lucene"
)

type benchDoc struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Price       float64  `json:"price"`
	Tags        []string `json:"tags"`
	Category    string   `json:"category"`
}

func BenchmarkBuilder_New(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = lucene.New[benchDoc]()
	}
}

func BenchmarkQuery_Term(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.Term("status", "active")
	}
}

func BenchmarkQuery_Match(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.Match("title", "search term")
	}
}

func BenchmarkQuery_Match_WithOptions(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.Match("title", "search term").
			Fuzziness("AUTO").
			Operator("and").
			Analyzer("standard").
			Boost(1.5)
	}
}

func BenchmarkQuery_Bool_Simple(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.Bool().
			Must(builder.Term("status", "active")).
			Filter(builder.Range("price").Gte(10))
	}
}

func BenchmarkQuery_Bool_Complex(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.Bool().
			Must(
				builder.Match("title", "search term"),
				builder.Term("status", "active"),
			).
			Should(
				builder.Match("description", "bonus content"),
				builder.Terms("category", "electronics", "computers"),
			).
			MustNot(
				builder.Term("status", "deleted"),
			).
			Filter(
				builder.Range("price").Gte(10).Lte(1000),
				builder.Exists("tags"),
			).
			MinimumShouldMatch(1).
			Boost(1.5)
	}
}

func BenchmarkQuery_Range(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.Range("price").Gte(10).Lt(100)
	}
}

func BenchmarkQuery_MultiMatch(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.MultiMatch("search term", "title", "description").
			Type("best_fields").
			TieBreaker(0.3)
	}
}

func BenchmarkQuery_Nested_Deep(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.Bool().
			Must(
				builder.Bool().
					Must(
						builder.Bool().
							Must(builder.Term("status", "active")).
							Filter(builder.Range("price").Gte(10)),
					).
					Should(builder.Match("title", "search")),
			)
	}
}

func BenchmarkAggregation_Terms(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.TermsAgg("by_category", "category").Size(10)
	}
}

func BenchmarkAggregation_Complex(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = builder.TermsAgg("by_category", "category").
			Size(10).
			SubAgg(builder.Avg("avg_price", "price")).
			SubAgg(builder.Max("max_price", "price")).
			SubAgg(builder.Stats("price_stats", "price"))
	}
}

func BenchmarkSearch_Full(b *testing.B) {
	builder, _ := lucene.New[benchDoc]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lucene.NewSearch().
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
			From(0).
			Sort(lucene.SortField{Field: "price", Order: "desc"})
	}
}
