package conversions

import (
	"errors"

	"github.com/Robert-Safin/go-extra-types/option"
	"github.com/Robert-Safin/go-extra-types/result"
)

func Otr[T any](opt option.Option[T], customErr error) result.Result[T] {
	if opt.IsNone() {
		if customErr == nil {
			return result.NewErr[T](errors.New("None value"))
		} else {
			return result.NewErr[T](customErr)
		}
	}
	return result.NewOk(opt.Unwrap())
}

func Rto[T any](res result.Result[T]) option.Option[T] {
	if res.IsErr() {
		return option.NoneOption[T]()
	}
	return option.SomeOption[T](res.Unwrap())
}
