package util

func Ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func Times(times int, f func()) {
	for range times {
		f()
	}
}
