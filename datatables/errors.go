package datatables

import "errors"

// Common errors returned by the datatables package
var (
	// ErrInvalidColumnName is returned when a column name contains invalid characters
	ErrInvalidColumnName = errors.New("invalid column name: must contain only alphanumeric characters, underscores, and dots")

	// ErrColumnNotFound is returned when a requested orderable column is not found
	ErrColumnNotFound = errors.New("column not found in orderable map")

	// ErrInvalidData is returned when the provided data is not a valid slice
	ErrInvalidData = errors.New("invalid data: expected a slice")

	// ErrDefaultOrderColumn is returned when default ordering references a non-existent column
	ErrDefaultOrderColumn = errors.New("default order column does not exist in the database")
)

// ValidationError represents a validation error with additional context
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return "validation error on field '" + e.Field + "': " + e.Message
}
