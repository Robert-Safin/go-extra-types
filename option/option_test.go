package option_test

import (
	"testing"

	"github.com/Robert-Safin/go-extra-types/option"
)

func TestSomeOption(t *testing.T) {
	t.Run("creates some option with value", func(t *testing.T) {
		opt := option.SomeOption("hello")

		if !opt.IsSome() {
			t.Error("Expected option to be Some")
		}

		if opt.IsNone() {
			t.Error("Expected option to not be None")
		}

		if opt.Unwrap() != "hello" {
			t.Errorf("Expected value 'hello', got %s", opt.Unwrap())
		}
	})

	t.Run("works with different types", func(t *testing.T) {
		intOpt := option.SomeOption(42)
		if intOpt.Unwrap() != 42 {
			t.Errorf("Expected value 42, got %d", intOpt.Unwrap())
		}

		boolOpt := option.SomeOption(true)
		if boolOpt.Unwrap() != true {
			t.Errorf("Expected value true, got %t", boolOpt.Unwrap())
		}
	})
}

func TestNoneOption(t *testing.T) {
	t.Run("creates none option", func(t *testing.T) {
		opt := option.NoneOption[string]()

		if opt.IsSome() {
			t.Error("Expected option to not be Some")
		}

		if !opt.IsNone() {
			t.Error("Expected option to be None")
		}
	})

	t.Run("works with different types", func(t *testing.T) {
		intOpt := option.NoneOption[int]()
		if !intOpt.IsNone() {
			t.Error("Expected int option to be None")
		}

		stringOpt := option.NoneOption[string]()
		if !stringOpt.IsNone() {
			t.Error("Expected string option to be None")
		}
	})
}

func TestNewInfer(t *testing.T) {
	t.Run("infers Some from non-zero value", func(t *testing.T) {
		opt := option.NewInfer("hello")

		if !opt.IsSome() {
			t.Error("Expected option to be Some for non-zero value")
		}

		if opt.Unwrap() != "hello" {
			t.Errorf("Expected value 'hello', got %s", opt.Unwrap())
		}
	})

	t.Run("infers None from zero value", func(t *testing.T) {
		opt := option.NewInfer("")

		if !opt.IsNone() {
			t.Error("Expected option to be None for zero value")
		}
	})

	t.Run("infers None from zero int", func(t *testing.T) {
		opt := option.NewInfer(0)

		if !opt.IsNone() {
			t.Error("Expected option to be None for zero int")
		}
	})

	t.Run("explicit boolean override - true", func(t *testing.T) {
		// Even with zero value, should be Some if explicitly set to true
		opt := option.NewInfer("", true)

		if !opt.IsSome() {
			t.Error("Expected option to be Some when explicitly set to true")
		}

		if opt.Unwrap() != "" {
			t.Errorf("Expected empty string, got %s", opt.Unwrap())
		}
	})

	t.Run("explicit boolean override - false", func(t *testing.T) {
		// Even with non-zero value, should be None if explicitly set to false
		opt := option.NewInfer("hello", false)

		if !opt.IsNone() {
			t.Error("Expected option to be None when explicitly set to false")
		}
	})

	t.Run("multiple boolean values uses first", func(t *testing.T) {
		opt := option.NewInfer("hello", true, false, true)

		if !opt.IsSome() {
			t.Error("Expected option to be Some using first boolean value")
		}
	})
}

func TestDestructure(t *testing.T) {
	t.Run("destructures Some option", func(t *testing.T) {
		opt := option.SomeOption("hello")
		val, ok := opt.Destructure()

		if !ok {
			t.Error("Expected ok to be true for Some option")
		}

		if val != "hello" {
			t.Errorf("Expected value 'hello', got %s", val)
		}
	})

	t.Run("destructures None option", func(t *testing.T) {
		opt := option.NoneOption[string]()
		val, ok := opt.Destructure()

		if ok {
			t.Error("Expected ok to be false for None option")
		}

		if val != "" {
			t.Errorf("Expected zero value, got %s", val)
		}
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("unwraps Some option", func(t *testing.T) {
		opt := option.SomeOption("hello")
		val := opt.Unwrap()

		if val != "hello" {
			t.Errorf("Expected value 'hello', got %s", val)
		}
	})

	t.Run("panics on None option", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when unwrapping None option")
			}
		}()

		opt := option.NoneOption[string]()
		opt.Unwrap()
	})
}

func TestUnwrapOrDefault(t *testing.T) {
	t.Run("returns value for Some option", func(t *testing.T) {
		opt := option.SomeOption("hello")
		val := opt.UnwrapOrDefault("default")

		if val != "hello" {
			t.Errorf("Expected value 'hello', got %s", val)
		}
	})

	t.Run("returns default for None option", func(t *testing.T) {
		opt := option.NoneOption[string]()
		val := opt.UnwrapOrDefault("default")

		if val != "default" {
			t.Errorf("Expected default value 'default', got %s", val)
		}
	})
}

func TestUnwrapOrZero(t *testing.T) {
	t.Run("returns value for Some option", func(t *testing.T) {
		opt := option.SomeOption("hello")
		val := opt.UnwrapOrZero()

		if val != "hello" {
			t.Errorf("Expected value 'hello', got %s", val)
		}
	})

	t.Run("returns zero value for None option", func(t *testing.T) {
		opt := option.NoneOption[string]()
		val := opt.UnwrapOrZero()

		if val != "" {
			t.Errorf("Expected zero value, got %s", val)
		}
	})

	t.Run("returns zero for different types", func(t *testing.T) {
		intOpt := option.NoneOption[int]()
		if intOpt.UnwrapOrZero() != 0 {
			t.Errorf("Expected zero int, got %d", intOpt.UnwrapOrZero())
		}

		boolOpt := option.NoneOption[bool]()
		if boolOpt.UnwrapOrZero() != false {
			t.Errorf("Expected false, got %t", boolOpt.UnwrapOrZero())
		}
	})
}

func TestUnwrapOrFunc(t *testing.T) {
	t.Run("returns value for Some option", func(t *testing.T) {
		opt := option.SomeOption("hello")
		val := opt.UnwrapOrFunc(func(o option.Option[string]) string {
			return "from function"
		})

		if val != "hello" {
			t.Errorf("Expected value 'hello', got %s", val)
		}
	})

	t.Run("calls function for None option", func(t *testing.T) {
		opt := option.NoneOption[string]()
		val := opt.UnwrapOrFunc(func(o option.Option[string]) string {
			if !o.IsNone() {
				t.Error("Expected option passed to function to be None")
			}
			return "from function"
		})

		if val != "from function" {
			t.Errorf("Expected value 'from function', got %s", val)
		}
	})

	t.Run("function receives correct option", func(t *testing.T) {
		opt := option.NoneOption[int]()
		called := false

		opt.UnwrapOrFunc(func(o option.Option[int]) int {
			called = true
			if !o.IsNone() {
				t.Error("Expected received option to be None")
			}
			return 42
		})

		if !called {
			t.Error("Expected function to be called")
		}
	})
}

func TestOptionWithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("struct option", func(t *testing.T) {
		person := Person{Name: "Alice", Age: 30}
		opt := option.SomeOption(person)

		if !opt.IsSome() {
			t.Error("Expected option to be Some")
		}

		result := opt.Unwrap()
		if result.Name != "Alice" || result.Age != 30 {
			t.Errorf("Expected person Alice(30), got %+v", result)
		}
	})

	t.Run("slice option", func(t *testing.T) {
		slice := []int{1, 2, 3}
		opt := option.SomeOption(slice)

		if !opt.IsSome() {
			t.Error("Expected option to be Some")
		}

		result := opt.Unwrap()
		if len(result) != 3 || result[0] != 1 {
			t.Errorf("Expected slice [1,2,3], got %v", result)
		}
	})

	t.Run("map option", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2}
		opt := option.SomeOption(m)

		if !opt.IsSome() {
			t.Error("Expected option to be Some")
		}

		result := opt.Unwrap()
		if result["a"] != 1 || result["b"] != 2 {
			t.Errorf("Expected map with a:1, b:2, got %v", result)
		}
	})
}

func TestIsMethodsConsistency(t *testing.T) {
	t.Run("Some option methods are consistent", func(t *testing.T) {
		opt := option.SomeOption("hello")

		if !opt.IsSome() {
			t.Error("IsSome() should return true for Some option")
		}

		if opt.IsNone() {
			t.Error("IsNone() should return false for Some option")
		}

		// IsSome and IsNone should be opposites
		if opt.IsSome() == opt.IsNone() {
			t.Error("IsSome() and IsNone() should return opposite values")
		}
	})

	t.Run("None option methods are consistent", func(t *testing.T) {
		opt := option.NoneOption[string]()

		if opt.IsSome() {
			t.Error("IsSome() should return false for None option")
		}

		if !opt.IsNone() {
			t.Error("IsNone() should return true for None option")
		}

		// IsSome and IsNone should be opposites
		if opt.IsSome() == opt.IsNone() {
			t.Error("IsSome() and IsNone() should return opposite values")
		}
	})
}
