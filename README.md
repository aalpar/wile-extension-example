# wile-extension-example

Working examples for embedding [Wile](https://github.com/aalpar/wile) (a Scheme interpreter) in Go applications.

## Prerequisites

- Go 1.23+
- `github.com/aalpar/wile` v1.3.0+

## Quick Start

```bash
go run ./cmd/ffi-basics
```

## Examples

### `cmd/ffi-basics` — Scalar Types, Errors, Variadic

Demonstrates `RegisterFunc` with natural Go signatures:

| Go signature | Scheme | Shows |
|---|---|---|
| `func(int64) int64` | `(double 21)` → `42` | integer round-trip |
| `func(float64) float64` | `(circle-area 5.0)` | float64 |
| `func(string) string` | `(greet "world")` | string |
| `func(int64) bool` | `(even? 4)` | bool return |
| `func([]byte) int64` | `(byte-count #u8(1 2 3))` | bytevector param |
| `func(Value) Value` | `(identity '(1 2))` | pass-through |
| `func(float64, float64) (float64, error)` | `(safe-divide 10 0)` | error return |
| `func(string)` | `(log-message "hi")` | void return |
| `func(...int64) int64` | `(sum 1 2 3)` | variadic |
| `func(string, ...string) string` | `(join "-" "a" "b")` | prefix + variadic |

```bash
go run ./cmd/ffi-basics
```

### `cmd/ffi-collections` — Slices, Maps, Structs

Demonstrates `RegisterFunc` with composite Go types:

| Go type | Scheme representation | Direction |
|---|---|---|
| `[]int64` | proper list `'(1 2 3)` | Scheme → Go |
| `[]string` | proper list `'("a" "b")` | Go → Scheme |
| `map[string]int64` | hashtable | both directions |
| `struct{Name string; Age int64}` | alist `'((Name . "Alice") (Age . 30))` | both directions |

```bash
go run ./cmd/ffi-collections
```

### `cmd/ffi-callbacks` — Callbacks & Context

Demonstrates callback parameters and `context.Context` forwarding:

| Go signature | Scheme | Shows |
|---|---|---|
| `func(func(int64) int64, int64) int64` | `(apply-twice (lambda (x) (* x 2)) 3)` | basic callback |
| `func(func(int64), int64)` | `(do-n-times (lambda (i) (display i)) 3)` | void callback |
| `func(context.Context) bool` | `(has-deadline?)` | context forwarding |
| `func(context.Context, int64) int64` | `(ctx-double 21)` | context + args |

```bash
go run ./cmd/ffi-callbacks
```

### `cmd/custom-extension` — Writing an Extension

Demonstrates the full extension authoring pattern using a key-value store (`kvstore/`):

- Implements `registry.Extension` (adds primitives to the registry)
- Implements `registry.Closeable` (cleanup on `engine.Close()`)
- Stateful: the `*KVStore` holds a `map[string]string` that primitives read and write
- Uses `machine.ForeignFunction` signature with `MachineContext` for argument access

Primitives provided:

| Primitive | Args | Description |
|---|---|---|
| `kv-set!` | 2 | Set a key-value pair |
| `kv-get` | 1-2 | Get by key, optional default |
| `kv-delete!` | 1 | Delete a key |
| `kv-keys` | 0 | List all keys (sorted) |
| `kv-count` | 0 | Number of entries |
| `kv-clear!` | 0 | Remove all entries |

```bash
go run ./cmd/custom-extension
```

## Writing Your Own Extension

See `kvstore/` for a complete example. The pattern is:

```go
package myext

import "github.com/aalpar/wile/registry"

type MyExtension struct {
    // your state here
}

func New() *MyExtension {
    return &MyExtension{}
}

func (e *MyExtension) Name() string {
    return "my-extension"
}

func (e *MyExtension) AddToRegistry(r *registry.Registry) error {
    r.AddPrimitives([]registry.PrimitiveSpec{
        // your primitives here
    }, registry.PhaseRuntime)
    return nil
}

// Optional: implement registry.Closeable for cleanup
func (e *MyExtension) Close() error {
    return nil
}
```

Load it:

```go
engine, _ := wile.NewEngine(ctx, wile.WithExtension(myext.New()))
```

## Key Concepts

**`RegisterFunc` vs `RegisterPrimitive`**: `RegisterFunc` accepts natural Go signatures
and handles type conversion automatically via reflection. `RegisterPrimitive` gives direct
access to `MachineContext` for argument handling. Use `RegisterFunc` for application-level
bindings; use `RegisterPrimitive` (via extensions) when you need fine-grained control.

**Type mappings** (`RegisterFunc`):

| Go type | Scheme type |
|---|---|
| `int64`, `int` | exact integer |
| `float64` | inexact real |
| `string` | string |
| `bool` | boolean |
| `[]byte` | bytevector |
| `[]T` | proper list |
| `map[K]V` | hashtable |
| `struct` | alist `((Field . value) ...)` |
| `func(...)` | procedure (lambda) |
| `wile.Value` | any Scheme value |
| `context.Context` | forwarded from VM (first param only) |
| `error` | runtime error (last return only) |

## Links

- [Wile repository](https://github.com/aalpar/wile)
