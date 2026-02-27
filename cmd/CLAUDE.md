# cmd/ — Example programs

Each subdirectory is a standalone `main` package demonstrating a specific wile embedding pattern. All share `internal/display` for output formatting and follow the same structure: `registerFunctions()` + `runExamples()`.

## Examples by feature

### ffi-basics (`cmd/ffi-basics/main.go`)
**RegisterFuncs with scalar types.** Covers the full type mapping surface:
- `int64`, `float64`, `string`, `bool`, `[]byte`, `wile.Value` — all scalar types
- `(T, error)` return — error propagates as Scheme runtime error
- void return (no return value → `(void)`)
- variadic: `...int64`, `string + ...string`

### ffi-collections (`cmd/ffi-collections/main.go`)
**RegisterFuncs with composite types.** Three composite Go types:
- `[]T` ↔ Scheme proper list (both directions)
- `map[K]V` ↔ Scheme hashtable (both directions)
- `struct` ↔ Scheme alist `((Field . value) ...)` (both directions)

Defines a local `User` struct to demonstrate struct round-trips.

### ffi-callbacks (`cmd/ffi-callbacks/main.go`)
**RegisterFuncs with function parameters and context.Context.** Two categories:
- Callbacks: `func(int64) int64`, `func(int64)` (void), callbacks + slices, local state capture
- Context forwarding: `context.Context` as first parameter is auto-injected from the VM. Demonstrates deadline detection and time-remaining queries.

### scheme-library (`cmd/scheme-library/main.go`)
**R7RS library system.** Demonstrates `wile.WithLibraryPaths()`:
- Scheme `.sld` file at `cmd/scheme-library/lib/stats.sld`
- Library exports: `mean`, `variance`, `describe`
- Composition: library exports + Go-registered `format-result` used together

**Must be run from repo root** (library path is relative): `go run ./cmd/scheme-library`

### custom-extension (`cmd/custom-extension/main.go`)
**Full Extension pattern.** Loads `kvstore.New()` via `wile.WithExtension()`:
- Exercises all 6 kvstore primitives
- Demonstrates `engine.Close()` triggering `kvstore.Close()`
- Shows Scheme-side composition: `string-append` with `kv-get`
