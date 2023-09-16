package jsong

import "fmt"

type Mapper interface {
	Map(v any) any
}

type Ident struct{}

func Map(v any) any { return v }

type MapSeq []Mapper

func (a MapSeq) Map(v any) any {
	for _, e := range a {
		v = e.Map(v)
	}
	return v
}

type MulScalar struct {
	M any
}

func (a MulScalar) Map(v any) any {
	return v.(num) * a.M.(num)
}

type AddScalar struct {
	C any
}

func (a AddScalar) Map(v any) any {
	if s, ok := v.(str); ok {
		return s + a.C.(str)
	}
	if n, ok := v.(num); ok {
		return num(n) + a.C.(num)
	}
	return v
}

type MathMapper struct {
	Fn func(float64) float64
}

func (a MathMapper) Map(v any) any {
	return num(a.Fn(float64(v.(num))))
}

type Math2Mapper struct {
	Fn2 func(float64, float64) float64
}

func (a Math2Mapper) Map(v any) any {
	es := v.(array)
	return num(a.Fn2(float64(es[0].(num)), float64(es[1].(num))))
}

type ObjectMapper map[string]Mapper

func (a ObjectMapper) Map(v any) any {
	val := ValueOf(v).(object)
	for k, m := range a {
		e := val.At(k)
		val = Merge(val, m.Map(e), k, "").(object)
	}
	return v
}

type ArrayMapper []Mapper

func (a ArrayMapper) Map(v any) any {
	val := ValueOf(v).(array)
	for i, m := range a {
		val[i] = m.Map(val[i])
	}
	return val
}

type ArrayRemapper []any

func (a ArrayRemapper) Map(v any) any {
	src := ValueOf(v)
	dst := make(array, len(a))
	for i, e := range a {
		switch e := e.(type) {
		case valueInterface:
			dst[i] = e
		case Mapper:
			dst[i] = e.Map(src.(array).At(i))
		case string:
			dst[i] = Extract(src, e)
		default:
			panic(fmt.Errorf("ArrayRemapper: unexpected element at %d: %T", i, e))
		}
	}
	return dst
}

type ObjectRemapper map[string]any

func (a ObjectRemapper) Map(v any) any {
	src := ValueOf(v)
	dst := make(object, len(a))
	for k, e := range a {
		switch e := e.(type) {
		case valueInterface:
			dst[k] = e
		case Mapper:
			dst[k] = e.Map(src.(object).At(k))
		case string:
			dst[k] = Extract(src, k)
		default:
			panic(fmt.Errorf("ObjectRemapper: unexpected entry at %q: %T", k, e))
		}
	}
	return dst
}
