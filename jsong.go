// Package jsong implements merging JSON objects.
package jsong

import (
	"fmt"
	"slices"

	"golang.org/x/exp/maps"
)

type valueInterface interface {
	Get(k any) (valueInterface, bool)
	Put(k any, v valueInterface)
	Delete(k any)
	Each(func(k, v any) bool)
	compare(other valueInterface) int
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Errorf("assertion failed: %w", err))
	}
	return v
}

func MustOk[T any](v T, ok bool) T {
	if !ok {
		panic(fmt.Errorf("assertion failed on %T", v))
	}
	return v
}

func IsNull(v any) bool {
	_, ok := v.(null)
	return ok
}

func Bool(v any) (bool, bool) {
	b, ok := v.(boolean)
	if !ok {
		return false, false
	}
	return bool(b), true
}

func Float64(v any) (float64, bool) {
	n, ok := v.(num)
	if !ok {
		return 0, false
	}
	return float64(n), true
}

func String(v any) (string, bool) {
	s, ok := v.(str)
	if !ok {
		return "", false
	}
	return string(s), true
}

func Array(v any) ([]any, bool) {
	a, ok := v.(array)
	if !ok {
		return nil, false
	}
	return ([]any)(a), true
}

func Object(v any) (map[string]any, bool) {
	m, ok := v.(object)
	if !ok {
		return nil, false
	}
	return (map[string]any)(m), true
}

type null struct{}

func (null) Get(k any) (valueInterface, bool) { return nil, false }
func (null) Put(k any, v valueInterface)      {}
func (null) Delete(k any)                     {}
func (null) Each(func(any, any) bool)         {}

type boolean bool

func (boolean) Get(k any) (valueInterface, bool) { return nil, false }
func (boolean) Put(k any, v valueInterface)      {}
func (boolean) Delete(k any)                     {}
func (boolean) Each(func(any, any) bool)         {}

type num float64

func (num) Get(k any) (valueInterface, bool) { return nil, false }
func (num) Put(k any, v valueInterface)      {}
func (num) Delete(k any)                     {}
func (num) Each(func(any, any) bool)         {}

type str string

func (str) Get(k any) (valueInterface, bool) { return nil, false }
func (str) Put(k any, v valueInterface)      {}
func (str) Delete(k any)                     {}
func (str) Each(func(any, any) bool)         {}

type array []any // []valueInterface

func (a array) At(i int) valueInterface {
	return a[i].(valueInterface)
}

func (a array) Get(k any) (valueInterface, bool) {
	if i, ok := k.(int64); ok && i < int64(len(a)) {
		return a[i].(valueInterface), true
	}
	return nil, false
}

// Put puts the value v in the array index k.
// Put panics if k's index is out of bounds.
func (a array) Put(k any, v valueInterface) {
	if i, ok := k.(int64); ok {
		a[i] = v
	}
}

func (a array) Delete(k any) {
	if i, ok := k.(int64); ok && i < int64(len(a)) {
		a[i] = nil
	}
}

func (a array) Each(fn func(k, v any) bool) {
	for i, e := range a {
		if !fn(int64(i), e) {
			break
		}
	}
}

func (a array) Clone() array { return slices.Clone(a) }

type object map[string]any // map[string]valueInterface

func (a object) At(k string) valueInterface {
	return a[k].(valueInterface)
}

func (a object) Get(k any) (valueInterface, bool) {
	if k, ok := k.(string); ok {
		v, ok := a[k]
		if !ok {
			return nil, false
		}
		return v.(valueInterface), true
	}
	return nil, false
}

func (a object) Put(k any, v valueInterface) {
	if k, ok := k.(string); ok {
		a[k] = v
	}
}

func (a object) Delete(k any) {
	if k, ok := k.(string); ok {
		delete(a, k)
	}
}

func (a object) Each(fn func(k, v any) bool) {
	for k, v := range a {
		if !fn(k, v) {
			break
		}
	}
}

func (a object) Len() int           { return len(a) }
func (a object) Merge(other object) { maps.Copy(a, other) }
func (a object) Clear()             { maps.Clear(a) }
