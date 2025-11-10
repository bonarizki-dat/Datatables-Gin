package datatables

import (
	"net/http"

	"github.com/bonarizki-dat/Datatables-Gin/datatables/dto"
	"github.com/gin-gonic/gin"
)

// JSON is a convenience helper that sends a standardized DataTables response.
// It wraps the DataTables result in a SuccessResponse structure and sends it
// as JSON with HTTP 200 OK status.
//
// Parameters:
//   - c: Gin context
//   - res: DataTables response containing draw, records total/filtered, and data
//
// Example:
//   result, err := datatables.OfReturn(c, query, &users, searchable, orderable, opts)
//   if err != nil {
//       c.JSON(500, gin.H{"error": err.Error()})
//       return
//   }
//   datatables.JSON(c, result)
func JSON(c *gin.Context, res dto.Datatables) {
	dto.ResponseDatatables(c, http.StatusOK, res, "success")
}

// JSONError is a convenience helper for sending error responses in a consistent format.
//
// Parameters:
//   - c: Gin context
//   - statusCode: HTTP status code (e.g., 400, 500)
//   - message: Error message to display
func JSONError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, dto.SuccessResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Errors:  message,
	})
}
