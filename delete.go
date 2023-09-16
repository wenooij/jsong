package jsong

// Delete the path from the value v and return the result.
//
// The empty path returns nil.
func Delete(v any, path string) any {
	if path == "" {
		return nil
	}
	var rv valueInterface
	if x, ok := v.(valueInterface); ok {
		rv = x
	} else {
		rv = ValueOf(v).(valueInterface)
	}
	return deleteRec(rv, path)
}

func deleteRec(rv valueInterface, path string) valueInterface {
	head, tail, leaf := Cut(path)
	if !leaf && tail == "" {
		return nil
	}
	if leaf {
		rv.Delete(head)
		return rv
	}
	next, ok := rv.Get(head)
	if !ok {
		return nil
	}
	if o, ok := rv.(object); ok {
		o[head.(string)] = deleteRec(next, tail)
		return o
	}
	if a, ok := rv.(array); ok {
		if i := head.(int64); i < int64(len(a)) {
			a[head.(int64)] = deleteRec(next, tail)
		}
	}
	return rv
}
