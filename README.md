# 🧮 Datatables-Gin

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](datatables/)

A **reusable, secure, and powerful DataTables-style server-side processor** for **Gin + GORM** in Go.
Inspired by [Yajra DataTables](https://github.com/yajra/laravel-datatables) (Laravel), designed for clean, composable, and production-ready data handling in APIs.

---

## ✨ Features

- 🔁 **Automatic Processing**: Handles pagination, ordering, and search out of the box
- 🛡️ **Security First**: Built-in SQL injection prevention and column validation
- 🧱 **Simple API**: Similar to Yajra DataTables, easy to learn and use
- ⚙️ **Highly Customizable**: Add, edit, or remove columns dynamically
- 📊 **Gin + GORM**: Works seamlessly with Gin framework and GORM ORM
- 💡 **Type-Safe**: Generic implementation (`OfReturn[T any]`) for compile-time safety
- 🧪 **Well-Tested**: Comprehensive unit tests with >90% coverage
- 📦 **Production-Ready**: Configurable defaults, error handling, and validation

---

## 📦 Installation

```bash
go get github.com/bonarizki-dat/Datatables-Gin
```

**Requirements:**
- Go 1.23 or higher
- Gin Web Framework
- GORM v2

---

## 🚀 Quick Start

### Basic Example

```go
package main

import (
    "github.com/bonarizki-dat/Datatables-Gin/datatables"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type User struct {
    ID        uint   `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

func GetUsers(c *gin.Context, db *gorm.DB) {
    var users []User

    // Define searchable columns
    searchable := []string{"name", "email"}

    // Map frontend column names to database columns
    orderable := map[string]string{
        "name":    "name",
        "email":   "email",
        "created": "created_at",
    }

    // Configure options
    opts := datatables.NewOptions().
        WithDefaultOrder("created_at DESC")

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
        datatables.JSONError(c, 500, err.Error())
        return
    }

    datatables.JSON(c, result)
}
```

### Frontend Integration

```html
<table id="users-table">
    <thead>
        <tr>
            <th data-data="DT_RowIndex">No</th>
            <th data-data="name">Name</th>
            <th data-data="email">Email</th>
            <th data-data="created">Created</th>
        </tr>
    </thead>
</table>

<script>
$('#users-table').DataTable({
    processing: true,
    serverSide: true,
    ajax: '/api/users',
    columns: [
        { data: 'DT_RowIndex', orderable: false, searchable: false },
        { data: 'name' },
        { data: 'email' },
        { data: 'created', name: 'created' }
    ]
});
</script>
```

---

## 🎯 Advanced Usage

### Custom Columns

Add computed columns dynamically:

```go
opts := datatables.NewOptions().
    Add("full_name", func(row map[string]interface{}) interface{} {
        return fmt.Sprintf("%s %s",
            row["first_name"],
            row["last_name"])
    }).
    Add("status_badge", func(row map[string]interface{}) interface{} {
        status := row["status"].(string)
        if status == "active" {
            return `<span class="badge badge-success">Active</span>`
        }
        return `<span class="badge badge-danger">Inactive</span>`
    })
```

### Transform Existing Columns

Edit column values before sending to frontend:

```go
opts := datatables.NewOptions().
    Edit("email", func(value interface{}, row map[string]interface{}) interface{} {
        // Lowercase all emails
        return strings.ToLower(value.(string))
    }).
    Edit("created_at", func(value interface{}, row map[string]interface{}) interface{} {
        // Format timestamps
        t := value.(time.Time)
        return t.Format("2006-01-02 15:04:05")
    })
```

### Remove Sensitive Data

Hide columns from the response:

```go
opts := datatables.NewOptions().
    Remove("password", "internal_id", "deleted_at")
```

### Custom Index Column

Configure row numbering:

```go
// Continuous numbering across pages
opts := datatables.NewOptions().
    WithIndex("row_number", false)

// Reset numbering on each page
opts := datatables.NewOptions().
    WithIndex("row_number", true)
```

### Complete Example

```go
func GetProducts(c *gin.Context, db *gorm.DB) {
    var products []Product

    searchable := []string{"name", "category", "sku"}

    orderable := map[string]string{
        "name":     "name",
        "price":    "price",
        "stock":    "stock",
        "category": "category",
        "created":  "created_at",
    }

    opts := datatables.NewOptions().
        WithIndex("row_num", true).
        WithDefaultOrder("created_at DESC").
        Add("price_formatted", func(row map[string]interface{}) interface{} {
            price := row["price"].(float64)
            return fmt.Sprintf("$%.2f", price)
        }).
        Add("stock_status", func(row map[string]interface{}) interface{} {
            stock := row["stock"].(int)
            if stock > 10 {
                return "In Stock"
            } else if stock > 0 {
                return "Low Stock"
            }
            return "Out of Stock"
        }).
        Edit("category", func(value interface{}, row map[string]interface{}) interface{} {
            return strings.ToUpper(value.(string))
        }).
        Remove("internal_notes", "cost_price")

    result, err := datatables.OfReturn(
        c,
        db.Model(&Product{}),
        &products,
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
```

---

## 🔒 Security Features

### SQL Injection Prevention

All column names are validated against a whitelist pattern before being used in queries:

```go
// ✅ Valid column names
"user_id"
"users.name"
"created_at"

// ❌ Blocked - prevents SQL injection
"id; DROP TABLE users--"
"name' OR '1'='1"
"id--"
```

### Validation Errors

The package returns descriptive errors for security violations:

```go
result, err := datatables.OfReturn(c, query, &users, searchable, orderable, opts)
if err != nil {
    // Handle validation errors
    log.Printf("DataTables error: %v", err)
    datatables.JSONError(c, 400, "Invalid column name")
    return
}
```

### Safe Search Queries

All search values use parameterized queries:

```go
// Automatically converts to safe parameterized query:
// WHERE LOWER(name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?)
// Parameters: ["%search_value%", "%search_value%"]
```

---

## 📚 API Reference

### Core Functions

#### `OfReturn[T any]()`

Main function for processing DataTables requests.

```go
func OfReturn[T any](
    c *gin.Context,
    query *gorm.DB,
    dest *[]T,
    searchable []string,
    orderable map[string]string,
    opts Options,
) (dto.Datatables, error)
```

**Parameters:**
- `c` - Gin context containing request parameters
- `query` - GORM query builder (can include WHERE, JOIN, etc.)
- `dest` - Pointer to slice where results will be stored
- `searchable` - Columns that support global search
- `orderable` - Map of frontend column names to database columns
- `opts` - Configuration options

**Returns:**
- `dto.Datatables` - Response compatible with DataTables JSON format
- `error` - Validation or database errors

#### `JSON()`

Sends a standardized DataTables JSON response.

```go
datatables.JSON(c, result)
```

#### `JSONError()`

Sends a consistent error response.

```go
datatables.JSONError(c, 500, "Database error")
```

### Options Builder

#### `NewOptions()`

Creates a new Options instance with defaults.

```go
opts := datatables.NewOptions()
```

#### `WithIndex(col string, reset bool)`

Configures the index column.

```go
opts.WithIndex("DT_RowIndex", false)
```

#### `WithDefaultOrder(order string)`

Sets default ordering when none is specified.

```go
opts.WithDefaultOrder("created_at DESC")
```

#### `Add(col string, fn func)`

Adds a custom computed column.

```go
opts.Add("full_name", func(row map[string]interface{}) interface{} {
    return row["first_name"].(string) + " " + row["last_name"].(string)
})
```

#### `Edit(col string, fn func)`

Transforms an existing column value.

```go
opts.Edit("email", func(value interface{}, row map[string]interface{}) interface{} {
    return strings.ToLower(value.(string))
})
```

#### `Remove(cols ...string)`

Removes columns from the response.

```go
opts.Remove("password", "internal_id")
```

---

## 🧪 Testing

Run the test suite:

```bash
go test ./datatables/... -v
```

Run with coverage:

```bash
go test ./datatables/... -cover
```

---

## 📂 Project Structure

```
datatables/
├── datatables.go      # Package documentation
├── processor.go       # Core OfReturn logic
├── options.go         # Options builder
├── parser.go          # Request parameter parsing
├── converter.go       # Struct to map conversion
├── transformer.go     # Column transformations
├── validation.go      # Security validation
├── errors.go          # Error types
├── response.go        # JSON response helpers
└── dto/
    ├── request.go     # Request DTOs
    └── response.go    # Response DTOs
```

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Inspired by [Yajra DataTables](https://github.com/yajra/laravel-datatables) for Laravel
- Built with [Gin](https://github.com/gin-gonic/gin) and [GORM](https://gorm.io)

---

## 📞 Support

If you encounter any issues or have questions:

- 📝 [Open an issue](https://github.com/bonarizki-dat/Datatables-Gin/issues)
- 💬 [Discussions](https://github.com/bonarizki-dat/Datatables-Gin/discussions)

---

<p align="center">Made with ❤️ for the Go community</p>
