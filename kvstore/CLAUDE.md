# kvstore

In-memory key-value store implementing wile's full extension pattern. This is the reference example for writing a wile extension.

## Design

`KVStore` implements two wile interfaces:
- **`registry.Extension`** — `Name() string` + `AddToRegistry(*registry.Registry) error`
- **`registry.Closeable`** — `Close() error` (called by `engine.Close()`)

State is a `sync.RWMutex`-protected `map[string]string`. Each primitive method is a receiver on `*KVStore`, capturing state through the receiver rather than closures.

## Extension loading flow

```
kvstore.New() → wile.NewEngine(ctx, wile.WithExtension(store))
                  ↓
                AddToRegistry(r) → r.AddPrimitives(specs, PhaseRuntime)
                  ↓
                Engine ready — primitives available in Scheme
                  ↓
                engine.Close() → kvstore.Close()
```

## Primitives

| Scheme name | Params | Variadic | Description |
|---|---|---|---|
| `kv-set!` | 2 (key, value) | no | Store a string key-value pair |
| `kv-get` | 1-2 (key [default]) | yes | Lookup; error if missing without default |
| `kv-delete!` | 1 (key) | no | Remove a key |
| `kv-keys` | 0 | no | Sorted list of all keys |
| `kv-count` | 0 | no | Number of entries |
| `kv-clear!` | 0 | no | Remove all entries |

## Error handling

- `ErrKeyNotFound` — sentinel error via `values.NewStaticError("key not found")`
- Errors wrapped with `values.WrapForeignErrorf(sentinel, format, args...)` for context
- Type errors use `values.ErrNotAString` sentinel

## ForeignFunction signature

As of wile v1.4.0, `machine.ForeignFunction` is `func(ctx context.Context, mc *MachineContext) error`. All primitive methods accept `context.Context` as their first parameter (currently unused, named `_`).

## Source index

| File | Key symbols |
|---|---|
| `extension.go:19` | `ErrKeyNotFound` — sentinel error |
| `extension.go:23` | `KVStore` struct (mu, data) |
| `extension.go:29` | `New()` constructor |
| `extension.go:41` | `AddToRegistry()` — registers primitives at PhaseRuntime |
| `extension.go:47` | `Close()` — prints summary, nils map |
| `primitives.go:13` | `primitiveSpecs()` — returns all 6 `PrimitiveSpec` entries |
| `primitives.go:65` | `primSet` |
| `primitives.go:84` | `primGet` — variadic, handles optional default |
| `primitives.go:120` | `primDelete` |
| `primitives.go:135` | `primKeys` — sorts keys, builds Scheme list |
| `primitives.go:154` | `primCount` |
| `primitives.go:164` | `primClear` |
| `primitives.go:174` | `requireString()` — shared argument extraction helper |
