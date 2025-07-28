package result

import (
	"fmt"
)

type Result[T any] struct {
	val T
	err error
}

func NewInfer[T any](val T, err error) Result[T] {
	var zeroValue T
	if err != nil {
		return Result[T]{val: zeroValue, err: err}
	}
	return Result[T]{val: val, err: nil}

}

func NewOk[T any](val T) Result[T] {
	return Result[T]{
		val: val,
		err: nil,
	}
}
func NewErr[T any](err error) Result[T] {
	var zeroValue T
	return Result[T]{
		val: zeroValue,
		err: err,
	}
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func (r Result[T]) Error() error {
	return r.err
}

func (r Result[T]) Destructure() (T, error) {
	return r.val, r.err
}

func (r Result[T]) Unwrap() T {
	if r.err == nil {
		return r.val
	}
	panic(fmt.Sprintf("Unwrap on error: %v", r.err))
}

func (r Result[T]) UnwrapOrDefault(def T) T {
	if r.err == nil {
		return r.val
	}
	return def
}

func (r Result[T]) UnwrapOrZero() T {
	if r.err == nil {
		return r.val
	}
	var zeroValue T
	return zeroValue
}

func (r Result[T]) UnwrapOrFunc(f func(r Result[T]) T) T {
	if r.err == nil {
		return r.val
	}
	return f(r)
}
