package conversions

import (
	"github.com/Robert-Safin/go-lib/option"
	"github.com/Robert-Safin/go-lib/result"
)

func Otr[T any](opt option.Option[T], customErr error) result.Result[T] {
	if opt.IsNone() {
		return result.NewErr[T](customErr)
	}
	return result.NewOk(opt.Unwrap())
}

func Rto[T any](res result.Result[T]) option.Option[T] {
	if res.IsErr() {
		return option.NoneOption[T]()
	}
	return option.SomeOption[T](res.Unwrap())
}
