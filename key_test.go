package jsong

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJoinKey(t *testing.T) {
	got := JoinKey("foo.bar", "a", "b", int64(0))

	want := "foo.bar.a.b.0"

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("JoinKey(): got diff:\n%s", diff)
	}
}

func TestJoinKeyQuoteReserved(t *testing.T) {
	got := JoinKey("foo.bar", `.`, "*", "")

	want := `foo.bar."."."*".`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("JoinKey(): got diff:\n%s", diff)
	}
}
