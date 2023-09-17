package jsong

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGlobNestedMap(t *testing.T) {
	m := map[string]any{
		"a": map[string]any{
			"k1": []any{"a", "b", "c"},
			"k2": []any{"a", "b", "c"},
			"k3": []any{"a", "b", "c"},
		},
	}

	got := GlobKey(m, "a.*")

	want := []string{
		"a.k1",
		"a.k2",
		"a.k3",
	}

	lessFunc := func(a, b string) bool { return a < b }
	if diff := cmp.Diff(want, got, cmpopts.SortSlices(lessFunc)); diff != "" {
		t.Errorf("Glob(): got diff:\n%s", diff)
	}
}

func TestGlobMultipleNestedMap(t *testing.T) {
	m := map[string]any{
		"a": map[string]any{
			"k1": []any{"a", "b", "c"},
			"k2": []any{"a", "b", "c"},
			"k3": []any{"a", "b", "c"},
		},
	}
	got := GlobKey(m, "a.*.*")

	want := []string{
		"a.k1.0",
		"a.k1.1",
		"a.k1.2",
		"a.k2.0",
		"a.k2.1",
		"a.k2.2",
		"a.k3.0",
		"a.k3.1",
		"a.k3.2",
	}

	lessFunc := func(a, b string) bool { return a < b }
	if diff := cmp.Diff(want, got, cmpopts.SortSlices(lessFunc)); diff != "" {
		t.Errorf("Glob(): got diff:\n%s", diff)
	}
}

func TestGlobDoubleStarNestedMap(t *testing.T) {
	m := map[string]any{
		"a": map[string]any{
			"k1": []any{"a", "b", "c"},
			"k2": []any{"a", "b", "c"},
			"k3": []any{"a", "b", "c"},
		},
	}
	got := GlobKey(m, "**")

	want := []string{
		"a",
		"a.k1",
		"a.k2",
		"a.k3",
		"a.k1.0",
		"a.k1.1",
		"a.k1.2",
		"a.k2.0",
		"a.k2.1",
		"a.k2.2",
		"a.k3.0",
		"a.k3.1",
		"a.k3.2",
	}

	lessFunc := func(a, b string) bool { return a < b }
	if diff := cmp.Diff(want, got, cmpopts.SortSlices(lessFunc)); diff != "" {
		t.Errorf("Glob(): got diff:\n%s", diff)
	}
}

func TestGlobDoubleStarSuffixNestedMap(t *testing.T) {
	m := map[string]any{
		"a": map[string]any{
			"k1": []any{"a", "b", "c"},
			"k2": []any{"a", "b", "c"},
			"k3": []any{"a", "b", "c"},
		},
	}
	got := GlobKey(m, "**.0")

	want := []string{
		"a.k1.0",
		"a.k2.0",
		"a.k3.0",
	}

	lessFunc := func(a, b string) bool { return a < b }
	if diff := cmp.Diff(want, got, cmpopts.SortSlices(lessFunc)); diff != "" {
		t.Errorf("Glob(): got diff:\n%s", diff)
	}
}
