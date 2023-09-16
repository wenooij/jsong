package jsong

import (
	"strings"
	"sync"
)

type Reducer interface {
	Add(any)
	Value() any
}

type HashReducer[T comparable] struct {
	New        func() Reducer
	Hash       func(any) T
	Partitions map[T]Reducer
}

func (a *HashReducer[T]) Add(x any) {
	h := a.Hash(x)
	r, ok := a.Partitions[h]
	if !ok {
		r = a.New()
		a.Partitions[h] = r
	}
	r.Add(x)
}

func (a *HashReducer[T]) Value() any {
	res := make(array, 0, len(a.Partitions))
	for _, r := range a.Partitions {
		res = append(res, r.Value())
	}
	return res
}

type PartitionReducer struct {
	New        func() Reducer
	Key        string
	partitions map[valueInterface]Reducer
	once       sync.Once
}

func (a *PartitionReducer) Reset() {
	a.partitions = nil
	a.once = sync.Once{}
}

func (a *PartitionReducer) Add(x any) {
	h := Extract(x, a.Key).(valueInterface)
	a.once.Do(func() { a.partitions = make(map[valueInterface]Reducer) })
	r, ok := a.partitions[h]
	if !ok {
		r = a.New()
		a.partitions[h] = r
	}
	r.Add(x)
}

func (a *PartitionReducer) Value() any {
	res := make(array, 0, len(a.partitions))
	for _, r := range a.partitions {
		res = append(res, r.Value())
	}
	return res
}

type ObjectReducer map[string]Reducer

func (a ObjectReducer) Add(x any) {
	val := ValueOf(x).(object)
	for k, r := range a {
		r.Add(val[k])
	}
}

func (a ObjectReducer) Value() any {
	res := make(object, len(a))
	for k, r := range a {
		res[k] = r.Value()
	}
	return res
}

type ArrayReducer []Reducer

func (a ArrayReducer) Add(x any) {
	val := ValueOf(x).(array)
	for i, r := range a {
		r.Add(val[i])
	}
}

func (a ArrayReducer) Value() any {
	res := make(array, len(a))
	for i, r := range a {
		res[i] = r.Value()
	}
	return res
}

type StringAgg struct {
	strings.Builder
}

func (a *StringAgg) Add(x any) {
	a.Builder.WriteString(x.(string))
}

type NullReducer struct{}

func (NullReducer) Add(any)    {}
func (NullReducer) Value() any { return nil }

type ReduceOp string

const (
	ReduceUndefined ReduceOp = ""
	ReduceSum       ReduceOp = "sum"
	ReduceMin       ReduceOp = "min"
	ReduceMax       ReduceOp = "max"
	ReduceAny       ReduceOp = "any"
	ReduceMean      ReduceOp = "avg"
)

type NumericReducer struct {
	Op  ReduceOp
	val float64
	cnt int
	set bool
}

func (a *NumericReducer) Add(x any) {
	v := x.(float64)
	switch a.Op {
	case ReduceUndefined, ReduceSum:
		a.set = true
		a.val += v
	case ReduceMin:
		if !a.set || v < a.val {
			a.set = true
			a.val = v
		}
	case ReduceMax:
		if !a.set || v > a.val {
			a.set = true
			a.val = v
		}
	case ReduceAny:
		if !a.set {
			a.set = true
			a.val = v
		}
	case ReduceMean:
		a.set = true
		a.val += v
		a.cnt++
	}
}

func (a *NumericReducer) Value() any {
	if a.Op == ReduceMean {
		return a.val / float64(a.cnt)
	}
	if !a.set {
		return (*num)(nil)
	}
	return num(a.val)
}

type SumReducer struct {
	sum float64
}

func (a *SumReducer) Add(x any) {
	a.sum += Must(Float64(x))
}

type TrueCounter struct {
	count int
}

func (a *TrueCounter) Add(x any) {
	if x.(bool) {
		a.count++
	}
}

func (a *StringAgg) Value() any   { return str(a.String()) }
func (a *SumReducer) Value() any  { return num(a.sum) }
func (a *TrueCounter) Value() any { return num(a.count) }

type AnyReducer struct{ V any }

func (a *AnyReducer) Add(v any)  { a.V = v }
func (a *AnyReducer) Value() any { return a.V }
