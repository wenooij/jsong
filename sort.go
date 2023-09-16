package jsong

import (
	"sort"
)

func Sort(vs any) any {
	if _, ok := vs.(valueInterface); !ok {
		vs = ValueOf(vs)
	}
	a, ok := vs.(array)
	if !ok {
		return vs
	}
	sort.Sort(a)
	return a
}

func (a array) Len() int           { return len(a) }
func (a array) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a array) Less(i, j int) bool { return compare(a.At(i), a.At(j)) < 0 }

// SortByKey sorts the values by extracting the key using jsong.
func SortByKey(vs any, key string) any {
	if key == "" {
		return Sort(vs)
	}
	if _, ok := vs.(valueInterface); !ok {
		vs = ValueOf(vs)
	}
	a, ok := vs.(array)
	if !ok {
		return vs
	}
	sort.Slice(a, func(i, j int) bool {
		e1 := extractRec(a.At(i), key)
		e2 := extractRec(a.At(j), key)
		return compare(e1, e2) < 0
	})
	return a
}
