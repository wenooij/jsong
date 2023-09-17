package jsong

// Delete the path from the value v and return the result.
//
// The empty path returns nil.
func Delete(v any, path string) any {
	if path == "" || v == nil || v == (null{}) {
		return null{}
	}
	rv, ok := v.(valueInterface)
	if !ok {
		rv = ValueOf(v).(valueInterface)
	}
	return deleteImpl(rv, path)
}

func deleteImpl(rv valueInterface, path string) valueInterface {
	if head, tail, leaf := CutKey(path); leaf {
		rv.Delete(head)
	} else if next, ok := rv.Get(head); ok {
		rv.Put(head, deleteImpl(next, tail))
	}
	return rv
}
