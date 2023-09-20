package jsong

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testVisit struct {
	K string
	V any
}

func TestVisitNestedMap(t *testing.T) {
	m := map[string]any{
		"a": 1,
		"b": 2,
		"c": []any{
			false,
			true,
			map[string]bool{"d": false},
		},
	}

	var gotVisits []testVisit
	Visit(m, func(k string, v any) error {
		gotVisits = append(gotVisits, testVisit{K: k, V: v})
		return nil
	})

	wantVisits := []testVisit{
		{K: "", V: ValueOf(m)},
		{K: "a", V: ValueOf(m["a"])},
		{K: "b", V: ValueOf(m["b"])},
		{K: "c", V: ValueOf(m["c"])},
		{K: "c.0", V: ValueOf(m["c"].([]any)[0])},
		{K: "c.1", V: ValueOf(m["c"].([]any)[1])},
		{K: "c.2", V: ValueOf(m["c"].([]any)[2])},
		{K: "c.2.d", V: ValueOf(m["c"].([]any)[2].(map[string]bool)["d"])},
	}

	less := func(a, b testVisit) bool { return a.K < b.K }
	if diff := cmp.Diff(wantVisits, gotVisits, cmpopts.SortSlices(less)); diff != "" {
		t.Errorf("Visit(): got diff:\n%s", diff)
	}
}
