package enum

import (
	"fmt"
)

type Enum[T any] struct {
	name     string
	variants map[string]T
}

type Variant[T any] struct {
	enum    *Enum[T]
	variant string
}

func NewEnum[T any](name string, input ...any) Enum[T] {
	m := make(map[string]T, len(input)/2)

	if len(input)%2 != 0 {
		panic("Received odd number of inputs")
	}

	for i := 0; i < len(input); i += 2 {
		key, ok := input[i].(string)
		if !ok {
			panic("Enum variant names should be strings")
		}
		value, ok := input[i+1].(T)
		if !ok {
			panic("Recieved mixed types for Enum values")
		}
		m[key] = value
	}
	return Enum[T]{
		name:     name,
		variants: m,
	}
}

func (e Enum[T]) Instance(key string) (Variant[T], error) {
	_, ok := e.variants[key]
	var zeroValue Variant[T]
	if !ok {
		return zeroValue, fmt.Errorf("Enum %v does not have variant %v", e.name, key)
	}
	return Variant[T]{enum: &e, variant: key}, nil
}

func (v Variant[T]) Cmp(other Variant[T]) bool {
	return v.variant == other.variant && v.enum.name == other.enum.name
}
