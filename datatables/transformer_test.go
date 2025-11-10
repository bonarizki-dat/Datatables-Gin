package datatables

import (
	"strings"
	"testing"
)

func TestApplyOptions(t *testing.T) {
	t.Run("Add index column without reset", func(t *testing.T) {
		data := []map[string]interface{}{
			{"id": 1, "name": "John"},
			{"id": 2, "name": "Jane"},
		}

		opts := NewOptions().WithIndex("DT_RowIndex", false)
		result := applyOptions(data, opts, 10) // Start from page offset 10

		if result[0]["DT_RowIndex"] != 11 {
			t.Errorf("Expected DT_RowIndex=11, got %v", result[0]["DT_RowIndex"])
		}
		if result[1]["DT_RowIndex"] != 12 {
			t.Errorf("Expected DT_RowIndex=12, got %v", result[1]["DT_RowIndex"])
		}
	})

	t.Run("Add index column with reset", func(t *testing.T) {
		data := []map[string]interface{}{
			{"id": 1, "name": "John"},
			{"id": 2, "name": "Jane"},
		}

		opts := NewOptions().WithIndex("row_num", true)
		result := applyOptions(data, opts, 50) // Start offset ignored when reset=true

		if result[0]["row_num"] != 1 {
			t.Errorf("Expected row_num=1, got %v", result[0]["row_num"])
		}
		if result[1]["row_num"] != 2 {
			t.Errorf("Expected row_num=2, got %v", result[1]["row_num"])
		}
	})

	t.Run("Add custom column", func(t *testing.T) {
		data := []map[string]interface{}{
			{"first_name": "John", "last_name": "Doe"},
			{"first_name": "Jane", "last_name": "Smith"},
		}

		opts := NewOptions().Add("full_name", func(row map[string]interface{}) interface{} {
			return row["first_name"].(string) + " " + row["last_name"].(string)
		})

		result := applyOptions(data, opts, 0)

		if result[0]["full_name"] != "John Doe" {
			t.Errorf("Expected full_name='John Doe', got %v", result[0]["full_name"])
		}
		if result[1]["full_name"] != "Jane Smith" {
			t.Errorf("Expected full_name='Jane Smith', got %v", result[1]["full_name"])
		}
	})

	t.Run("Edit existing column", func(t *testing.T) {
		data := []map[string]interface{}{
			{"id": 1, "email": "JOHN@EXAMPLE.COM"},
			{"id": 2, "email": "JANE@EXAMPLE.COM"},
		}

		opts := NewOptions().Edit("email", func(value interface{}, row map[string]interface{}) interface{} {
			return strings.ToLower(value.(string))
		})

		result := applyOptions(data, opts, 0)

		if result[0]["email"] != "john@example.com" {
			t.Errorf("Expected email='john@example.com', got %v", result[0]["email"])
		}
		if result[1]["email"] != "jane@example.com" {
			t.Errorf("Expected email='jane@example.com', got %v", result[1]["email"])
		}
	})

	t.Run("Remove columns", func(t *testing.T) {
		data := []map[string]interface{}{
			{"id": 1, "name": "John", "password": "secret", "internal_id": 999},
		}

		opts := NewOptions().Remove("password", "internal_id")
		result := applyOptions(data, opts, 0)

		if _, exists := result[0]["password"]; exists {
			t.Error("password should be removed")
		}
		if _, exists := result[0]["internal_id"]; exists {
			t.Error("internal_id should be removed")
		}
		if _, exists := result[0]["name"]; !exists {
			t.Error("name should still exist")
		}
	})

	t.Run("Complex transformation", func(t *testing.T) {
		data := []map[string]interface{}{
			{"first_name": "John", "last_name": "Doe", "email": "JOHN@EXAMPLE.COM", "password": "secret"},
		}

		opts := NewOptions().
			WithIndex("row_num", true).
			Add("full_name", func(row map[string]interface{}) interface{} {
				return row["first_name"].(string) + " " + row["last_name"].(string)
			}).
			Edit("email", func(value interface{}, row map[string]interface{}) interface{} {
				return strings.ToLower(value.(string))
			}).
			Remove("password")

		result := applyOptions(data, opts, 0)

		// Check index
		if result[0]["row_num"] != 1 {
			t.Errorf("Expected row_num=1, got %v", result[0]["row_num"])
		}

		// Check added column
		if result[0]["full_name"] != "John Doe" {
			t.Errorf("Expected full_name='John Doe', got %v", result[0]["full_name"])
		}

		// Check edited column
		if result[0]["email"] != "john@example.com" {
			t.Errorf("Expected email='john@example.com', got %v", result[0]["email"])
		}

		// Check removed column
		if _, exists := result[0]["password"]; exists {
			t.Error("password should be removed")
		}
	})

	t.Run("Nil data", func(t *testing.T) {
		opts := NewOptions()
		result := applyOptions(nil, opts, 0)

		if result != nil {
			t.Error("Expected nil result for nil input")
		}
	})

	t.Run("Empty data", func(t *testing.T) {
		data := []map[string]interface{}{}
		opts := NewOptions()
		result := applyOptions(data, opts, 0)

		if len(result) != 0 {
			t.Errorf("Expected empty result, got %d items", len(result))
		}
	})
}
