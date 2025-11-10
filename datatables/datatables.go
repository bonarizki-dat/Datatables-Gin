// Package datatables provides a reusable DataTables-style server-side processor for Gin + GORM.
//
// Inspired by Yajra DataTables (Laravel), this package offers a clean and powerful API
// for handling server-side DataTables operations including pagination, searching, ordering,
// and dynamic column manipulation.
//
// Key Features:
//   - Automatic pagination, ordering, and search handling
//   - Simple API similar to Yajra DataTables
//   - Fully customizable columns (add, edit, remove)
//   - Type-safe implementation with Go generics
//   - SQL injection prevention with column name validation
//   - Configurable default ordering
//
// Basic Usage:
//
//	var users []User
//	result, err := datatables.OfReturn(
//	    c,
//	    db.Model(&User{}),
//	    &users,
//	    []string{"name", "email"},
//	    map[string]string{
//	        "name": "name",
//	        "email": "email",
//	        "created": "created_at",
//	    },
//	    datatables.NewOptions().WithDefaultOrder("created_at DESC"),
//	)
//	if err != nil {
//	    c.JSON(500, gin.H{"error": err.Error()})
//	    return
//	}
//	datatables.JSON(c, result)
//
// Advanced Usage with Column Customization:
//
//	opts := datatables.NewOptions().
//	    WithIndex("row_num", true).
//	    WithDefaultOrder("id DESC").
//	    Add("full_name", func(row map[string]interface{}) interface{} {
//	        return fmt.Sprintf("%s %s", row["first_name"], row["last_name"])
//	    }).
//	    Edit("email", func(value interface{}, row map[string]interface{}) interface{} {
//	        return strings.ToLower(value.(string))
//	    }).
//	    Remove("password", "internal_id")
//
// Security:
//   - All column names are validated against a whitelist pattern
//   - Search queries use parameterized statements
//   - Maximum page size is enforced (500 records)
//
// For more information, visit: https://github.com/bonarizki-dat/Datatables-Gin
package datatables
