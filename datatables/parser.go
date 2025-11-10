package datatables

import (
	"strconv"
	"strings"

	"github.com/bonarizki-dat/Datatables-Gin/datatables/dto"
	"github.com/gin-gonic/gin"
)

// ParseParams reads and normalizes query parameters used by the DataTables frontend.
// It extracts pagination, sorting, and search information into a standardized dto.Params struct.
//
// Supported DataTables parameters:
//   - draw: Draw counter for synchronization
//   - start: Record offset for pagination
//   - length: Number of records per page (max 500)
//   - search[value]: Global search value
//   - order[0][column]: Column to order by
//   - order[0][dir]: Order direction (asc/desc)
//
// Returns a dto.Params struct with parsed values and sensible defaults.
func ParseParams(c *gin.Context) dto.Params {
	// Parse draw counter (used by DataTables for synchronization)
	draw, _ := strconv.ParseInt(c.DefaultQuery("draw", "1"), 10, 64)

	// Parse pagination parameters
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	length, _ := strconv.Atoi(c.DefaultQuery("length", "10"))

	// Parse search value
	search := c.DefaultQuery("search[value]", "")

	// Try to get order column from different possible sources
	orderColumn := c.DefaultQuery("order[0][column]", "")
	order := ""

	// First try: direct column name from order[0][column]
	if orderColumn != "" {
		order = orderColumn
	} else {
		// Fallback: try the old DataTables format (column index)
		columnIndex := c.DefaultQuery("order[0][column]", "0")
		order = c.DefaultQuery("columns["+columnIndex+"][data]", "")
	}

	// Parse and validate order direction
	dir := strings.ToLower(c.DefaultQuery("order[0][dir]", "asc"))
	if dir != "asc" && dir != "desc" {
		dir = "asc" // Default to ascending if invalid
	}

	// Enforce maximum page size to prevent abuse
	// -1 means "all records" and is allowed
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
