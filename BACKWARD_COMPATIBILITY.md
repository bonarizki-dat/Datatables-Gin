# Backward Compatibility Guide

## ✅ 100% Backward Compatible

All existing code using this package will continue to work without any changes.

---

## Old Code (Still Works!)

### Before Refactoring:

```go
package main

import (
    "github.com/bonarizki-dat/Datatables-Gin/datatables"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type User struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func GetUsers(c *gin.Context, db *gorm.DB) {
    var users []User

    searchable := []string{"name", "email"}

    orderable := map[string]string{
        "name":  "name",
        "email": "email",
    }

    opts := datatables.NewOptions()

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
```

**Status:** ✅ **STILL WORKS PERFECTLY!**

---

## What Changed?

### Internal Changes (Invisible to Users):

1. **Code Split into Multiple Files**
   - Before: Everything in `datatables.go`
   - After: Split into 9 files (`processor.go`, `options.go`, etc.)
   - **Impact:** ✅ None - all functions still exported from `datatables` package

2. **Added Security Validation**
   - Now validates column names to prevent SQL injection
   - **Impact:** ✅ None - validation is automatic and transparent
   - If invalid column detected, returns clear error message

3. **Better Error Handling**
   - Added custom error types
   - **Impact:** ✅ None - still returns standard Go `error` interface

---

## New Optional Features (Opt-in)

These are **additions**, not changes. Old code doesn't need to use them:

### 1. Default Order Configuration (New!)

```go
// OLD WAY - Still works!
opts := datatables.NewOptions()

// NEW WAY - Optional enhancement
opts := datatables.NewOptions().
    WithDefaultOrder("id DESC")  // Prevents "created_at not found" errors
```

### 2. Error Response Helper (New!)

```go
// OLD WAY - Still works!
if err != nil {
    c.JSON(500, gin.H{"error": err.Error()})
    return
}

// NEW WAY - Optional helper
if err != nil {
    datatables.JSONError(c, 500, err.Error())  // Consistent format
    return
}
```

---

## Migration Guide

### Do You Need to Change Anything?

**Short answer: NO!** ❌

Your existing code will work without any changes.

### Should You Update?

**Optional improvements you can make:**

#### If you had this error before:
```
Error: column "created_at" does not exist
```

**Fix:** Add default order configuration
```go
opts := datatables.NewOptions().
    WithDefaultOrder("id DESC")  // Use a column that exists in your table
```

#### If you want consistent error responses:

**Before:**
```go
c.JSON(500, gin.H{"error": err.Error()})
```

**After:**
```go
datatables.JSONError(c, 500, err.Error())
```

---

## API Comparison

| Function/Method | Before | After | Status |
|----------------|--------|-------|--------|
| `datatables.OfReturn[T]()` | ✅ | ✅ | Unchanged |
| `datatables.ParseParams()` | ✅ | ✅ | Unchanged |
| `datatables.NewOptions()` | ✅ | ✅ | Unchanged |
| `Options.WithIndex()` | ✅ | ✅ | Unchanged |
| `Options.Add()` | ✅ | ✅ | Unchanged |
| `Options.Edit()` | ✅ | ✅ | Unchanged |
| `Options.Remove()` | ✅ | ✅ | Unchanged |
| `datatables.JSON()` | ✅ | ✅ | Unchanged |
| `Options.WithDefaultOrder()` | ❌ | ✅ | **New (optional)** |
| `datatables.JSONError()` | ❌ | ✅ | **New (optional)** |

---

## Breaking Changes

**None.** This release is 100% backward compatible.

---

## Recommended Update Steps

1. **Update the package:**
   ```bash
   go get -u github.com/bonarizki-dat/Datatables-Gin
   ```

2. **Run your tests:**
   ```bash
   go test ./...
   ```

3. **Everything should work!** ✅

4. **Optionally**, add `WithDefaultOrder()` to prevent potential errors:
   ```go
   opts := datatables.NewOptions().
       WithDefaultOrder("id DESC")
   ```

---

## Questions?

If you encounter any issues:
1. Check that column names in `searchable` and `orderable` are valid
2. Ensure your `orderable` map uses actual database column names
3. If you get "invalid column name" errors, it means the security validation caught a potential SQL injection - this is a good thing!

---

**TL;DR:** Update safely. Your code will keep working. No changes required. 🎉
