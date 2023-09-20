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

func TestCompileKeyMatcherLit(t *testing.T) {
	m := Must(CompileKeyMatcher("a.b.c.d"))

	got := m.r.String()

	want := `^a\.b\.c\.d$`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("CompileKeyMatcher(): got diff:\n%s", diff)
	}
}

func TestCompileKeyMatcherStar(t *testing.T) {
	m := Must(CompileKeyMatcher("a.*.c"))

	got := m.r.String()

	want := `^a\.[^.]*\.c$`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("CompileKeyMatcher(): got diff:\n%s", diff)
	}
}

func TestCompileKeyMatcherDoubleStar(t *testing.T) {
	m := Must(CompileKeyMatcher("a.**.c"))

	got := m.r.String()

	want := `^a\..*\.c$`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("CompileKeyMatcher(): got diff:\n%s", diff)
	}
}
