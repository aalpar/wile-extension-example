// Package display provides output formatting helpers for example programs.
package display

import (
	"context"
	"fmt"

	"github.com/aalpar/wile"
)

// Section prints a section header.
func Section(name string) {
	fmt.Printf("\n=== %s ===\n", name)
}

// Run evaluates code, printing the label and result.
func Run(engine *wile.Engine, label, code string) {
	ctx := context.Background()
	result, err := engine.Eval(ctx, code)
	if err != nil {
		fmt.Printf("  %-30s ERROR: %v\n", label, err)
		return
	}
	if result.IsVoid() {
		fmt.Printf("  %-30s => (void)\n", label)
		return
	}
	fmt.Printf("  %-30s => %s\n", label, result.SchemeString())
}

// RunMultiple evaluates multiple expressions, printing the label and last result.
func RunMultiple(engine *wile.Engine, label, code string) {
	ctx := context.Background()
	result, err := engine.EvalMultiple(ctx, code)
	if err != nil {
		fmt.Printf("  %-30s ERROR: %v\n", label, err)
		return
	}
	if result.IsVoid() {
		fmt.Printf("  %-30s => (void)\n", label)
		return
	}
	fmt.Printf("  %-30s => %s\n", label, result.SchemeString())
}

// RunExpectError evaluates code that should fail, printing the error.
func RunExpectError(engine *wile.Engine, label, code string) {
	ctx := context.Background()
	_, err := engine.Eval(ctx, code)
	if err != nil {
		fmt.Printf("  %-30s => error: %v\n", label, err)
		return
	}
	fmt.Printf("  %-30s => unexpected success\n", label)
}
