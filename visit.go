package jsong

import "errors"

var (
	ErrSkip = errors.New("skip")
	ErrStop = errors.New("stop")
)

func Visit(v any, visitFn func(k string, v any) error) error {
	val, ok := v.(valueInterface)
	if !ok {
		val = ValueOf(v).(valueInterface)
	}
	if err := visit("", val, visitFn); err != nil && err != ErrStop {
		return err
	}
	return nil
}

func visit(path string, v valueInterface, visitFn func(k string, v any) error) error {
	if err := visitFn(path, v); err != nil {
		if err == ErrSkip {
			return nil
		}
		return err
	}
	var err error
	v.Each(func(k, v any) bool {
		if err = visit(JoinKey(path, k), v.(valueInterface), visitFn); err != nil {
			return false
		}
		return true
	})
	return err
}
