package jsong

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const dot = rune('.')

const reserved = ".*"

func quoteHint(k string) bool { return strings.HasPrefix(k, `"`) }

func CutKey(k string) (head any, tail string, leaf bool) {
	if quoteHint(k) {
		panic("CutKey: unquote is not implemented yet!")
	}
	s, tail, found := strings.Cut(k, string(dot))
	leaf = !found
	if i, ok := index(s); ok {
		return i, tail, leaf
	}
	return s, tail, leaf
}

func IsLeaf(k string) bool {
	return !strings.ContainsRune(k, dot)
}

func indexHint(k string) bool {
	return len(k) > 0 && '0' <= k[0] && k[0] <= '9'
}

func index(k string) (int64, bool) {
	if indexHint(k) {
		if i, err := strconv.ParseInt(k, 10, 64); err == nil {
			return i, true
		}
	}
	return 0, false
}

// JoinKey appends the key args to the base.
//
// Args should be either string or int64
// or else JoinKey panics.
func JoinKey(base string, as ...any) string {
	var sb strings.Builder
	sb.WriteString(base)
	for i, a := range as {
		if i > 0 || base != "" {
			sb.WriteRune(dot)
		}
		switch a := a.(type) {
		case int64:
			fmt.Fprint(&sb, a)
		case string:
			if indexHint(a) || strings.ContainsAny(a, reserved) {
				sb.WriteString(strconv.Quote(a))
				continue
			}
			sb.WriteString(a)
		default:
			panic(fmt.Errorf("JoinKey: unexpected type in key at %T", a))
		}
	}
	return sb.String()
}

type KeyMatcher struct{ r *regexp.Regexp }

func CompileKeyMatcher(glob string) (*KeyMatcher, error) {
	glob = regexp.QuoteMeta(glob)
	glob = strings.ReplaceAll(glob, `\*\*`, ".*")
	glob = strings.ReplaceAll(glob, `\*`, "[^.]*")
	r, err := regexp.Compile(fmt.Sprint("^", glob, "$"))
	if err != nil {
		return nil, err
	}
	return &KeyMatcher{r: r}, nil
}

func (m *KeyMatcher) MatchKey(k string) bool {
	return m.r.MatchString(k)
}
