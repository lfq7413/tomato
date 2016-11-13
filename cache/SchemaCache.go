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
// singleCache 默认为 false
func NewSchemaCache(ttl int, singleCache bool) *SchemaCache {
	if adapter == nil {
		adapter = newInMemoryCacheAdapter(5)
	}
	prefix := schemaCachePrefix
	if singleCache == false {
		prefix = prefix + utils.CreateToken()
	}
	return &SchemaCache{
		ttl:    ttl,
		prefix: prefix,
	}
}

// Put ...
func (s *SchemaCache) Put(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var keys map[string]interface{}
	v := get(s.prefix + allKeys)
	if v == nil {
		keys = map[string]interface{}{}
	} else {
		if m, ok := v.(map[string]interface{}); ok {
			keys = m
		} else {
			keys = map[string]interface{}{}
		}
	}
	if _, ok := keys[key]; ok == false {
		keys[key] = true
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
	} else if r, ok := v.([]interface{}); ok {
		res := []types.M{}
		for _, m := range r {
			if subv := utils.M(m); subv != nil {
				res = append(res, subv)
			}
		}
		return res
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
	schema := utils.M(v)
	if schema != nil {
		return schema
	}
	// 从 mainSchema 中查找
	cachedSchemas := []types.M{}
	v = get(s.prefix + mainSchema)
	if r, ok := v.([]types.M); ok {
		cachedSchemas = r
	} else if r, ok := v.([]interface{}); ok {
		res := []types.M{}
		for _, m := range r {
			if subv := utils.M(m); subv != nil {
				res = append(res, subv)
			}
		}
		cachedSchemas = res
	}

	for _, s := range cachedSchemas {
		if utils.S(s["className"]) == className {
			schema = s
			break
		}
	}

	if schema != nil {
		return schema
	}
	return nil
}

// Clear ...
func (s *SchemaCache) Clear() {
	var keys map[string]interface{}
	v := get(s.prefix + allKeys)
	if v == nil {
		return
	}
	if m, ok := v.(map[string]interface{}); ok {
		keys = m
	} else {
		return
	}

	for key := range keys {
		del(key)
	}
	del(s.prefix + allKeys)
}
