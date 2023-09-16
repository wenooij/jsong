package jsong

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExtractNil(t *testing.T) {
	got := Extract(nil, "a.simple.path")

	var want any = nil

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestExtractNil(): got diff:\n%v", diff)
	}
}
