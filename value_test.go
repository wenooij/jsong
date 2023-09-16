package jsong

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValueOfNil(t *testing.T) {
	var a any

	got := ValueOf(a)

	want := null{}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestValueOfNil(): got diff:\n%v", diff)
	}
}

func TestValueOfBool(t *testing.T) {
	var b bool

	got := ValueOf(b)

	want := boolean(false)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestValueOfBool(): got diff:\n%v", diff)
	}
}

func TestValueOfArray(t *testing.T) {
	m := []any{nil}

	got := ValueOf(m)

	want := array{nil}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestValueOfArray(): got diff:\n%v", diff)
	}
}

func TestValueOfStruct(t *testing.T) {
	x := struct{ V int }{V: 1}

	got := ValueOf(x)

	want := object{"V": num(1)}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestValueOfStruct(): got diff:\n%v", diff)
	}
}

func TestValueOfEmbededSlice(t *testing.T) {
	m := map[string]any{"a": []any{nil}}

	got := ValueOf(m)

	want := object{"a": array{nil}}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestValueOfEmbededSlice(): got diff:\n%v", diff)
	}
}

func TestValueOfArrayCycle(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("TestValueOfArrayCycle(): want panic = true, got panic = false")
		}
	}()

	a := []any{nil, nil, nil}
	a[1] = a

	ValueOf(a)
}

func TestValueOfMapCycle(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("TestValueOfArrayCycle(): want panic = true, got panic = false")
		}
	}()

	a := map[string]any{"a": nil, "b": nil, "c": nil}
	a["b"] = a

	ValueOf(a)
}
