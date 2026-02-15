// Command ffi-callbacks demonstrates RegisterFunc with callback parameters
// and context.Context forwarding.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

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
	// Callback functions
	must(engine.RegisterFuncs(map[string]any{
		// func(int64) int64 callback: Scheme lambda invoked from Go
		"apply-twice": func(f func(int64) int64, n int64) int64 {
			return f(f(n))
		},
		// func(int64) callback: Go calls a Scheme lambda for side effects
		"do-n-times": func(f func(int64), n int64) {
			for i := range n {
				f(int64(i))
			}
		},
		// Apply a callback to each element and collect results
		"map-ints": func(f func(int64) int64, ns []int64) []int64 {
			result := make([]int64, len(ns))
			for i, n := range ns {
				result[i] = f(n)
			}
			return result
		},
	}))

	// Accumulate callback results in Go (captures local state)
	var collected []int64
	must(engine.RegisterFunc("collect",
		func(f func(int64) int64, ns []int64) []int64 {
			collected = collected[:0]
			for _, n := range ns {
				collected = append(collected, f(n))
			}
			result := make([]int64, len(collected))
			copy(result, collected)
			return result
		}))

	// Context forwarding functions
	must(engine.RegisterFuncs(map[string]any{
		// context.Context as first param: automatically forwarded from VM
		"has-deadline?": func(ctx context.Context) bool {
			_, ok := ctx.Deadline()
			return ok
		},
		// context.Context + regular args
		"ctx-double": func(ctx context.Context, n int64) int64 {
			_ = ctx // available for cancellation checks, tracing, etc.
			return n * 2
		},
		// context.Context with deadline check
		"time-left-ms": func(ctx context.Context) int64 {
			deadline, ok := ctx.Deadline()
			if !ok {
				return -1
			}
			return time.Until(deadline).Milliseconds()
		},
	}))
}

func runExamples(engine *wile.Engine) {
	display.Section("Basic callback")
	display.Run(engine, "(apply-twice double 3)",
		"(apply-twice (lambda (x) (* x 2)) 3)")
	display.Run(engine, "(apply-twice inc 10)",
		"(apply-twice (lambda (x) (+ x 1)) 10)")

	display.Section("Void callback")
	display.Run(engine, "do-n-times (accumulate)",
		"(do-n-times (lambda (i) (+ i 100)) 3)")

	display.Section("Callback collecting results")
	display.Run(engine, "collect squares",
		"(collect (lambda (x) (* x x)) '(1 2 3 4 5))")

	display.Section("Callback with collection")
	display.Run(engine, "map-ints (square) over list",
		"(map-ints (lambda (x) (* x x)) '(1 2 3 4 5))")

	display.Section("Context forwarding")
	display.Run(engine, "(has-deadline?)", "(has-deadline?)")
	display.Run(engine, "(ctx-double 21)", "(ctx-double 21)")

	display.Section("Context with deadline")
	runWithDeadline(engine)

	display.Section("Scheme-defined function as callback")
	display.RunMultiple(engine, "define + callback", `
		(define (square x) (* x x))
		(apply-twice square 3)
	`)
}

func runWithDeadline(engine *wile.Engine) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := engine.Eval(ctx, "(has-deadline?)")
	if err != nil {
		fmt.Printf("  %-30s ERROR: %v\n", "with-deadline", err)
		return
	}
	fmt.Printf("  %-30s => %s\n", "(has-deadline?) with timeout", result.SchemeString())

	result, err = engine.Eval(ctx, "(time-left-ms)")
	if err != nil {
		fmt.Printf("  %-30s ERROR: %v\n", "time-left-ms", err)
		return
	}
	fmt.Printf("  %-30s => %s ms\n", "(time-left-ms)", result.SchemeString())
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
