package jsong

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestArrayMarshalJSON(t *testing.T) {
	m := array{num(1), str("2"), boolean(true), nil}

	got, gotErr := json.Marshal(m)

	if gotErr != nil {
		t.Errorf("TestObjectMarshalJSON(): got err: %v", gotErr)
	}

	want := []byte(`[1,"2",true,null]`)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestObjectMarshalJSON() got diff:\n%v", diff)
	}
}

func TestObjectMarshalJSON(t *testing.T) {
	m := object{"a": num(1), "b": str("2"), "c": boolean(true), "d": nil}

	got, gotErr := json.Marshal(m)

	if gotErr != nil {
		t.Errorf("TestObjectMarshalJSON(): got err: %v", gotErr)
	}

	want := []byte(`{"a":1,"b":"2","c":true,"d":null}`)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestObjectMarshalJSON() got diff:\n%v", diff)
	}
}
