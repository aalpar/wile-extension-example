// Package kvstore implements an in-memory key-value store as a Wile extension.
//
// It demonstrates the full extension authoring pattern:
//   - Implementing registry.Extension and registry.Closeable
//   - Stateful extension (the KVStore holds a map that primitives read/write)
//   - Registering ForeignFunction primitives via AddPrimitives
//   - Proper error handling with sentinel errors and WrapForeignErrorf
package kvstore

import (
	"fmt"
	"sync"

	"github.com/aalpar/wile/registry"
	"github.com/aalpar/wile/values"
)

// ErrKeyNotFound is returned when a key lookup fails without a default value.
var ErrKeyNotFound = values.NewStaticError("key not found")

// KVStore is an in-memory key-value store extension.
// It implements both registry.Extension and registry.Closeable.
type KVStore struct {
	mu   sync.RWMutex
	data map[string]string
}

// New creates a new KVStore extension.
func New() *KVStore {
	return &KVStore{
		data: make(map[string]string),
	}
}

// Name returns the extension name.
func (kv *KVStore) Name() string {
	return "kvstore"
}

// AddToRegistry registers all kvstore primitives.
func (kv *KVStore) AddToRegistry(r *registry.Registry) error {
	r.AddPrimitives(kv.primitiveSpecs(), registry.PhaseRuntime)
	return nil
}

// Close cleans up the store. Implements registry.Closeable.
func (kv *KVStore) Close() error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	count := len(kv.data)
	kv.data = nil
	fmt.Printf("[kvstore] closed (had %d entries)\n", count)
	return nil
}
