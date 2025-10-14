package dto

import (
	"github.com/gin-gonic/gin"
)

// ========================
// DataTables Response Format
// ========================

// Datatables represents the standard response structure used by the
// jQuery DataTables plugin. It includes pagination metadata and data rows.
type Datatables struct {
	Draw            int64       `json:"draw"`            // Draw counter to synchronize client-side and server-side data
	RecordsTotal    int64       `json:"recordsTotal"`    // Total number of records available
	RecordsFiltered int64       `json:"recordsFiltered"` // Number of records after applying filters
	Data            interface{} `json:"data"`            // Actual data rows to be displayed in the DataTable
}

// ========================
// Generic Success Response
// ========================

// SuccessResponse defines a consistent structure for successful API responses.
// It can be used for both standard API calls and DataTables integrations.
type SuccessResponse struct {
	Success bool        `json:"success"` // Indicates whether the operation was successful
	Message string      `json:"message"` // Optional message describing the result
	Data    interface{} `json:"data"`    // Returned data payload
	Errors  interface{} `json:"errors"`  // Optional error details (usually nil for success)
}

// ========================
// Response Helper Function
// ========================

// ResponseDatatables sends a standardized JSON success response to the client.
// It wraps the data in a SuccessResponse structure. If the message is empty,
// it will default to "Success".
//
// Example:
//   ResponseDatatables(c, http.StatusOK, result, "")
func ResponseDatatables(c *gin.Context, code int, data interface{}, message string) {
	if message == "" {
		message = "Success"
	}

	c.JSON(code, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Errors:  nil,
	})
}