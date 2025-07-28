package enum

import (
	"fmt"
	"maps"
)

type Enum[T any] struct {
	name     string
	variants map[string]T
}

type Variant[T any] struct {
	enum  string
	name  string
	value T
}

func NewEnum[T any](name string, variants map[string]T) Enum[T] {
	if name == "" {
		panic("Enum name cannot be empty")
	}
	if len(variants) == 0 {
		panic("Enum must have at least one variant")
	}
	copy := map[string]T{}
	maps.Copy(copy, variants)
	return Enum[T]{
		name:     name,
		variants: copy,
	}
}

func (e Enum[T]) NewInstance(name string) Variant[T] {
	if name == "" {
		panic("Variant name cannot be empty")
	}
	v, ok := e.variants[name]
	if !ok {
		panic(fmt.Sprintf("Enum %v does not have variant %v\n", e.name, name))
	}
	return Variant[T]{enum: e.name, name: name, value: v}
}

func (e Enum[T]) VariantNames() []string {
	names := make([]string, 0, len(e.variants))
	for name := range e.variants {
		names = append(names, name)
	}
	return names
}

func (e Enum[T]) String() string {
	return fmt.Sprintf("Enum{name: %s, variant count: %v, variants: %v}", e.name, len(e.variants), e.variants)
}

func (v Variant[T]) IsInstanceOf(enum Enum[T]) bool {
	_, ok := enum.variants[v.name]
	return ok && enum.name == v.enum
}

func (v Variant[T]) Value() T {
	return v.value
}
func (v Variant[T]) Name() string {
	return v.name
}

func (v Variant[T]) String() string {
	return fmt.Sprintf("Variant{enum: %s, name: %s}", v.enum, v.name)
}
