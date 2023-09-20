package jsong

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecodeNull(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`null`)).Decode()

	want := null{}
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestDecodeFalse(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`false`)).Decode()

	want := boolean(false)
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestDecodeTrue(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`true`)).Decode()

	want := boolean(true)
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestDecodeZero(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`0`)).Decode()

	want := num(0)
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestDecodeNumber(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`-1.2e+5`)).Decode()

	want := num(-1.2e+5)
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestDecodeString(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`"abc"`)).Decode()

	want := str("abc")
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestDecodeEscapeString(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`"\"abc\""`)).Decode()

	want := str("\"abc\"")
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestDecodeEmptyArray(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`[]`)).Decode()

	want := array{}
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestDecodeArray(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`[0, 1, 2]`)).Decode()

	want := array{num(0), num(1), num(2)}
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestDecodeArrayNull(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`[null, null]`)).Decode()

	want := array{null{}, null{}}
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestDecodeArrayBool(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`[false,false]`)).Decode()

	want := array{boolean(false), boolean(false)}
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestDecodeNestedArray(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`[[[["nested"]]]]`)).Decode()

	want := array{array{array{array{str("nested")}}}}
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestDecodeEmptyObject(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`{}`)).Decode()

	want := object{}
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}
func TestDecodeObject(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`{"a": 0, "b": 1, "c": 2}`)).Decode()

	want := object{"a": num(0), "b": num(1), "c": num(2)}
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): want err = %v, got err = %v", wantErr, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestDecodeNestedObject(t *testing.T) {
	got, err := NewDecoder(strings.NewReader(`{"a": {}}`)).Decode()

	want := object{"a": object{}}
	wantErr := false

	gotErr := err != nil
	if wantErr != gotErr {
		t.Fatalf("Decode(): got err = %v, want err = %v", err, wantErr)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Decode(): got diff:\n%s", diff)
	}
}

func TestUnquoteInPlace(t *testing.T) {
	got, err := unquoteInPlace([]byte(`"hello \"world\"!"`))
	want := []byte(`hello "world"!`)

	gotErr := err != nil
	wantErr := false
	if gotErr != wantErr {
		t.Fatalf("UnquoteInPlace(): got err = %v, want err = %v", err, wantErr)
	}

	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("UnquoteInPlace(): got diff:\n%s", diff)
	}
}

func TestUnquoteInPlaceQuotedEmpty(t *testing.T) {
	got, err := unquoteInPlace([]byte(`""`))
	want := []byte{}

	gotErr := err != nil
	wantErr := false
	if gotErr != wantErr {
		t.Fatalf("UnquoteInPlace(): got err = %v, want err = %v", err, wantErr)
	}

	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("UnquoteInPlace(): got diff:\n%s", diff)
	}
}
