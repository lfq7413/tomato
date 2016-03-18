package rest

import "github.com/lfq7413/tomato/utils"

// Find ...
func Find(
	auth *Auth,
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
	auth *Auth,
	className string,
	objectID string,
) map[string]interface{} {

	if className == "_User" && auth.CouldUpdateUserID(objectID) == false {
		// TODO 权限不足
	}

	enforceRoleSecurity("delete", className, auth)

	var inflatedObject map[string]interface{}

	if TriggerExists(TypeBeforeDelete, className) ||
		TriggerExists(TypeAfterDelete, className) ||
		className == "_Session" {
		response := Find(auth, className, map[string]interface{}{"objectId": objectID}, map[string]interface{}{})
		if utils.HasResults(response) == false {
			// TODO 未找到要删除的对象
		}

		result := utils.SliceInterface(response["results"])
		inflatedObject = utils.MapInterface(result[0])
		if inflatedObject == nil {
			// TODO 未找到要删除的对象
		}
	}

	destroy := NewDestroy(auth, className, map[string]interface{}{"objectId": objectID}, inflatedObject)

	return destroy.Execute()
}

// Create ...
func Create(
	auth *Auth,
	className string,
	object map[string]interface{},
) map[string]interface{} {

	enforceRoleSecurity("create", className, auth)
	write := NewWrite(auth, className, nil, object, nil)

	return write.Execute()
}

// Update ...
func Update(
	auth *Auth,
	className string,
	objectID string,
	object map[string]interface{},
) map[string]interface{} {

	enforceRoleSecurity("update", className, auth)

	var originalRestObject map[string]interface{}

	var response map[string]interface{}
	if TriggerExists(TypeBeforeSave, className) ||
		TriggerExists(TypeAfterSave, className) {
		response = Find(auth, className, map[string]interface{}{"objectId": objectID}, map[string]interface{}{})

		if utils.HasResults(response) == false {
			// TODO 未找到要更新的对象
		}

		result := utils.SliceInterface(response["results"])
		originalRestObject = utils.MapInterface(result[0])
		if originalRestObject == nil {
			// TODO 未找到要更新的对象
		}
	}

	write := NewWrite(auth, className, map[string]interface{}{"objectId": objectID}, object, originalRestObject)

	return write.Execute()
}

func enforceRoleSecurity(method string, className string, auth *Auth) {
	if className == "_Role" && auth.IsMaster == false {
		// TODO 权限不足
	}
	if method == "delete" && className == "_Installation" && auth.IsMaster == false {
		// TODO 权限不足
	}
}
