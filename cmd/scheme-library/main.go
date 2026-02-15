// Command scheme-library demonstrates the hybrid Go + Scheme library pattern:
//
//	┌──────────────┐     ┌──────────────────┐     ┌─────────────────┐
//	│ Go: register  │────→│ .sld: compose    │────→│ Scheme: import  │
//	│ native funcs  │     │ into library API │     │ and use         │
//	└──────────────┘     └──────────────────┘     └─────────────────┘
//
// WithLibraryPaths enables the R7RS library system so embedders can ship
// Scheme libraries (.sld files) alongside Go code. Library environments are
// isolated — they see registry primitives (+, *, map, etc.) but not
// RegisterFunc bindings. Go functions and library exports are composed at
// the engine level.
//
// NOTE: Run from the repository root: go run ./cmd/scheme-library
// The library search path is relative to the working directory.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aalpar/wile"
	"github.com/aalpar/wile-extension-example/internal/display"
)

func main() {
	ctx := context.Background()

	// Create engine with library support.
	// The search path points to our lib/ directory where stats.sld lives.
	// This path is relative to the working directory (run from repo root).
	engine, err := wile.NewEngine(ctx,
		wile.WithLibraryPaths("./cmd/scheme-library/lib"),
	)
	if err != nil {
		log.Fatal(err)
	}

	registerFunctions(engine)
	runExamples(engine)
}

func registerFunctions(engine *wile.Engine) {
	// Register a Go function to demonstrate composition with library exports.
	must(engine.RegisterFunc("format-result", func(label string, value float64) string {
		return fmt.Sprintf("%s: %.4f", label, value)
	}))
}

func runExamples(engine *wile.Engine) {
	display.Section("Import (stats) library")
	display.RunMultiple(engine, "(import (stats))", "(import (stats))")

	display.Section("Mean")
	display.RunMultiple(engine, "(mean '(10 20 30))", `
		(import (stats))
		(mean '(10 20 30))
	`)
	display.RunMultiple(engine, "(mean '(1 2 3 4 5))", `
		(import (stats))
		(mean '(1 2 3 4 5))
	`)

	display.Section("Variance")
	display.RunMultiple(engine, "(variance '(2 4 4 4 5 5 7 9))", `
		(import (stats))
		(variance '(2 4 4 4 5 5 7 9))
	`)

	display.Section("Describe")
	display.RunMultiple(engine, "(describe '(2 4 4 4 5 5 7 9))", `
		(import (stats))
		(describe '(2 4 4 4 5 5 7 9))
	`)

	display.Section("Composing library exports with Go functions")
	display.RunMultiple(engine, "format-result + mean", `
		(import (stats))
		(format-result "average" (mean '(10 20 30)))
	`)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
