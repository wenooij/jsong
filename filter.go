package jsong

type Filter interface {
	Filter(any) bool
}

// ObjectFieldFilter returns a Filter which returns true
// only if the object has all the given fields.
type ObjectFieldFilter map[string]struct{}

func (f ObjectFieldFilter) Filter(v any) bool {
	val, ok := ValueOf(v).(object)
	if !ok {
		return false
	}
	for k := range f {
		if _, ok := val.Get(k); !ok {
			return false
		}
	}
	return true
}

type GlobFilter struct {
	Glob string
}

func (f GlobFilter) Filter(v any) bool {
	var match bool
	Glob(v, f.Glob, func(string, any) { match = true })
	return match
}
