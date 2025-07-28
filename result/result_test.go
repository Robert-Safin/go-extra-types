package result_test

import (
	"errors"
	"testing"

	"github.com/Robert-Safin/go-extra-types/result"
)

func TestNewInfer(t *testing.T) {
	t.Run("creates Ok result when error is nil", func(t *testing.T) {
		res := result.NewInfer("hello", nil)

		if !res.IsOk() {
			t.Error("Expected result to be Ok when error is nil")
		}

		if res.IsErr() {
			t.Error("Expected result to not be Err when error is nil")
		}

		if res.Unwrap() != "hello" {
			t.Errorf("Expected value 'hello', got %s", res.Unwrap())
		}

		if res.Error() != nil {
			t.Error("Expected error to be nil for Ok result")
		}
	})

	t.Run("creates Err result when error is present", func(t *testing.T) {
		err := errors.New("test error")
		res := result.NewInfer("hello", err)

		if res.IsOk() {
			t.Error("Expected result to not be Ok when error is present")
		}

		if !res.IsErr() {
			t.Error("Expected result to be Err when error is present")
		}

		if res.Error() != err {
			t.Errorf("Expected error to be %v, got %v", err, res.Error())
		}
	})

	t.Run("stores zero value when error is present", func(t *testing.T) {
		err := errors.New("test error")
		res := result.NewInfer("hello", err)

		val, resErr := res.Destructure()
		if val != "" {
			t.Errorf("Expected zero value (empty string), got %s", val)
		}

		if resErr != err {
			t.Errorf("Expected error to be %v, got %v", err, resErr)
		}
	})
}

func TestNewOk(t *testing.T) {
	t.Run("creates Ok result with value", func(t *testing.T) {
		res := result.NewOk("hello")

		if !res.IsOk() {
			t.Error("Expected result to be Ok")
		}

		if res.IsErr() {
			t.Error("Expected result to not be Err")
		}

		if res.Unwrap() != "hello" {
			t.Errorf("Expected value 'hello', got %s", res.Unwrap())
		}

		if res.Error() != nil {
			t.Error("Expected error to be nil for Ok result")
		}
	})

	t.Run("works with different types", func(t *testing.T) {
		intRes := result.NewOk(42)
		if intRes.Unwrap() != 42 {
			t.Errorf("Expected value 42, got %d", intRes.Unwrap())
		}

		boolRes := result.NewOk(true)
		if boolRes.Unwrap() != true {
			t.Errorf("Expected value true, got %t", boolRes.Unwrap())
		}
	})
}

func TestNewErr(t *testing.T) {
	t.Run("creates Err result with error", func(t *testing.T) {
		err := errors.New("test error")
		res := result.NewErr[string](err)

		if res.IsOk() {
			t.Error("Expected result to not be Ok")
		}

		if !res.IsErr() {
			t.Error("Expected result to be Err")
		}

		if res.Error() != err {
			t.Errorf("Expected error to be %v, got %v", err, res.Error())
		}
	})

	t.Run("stores zero value for different types", func(t *testing.T) {
		err := errors.New("test error")

		stringRes := result.NewErr[string](err)
		val, _ := stringRes.Destructure()
		if val != "" {
			t.Errorf("Expected zero value (empty string), got %s", val)
		}

		intRes := result.NewErr[int](err)
		intVal, _ := intRes.Destructure()
		if intVal != 0 {
			t.Errorf("Expected zero value (0), got %d", intVal)
		}
	})
}

func TestDestructure(t *testing.T) {
	t.Run("destructures Ok result", func(t *testing.T) {
		res := result.NewOk("hello")
		val, err := res.Destructure()

		if val != "hello" {
			t.Errorf("Expected value 'hello', got %s", val)
		}

		if err != nil {
			t.Errorf("Expected error to be nil, got %v", err)
		}
	})

	t.Run("destructures Err result", func(t *testing.T) {
		expectedErr := errors.New("test error")
		res := result.NewErr[string](expectedErr)
		val, err := res.Destructure()

		if val != "" {
			t.Errorf("Expected zero value, got %s", val)
		}

		if err != expectedErr {
			t.Errorf("Expected error to be %v, got %v", expectedErr, err)
		}
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("unwraps Ok result", func(t *testing.T) {
		res := result.NewOk("hello")
		val := res.Unwrap()

		if val != "hello" {
			t.Errorf("Expected value 'hello', got %s", val)
		}
	})

	t.Run("panics on Err result", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when unwrapping Err result")
			}
		}()

		err := errors.New("test error")
		res := result.NewErr[string](err)
		res.Unwrap()
	})
}

func TestUnwrapOrDefault(t *testing.T) {
	t.Run("returns value for Ok result", func(t *testing.T) {
		res := result.NewOk("hello")
		val := res.UnwrapOrDefault("default")

		if val != "hello" {
			t.Errorf("Expected value 'hello', got %s", val)
		}
	})

	t.Run("returns default for Err result", func(t *testing.T) {
		err := errors.New("test error")
		res := result.NewErr[string](err)
		val := res.UnwrapOrDefault("default")

		if val != "default" {
			t.Errorf("Expected default value 'default', got %s", val)
		}
	})
}

func TestUnwrapOrZero(t *testing.T) {
	t.Run("returns value for Ok result", func(t *testing.T) {
		res := result.NewOk("hello")
		val := res.UnwrapOrZero()

		if val != "hello" {
			t.Errorf("Expected value 'hello', got %s", val)
		}
	})

	t.Run("returns zero value for Err result", func(t *testing.T) {
		err := errors.New("test error")
		res := result.NewErr[string](err)
		val := res.UnwrapOrZero()

		if val != "" {
			t.Errorf("Expected zero value, got %s", val)
		}
	})

	t.Run("returns zero for different types", func(t *testing.T) {
		err := errors.New("test error")

		intRes := result.NewErr[int](err)
		if intRes.UnwrapOrZero() != 0 {
			t.Errorf("Expected zero int, got %d", intRes.UnwrapOrZero())
		}

		boolRes := result.NewErr[bool](err)
		if boolRes.UnwrapOrZero() != false {
			t.Errorf("Expected false, got %t", boolRes.UnwrapOrZero())
		}
	})
}

func TestUnwrapOrFunc(t *testing.T) {
	t.Run("returns value for Ok result", func(t *testing.T) {
		res := result.NewOk("hello")
		val := res.UnwrapOrFunc(func(r result.Result[string]) string {
			return "from function"
		})

		if val != "hello" {
			t.Errorf("Expected value 'hello', got %s", val)
		}
	})

	t.Run("calls function for Err result", func(t *testing.T) {
		err := errors.New("test error")
		res := result.NewErr[string](err)
		val := res.UnwrapOrFunc(func(r result.Result[string]) string {
			if !r.IsErr() {
				t.Error("Expected result passed to function to be Err")
			}
			if r.Error() != err {
				t.Errorf("Expected error to be %v, got %v", err, r.Error())
			}
			return "from function"
		})

		if val != "from function" {
			t.Errorf("Expected value 'from function', got %s", val)
		}
	})

	t.Run("function receives correct result", func(t *testing.T) {
		err := errors.New("test error")
		res := result.NewErr[int](err)
		called := false

		res.UnwrapOrFunc(func(r result.Result[int]) int {
			called = true
			if !r.IsErr() {
				t.Error("Expected received result to be Err")
			}
			if r.Error() != err {
				t.Errorf("Expected error to be %v, got %v", err, r.Error())
			}
			return 42
		})

		if !called {
			t.Error("Expected function to be called")
		}
	})
}

func TestResultWithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("struct result", func(t *testing.T) {
		person := Person{Name: "Alice", Age: 30}
		res := result.NewOk(person)

		if !res.IsOk() {
			t.Error("Expected result to be Ok")
		}

		result := res.Unwrap()
		if result.Name != "Alice" || result.Age != 30 {
			t.Errorf("Expected person Alice(30), got %+v", result)
		}
	})

	t.Run("slice result", func(t *testing.T) {
		slice := []int{1, 2, 3}
		res := result.NewOk(slice)

		if !res.IsOk() {
			t.Error("Expected result to be Ok")
		}

		result := res.Unwrap()
		if len(result) != 3 || result[0] != 1 {
			t.Errorf("Expected slice [1,2,3], got %v", result)
		}
	})

	t.Run("map result", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2}
		res := result.NewOk(m)

		if !res.IsOk() {
			t.Error("Expected result to be Ok")
		}

		result := res.Unwrap()
		if result["a"] != 1 || result["b"] != 2 {
			t.Errorf("Expected map with a:1, b:2, got %v", result)
		}
	})
}

func TestIsMethodsConsistency(t *testing.T) {
	t.Run("Ok result methods are consistent", func(t *testing.T) {
		res := result.NewOk("hello")

		if !res.IsOk() {
			t.Error("IsOk() should return true for Ok result")
		}

		if res.IsErr() {
			t.Error("IsErr() should return false for Ok result")
		}

		// IsOk and IsErr should be opposites
		if res.IsOk() == res.IsErr() {
			t.Error("IsOk() and IsErr() should return opposite values")
		}
	})

	t.Run("Err result methods are consistent", func(t *testing.T) {
		err := errors.New("test error")
		res := result.NewErr[string](err)

		if res.IsOk() {
			t.Error("IsOk() should return false for Err result")
		}

		if !res.IsErr() {
			t.Error("IsErr() should return true for Err result")
		}

		// IsOk and IsErr should be opposites
		if res.IsOk() == res.IsErr() {
			t.Error("IsOk() and IsErr() should return opposite values")
		}
	})
}

func TestChainedOperations(t *testing.T) {
	t.Run("multiple operations on Ok result", func(t *testing.T) {
		res := result.NewOk("hello")

		// Should be able to call multiple methods
		if !res.IsOk() || res.IsErr() {
			t.Error("Result state should be consistent")
		}

		val1 := res.Unwrap()
		val2 := res.UnwrapOrDefault("default")
		val3 := res.UnwrapOrZero()

		if val1 != "hello" || val2 != "hello" || val3 != "hello" {
			t.Error("All unwrap methods should return the same value for Ok result")
		}
	})

	t.Run("multiple operations on Err result", func(t *testing.T) {
		err := errors.New("test error")
		res := result.NewErr[string](err)

		// Should be able to call multiple methods
		if res.IsOk() || !res.IsErr() {
			t.Error("Result state should be consistent")
		}

		if res.Error() != err {
			t.Error("Error should be consistent across calls")
		}

		val1 := res.UnwrapOrDefault("default")
		val2 := res.UnwrapOrZero()

		if val1 != "default" || val2 != "" {
			t.Error("Unwrap methods should return expected fallback values")
		}
	})
}
