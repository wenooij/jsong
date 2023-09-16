package jsong

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDeleteMap(t *testing.T) {
	m := map[string]any{"a": 1, "b": 2, "c": 3}

	got := Delete(m, "a")

	want := object{"b": num(2), "c": num(3)}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Delete(): got diff:\n%v", diff)
	}
}

func TestDeleteSlice(t *testing.T) {
	m := []any{1, 2, 3}

	got := Delete(m, "1")

	want := array{num(1), nil, num(3)}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Delete(): got diff:\n%v", diff)
	}
}

func TestDeletePath(t *testing.T) {
	m := map[string]any{"a": []any{1.0}}

	got := Delete(m, "a.0")

	want := object{"a": array{nil}}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Delete(): got diff:\n%v", diff)
	}
}

func TestDeleteEmbededPath(t *testing.T) {
	m := map[string]any{"a": []any{nil, map[string]any{"b": nil}, nil}}

	got := Delete(m, "a.1.b")

	want := object{"a": array{nil, object{}, nil}}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Delete(): got diff:\n%v", diff)
	}
}
