package enum_test

import (
	"testing"

	"github.com/Robert-Safin/go-extra-types/enum"
)

func TestNewEnum(t *testing.T) {
	t.Run("valid enum creation", func(t *testing.T) {
		variants := map[string]string{
			"Red":   "#FF0000",
			"Green": "#00FF00",
			"Blue":  "#0000FF",
		}

		e := enum.NewEnum("Colors", variants)

		if e.String() == "" {
			t.Error("Expected non-empty string representation")
		}

		names := e.VariantNames()
		if len(names) != 3 {
			t.Errorf("Expected 3 variant names, got %d", len(names))
		}

		// Check that all expected names are present
		nameMap := make(map[string]bool)
		for _, name := range names {
			nameMap[name] = true
		}

		for expectedName := range variants {
			if !nameMap[expectedName] {
				t.Errorf("Expected variant name %s not found", expectedName)
			}
		}
	})

	t.Run("empty enum name panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for empty enum name")
			}
		}()

		variants := map[string]string{"Red": "#FF0000"}
		enum.NewEnum("", variants)
	})

	t.Run("empty variants map panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for empty variants map")
			}
		}()

		variants := map[string]string{}
		enum.NewEnum("Colors", variants)
	})

	t.Run("enum is immutable after creation", func(t *testing.T) {
		variants := map[string]string{
			"Red": "#FF0000",
		}

		e := enum.NewEnum("Colors", variants)

		// Modify original map
		variants["Blue"] = "#0000FF"
		delete(variants, "Red")

		// Enum should still have only Red variant
		names := e.VariantNames()
		if len(names) != 1 {
			t.Errorf("Expected 1 variant, got %d", len(names))
		}

		if names[0] != "Red" {
			t.Errorf("Expected Red variant, got %s", names[0])
		}
	})
}

func TestNewInstance(t *testing.T) {
	variants := map[string]string{
		"Red":   "#FF0000",
		"Green": "#00FF00",
		"Blue":  "#0000FF",
	}
	e := enum.NewEnum("Colors", variants)

	t.Run("valid variant creation", func(t *testing.T) {
		v := e.NewInstance("Red")

		if v.Name() != "Red" {
			t.Errorf("Expected variant name Red, got %s", v.Name())
		}

		if v.Value() != "#FF0000" {
			t.Errorf("Expected variant value #FF0000, got %s", v.Value())
		}

		if !v.IsInstanceOf(e) {
			t.Error("Expected variant to be instance of its enum")
		}

		if v.String() == "" {
			t.Error("Expected non-empty string representation")
		}
	})

	t.Run("empty variant name panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for empty variant name")
			}
		}()

		e.NewInstance("")
	})

	t.Run("invalid variant name panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for invalid variant name")
			}
		}()

		e.NewInstance("Yellow")
	})
}

func TestVariantIsInstanceOf(t *testing.T) {
	colors := map[string]string{
		"Red": "#FF0000",
	}
	fruits := map[string]string{
		"Red": "Apple",
	}

	colorsEnum := enum.NewEnum("Colors", colors)
	fruitsEnum := enum.NewEnum("Fruits", fruits)

	redColor := colorsEnum.NewInstance("Red")
	redFruit := fruitsEnum.NewInstance("Red")

	t.Run("variant belongs to correct enum", func(t *testing.T) {
		if !redColor.IsInstanceOf(colorsEnum) {
			t.Error("Expected redColor to be instance of colorsEnum")
		}

		if !redFruit.IsInstanceOf(fruitsEnum) {
			t.Error("Expected redFruit to be instance of fruitsEnum")
		}
	})

	t.Run("variant does not belong to different enum", func(t *testing.T) {
		if redColor.IsInstanceOf(fruitsEnum) {
			t.Error("Expected redColor to NOT be instance of fruitsEnum")
		}

		if redFruit.IsInstanceOf(colorsEnum) {
			t.Error("Expected redFruit to NOT be instance of colorsEnum")
		}
	})
}

func TestEnumWithDifferentTypes(t *testing.T) {
	t.Run("integer enum", func(t *testing.T) {
		variants := map[string]int{
			"One":   1,
			"Two":   2,
			"Three": 3,
		}

		e := enum.NewEnum("Numbers", variants)
		v := e.NewInstance("Two")

		if v.Value() != 2 {
			t.Errorf("Expected value 2, got %d", v.Value())
		}
	})

	t.Run("struct enum", func(t *testing.T) {
		type Config struct {
			Host string
			Port int
		}

		variants := map[string]Config{
			"Dev":  {Host: "localhost", Port: 3000},
			"Prod": {Host: "example.com", Port: 443},
		}

		e := enum.NewEnum("Environments", variants)
		v := e.NewInstance("Dev")

		config := v.Value()
		if config.Host != "localhost" || config.Port != 3000 {
			t.Errorf("Expected dev config, got %+v", config)
		}
	})

	t.Run("function enum", func(t *testing.T) {
		variants := map[string]func() string{
			"Hello": func() string { return "Hello, World!" },
			"Bye":   func() string { return "Goodbye!" },
		}

		e := enum.NewEnum("Greetings", variants)
		v := e.NewInstance("Hello")

		fn := v.Value()
		if fn() != "Hello, World!" {
			t.Errorf("Expected 'Hello, World!', got %s", fn())
		}
	})
}

func TestEnumVariantNames(t *testing.T) {
	variants := map[string]string{
		"C": "Third",
		"A": "First",
		"B": "Second",
	}

	e := enum.NewEnum("Letters", variants)
	names := e.VariantNames()

	if len(names) != 3 {
		t.Errorf("Expected 3 names, got %d", len(names))
	}

	// Verify all names are present (order doesn't matter for maps)
	nameSet := make(map[string]bool)
	for _, name := range names {
		nameSet[name] = true
	}

	expectedNames := []string{"A", "B", "C"}
	for _, expected := range expectedNames {
		if !nameSet[expected] {
			t.Errorf("Expected name %s not found in %v", expected, names)
		}
	}
}

func TestEnumString(t *testing.T) {
	variants := map[string]string{
		"Red": "#FF0000",
	}

	e := enum.NewEnum("Colors", variants)
	str := e.String()

	if str == "" {
		t.Error("Expected non-empty string representation")
	}

	// Should contain enum name and variant count
	if !contains(str, "Colors") {
		t.Error("Expected string to contain enum name")
	}

	if !contains(str, "1") {
		t.Error("Expected string to contain variant count")
	}
}

func TestVariantString(t *testing.T) {
	variants := map[string]string{
		"Red": "#FF0000",
	}

	e := enum.NewEnum("Colors", variants)
	v := e.NewInstance("Red")
	str := v.String()

	if str == "" {
		t.Error("Expected non-empty string representation")
	}

	// Should contain enum name and variant name
	if !contains(str, "Colors") {
		t.Error("Expected string to contain enum name")
	}

	if !contains(str, "Red") {
		t.Error("Expected string to contain variant name")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && someContains(s, substr)))
}

func someContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
