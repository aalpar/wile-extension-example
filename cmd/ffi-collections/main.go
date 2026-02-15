// Command ffi-collections demonstrates RegisterFunc with composite types:
// slices, maps, and structs.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aalpar/wile"
	"github.com/aalpar/wile-extension-example/internal/display"
)

// User is an example struct for struct round-trip demonstrations.
type User struct {
	Name string
	Age  int64
}

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
	// --- Slice parameters ---

	// []int64 parameter: Scheme list → Go slice
	must(engine.RegisterFunc("sum-list", func(ns []int64) int64 {
		var total int64
		for _, n := range ns {
			total += n
		}
		return total
	}))

	// []string return: Go slice → Scheme list
	must(engine.RegisterFunc("make-tags", func(n int) []string {
		tags := make([]string, n)
		for i := range tags {
			tags[i] = fmt.Sprintf("tag-%d", i+1)
		}
		return tags
	}))

	// --- Map parameters ---

	// map[string]int64 parameter: Scheme hashtable → Go map
	must(engine.RegisterFunc("total-score", func(scores map[string]int64) int64 {
		var total int64
		for _, v := range scores {
			total += v
		}
		return total
	}))

	// map[string]int64 return: Go map → Scheme hashtable
	must(engine.RegisterFunc("default-config", func() map[string]int64 {
		return map[string]int64{
			"timeout": 30,
			"retries": 3,
			"port":    8080,
		}
	}))

	// --- Struct parameters ---

	// struct parameter: Scheme alist → Go struct
	must(engine.RegisterFunc("greet-user", func(u User) string {
		return fmt.Sprintf("Hello, %s (age %d)!", u.Name, u.Age)
	}))

	// struct return: Go struct → Scheme alist
	must(engine.RegisterFunc("make-user", func(name string, age int64) User {
		return User{Name: name, Age: age}
	}))
}

func runExamples(engine *wile.Engine) {
	display.Section("Slice parameter (list → []int64)")
	display.Run(engine, "(sum-list '(10 20 30))", "(sum-list '(10 20 30))")
	display.Run(engine, "(sum-list '())", "(sum-list '())")

	display.Section("Slice return ([]string → list)")
	display.Run(engine, "(make-tags 3)", "(make-tags 3)")

	display.Section("Map parameter (hashtable → map)")
	display.RunMultiple(engine, "total-score", `
		(let ((ht (make-hashtable)))
			(hashtable-set! ht "math" 95)
			(hashtable-set! ht "science" 87)
			(hashtable-set! ht "english" 92)
			(total-score ht))
	`)

	display.Section("Map return (map → hashtable)")
	display.Run(engine, "(hashtable-size (default-config))", "(hashtable-size (default-config))")
	display.RunMultiple(engine, "get timeout value", `
		(let ((cfg (default-config)))
			(hashtable-ref cfg "timeout"))
	`)

	display.Section("Struct parameter (alist → struct)")
	display.Run(engine, "greet-user",
		`(greet-user '((Name . "Alice") (Age . 30)))`)

	display.Section("Struct return (struct → alist)")
	display.Run(engine, `(make-user "Bob" 25)`,
		`(make-user "Bob" 25)`)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
