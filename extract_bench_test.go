package jsong

import (
	"fmt"
	"testing"
)

var benchExtractValues = []any{
	"",
	"abc",
	[]int{1, 2, 3},
	map[string]int{"a": 1, "b": 2, "c": 3},
	struct {
		X int
		Y struct{ Z *int }
	}{X: 5, Y: struct{ Z *int }{}},
	make([]int, 100),
	make([]byte, 20),
	string(make([]byte, 20)),
	make([]struct{ X, Y, Z **int }, 20),
}

var benchExtractPaths = []string{
	"",
	"0",
	"a",
	"y.z",
	"15",
	"z",
}

func BenchmarkExtract(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, v := range benchExtractValues {
			for _, path := range benchExtractPaths {
				Extract(v, path)
			}
		}
	}
}

func fallbackExtract(x any, path string) (any, error) {
	v := ValueOf(x)
	if path == "" {
		return v, nil
	}
	for head, tail, _ := Cut(path); ; head, tail, _ = Cut(tail) {
		var t any
		switch v := v.(type) {
		case []any:
			i, ok := head.(int64)
			if !ok {
				return nil, fmt.Errorf("not an array index: %v", head)
			}
			if int64(len(v)) <= i {
				return nil, fmt.Errorf("array index %d out of bounds: %d", i, len(v))
			}
			t = v[i]
		case map[string]any:
			var ok bool
			t, ok = v[head.(string)]
			if !ok {
				return nil, fmt.Errorf("key error: %v", head)
			}
		default:
			return nil, fmt.Errorf("not a JSON object or array at %v: %T", head, v)
		}
		if tail == "" {
			return t, nil
		}
		v = t
	}
}

func BenchmarkExtractBaseline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, v := range benchExtractValues {
			for _, path := range benchExtractPaths {
				fallbackExtract(v, path)
			}
		}
	}
}
