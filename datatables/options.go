package datatables

// Options provides customization similar to Yajra DataTables.
// It allows adding, editing, and removing columns dynamically,
// as well as controlling the row index column and default ordering.
type Options struct {
	// IndexColumn is the name of the index column to be added (e.g., "DT_RowIndex")
	IndexColumn string

	// ResetIndex determines whether to reset index numbering to start from 1
	ResetIndex bool

	// DefaultOrder specifies the default ordering when no order is provided
	// Format: "column_name direction" (e.g., "created_at DESC")
	// Leave empty to skip default ordering
	DefaultOrder string

	// AddColumns contains custom columns to add, computed from existing row data
	AddColumns map[string]func(row map[string]interface{}) interface{}

	// EditColumns contains custom transformations for existing columns
	EditColumns map[string]func(value interface{}, row map[string]interface{}) interface{}

	// RemoveColumns is a list of columns to be removed from the final output
	RemoveColumns []string
}

// NewOptions returns a new Options instance with sensible defaults.
// By default, it creates an index column named "DT_RowIndex" without resetting indices.
func NewOptions() Options {
	return Options{
		IndexColumn:   "DT_RowIndex",
		ResetIndex:    false,
		DefaultOrder:  "",
		AddColumns:    make(map[string]func(row map[string]interface{}) interface{}),
		EditColumns:   make(map[string]func(value interface{}, row map[string]interface{}) interface{}),
		RemoveColumns: []string{},
	}
}

// WithIndex configures the name and behavior of the index column.
//
// Parameters:
//   - col: The column name for the index (e.g., "DT_RowIndex", "row_number")
//   - reset: If true, index starts from 1 on each page; if false, continues from previous pages
//
// Example:
//   opts.WithIndex("row_num", true)
func (o Options) WithIndex(col string, reset bool) Options {
	o.IndexColumn = col
	o.ResetIndex = reset
	return o
}

// WithDefaultOrder sets the default ordering clause to use when no order is specified.
// This prevents errors when tables don't have a "created_at" column.
//
// Parameters:
//   - order: The default order clause (e.g., "id DESC", "name ASC")
//
// Example:
//   opts.WithDefaultOrder("id DESC")
func (o Options) WithDefaultOrder(order string) Options {
	o.DefaultOrder = order
	return o
}

// Add registers a new column to be added dynamically using a callback function.
// The callback receives the entire row data and should return the value for the new column.
//
// Parameters:
//   - col: The name of the new column
//   - fn: A function that computes the column value from row data
//
// Example:
//   opts.Add("full_name", func(row map[string]interface{}) interface{} {
//       return row["first_name"].(string) + " " + row["last_name"].(string)
//   })
func (o Options) Add(col string, fn func(row map[string]interface{}) interface{}) Options {
	o.AddColumns[col] = fn
	return o
}

// Edit registers a callback function to modify an existing column's value.
// The callback receives both the current value and the entire row data.
//
// Parameters:
//   - col: The name of the column to edit
//   - fn: A function that transforms the column value
//
// Example:
//   opts.Edit("email", func(value interface{}, row map[string]interface{}) interface{} {
//       return strings.ToLower(value.(string))
//   })
func (o Options) Edit(col string, fn func(value interface{}, row map[string]interface{}) interface{}) Options {
	o.EditColumns[col] = fn
	return o
}

// Remove specifies one or more columns to be removed from the final output.
// This is useful for hiding sensitive data or reducing payload size.
//
// Parameters:
//   - cols: One or more column names to remove
//
// Example:
//   opts.Remove("password", "internal_id", "deleted_at")
func (o Options) Remove(cols ...string) Options {
	o.RemoveColumns = append(o.RemoveColumns, cols...)
	return o
}
