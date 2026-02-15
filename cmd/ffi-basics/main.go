// Command ffi-basics demonstrates RegisterFunc with scalar types,
// error returns, void functions, and variadic parameters.
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/aalpar/wile"
	"github.com/aalpar/wile-extension-example/internal/display"
)

func main() {
	ctx := context.Background()
	engine, err := wile.NewEngine(ctx)
	if err != nil {
		log.Fatal(err)
	}

	registerFunctions(engine)
	runExamples(engine)
}

func registerFunctions(engine *wile.Engine) {
	// int64 → int64
	must(engine.RegisterFunc("double", func(n int64) int64 {
		return n * 2
	}))

	// float64 → float64
	must(engine.RegisterFunc("circle-area", func(r float64) float64 {
		return math.Pi * r * r
	}))

	// string → string
	must(engine.RegisterFunc("greet", func(s string) string {
		return "Hello, " + s + "!"
	}))

	// int64 → bool
	must(engine.RegisterFunc("even?", func(n int64) bool {
		return n%2 == 0
	}))

	// []byte → int64
	must(engine.RegisterFunc("byte-count", func(data []byte) int64 {
		return int64(len(data))
	}))

	// Value → Value (pass-through)
	must(engine.RegisterFunc("identity", func(v wile.Value) wile.Value {
		return v
	}))

	// (float64, float64) → (float64, error)
	must(engine.RegisterFunc("safe-divide", func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	}))

	// string → void (side effect only)
	must(engine.RegisterFunc("log-message", func(msg string) {
		fmt.Printf("    [LOG] %s\n", msg)
	}))

	// variadic int64
	must(engine.RegisterFunc("sum", func(nums ...int64) int64 {
		var total int64
		for _, n := range nums {
			total += n
		}
		return total
	}))

	// fixed prefix + variadic rest
	must(engine.RegisterFunc("join", func(sep string, parts ...string) string {
		return strings.Join(parts, sep)
	}))
}

func runExamples(engine *wile.Engine) {
	display.Section("Integer round-trip")
	display.Run(engine, "(double 21)", "(double 21)")
	display.Run(engine, "(double -5)", "(double -5)")

	display.Section("Float computation")
	display.Run(engine, "(circle-area 5.0)", "(circle-area 5.0)")

	display.Section("String transformation")
	display.Run(engine, `(greet "world")`, `(greet "world")`)

	display.Section("Boolean return")
	display.Run(engine, "(even? 4)", "(even? 4)")
	display.Run(engine, "(even? 7)", "(even? 7)")

	display.Section("Bytevector parameter")
	display.Run(engine, "(byte-count #u8(1 2 3))", "(byte-count #u8(1 2 3))")

	display.Section("Value pass-through")
	display.Run(engine, "(identity '(1 2 3))", "(identity '(1 2 3))")
	display.Run(engine, "(identity 'hello)", "(identity 'hello)")

	display.Section("Error return")
	display.Run(engine, "(safe-divide 10.0 3.0)", "(safe-divide 10.0 3.0)")
	display.RunExpectError(engine, "(safe-divide 10.0 0.0)", "(safe-divide 10.0 0.0)")

	display.Section("Void return (side effect)")
	display.Run(engine, "(log-message \"startup\")", `(log-message "startup")`)

	display.Section("Variadic — sum")
	display.Run(engine, "(sum)", "(sum)")
	display.Run(engine, "(sum 1 2 3 4 5)", "(sum 1 2 3 4 5)")

	display.Section("Fixed prefix + variadic")
	display.Run(engine, `(join "-" "a" "b" "c")`, `(join "-" "a" "b" "c")`)
	display.Run(engine, `(join ", " "x")`, `(join ", " "x")`)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
