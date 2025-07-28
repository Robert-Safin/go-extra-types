package conversions_test

import (
	"errors"
	"testing"

	"github.com/Robert-Safin/go-extra-types/conversions"
	"github.com/Robert-Safin/go-extra-types/option"
	"github.com/Robert-Safin/go-extra-types/result"
)

func TestOtr(t *testing.T) {
	t.Run("converts Some option to Ok result", func(t *testing.T) {
		opt := option.SomeOption("hello")
		customErr := errors.New("custom error")
		res := conversions.Otr(opt, customErr)

		if !res.IsOk() {
			t.Error("Expected result to be Ok when option is Some")
		}

		if res.IsErr() {
			t.Error("Expected result to not be Err when option is Some")
		}

		if res.Unwrap() != "hello" {
			t.Errorf("Expected value 'hello', got %s", res.Unwrap())
		}

		if res.Error() != nil {
			t.Error("Expected error to be nil for Ok result")
		}
	})

	t.Run("converts None option to Err result", func(t *testing.T) {
		opt := option.NoneOption[string]()
		customErr := errors.New("custom error")
		res := conversions.Otr(opt, customErr)

		if res.IsOk() {
			t.Error("Expected result to not be Ok when option is None")
		}

		if !res.IsErr() {
			t.Error("Expected result to be Err when option is None")
		}

		if res.Error() != customErr {
			t.Errorf("Expected error to be %v, got %v", customErr, res.Error())
		}
	})

	t.Run("preserves value type", func(t *testing.T) {
		intOpt := option.SomeOption(42)
		customErr := errors.New("custom error")
		res := conversions.Otr(intOpt, customErr)

		if !res.IsOk() {
			t.Error("Expected result to be Ok")
		}

		if res.Unwrap() != 42 {
			t.Errorf("Expected value 42, got %d", res.Unwrap())
		}
	})

	t.Run("works with complex types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		person := Person{Name: "Alice", Age: 30}
		opt := option.SomeOption(person)
		customErr := errors.New("person not found")
		res := conversions.Otr(opt, customErr)

		if !res.IsOk() {
			t.Error("Expected result to be Ok")
		}

		result := res.Unwrap()
		if result.Name != "Alice" || result.Age != 30 {
			t.Errorf("Expected person Alice(30), got %+v", result)
		}
	})

	t.Run("uses provided custom error", func(t *testing.T) {
		opt := option.NoneOption[string]()
		customErr := errors.New("specific error message")
		res := conversions.Otr(opt, customErr)

		if res.Error() != customErr {
			t.Errorf("Expected custom error %v, got %v", customErr, res.Error())
		}

		if res.Error().Error() != "specific error message" {
			t.Errorf("Expected error message 'specific error message', got %s", res.Error().Error())
		}
	})

	t.Run("works with nil custom error", func(t *testing.T) {
		opt := option.NoneOption[string]()
		res := conversions.Otr(opt, nil)

		if !res.IsErr() {
			t.Error("Expected result to be Err")
		}

		if res.Error() == nil {
			t.Errorf("Expected error to be nil, got %v", res.Error())
		}
	})
}

func TestRto(t *testing.T) {
	t.Run("converts Ok result to Some option", func(t *testing.T) {
		res := result.NewOk("hello")
		opt := conversions.Rto(res)

		if !opt.IsSome() {
			t.Error("Expected option to be Some when result is Ok")
		}

		if opt.IsNone() {
			t.Error("Expected option to not be None when result is Ok")
		}

		if opt.Unwrap() != "hello" {
			t.Errorf("Expected value 'hello', got %s", opt.Unwrap())
		}
	})

	t.Run("converts Err result to None option", func(t *testing.T) {
		err := errors.New("test error")
		res := result.NewErr[string](err)
		opt := conversions.Rto(res)

		if opt.IsSome() {
			t.Error("Expected option to not be Some when result is Err")
		}

		if !opt.IsNone() {
			t.Error("Expected option to be None when result is Err")
		}
	})

	t.Run("preserves value type", func(t *testing.T) {
		res := result.NewOk(42)
		opt := conversions.Rto(res)

		if !opt.IsSome() {
			t.Error("Expected option to be Some")
		}

		if opt.Unwrap() != 42 {
			t.Errorf("Expected value 42, got %d", opt.Unwrap())
		}
	})

	t.Run("works with complex types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		person := Person{Name: "Bob", Age: 25}
		res := result.NewOk(person)
		opt := conversions.Rto(res)

		if !opt.IsSome() {
			t.Error("Expected option to be Some")
		}

		result := opt.Unwrap()
		if result.Name != "Bob" || result.Age != 25 {
			t.Errorf("Expected person Bob(25), got %+v", result)
		}
	})

	t.Run("discards error information", func(t *testing.T) {
		err := errors.New("important error")
		res := result.NewErr[string](err)
		opt := conversions.Rto(res)

		if !opt.IsNone() {
			t.Error("Expected option to be None")
		}

		// Error information is lost in conversion - this is expected behavior
		// We can only test that the option is None
	})
}

func TestRoundTripConversions(t *testing.T) {
	t.Run("Some -> Ok -> Some preserves value", func(t *testing.T) {
		original := option.SomeOption("hello")
		customErr := errors.New("error")

		// Some -> Ok
		res := conversions.Otr(original, customErr)
		// Ok -> Some
		final := conversions.Rto(res)

		if !final.IsSome() {
			t.Error("Expected final option to be Some")
		}

		if final.Unwrap() != "hello" {
			t.Errorf("Expected value 'hello', got %s", final.Unwrap())
		}
	})

	t.Run("None -> Err -> None preserves None state", func(t *testing.T) {
		original := option.NoneOption[string]()
		customErr := errors.New("error")

		// None -> Err
		res := conversions.Otr(original, customErr)
		// Err -> None
		final := conversions.Rto(res)

		if !final.IsNone() {
			t.Error("Expected final option to be None")
		}
	})

	t.Run("Ok -> Some -> Ok with error loses Ok state", func(t *testing.T) {
		original := result.NewOk("hello")
		customErr := errors.New("error")

		// Ok -> Some
		opt := conversions.Rto(original)
		// Some -> Ok (but this will always be Ok since opt is Some)
		final := conversions.Otr(opt, customErr)

		if !final.IsOk() {
			t.Error("Expected final result to be Ok")
		}

		if final.Unwrap() != "hello" {
			t.Errorf("Expected value 'hello', got %s", final.Unwrap())
		}
	})

	t.Run("Err -> None -> Err changes error", func(t *testing.T) {
		originalErr := errors.New("original error")
		newErr := errors.New("new error")
		original := result.NewErr[string](originalErr)

		// Err -> None
		opt := conversions.Rto(original)
		// None -> Err
		final := conversions.Otr(opt, newErr)

		if !final.IsErr() {
			t.Error("Expected final result to be Err")
		}

		if final.Error() != newErr {
			t.Errorf("Expected error to be %v, got %v", newErr, final.Error())
		}

		// Original error is lost - this is expected
		if final.Error() == originalErr {
			t.Error("Did not expect original error to be preserved")
		}
	})
}

func TestConversionsWithDifferentTypes(t *testing.T) {
	t.Run("integer conversions", func(t *testing.T) {
		// Option -> Result
		intOpt := option.SomeOption(42)
		err := errors.New("int error")
		intRes := conversions.Otr(intOpt, err)

		if intRes.Unwrap() != 42 {
			t.Errorf("Expected value 42, got %d", intRes.Unwrap())
		}

		// Result -> Option
		backToOpt := conversions.Rto(intRes)
		if backToOpt.Unwrap() != 42 {
			t.Errorf("Expected value 42, got %d", backToOpt.Unwrap())
		}
	})

	t.Run("boolean conversions", func(t *testing.T) {
		// Option -> Result
		boolOpt := option.SomeOption(true)
		err := errors.New("bool error")
		boolRes := conversions.Otr(boolOpt, err)

		if boolRes.Unwrap() != true {
			t.Errorf("Expected value true, got %t", boolRes.Unwrap())
		}

		// Result -> Option
		backToOpt := conversions.Rto(boolRes)
		if backToOpt.Unwrap() != true {
			t.Errorf("Expected value true, got %t", backToOpt.Unwrap())
		}
	})

	t.Run("slice conversions", func(t *testing.T) {
		slice := []string{"a", "b", "c"}

		// Option -> Result
		sliceOpt := option.SomeOption(slice)
		err := errors.New("slice error")
		sliceRes := conversions.Otr(sliceOpt, err)

		result := sliceRes.Unwrap()
		if len(result) != 3 || result[0] != "a" {
			t.Errorf("Expected slice [a,b,c], got %v", result)
		}

		// Result -> Option
		backToOpt := conversions.Rto(sliceRes)
		optResult := backToOpt.Unwrap()
		if len(optResult) != 3 || optResult[0] != "a" {
			t.Errorf("Expected slice [a,b,c], got %v", optResult)
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("zero value conversions", func(t *testing.T) {
		// Some with zero value
		zeroOpt := option.SomeOption("")
		err := errors.New("error")
		res := conversions.Otr(zeroOpt, err)

		if !res.IsOk() {
			t.Error("Expected result to be Ok even with zero value")
		}

		if res.Unwrap() != "" {
			t.Errorf("Expected empty string, got %s", res.Unwrap())
		}
	})

	t.Run("multiple conversions maintain consistency", func(t *testing.T) {
		original := option.SomeOption("test")
		err := errors.New("error")

		// Multiple round trips
		res1 := conversions.Otr(original, err)
		opt1 := conversions.Rto(res1)
		res2 := conversions.Otr(opt1, err)
		opt2 := conversions.Rto(res2)

		if !opt2.IsSome() {
			t.Error("Expected final option to be Some")
		}

		if opt2.Unwrap() != "test" {
			t.Errorf("Expected value 'test', got %s", opt2.Unwrap())
		}
	})
}
