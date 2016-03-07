package rest

import (
	"github.com/lfq7413/tomato/auth"
)

// Find ...
func Find(
	auth *auth.Auth,
	className string,
	where map[string]interface{},
	options map[string]interface{},
) []map[string]interface{} {
	return []map[string]interface{}{}
}

// Delete ...
func Delete(
	auth *auth.Auth,
	className string,
	objectID string,
) map[string]interface{} {
	return map[string]interface{}{}
}

// Create ...
func Create(
	auth *auth.Auth,
	className string,
	object map[string]interface{},
) map[string]interface{} {
	return map[string]interface{}{}
}

// Update ...
func Update(
	auth *auth.Auth,
	className string,
	objectID string,
	object map[string]interface{},
) map[string]interface{} {
	return map[string]interface{}{}
}
