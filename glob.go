package jsong

const (
	star       = '*'
	doubleStar = "**"
)

func GlobKey(v any, k string) []string {
	var results []string
	Glob(v, k, func(k string, _ any) { results = append(results, k) })
	return results
}

func GlobValues(v any, k string) []any {
	var results []any
	Glob(v, k, func(_ string, v any) { results = append(results, v) })
	return results
}

func Glob(v any, glob string, visitFn func(k string, v any)) {
	val, ok := v.(valueInterface)
	if !ok {
		val = ValueOf(v).(valueInterface)
	}
	m := Must(CompileKeyMatcher(glob))
	visit("", val, func(k string, v any) error {
		if !m.MatchKey(k) {
			return ErrSkip
		}
		return nil
	})
}
