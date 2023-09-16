package jsong

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type mergeTestCase struct {
	name      string
	dst       any
	src       any
	dstField  string
	srcField  string
	want      any
	wantPanic bool
}

func (tc mergeTestCase) runTest(t *testing.T) {
	t.Helper()
	defer func() {
		if gotPanic := recover(); gotPanic != nil {
			gotErr := gotPanic.(error)
			if gotErr == nil && tc.wantPanic {
				t.Errorf("Merge(%q): want panic = true, got panic = false", tc.name)
			} else if gotErr != nil && !tc.wantPanic {
				t.Errorf("Merge(%q): want panic = false, got panic = %v", tc.name, gotErr)
			}
		}
	}()
	got := Merge(tc.dst, tc.src, tc.dstField, tc.srcField)
	if diff := cmp.Diff(tc.want, got); diff != "" {
		t.Errorf("Merge(%q): got diff:\n%s", tc.name, diff)
	}
}

func TestMerge(t *testing.T) {
	for _, tc := range []mergeTestCase{{
		name: "empty",
		want: null{},
	}, {
		name: "merge empty map into nil",
		dst:  nil,
		src:  map[string]any{},
		want: object{},
	}, {
		name: "merge map into nil",
		dst:  nil,
		src:  map[string]int{"a": 1, "b": 2, "c": 3},
		want: object{"a": num(1), "b": num(2), "c": num(3)},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			tc.runTest(t)
		})
	}
}
