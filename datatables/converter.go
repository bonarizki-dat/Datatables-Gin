package datatables

import (
	"reflect"
	"strings"
)

// structToMapSlice converts a slice of structs into a slice of map[string]interface{}.
// It uses reflection to read exported fields and respects JSON struct tags.
//
// Supported struct tag formats:
//   - `json:"field_name"`: Uses "field_name" as the map key
//   - `json:"field_name,omitempty"`: Uses "field_name" (options are ignored)
//   - `json:"-"`: Field is excluded from output
//   - No tag: Uses the field name as-is
//
// Parameters:
//   - data: Pointer to a slice of structs (e.g., *[]User)
//
// Returns a slice of maps where each map represents one struct instance.
// Returns nil if the input is not a valid slice.
//
// Example:
//   type User struct {
//       ID   int    `json:"id"`
//       Name string `json:"name"`
//   }
//   users := []User{{ID: 1, Name: "John"}}
//   result := structToMapSlice(&users)
//   // result: [{"id": 1, "name": "John"}]
func structToMapSlice(data interface{}) []map[string]interface{} {
	v := reflect.ValueOf(data)

	// Dereference pointer if necessary
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Validate that we have a slice
	if v.Kind() != reflect.Slice {
		return nil
	}

	result := make([]map[string]interface{}, 0, v.Len())

	// Iterate through each element in the slice
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)

		// Dereference pointer if necessary
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		// Convert struct to map
		m := structToMap(item)
		result = append(result, m)
	}

	return result
}

// structToMap converts a single struct value to a map[string]interface{}.
// It processes all exported fields and respects JSON tags.
func structToMap(v reflect.Value) map[string]interface{} {
	m := make(map[string]interface{})

	// Iterate through all fields in the struct
	for j := 0; j < v.NumField(); j++ {
		field := v.Type().Field(j)
		fieldValue := v.Field(j)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Determine the map key from JSON tag or field name
		col := getFieldName(field)

		// Skip fields marked with json:"-"
		if col == "" {
			continue
		}

		// Add field to map
		m[col] = fieldValue.Interface()
	}

	return m
}

// getFieldName extracts the field name from the JSON struct tag.
// Returns an empty string if the field should be excluded (json:"-").
func getFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")

	// If no JSON tag, use the field name
	if jsonTag == "" {
		return field.Name
	}

	// Handle json:"-" (exclude field)
	if jsonTag == "-" {
		return ""
	}

	// Extract field name (before comma, if present)
	// e.g., "field_name,omitempty" -> "field_name"
	parts := strings.Split(jsonTag, ",")
	return parts[0]
}
