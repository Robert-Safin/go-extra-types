package iter

import "math/rand"

type Iter[T any] []T

// Constructors and Conversions and otherwise functions
func NewIter[T any](slice []T) Iter[T] {
	dst := make([]T, len(slice))
	copy(dst, slice)
	return Iter[T](dst)
}

func Map[T, U any](i Iter[T], applyFunc func(item T) U) Iter[U] {
	res := make([]U, len(i))
	for idx, item := range i {
		res[idx] = applyFunc(item)
	}
	return Iter[U](res)
}

func ToMap[K comparable, V any, T any](iterator Iter[T], conversionFunc func(index int, item T) (K, V)) map[K]V {
	m := make(map[K]V, len(iterator))
	for i, v := range iterator {
		key, value := conversionFunc(i, v)
		m[key] = value
	}
	return m
}
func Reduce[T any, U any](iterator Iter[T], reduceFunc func(acc U, item T) U, initial U) U {
	for _, item := range iterator {
		initial = reduceFunc(initial, item)
	}
	return initial
}

func Deduped[T comparable](i Iter[T]) Iter[T] {
	if len(i) <= 1 {
		return NewIter(i) // Return copy
	}

	result := make([]T, 0, len(i))
	seen := make(map[T]bool)

	for _, item := range i {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return Iter[T](result)
}

type Enumerated[T any] struct {
	Index int
	Value T
}

func Enumerate[T any](iter Iter[T]) Iter[Enumerated[T]] {
	res := make([]Enumerated[T], len(iter))
	for i, v := range iter {
		res[i] = Enumerated[T]{Index: i, Value: v}
	}
	return Iter[Enumerated[T]](res)
}

type Pair[T, U any] struct {
	First  T
	Second U
}

func Zip[T, U any](i Iter[T], other Iter[U]) Iter[Pair[T, U]] {
	minLen := min(len(i), len(other))
	result := make([]Pair[T, U], minLen)

	for idx := range minLen {
		result[idx] = Pair[T, U]{First: i[idx], Second: other[idx]}
	}

	return Iter[Pair[T, U]](result)
}
func Sum[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](i Iter[T]) T {
	var sum T
	for _, v := range i {
		sum += v
	}
	return sum
}

func Max[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](i Iter[T]) (T, bool) {
	if len(i) == 0 {
		var zero T
		return zero, false
	}
	max := i[0]
	for _, v := range i {
		if v > max {
			max = v
		}
	}
	return max, true
}

func Min[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64](i Iter[T]) (T, bool) {
	if len(i) == 0 {
		var zero T
		return zero, false
	}
	min := i[0]
	for _, v := range i {
		if v < min {
			min = v
		}
	}
	return min, true
}

// Navigation & Control
func (i *Iter[T]) Next() bool {
	if len(*i) == 0 {
		return false
	}
	*i = (*i)[1:]
	return true
}

func (i *Iter[T]) Skip(n uint) {
	if n >= uint(len(*i)) {
		*i = (*i)[:0]
		return
	}
	*i = (*i)[n:]
}

func (i *Iter[T]) Take(n uint) {
	if n > uint(len(*i)) {
		return
	}
	*i = (*i)[:n]
}

func (i Iter[T]) First() (T, bool) {
	if len(i) == 0 {
		var zero T
		return zero, false
	}
	return i[0], true
}
func (i Iter[T]) Last() (T, bool) {
	if len(i) == 0 {
		var zero T
		return zero, false
	}
	return i[len(i)-1], true
}

func (i Iter[T]) Nth(n uint) (T, bool) {
	var zero T
	if len(i) == 0 {
		return zero, false
	}
	if n > uint(len(i)-1) {

		return zero, false
	}
	return i[n], true
}

// Transformation
func (i *Iter[T]) ForEach(applyFunc func(item *T)) {
	for idx := range *i {
		applyFunc(&(*i)[idx])
	}
}

type Group[T any] struct {
	Key   string
	Items []T
}

// maintains order
func (i Iter[T]) GroupBy(groupFunc func(item T) string) []Group[T] {
	groups := make(map[string][]T)
	for _, item := range i {
		key := groupFunc(item)
		groups[key] = append(groups[key], item)
	}
	var result []Group[T]
	for key, items := range groups {
		result = append(result, Group[T]{Key: key, Items: items})
	}
	return result

}

// Filtering, Searching, Slicing
func (i *Iter[T]) Filter(filterFunc func(item T) bool) {
	res := []T{}
	for _, item := range *i {
		if filterFunc(item) {
			res = append(res, item)
		}
	}
	*i = res
}
func (i Iter[T]) FindOne(searchFunc func(item T) bool) (int, T, bool) {
	for index, item := range i {
		if searchFunc(item) {
			return index, item, true
		}
	}
	var zero T
	return 0, zero, false
}

type Result[T any] struct {
	Index int
	Item  T
}

func (i Iter[T]) FindMany(searchFunc func(item T) bool) []Result[T] {
	res := make([]Result[T], 0, len(i))
	for index, item := range i {
		if searchFunc(item) {
			res = append(res, Result[T]{Index: index, Item: item})
		}
	}
	return res
}
func (i Iter[T]) All(condFunc func(item T) bool) bool {
	for _, item := range i {
		if !condFunc(item) {
			return false
		}
	}
	return true
}
func (i Iter[T]) Any(condFunc func(item T) bool) bool {
	for _, item := range i {
		if condFunc(item) {
			return true
		}
	}
	return false
}

func (i Iter[T]) ContainsFunc(targetFunc func(item T) bool) bool {
	for _, item := range i {
		if targetFunc(item) {
			return true
		}
	}
	return false
}
func (i *Iter[T]) Reverse() {
	for left, right := 0, len(*i)-1; left < right; left, right = left+1, right-1 {
		(*i)[left], (*i)[right] = (*i)[right], (*i)[left]
	}
}
func (i *Iter[T]) Shuffle() {
	for j := len(*i) - 1; j > 0; j-- {
		k := rand.Intn(j + 1)
		(*i)[j], (*i)[k] = (*i)[k], (*i)[j]
	}
}

// Collection & Aggregation
func (i Iter[T]) Collect() []T {
	dst := make([]T, len(i))
	copy(dst, i)
	return dst
}

func (i Iter[T]) SumFunc(sumFunc func(item T) int) int {
	var sum int
	for _, v := range i {
		sum += sumFunc(v)
	}
	return sum
}

func (i *Iter[T]) Chain(other Iter[T]) {
	for _, item := range other {
		*i = append(*i, item)
	}
}

// Stateful or Advanced Patterns
func (i *Iter[T]) Windows(size uint, applyFunc func(window []T)) {
	if size == 0 || size > uint(len(*i)) {
		return
	}

	for start := 0; start <= len(*i)-int(size); start++ {
		window := (*i)[start : start+int(size)]
		applyFunc(window)
	}
}
func (i *Iter[T]) Chunks(size uint, applyFunc func(chunk []T)) {
	if size == 0 {
		return
	}

	if size > uint(len(*i)) {
		applyFunc(*i)
		return
	}

	for start := 0; start < len(*i); start += int(size) {
		end := min(start+int(size), len(*i))
		chunk := (*i)[start:end]
		applyFunc(chunk)
	}
}

func (i *Iter[T]) Cycle(applyFunc func(item T), breakFunc func() bool) {
	if len(*i) == 0 {
		return
	}
	index := 0
	for {
		if breakFunc() {
			break
		}
		applyFunc((*i)[index])
		index = (index + 1) % len(*i)
	}
}
