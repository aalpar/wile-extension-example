package kvstore

import (
	"context"
	"sort"

	"github.com/aalpar/wile/machine"
	"github.com/aalpar/wile/registry"
	"github.com/aalpar/wile/values"
)

// primitiveSpecs returns the PrimitiveSpec slice for all kvstore operations.
// Each spec's Impl is a method on *KVStore, capturing state via the receiver.
func (kv *KVStore) primitiveSpecs() []registry.PrimitiveSpec {
	return []registry.PrimitiveSpec{
		{
			Name:       "kv-set!",
			ParamCount: 2,
			Impl:       kv.primSet,
			Doc:        "Set a key-value pair (both strings).",
			ParamNames: []string{"key", "value"},
			Category:   "kvstore",
		},
		{
			Name:       "kv-get",
			ParamCount: 2,
			IsVariadic: true,
			Impl:       kv.primGet,
			Doc:        "Get value by key. Optional default if key missing.",
			ParamNames: []string{"key", "default"},
			Category:   "kvstore",
		},
		{
			Name:       "kv-delete!",
			ParamCount: 1,
			Impl:       kv.primDelete,
			Doc:        "Delete a key.",
			ParamNames: []string{"key"},
			Category:   "kvstore",
		},
		{
			Name:       "kv-keys",
			ParamCount: 0,
			Impl:       kv.primKeys,
			Doc:        "Return a sorted list of all keys.",
			Category:   "kvstore",
		},
		{
			Name:       "kv-count",
			ParamCount: 0,
			Impl:       kv.primCount,
			Doc:        "Return the number of entries.",
			Category:   "kvstore",
		},
		{
			Name:       "kv-clear!",
			ParamCount: 0,
			Impl:       kv.primClear,
			Doc:        "Remove all entries.",
			Category:   "kvstore",
		},
	}
}

// primSet implements (kv-set! key value).
func (kv *KVStore) primSet(_ context.Context, mc *machine.MachineContext) error {
	key, err := requireString(mc, 0, "kv-set!")
	if err != nil {
		return err
	}
	val, err := requireString(mc, 1, "kv-set!")
	if err != nil {
		return err
	}

	kv.mu.Lock()
	kv.data[key] = val
	kv.mu.Unlock()

	mc.SetValue(values.Void)
	return nil
}

// primGet implements (kv-get key [default]).
func (kv *KVStore) primGet(_ context.Context, mc *machine.MachineContext) error {
	key, err := requireString(mc, 0, "kv-get")
	if err != nil {
		return err
	}

	// Parse optional default from rest args.
	rest := mc.Arg(1)
	var defaultVal values.Value
	hasDefault := false
	if !values.IsEmptyList(rest) {
		tuple, ok := rest.(values.Tuple)
		if ok {
			defaultVal = tuple.Car()
			hasDefault = true
		}
	}

	kv.mu.RLock()
	val, found := kv.data[key]
	kv.mu.RUnlock()

	if !found {
		if hasDefault {
			mc.SetValue(defaultVal)
			return nil
		}
		return values.WrapForeignErrorf(ErrKeyNotFound,
			"kv-get: key %q not found", key)
	}

	mc.SetValue(values.NewString(val))
	return nil
}

// primDelete implements (kv-delete! key).
func (kv *KVStore) primDelete(_ context.Context, mc *machine.MachineContext) error {
	key, err := requireString(mc, 0, "kv-delete!")
	if err != nil {
		return err
	}

	kv.mu.Lock()
	delete(kv.data, key)
	kv.mu.Unlock()

	mc.SetValue(values.Void)
	return nil
}

// primKeys implements (kv-keys) → sorted list of all keys.
func (kv *KVStore) primKeys(_ context.Context, mc *machine.MachineContext) error {
	kv.mu.RLock()
	keys := make([]string, 0, len(kv.data))
	for k := range kv.data {
		keys = append(keys, k)
	}
	kv.mu.RUnlock()

	sort.Strings(keys)

	elems := make([]values.Value, len(keys))
	for i, k := range keys {
		elems[i] = values.NewString(k)
	}
	mc.SetValue(values.List(elems...))
	return nil
}

// primCount implements (kv-count) → number of entries.
func (kv *KVStore) primCount(_ context.Context, mc *machine.MachineContext) error {
	kv.mu.RLock()
	n := len(kv.data)
	kv.mu.RUnlock()

	mc.SetValue(values.NewInteger(int64(n)))
	return nil
}

// primClear implements (kv-clear!) → removes all entries.
func (kv *KVStore) primClear(_ context.Context, mc *machine.MachineContext) error {
	kv.mu.Lock()
	clear(kv.data)
	kv.mu.Unlock()

	mc.SetValue(values.Void)
	return nil
}

// requireString extracts a string argument from the given index.
func requireString(mc *machine.MachineContext, index int, name string) (string, error) {
	v := mc.Arg(index)
	s, ok := v.(*values.String)
	if !ok {
		return "", values.WrapForeignErrorf(values.ErrNotAString,
			"%s: expected string at argument %d but got %T", name, index+1, v)
	}
	return s.Value, nil
}
