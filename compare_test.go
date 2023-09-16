package jsong

import "testing"

func TestCompareCycle(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("TestCompareCycle(): got panic = false, want panic = true")
		}
	}()

	a := []any{1, 2, 3}
	a[1] = a

	b := []any{}

	Compare(a, b)
}
