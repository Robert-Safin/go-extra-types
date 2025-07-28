package option

import (
	"reflect"
)

type Option[T any] struct {
	val T
	ok  bool
}

func SomeOption[T any](val T) Option[T] {
	return Option[T]{val: val, ok: true}
}

func NoneOption[T any]() Option[T] {
	var zeroValue T
	return Option[T]{val: zeroValue, ok: false}
}

func NewInfer[T any](val T, infer ...bool) Option[T] {
	var zeroValue T
	if len(infer) == 0 {
		if reflect.ValueOf(val).IsZero() {
			return Option[T]{val: zeroValue, ok: false}
		} else {
			return Option[T]{val: val, ok: true}
		}
	}
	if infer[0] == false {
		return Option[T]{val: zeroValue, ok: false}
	}
	return Option[T]{val: val, ok: infer[0]}

}

func (o Option[T]) Destructure() (T, bool) {
	return o.val, o.IsSome()
}

func (o Option[T]) IsNone() bool {
	return o.ok == false
}

func (o Option[T]) IsSome() bool {
	return o.ok != false
}

func (o Option[T]) Unwrap() T {
	if o.IsSome() {
		return o.val
	}
	panic("Unwrap on empty Option")
}

func (o Option[T]) UnwrapOrDefault(def T) T {
	if o.IsSome() {
		return o.val
	}
	return def
}

func (o Option[T]) UnwrapOrZero() T {
	if o.IsSome() {
		return o.val
	}
	var zeroValue T
	return zeroValue
}

func (o Option[T]) UnwrapOrFunc(f func(o Option[T]) T) T {
	if o.IsSome() {
		return o.val
	}
	return f(o)
}
