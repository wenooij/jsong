package jsong

import (
	"errors"
	"strconv"
	"strings"
)

const (
	Sep  = rune('.')
	Glob = rune('*')
)

var (
	ErrMaxDepth = errors.New("max depth reached")
	ErrMaxIndex = errors.New("max index reached")
)

type Limits struct {
	MaxDepth int
	MaxIndex int64
}

func DefaultLimits() Limits {
	return Limits{
		MaxDepth: 6,
		MaxIndex: 32,
	}
}

func Each(k string, lim Limits, visitFn func(head any)) error {
	for depth := 1; ; depth++ {
		if lim.MaxDepth > 0 && lim.MaxDepth < depth {
			return ErrMaxDepth
		}
		head, tail, leaf := Cut(k)
		if i, ok := head.(int64); ok {
			if lim.MaxIndex < i {
				return ErrMaxIndex
			}
		}
		visitFn(head)
		if leaf {
			break
		}
		k = tail
	}
	return nil
}

func Cut(k string) (head any, tail string, leaf bool) {
	s, tail, found := strings.Cut(k, string(Sep))
	leaf = !found
	if i, ok := index(s); ok {
		return i, tail, leaf
	}
	return s, tail, leaf
}

func Leaf(k string) bool {
	return !strings.ContainsRune(k, Sep)
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
