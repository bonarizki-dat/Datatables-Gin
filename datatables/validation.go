package datatables

import "regexp"

// columnNamePattern defines the allowed pattern for column names
// Allows: alphanumeric characters, underscores, and dots (for table.column notation)
var columnNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)

// isValidColumnName checks if a column name is safe to use in SQL queries.
// This helps prevent SQL injection attacks by validating column names
// before they are used in dynamic queries.
//
// Valid column names must contain only:
//   - Alphanumeric characters (a-z, A-Z, 0-9)
//   - Underscores (_)
//   - Dots (.) for table.column notation
//
// Returns true if the column name is valid, false otherwise.
func isValidColumnName(name string) bool {
	if name == "" {
		return false
	}
	return columnNamePattern.MatchString(name)
}

// validateSearchableColumns validates all searchable column names
// to ensure they are safe for SQL queries.
//
// Returns an error if any column name is invalid.
func validateSearchableColumns(columns []string) error {
	for _, col := range columns {
		if !isValidColumnName(col) {
			return &ValidationError{
				Field:   col,
				Message: "searchable column name contains invalid characters",
			}
		}
	}
	return nil
}

// validateOrderableColumns validates all orderable column mappings
// to ensure both keys and values are safe for SQL queries.
//
// Returns an error if any column name is invalid.
func validateOrderableColumns(columns map[string]string) error {
	for key, val := range columns {
		if !isValidColumnName(key) {
			return &ValidationError{
				Field:   key,
				Message: "orderable column key contains invalid characters",
			}
		}
		if !isValidColumnName(val) {
			return &ValidationError{
				Field:   val,
				Message: "orderable column value contains invalid characters",
			}
		}
	}
	return nil
}
