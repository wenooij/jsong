package jsong

// Merge the field of the JSON object data with the MergeOptions.
// It returns the resulting JSON data or any merge errors.
func Merge(dst, src any, dstPath, srcPath string) any {
	if _, ok := dst.(valueInterface); !ok {
		dst = ValueOf(dst)
	}
	src = Extract(src, srcPath)
	return mergeRec(dst.(valueInterface), src.(valueInterface), dstPath)
}

func mergeRec(dst, src valueInterface, dstPath string) valueInterface {
	if dstPath == "" {
		return src
	}
	head, tail, leaf := CutKey(dstPath)
	if head == "" {
		return dst
	}
	if !leaf && tail == "" {
		return dst
	}
	switch dst := dst.(type) {
	case null:
		switch head := head.(type) {
		case int64:
			return mergeRec(make(array, head+1), src, dstPath)
		case string:
			return mergeRec(make(object), src, dstPath)
		default:
			return dst
		}
	case array:
		i := head.(int64)
		if int64(len(dst)) <= i {
			return dst
		}
		if leaf {
			dst[i] = src
			return dst
		}
		return mergeRec(dst.At(int(i)), src, tail)
	case object:
		if leaf {
			dst[head.(string)] = src
			return dst
		}
		v, ok := dst[head.(string)]
		if !ok {
			v = make(object)
		}
		return mergeRec(v.(valueInterface), src, tail)
	default:
		return dst
	}
}
