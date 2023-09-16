package jsong

import (
	"slices"

	"golang.org/x/exp/maps"
)

func Compare(a, b any) int {
	if _, ok := a.(valueInterface); !ok {
		a = ValueOf(a)
	}
	if _, ok := b.(valueInterface); !ok {
		b = ValueOf(b)
	}
	return compare(a.(valueInterface), b.(valueInterface))
}

func fastCompare(a, b any) (int, bool) {
	if a == nil {
		if b == nil {
			return 0, true
		}
		return -1, true
	}
	if b == nil {
		return +1, true
	}
	if &a == &b {
		return 0, true
	}
	return 0, false
}

func compare(a, b valueInterface) int {
	if cmp, ok := fastCompare(a, b); ok {
		return cmp
	}
	switch a := a.(type) {
	case null:
		return a.compare(b)
	case boolean:
		switch b := b.(type) {
		case null:
			return +1
		case boolean:
			return a.compare(b)
		default:
			return -1
		}
	case num:
		switch b := b.(type) {
		case null, boolean:
			return +1
		case num:
			return a.compare(b)
		default:
			return -1
		}
	case str:
		switch b := b.(type) {
		case null, boolean, num:
			return +1
		case str:
			return a.compare(b)
		default:
			return -1
		}
	case array:
		switch b := b.(type) {
		case null, boolean, num, str:
			return +1
		case array:
			return a.compare(b)
		default:
			return -1
		}
	case object:
		switch b := b.(type) {
		case null, boolean, num, str, array:
			return +1
		case object:
			return a.compare(b)
		default:
			return -1
		}
	default:
		return -1
	}
}

func (a null) compare(other valueInterface) int {
	if IsNull(other) {
		return 0
	}
	return -1
}

func (a boolean) compare(other valueInterface) int {
	if b := other.(boolean); a != b {
		if !a {
			return -1
		}
		return +1
	}
	return 0
}

func (a num) compare(other valueInterface) int {
	if b := other.(num); a != b {
		if a < b {
			return -1
		}
		return +1
	}
	return 0
}

func (a str) compare(other valueInterface) int {
	if b := other.(str); a != b {
		if a < b {
			return -1
		}
		return +1
	}
	return 0
}

func (a array) compare(other valueInterface) int {
	b := other.(array)
	if a == nil {
		if b == nil {
			return 0
		}
		return -1
	}
	if b == nil {
		return +1
	}
	if &a == &b {
		return 0
	}
	if len(a) != len(b) {
		if len(a) < len(b) {
			return -1
		}
		return +1
	}
	for i := range a {
		if cmp := compare(a.At(i), b.At(i)); cmp != 0 {
			return cmp
		}
	}
	return 0
}

func (a object) compare(other valueInterface) int {
	b := other.(object)
	if a == nil {
		if b == nil {
			return 0
		}
		return -1
	}
	if b == nil {
		return +1
	}
	if &a == &b {
		return 0
	}
	if len(a) != len(b) {
		if len(a) < len(b) {
			return -1
		}
		return +1
	}
	ka := maps.Keys(a)
	kb := maps.Keys(b)
	slices.Sort(ka)
	slices.Sort(kb)
	if cmp := slices.Compare(ka, kb); cmp != 0 {
		return cmp
	}
	for _, k := range ka {
		if cmp := compare(a.At(k), b.At(k)); cmp != 0 {
			return cmp
		}
	}
	return 0
}
