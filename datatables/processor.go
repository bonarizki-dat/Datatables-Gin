package datatables

import (
	"github.com/bonarizki-dat/Datatables-Gin/datatables/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OfReturn executes the core DataTables server-side logic.
// It supports searching, ordering, pagination, and custom column manipulation.
//
// Security features:
//   - Validates all column names to prevent SQL injection
//   - Sanitizes search input with parameterized queries
//   - Enforces maximum page size limits
//
// Parameters:
//   - c: Gin context containing request parameters
//   - query: GORM query builder (can include WHERE clauses, JOINs, etc.)
//   - dest: Pointer to a slice where results will be stored
//   - searchable: List of columns that support global search
//   - orderable: Mapping between frontend column names and database columns
//   - opts: Optional column customizations (add/edit/remove/index/default order)
//
// Returns a dto.Datatables response structure compatible with DataTables JSON format,
// or an error if validation fails or database operations fail.
//
// Example:
//   var users []User
//   result, err := datatables.OfReturn(
//       c,
//       db.Model(&User{}),
//       &users,
//       []string{"name", "email"},
//       map[string]string{"name": "name", "email": "email", "created": "created_at"},
//       datatables.NewOptions().WithDefaultOrder("created_at DESC"),
//   )
func OfReturn[T any](
	c *gin.Context,
	query *gorm.DB,
	dest *[]T,
	searchable []string,
	orderable map[string]string,
	opts Options,
) (dto.Datatables, error) {
	// Validate column names to prevent SQL injection
	if err := validateSearchableColumns(searchable); err != nil {
		return dto.Datatables{}, err
	}
	if err := validateOrderableColumns(orderable); err != nil {
		return dto.Datatables{}, err
	}

	// Parse DataTables request parameters
	params := ParseParams(c)

	// Count total records (before filtering)
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return dto.Datatables{}, err
	}

	// Apply filtering (global search)
	filteredQuery := query.Session(&gorm.Session{})
	if params.Search != "" && len(searchable) > 0 {
		filteredQuery = applySearch(filteredQuery, searchable, params.Search)
	}

	// Count filtered records (after search, before pagination)
	var filtered int64
	if err := filteredQuery.Count(&filtered).Error; err != nil {
		return dto.Datatables{}, err
	}

	// Apply ordering
	filteredQuery = applyOrdering(filteredQuery, params, orderable, opts.DefaultOrder)

	// Apply pagination
	if params.Length > 0 {
		filteredQuery = filteredQuery.Offset(params.Start).Limit(params.Length)
	}

	// Fetch results from database
	if err := filteredQuery.Find(dest).Error; err != nil {
		return dto.Datatables{}, err
	}

	// Convert struct slice to []map[string]interface{}
	rows := structToMapSlice(dest)

	// Apply DataTables options (add/edit/remove columns, indexes)
	rows = applyOptions(rows, opts, params.Start)

	return dto.Datatables{
		Draw:            params.Draw,
		RecordsTotal:    total,
		RecordsFiltered: filtered,
		Data:            rows,
	}, nil
}

// applySearch adds global search conditions to the query.
// Uses OR conditions across all searchable columns with case-insensitive matching.
func applySearch(query *gorm.DB, searchable []string, searchValue string) *gorm.DB {
	for i, col := range searchable {
		searchPattern := "%" + searchValue + "%"
		if i == 0 {
			query = query.Where("LOWER("+col+") LIKE LOWER(?)", searchPattern)
		} else {
			query = query.Or("LOWER("+col+") LIKE LOWER(?)", searchPattern)
		}
	}
	return query
}

// applyOrdering adds ORDER BY clause to the query.
// Uses the orderable map to translate frontend column names to database columns.
// Falls back to defaultOrder if no order is specified.
func applyOrdering(query *gorm.DB, params dto.Params, orderable map[string]string, defaultOrder string) *gorm.DB {
	if params.Order != "" {
		// Check if the requested column is in the orderable map
		if col, ok := orderable[params.Order]; ok {
			return query.Order(col + " " + params.Dir)
		}
	}

	// Apply default ordering if specified and no valid order was provided
	if defaultOrder != "" {
		return query.Order(defaultOrder)
	}

	return query
}
