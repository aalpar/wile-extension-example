# wile-extension-example

Working examples for embedding [Wile](https://github.com/aalpar/wile) (an R7RS Scheme interpreter) in Go. Demonstrates the two extension mechanisms: `RegisterFunc` (reflection-based, natural Go signatures) and `registry.Extension` (low-level, direct `MachineContext` access).

## Architecture

```
cmd/                        Example programs (each is a standalone main)
  ffi-basics/                 RegisterFunc: scalars, errors, variadic
  ffi-collections/            RegisterFunc: slices, maps, structs
  ffi-callbacks/              RegisterFunc: callbacks, context.Context
  scheme-library/             R7RS library system (.sld) + Go composition
  custom-extension/           Full Extension pattern (kvstore)
kvstore/                    Reusable Extension package (registry.Extension + Closeable)
internal/display/           Output formatting helpers shared by all examples
```

Two extension axes in wile, both demonstrated here:

1. **RegisterFunc** (`cmd/ffi-*`) — reflection-based. Write `func(int64) string`, wile handles Scheme-to-Go conversion. Good for application-level bindings.
2. **Extension + PrimitiveSpec** (`kvstore/`) — implement `registry.Extension`, register `ForeignFunction` primitives via `AddPrimitives`. Direct `MachineContext` access for argument handling. Good for libraries needing state or fine-grained control.

## Dependency

Single dependency: `github.com/aalpar/wile` (currently v1.4.0). The `go.mod` is at the repo root.

## Build & Run

```bash
make build              # Build all examples to ./dist/{os}/{arch}/
make test               # Run tests
make ci                 # Full CI: lint + build + test + verify-mod
make run-examples       # Build + run all examples (30s timeout each)
go run ./cmd/ffi-basics # Run a single example directly
```

## Project Conventions

- **Commits:** No Co-Authored-By lines. Branch + PR workflow. No direct push to master.
- **Dependencies:** Prefer standard library. Single external dependency (wile).
- **Version:** v0.0.1 (see `VERSION`). Zero consumers — break freely.
- **Coverage:** 80% threshold enforced by `tools/sh/covercheck.sh`. Example programs (`cmd/*`) and `internal/display` are excluded.

## Source Index

### Core library package
| File | Purpose |
|---|---|
| `kvstore/extension.go` | `KVStore` struct, `New()`, `Name()`, `AddToRegistry()`, `Close()` — Extension + Closeable impl |
| `kvstore/primitives.go` | `primitiveSpecs()` returns `[]PrimitiveSpec`; 6 primitive methods + `requireString` helper |

### Example programs (cmd/)
| Directory | Wile feature demonstrated |
|---|---|
| `cmd/ffi-basics/` | `RegisterFuncs` with scalars, errors, void, variadic |
| `cmd/ffi-collections/` | `RegisterFuncs` with `[]T`, `map[K]V`, struct |
| `cmd/ffi-callbacks/` | `RegisterFuncs` with `func(...)` callbacks, `context.Context` forwarding |
| `cmd/scheme-library/` | `WithLibraryPaths` + R7RS `.sld` library loading |
| `cmd/custom-extension/` | `WithExtension` loading the kvstore package |

### Supporting files
| File | Purpose |
|---|---|
| `internal/display/display.go` | `Section()`, `Run()`, `RunMultiple()`, `RunExpectError()` — output helpers |
| `cmd/scheme-library/lib/stats.sld` | Pure Scheme R7RS library: `mean`, `variance`, `describe` |
| `Makefile` | Build, test, lint, coverage, Docker, version bumping |
| `VERSION` | Semver string read by Makefile for `LDFLAGS` |
| `docker/Dockerfile` | Go 1.23 bookworm image, builds all examples |
| `tools/sh/` | Shell scripts for CI/CD tasks (run-examples, smoke-test, covercheck, bump-version, docker-build/shell) |
