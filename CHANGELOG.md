# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased] - 2025-11-10

### ✅ Added

- **Security**: Column name validation to prevent SQL injection attacks
- **Feature**: `WithDefaultOrder()` method for configurable default ordering
- **Feature**: `JSONError()` helper function for consistent error responses
- **Error Types**: Custom error types (`ValidationError`, `ErrInvalidColumnName`, etc.)
- **Testing**: Comprehensive unit tests with 57.4% coverage
  - `validation_test.go` - Security validation tests
  - `converter_test.go` - Struct conversion tests
  - `options_test.go` - Options builder tests
  - `transformer_test.go` - Transformation tests
- **Documentation**: Extensive README with examples and API reference
- **Documentation**: `BACKWARD_COMPATIBILITY.md` guide
- **Examples**: `examples/basic_usage.go` demonstrating usage patterns

### 🔧 Changed

- **Refactor**: Split monolithic `datatables.go` into 9 modular files:
  - `datatables.go` - Package documentation (55 lines)
  - `processor.go` - Core OfReturn logic (133 lines)
  - `options.go` - Options builder (111 lines)
  - `parser.go` - Request parameter parsing (67 lines)
  - `converter.go` - Struct to map conversion (112 lines)
  - `transformer.go` - Column transformations (65 lines)
  - `validation.go` - Security validation (62 lines)
  - `errors.go` - Error types (28 lines)
  - `response.go` - JSON response helpers (42 lines)
- **Improved**: All code comments now in English
- **Improved**: Each file kept under 300 lines for better maintainability
- **Fixed**: Go version from invalid `1.25.3` to valid `1.23.0`
- **Fixed**: Removed hardcoded `created_at DESC` default ordering

### 🛡️ Security

- Added regex-based column name validation (`^[a-zA-Z0-9_.]+$`)
- All column names validated before SQL query construction
- Prevents SQL injection via column name manipulation
- Validates both searchable and orderable column mappings

### 🔄 Backward Compatibility

**✅ 100% backward compatible** - No breaking changes!

All existing code will continue to work without modifications:
- Function signatures unchanged
- Method names unchanged
- Default behavior unchanged
- All public APIs preserved

New features are opt-in and don't affect existing functionality.

### 📦 Dependencies

- Go 1.23.0 or higher
- `github.com/gin-gonic/gin` v1.11.0
- `gorm.io/gorm` v1.31.0

---

## Migration Guide

### Existing Users

**No action required!** Your code will continue to work.

### Recommended Updates

1. **Add default order configuration** (optional but recommended):
   ```go
   opts := datatables.NewOptions().
       WithDefaultOrder("id DESC")
   ```

2. **Use error helper** (optional):
   ```go
   if err != nil {
       datatables.JSONError(c, 500, err.Error())
       return
   }
   ```

### New Users

See `examples/basic_usage.go` for complete examples.

---

## What's Next?

Potential future enhancements (feedback welcome):
- [ ] Support for individual column search
- [ ] Advanced filtering options
- [ ] Export functionality (CSV, Excel)
- [ ] Caching layer for frequent queries
- [ ] Middleware for automatic setup

---

## Links

- [README](README.md) - Full documentation
- [Backward Compatibility Guide](BACKWARD_COMPATIBILITY.md)
- [Examples](examples/)
- [Report Issues](https://github.com/bonarizki-dat/Datatables-Gin/issues)
