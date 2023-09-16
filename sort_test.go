package jsong

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type sortTestCase struct {
	name        string
	inputValues any
	wantValues  valueInterface
	wantPanic   bool
}

func (tc sortTestCase) recover(t *testing.T, testName string) {
	t.Helper()
	if err := recover(); err != nil {
		if !tc.wantPanic {
			t.Errorf("%s(%q): recovered from unexpected panic: %v", testName, tc.name, err)
		}
	} else if tc.wantPanic {
		t.Errorf("%s(%q): test case expected a panic but got no panic", testName, tc.name)
	}
}

func (tc sortTestCase) check(t *testing.T, gotValues any, testName string) {
	t.Helper()
	if !tc.wantPanic {
		if _, ok := tc.inputValues.(array); ok {
			// Test diff in place.
			if diff := cmp.Diff(tc.wantValues, tc.inputValues); diff != "" {
				t.Errorf("TestSortByKey(%q): got (in place) diff:\n%s", tc.name, diff)
			}
		}
		if diff := cmp.Diff(tc.wantValues, gotValues); diff != "" {
			t.Errorf("TestSortByKey(%q): got diff:\n%s", tc.name, diff)
		}
	}
}

func (tc sortTestCase) runTest(t *testing.T) {
	t.Helper()
	defer tc.recover(t, "TestSort")

	gotValues := Sort(tc.inputValues)
	tc.check(t, gotValues, "TestSort")
}

type sortByKeyTestCase struct {
	sortTestCase
	inputKey string
}

func (tc sortByKeyTestCase) runTest(t *testing.T) {
	t.Helper()
	defer tc.recover(t, "TestSortByKey")

	gotValues := SortByKey(tc.inputValues, tc.inputKey)
	tc.check(t, gotValues, "TestSortByKey")
}

func TestSort(t *testing.T) {
	for _, tc := range []sortTestCase{{
		name:       "empty",
		wantValues: null{},
	}, {
		name:        "strings",
		inputValues: []any{"c", "b", "a"},
		wantValues:  array{str("a"), str("b"), str("c")},
	}, {
		name: "heterogeneous array sorted in place",
		inputValues: array{
			object{},
			object(nil),
			object{"a": num(1)},
			array{nil},
			array{},
			array(nil),
			boolean(false),
			num(1),
			num(0),
		},
		wantValues: array{
			boolean(false),
			num(0),
			num(1),
			array(nil),
			array{},
			array{nil},
			object(nil),
			object{},
			object{"a": num(1)},
		},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			tc.runTest(t)
		})
	}
}

func TestSortByKey(t *testing.T) {
	for _, tc := range []sortByKeyTestCase{{
		sortTestCase: sortTestCase{
			name:       "empty",
			wantValues: null{},
		},
	}, {
		sortTestCase: sortTestCase{
			name:        "int key sorts by slice index",
			inputValues: []any{[]any{"c"}, []any{"b"}, []any{"a"}},
			wantValues:  array{array{str("a")}, array{str("b")}, array{str("c")}},
		},
		inputKey: "0",
	}, {
		sortTestCase: sortTestCase{
			name: "field key sorts by field",
			inputValues: []any{
				struct{ V string }{V: "b"},
				struct{ V string }{V: "c"},
				struct{ V string }{V: "a"},
			},
			wantValues: array{
				object{"V": str("a")},
				object{"V": str("b")},
				object{"V": str("c")},
			},
		},
		inputKey: "V",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			tc.runTest(t)
		})
	}
}
