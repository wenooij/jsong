package jsong

// Extract the path from the value v or return nil if not present.
//
// JSON paths are field or array indices joined by the dot character.
// The empty path returns the input value processed by ValueOf.
func Extract(v any, path string) any {
	var rv valueInterface
	if x, ok := v.(valueInterface); ok {
		rv = x
	} else if v := ValueOf(v); v != nil {
		rv = v.(valueInterface)
	} else {
		return nil
	}
	return extractRec(rv, path)
}

func extractRec(rv valueInterface, path string) valueInterface {
	if path == "" {
		return rv
	}
	head, tail, leaf := Cut(path)
	if head == "" {
		return nil
	}
	if !leaf && tail == "" {
		return nil
	}
	rv, ok := rv.Get(head)
	if !ok {
		return nil
	}
	return extractRec(rv, tail)
}
