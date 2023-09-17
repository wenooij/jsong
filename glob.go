package jsong

const (
	star       = '*'
	doubleStar = "**"
)

type globStackEntry struct {
	all  bool // **
	path string
	tail string
	v    valueInterface
}

func GlobKey(v any, k string) []string {
	val, ok := v.(valueInterface)
	if !ok {
		val = ValueOf(v).(valueInterface)
	}
	return globKey(val, k)
}

func globKey(v valueInterface, k string) []string {
	var results []string
	var stack []globStackEntry
	stack = append(stack, globStackEntry{
		tail: k,
		v:    v,
	})
	for len(stack) > 0 {
		n := len(stack) - 1
		e := stack[n]
		stack = stack[:n]
		k := e.tail
		v := e.v
		if e.all {
			if k == "" {
				results = append(results, e.path)
			}
			v.Each(func(i any) bool {
				path := JoinKey(e.path, i)
				stack = append(stack, globStackEntry{
					all:  true,
					path: path,
					tail: k,
					v:    Must(v.Get(i)),
				})
				return true
			})
			head, tail, _ := CutKey(k)
			if head == doubleStar || head == string(star) { // Remove redundant stars.
				stack = append(stack, globStackEntry{
					all:  true,
					path: e.path,
					tail: tail,
					v:    v,
				})
				continue
			}
			if v, ok := v.Get(head); ok {
				stack = append(stack, globStackEntry{
					path: JoinKey(e.path, head),
					tail: tail,
					v:    v,
				})
			}
			continue
		}
		if k == "" {
			results = append(results, e.path)
			continue
		}
		head, tail, _ := CutKey(k)
		if head == string(star) {
			v.Each(func(i any) bool {
				stack = append(stack, globStackEntry{
					path: JoinKey(e.path, i),
					tail: tail,
					v:    Must(v.Get(i)),
				})
				return true
			})
			continue
		}
		if head == doubleStar {
			v.Each(func(i any) bool {
				stack = append(stack, globStackEntry{
					all:  true,
					path: JoinKey(e.path, i),
					tail: tail,
					v:    Must(v.Get(i)),
				})
				return true
			})
			continue
		}
		if v, ok := v.Get(head); ok {
			stack = append(stack, globStackEntry{
				path: JoinKey(e.path, head),
				tail: tail,
				v:    v,
			})
		}
	}
	return results
}
