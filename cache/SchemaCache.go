package cache

import (
	"sync"

	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

const mainSchema = "__MAIN_SCHEMA"
const schemaCachePrefix = "__SCHEMA"
const allKeys = "__ALL_KEYS"

// SchemaCache ...
type SchemaCache struct {
	ttl    int
	prefix string
	mu     sync.Mutex
}

// NewSchemaCache ...
func NewSchemaCache(ttl int) *SchemaCache {
	if adapter == nil {
		adapter = newInMemoryCacheAdapter(5)
	}
	return &SchemaCache{
		ttl:    ttl,
		prefix: schemaCachePrefix + utils.CreateToken(),
	}
}

// Put ...
func (s *SchemaCache) Put(key string, value interface{}) {
	var keys map[string]bool
	v := get(s.prefix + allKeys)
	if v == nil {
		keys = map[string]bool{}
	} else {
		if m, ok := v.(map[string]bool); ok {
			keys = m
		} else {
			keys = map[string]bool{}
		}
	}
	if _, ok := keys[key]; ok == false {
		s.mu.Lock()
		keys[key] = true
		s.mu.Unlock()
	}
	put(s.prefix+allKeys, keys, int64(s.ttl))
	put(key, value, int64(s.ttl))
}

// GetAllClasses ...
func (s *SchemaCache) GetAllClasses() []types.M {
	if s.ttl < 0 {
		return nil
	}
	v := get(s.prefix + mainSchema)
	if r, ok := v.([]types.M); ok {
		return r
	}
	return []types.M{}
}

// SetAllClasses ...
func (s *SchemaCache) SetAllClasses(schema []types.M) {
	if s.ttl < 0 {
		return
	}
	s.Put(s.prefix+mainSchema, schema)
}

// SetOneSchema ...
func (s *SchemaCache) SetOneSchema(className string, schema types.M) {
	if s.ttl < 0 {
		return
	}
	s.Put(s.prefix+className, schema)
}

// GetOneSchema ...
func (s *SchemaCache) GetOneSchema(className string) types.M {
	if s.ttl < 0 {
		return nil
	}
	v := get(s.prefix + className)
	return utils.M(v)
}

// Clear ...
func (s *SchemaCache) Clear() {
	var keys map[string]bool
	v := get(s.prefix + allKeys)
	if v == nil {
		return
	}
	if m, ok := v.(map[string]bool); ok {
		keys = m
	} else {
		return
	}

	for key := range keys {
		del(key)
	}
	del(s.prefix + allKeys)
}
