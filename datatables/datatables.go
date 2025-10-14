package datatables

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/bonarizki-dat/Datatables-Gin/datatables/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ========================
// Options → Customization for DataTables responses
// ========================

// Options provides customization similar to Yajra DataTables.
// It allows adding, editing, and removing columns dynamically,
// as well as controlling the row index column.
type Options struct {
	IndexColumn   string                                                     // The name of the index column to be added (e.g., "DT_RowIndex")
	ResetIndex    bool                                                       // Whether to reset index numbering to start from 1
	AddColumns    map[string]func(row map[string]interface{}) interface{}    // Custom columns to add, computed from existing row data
	EditColumns   map[string]func(value interface{}, row map[string]interface{}) interface{} // Custom transformations for existing columns
	RemoveColumns []string                                                   // Columns to be removed from the final output
}

// ========================
// Builder for easier options configuration
// ========================

// NewOptions returns a new Options instance with sensible defaults.
func NewOptions() Options {
	return Options{
		IndexColumn:   "DT_RowIndex",
		ResetIndex:    false,
		AddColumns:    make(map[string]func(row map[string]interface{}) interface{}),
		EditColumns:   make(map[string]func(value interface{}, row map[string]interface{}) interface{}),
		RemoveColumns: []string{},
	}
}

// WithIndex configures the name and behavior of the index column.
func (o Options) WithIndex(col string, reset bool) Options {
	o.IndexColumn = col
	o.ResetIndex = reset
	return o
}

// Add registers a new column to be added dynamically using a callback function.
func (o Options) Add(col string, fn func(row map[string]interface{}) interface{}) Options {
	o.AddColumns[col] = fn
	return o
}

// Edit registers a callback function to modify an existing column’s value.
func (o Options) Edit(col string, fn func(value interface{}, row map[string]interface{}) interface{}) Options {
	o.EditColumns[col] = fn
	return o
}

// Remove specifies one or more columns to be removed from the final output.
func (o Options) Remove(cols ...string) Options {
	o.RemoveColumns = append(o.RemoveColumns, cols...)
	return o
}

// ========================
// ParseParams → Extracts DataTables query parameters from request
// ========================

// ParseParams reads and normalizes query parameters used by the DataTables frontend.
// It extracts pagination, sorting, and search information into a standardized dto.Params struct.
func ParseParams(c *gin.Context) dto.Params {
	draw, _ := strconv.ParseInt(c.DefaultQuery("draw", "1"), 10, 64)
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	length, _ := strconv.Atoi(c.DefaultQuery("length", "10"))
	search := c.DefaultQuery("search[value]", "")

	// Try to get order column from different possible sources
	orderColumn := c.DefaultQuery("order[0][column]", "")
	order := ""

	// First try: direct column name from order[0][column]
	if orderColumn != "" {
		order = orderColumn
	} else {
		// Fallback: try the old DataTables format
		order = c.DefaultQuery("columns["+c.DefaultQuery("order[0][column]", "0")+"][data]", "")
	}

	dir := strings.ToLower(c.DefaultQuery("order[0][dir]", "asc"))

	if dir != "asc" && dir != "desc" {
		dir = "asc"
	}
	if length > 500 && length != -1 {
		length = 500
	}

	return dto.Params{
		Draw:   draw,
		Start:  start,
		Length: length,
		Search: search,
		Order:  order,
		Dir:    dir,
	}
}

// ========================
// OfReturn → Core DataTables processor
// ========================

// OfReturn executes the core DataTables server-side logic.
// It supports searching, ordering, pagination, and custom column manipulation.
//
// Parameters:
//   - c: Gin context
//   - query: GORM query builder
//   - dest: Pointer to a slice where results will be stored
//   - searchable: List of columns that are searchable
//   - orderable: Mapping between column names and actual database columns
//   - opts: Optional column customizations (add/edit/remove/index)
//
// Returns a dto.Datatables response structure compatible with DataTables JSON format.
func OfReturn[T any](
	c *gin.Context,
	query *gorm.DB,
	dest *[]T,
	searchable []string,
	orderable map[string]string,
	opts Options,
) (dto.Datatables, error) {
	params := ParseParams(c)

	// --- Count total records ---
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return dto.Datatables{}, err
	}

	// --- Filtering (search) ---
	filteredQuery := query.Session(&gorm.Session{})
	if params.Search != "" && len(searchable) > 0 {
		for i, col := range searchable {
			if i == 0 {
				filteredQuery = filteredQuery.Where("LOWER("+col+") LIKE LOWER(?)", "%"+params.Search+"%")
			} else {
				filteredQuery = filteredQuery.Or("LOWER("+col+") LIKE LOWER(?)", "%"+params.Search+"%")
			}
		}
	}

	// --- Count filtered records ---
	var filtered int64
	if err := filteredQuery.Count(&filtered).Error; err != nil {
		return dto.Datatables{}, err
	}

	// --- Ordering ---
	if params.Order != "" {
		if col, ok := orderable[params.Order]; ok {
			filteredQuery = filteredQuery.Order(col + " " + params.Dir)
		}
	} else {
		// Default ordering if no order specified
		filteredQuery = filteredQuery.Order("created_at DESC")
	}

	// --- Pagination ---
	if params.Length > 0 {
		filteredQuery = filteredQuery.Offset(params.Start).Limit(params.Length)
	}

	// --- Fetch results ---
	if err := filteredQuery.Find(dest).Error; err != nil {
		return dto.Datatables{}, err
	}

	// --- Convert struct slice to []map[string]interface{} ---
	rows := structToMapSliceReflect(dest)

	// --- Apply DataTables options ---
	rows = applyOptions(rows, opts, params.Start)

	return dto.Datatables{
		Draw:            params.Draw,
		RecordsTotal:    total,
		RecordsFiltered: filtered,
		Data:            rows,
	}, nil
}

// ========================
// structToMapSliceReflect → Converts struct slice to []map[string]interface{}
// ========================

// structToMapSliceReflect converts a slice of structs into a slice of map[string]interface{}.
// It uses reflection to read exported fields and respects JSON struct tags.
func structToMapSliceReflect(data interface{}) []map[string]interface{} {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return nil
	}

	var result []map[string]interface{}
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		m := make(map[string]interface{})
		for j := 0; j < item.NumField(); j++ {
			field := item.Type().Field(j)
			jsonTag := field.Tag.Get("json")

			col := field.Name
			if jsonTag != "" && jsonTag != "-" {
				col = strings.Split(jsonTag, ",")[0]
			}
			m[col] = item.Field(j).Interface()
		}
		result = append(result, m)
	}
	return result
}

// ========================
// applyOptions → Applies add/edit/remove/index options
// ========================

// applyOptions processes DataTables customization options such as adding new columns,
// editing existing ones, removing unwanted fields, and setting row indexes.
func applyOptions(data []map[string]interface{}, opts Options, start int) []map[string]interface{} {
	if data == nil {
		return nil
	}

	out := make([]map[string]interface{}, 0, len(data))
	for i, row := range data {
		newRow := make(map[string]interface{})
		for k, v := range row {
			newRow[k] = v
		}

		// --- Index column ---
		if opts.IndexColumn != "" {
			if opts.ResetIndex {
				newRow[opts.IndexColumn] = i + 1
			} else {
				newRow[opts.IndexColumn] = start + i + 1
			}
		}

		// --- Add custom columns ---
		for k, fn := range opts.AddColumns {
			newRow[k] = fn(row)
		}

		// --- Edit existing columns ---
		for k, fn := range opts.EditColumns {
			if val, ok := newRow[k]; ok {
				newRow[k] = fn(val, row)
			}
		}

		// --- Remove unwanted columns ---
		for _, col := range opts.RemoveColumns {
			delete(newRow, col)
		}

		out = append(out, newRow)
	}
	return out
}

// ========================
// JSON Helper → Sends standardized DataTables JSON response
// ========================

// JSON is a convenience helper that sends a standardized DataTables response
// using the gin context and dto.ResponseDatatables function.
func JSON(c *gin.Context, res dto.Datatables) {
	dto.ResponseDatatables(c, http.StatusOK, res, "success")
}
