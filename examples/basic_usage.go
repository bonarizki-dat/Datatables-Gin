package examples

import (
	"time"

	"github.com/bonarizki-dat/Datatables-Gin/datatables"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// User model example
type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// BasicUsage demonstrates the simplest way to use the package
// This is the SAME as before - no breaking changes!
func BasicUsage(c *gin.Context, db *gorm.DB) {
	var users []User

	// Define searchable columns
	searchable := []string{"name", "email"}

	// Map frontend column names to database columns
	orderable := map[string]string{
		"name":  "name",
		"email": "email",
	}

	// Create options (same as before)
	opts := datatables.NewOptions()

	// Process DataTables request
	result, err := datatables.OfReturn(
		c,
		db.Model(&User{}),
		&users,
		searchable,
		orderable,
		opts,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	datatables.JSON(c, result)
}

// ImprovedUsage shows new optional features
// Old code doesn't need to use these!
func ImprovedUsage(c *gin.Context, db *gorm.DB) {
	var users []User

	searchable := []string{"name", "email"}

	orderable := map[string]string{
		"name":    "name",
		"email":   "email",
		"created": "created_at",
	}

	// NEW: Optional improvements
	opts := datatables.NewOptions().
		WithDefaultOrder("created_at DESC"). // Prevent "column not found" errors
		WithIndex("row_num", true)           // Custom index column

	result, err := datatables.OfReturn(
		c,
		db.Model(&User{}),
		&users,
		searchable,
		orderable,
		opts,
	)

	if err != nil {
		// NEW: Optional error helper
		datatables.JSONError(c, 500, err.Error())
		return
	}

	datatables.JSON(c, result)
}

// AdvancedUsage shows all features including custom columns
func AdvancedUsage(c *gin.Context, db *gorm.DB) {
	var users []User

	searchable := []string{"name", "email"}

	orderable := map[string]string{
		"name":    "name",
		"email":   "email",
		"created": "created_at",
	}

	opts := datatables.NewOptions().
		WithDefaultOrder("created_at DESC").
		WithIndex("DT_RowIndex", false).
		Add("actions", func(row map[string]interface{}) interface{} {
			id := row["id"]
			return `<a href="/users/` + id.(string) + `/edit">Edit</a>`
		}).
		Remove("created_at") // Hide internal field

	result, err := datatables.OfReturn(
		c,
		db.Model(&User{}),
		&users,
		searchable,
		orderable,
		opts,
	)

	if err != nil {
		datatables.JSONError(c, 500, err.Error())
		return
	}

	datatables.JSON(c, result)
}
