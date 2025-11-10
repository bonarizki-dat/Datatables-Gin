package datatables

import (
	"reflect"
	"testing"
)

type TestUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"` // Should be excluded
}

type TestProduct struct {
	ProductID int     `json:"product_id"`
	Title     string  // No JSON tag, should use field name
	Price     float64 `json:"price,omitempty"`
}

func TestStructToMapSlice(t *testing.T) {
	t.Run("Convert user slice", func(t *testing.T) {
		users := []TestUser{
			{ID: 1, Name: "John Doe", Email: "john@example.com", Password: "secret"},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Password: "secret2"},
		}

		result := structToMapSlice(&users)

		if len(result) != 2 {
			t.Errorf("Expected 2 results, got %d", len(result))
		}

		// Check first user
		if result[0]["id"] != 1 {
			t.Errorf("Expected id=1, got %v", result[0]["id"])
		}
		if result[0]["name"] != "John Doe" {
			t.Errorf("Expected name='John Doe', got %v", result[0]["name"])
		}
		// Password should be excluded due to json:"-"
		if _, exists := result[0]["Password"]; exists {
			t.Error("Password field should be excluded")
		}
	})

	t.Run("Convert product slice with mixed tags", func(t *testing.T) {
		products := []TestProduct{
			{ProductID: 1, Title: "Laptop", Price: 999.99},
		}

		result := structToMapSlice(&products)

		if len(result) != 1 {
			t.Errorf("Expected 1 result, got %d", len(result))
		}

		// Check JSON tag is respected
		if result[0]["product_id"] != 1 {
			t.Errorf("Expected product_id=1, got %v", result[0]["product_id"])
		}

		// Check field without tag uses field name
		if result[0]["Title"] != "Laptop" {
			t.Errorf("Expected Title='Laptop', got %v", result[0]["Title"])
		}

		// Check omitempty is ignored (field still included)
		if result[0]["price"] != 999.99 {
			t.Errorf("Expected price=999.99, got %v", result[0]["price"])
		}
	})

	t.Run("Empty slice", func(t *testing.T) {
		users := []TestUser{}
		result := structToMapSlice(&users)

		if len(result) != 0 {
			t.Errorf("Expected 0 results, got %d", len(result))
		}
	})

	t.Run("Invalid input - not a slice", func(t *testing.T) {
		notASlice := "invalid"
		result := structToMapSlice(&notASlice)

		if result != nil {
			t.Error("Expected nil for invalid input")
		}
	})
}

func TestStructToMap(t *testing.T) {
	t.Run("Convert single struct", func(t *testing.T) {
		user := TestUser{
			ID:       1,
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "secret",
		}

		v := reflect.ValueOf(user)
		result := structToMap(v)

		if result["id"] != 1 {
			t.Errorf("Expected id=1, got %v", result["id"])
		}
		if result["name"] != "John Doe" {
			t.Errorf("Expected name='John Doe', got %v", result["name"])
		}
		if result["email"] != "john@example.com" {
			t.Errorf("Expected email='john@example.com', got %v", result["email"])
		}

		// Password should not be in map due to json:"-"
		if _, exists := result["Password"]; exists {
			t.Error("Password should be excluded")
		}
	})
}

func TestGetFieldName(t *testing.T) {
	tests := []struct {
		name     string
		jsonTag  string
		expected string
	}{
		{"Simple tag", "user_id", "user_id"},
		{"Tag with omitempty", "email,omitempty", "email"},
		{"Exclude field", "-", ""},
		{"No tag", "", ""},
		{"Multiple options", "name,omitempty,string", "name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := reflect.StructField{
				Name: "TestField",
				Tag:  reflect.StructTag(`json:"` + tt.jsonTag + `"`),
			}

			result := getFieldName(field)

			// Special case: empty tag should return field name
			expected := tt.expected
			if tt.jsonTag == "" {
				expected = "TestField"
			}

			if result != expected {
				t.Errorf("getFieldName() = %q, want %q", result, expected)
			}
		})
	}
}
