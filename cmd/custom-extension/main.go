// Command custom-extension demonstrates loading a custom extension
// that implements registry.Extension and registry.Closeable.
package main

import (
	"context"
	"log"

	"github.com/aalpar/wile"
	"github.com/aalpar/wile-extension-example/internal/display"
	"github.com/aalpar/wile-extension-example/kvstore"
)

func main() {
	ctx := context.Background()

	// Create an engine with the kvstore extension loaded.
	store := kvstore.New()
	engine, err := wile.NewEngine(ctx, wile.WithExtension(store))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		// Close calls kvstore's Close(), which prints a summary.
		closeErr := engine.Close()
		if closeErr != nil {
			log.Fatal(closeErr)
		}
	}()

	display.Section("kv-set! and kv-get")
	display.Run(engine, `(kv-set! "host" "localhost")`, `(kv-set! "host" "localhost")`)
	display.Run(engine, `(kv-set! "port" "8080")`, `(kv-set! "port" "8080")`)
	display.Run(engine, `(kv-get "host")`, `(kv-get "host")`)
	display.Run(engine, `(kv-get "port")`, `(kv-get "port")`)

	display.Section("kv-get with default")
	display.Run(engine, `(kv-get "missing" "N/A")`, `(kv-get "missing" "N/A")`)

	display.Section("kv-get without default (error)")
	display.RunExpectError(engine, `(kv-get "missing")`, `(kv-get "missing")`)

	display.Section("kv-count and kv-keys")
	display.Run(engine, "(kv-count)", "(kv-count)")
	display.Run(engine, "(kv-keys)", "(kv-keys)")

	display.Section("kv-delete!")
	display.Run(engine, `(kv-delete! "port")`, `(kv-delete! "port")`)
	display.Run(engine, "(kv-count)", "(kv-count)")
	display.Run(engine, "(kv-keys)", "(kv-keys)")

	display.Section("kv-clear!")
	display.Run(engine, "(kv-clear!)", "(kv-clear!)")
	display.Run(engine, "(kv-count)", "(kv-count)")

	display.Section("Use from Scheme")
	display.RunMultiple(engine, "store and retrieve", `
		(kv-set! "greeting" "hello")
		(string-append (kv-get "greeting") ", world!")
	`)
}
