package lucene

import "testing"

func TestResolveFieldKind(t *testing.T) {
	tests := []struct {
		typeName string
		want     FieldKind
	}{
		{"string", KindString},
		{"*string", KindString},
		{"bool", KindBool},
		{"int", KindInt},
		{"int64", KindInt},
		{"uint32", KindInt},
		{"float32", KindFloat},
		{"float64", KindFloat},
		{"time.Time", KindTime},
		{"*time.Time", KindTime},
		{"[]string", KindSlice},
		{"[]int", KindSlice},
		{"[]float32", KindVector},
		{"[]float64", KindVector},
		{"SomeStruct", KindUnknown},
		{"map[string]any", KindUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			got := resolveFieldKind(tt.typeName)
			if got != tt.want {
				t.Errorf("resolveFieldKind(%q) = %v, want %v", tt.typeName, got, tt.want)
			}
		})
	}
}
