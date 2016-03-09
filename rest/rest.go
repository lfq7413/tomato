package rest

import (
	"fmt"

	"github.com/lfq7413/tomato/auth"
)

// Find ...
func Find(
	auth *auth.Auth,
	className string,
	where map[string]interface{},
	options map[string]interface{},
) map[string]interface{} {

	enforceRoleSecurity("find", className, auth)
	query := NewQuery(auth, className, where, options)

	return query.Execute()
}

// Delete ...
func Delete(
	auth *auth.Auth,
	className string,
	objectID string,
) map[string]interface{} {

	if className == "_User" && auth.CouldUpdateUserID(objectID) == false {
		// TODO 权限不足
	}

	enforceRoleSecurity("delete", className, auth)

	var inflatedObject map[string]interface{}

	// TODO 获取要删除的对象

	destroy := NewDestroy(auth, className, map[string]interface{}{"objectId": objectID}, inflatedObject)

	return destroy.Execute()
}

// Create ...
func Create(
	auth *auth.Auth,
	className string,
	object map[string]interface{},
) map[string]interface{} {
	fmt.Println("object", object)
	return map[string]interface{}{}
}

// Update ...
func Update(
	auth *auth.Auth,
	className string,
	objectID string,
	object map[string]interface{},
) map[string]interface{} {
	fmt.Println("object", object)
	return map[string]interface{}{}
}

func enforceRoleSecurity(method string, className string, auth *auth.Auth) {
	if className == "_Role" && auth.IsMaster == false {
		// TODO 权限不足
	}
	if method == "delete" && className == "_Installation" && auth.IsMaster == false {
		// TODO 权限不足
	}
}
