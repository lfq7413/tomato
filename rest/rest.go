package rest

import (
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// Find ...
func Find(
	auth *Auth,
	className string,
	where types.M,
	options types.M,
) (types.M, error) {

	enforceRoleSecurity("find", className, auth)
	query := NewQuery(auth, className, where, options)

	return query.Execute()
}

// Delete ...
func Delete(
	auth *Auth,
	className string,
	objectID string,
) types.M {

	if className == "_User" && auth.CouldUpdateUserID(objectID) == false {
		// TODO 权限不足
	}

	enforceRoleSecurity("delete", className, auth)

	var inflatedObject types.M

	if TriggerExists(TypeBeforeDelete, className) ||
		TriggerExists(TypeAfterDelete, className) ||
		className == "_Session" {
		// TODO 处理错误
		response, _ := Find(auth, className, types.M{"objectId": objectID}, types.M{})
		if utils.HasResults(response) == false {
			// TODO 未找到要删除的对象
		}

		result := utils.SliceInterface(response["results"])
		inflatedObject = utils.MapInterface(result[0])
		if inflatedObject == nil {
			// TODO 未找到要删除的对象
		}
	}

	destroy := NewDestroy(auth, className, types.M{"objectId": objectID}, inflatedObject)

	return destroy.Execute()
}

// Create ...
func Create(
	auth *Auth,
	className string,
	object types.M,
) (types.M, error) {

	enforceRoleSecurity("create", className, auth)
	write := NewWrite(auth, className, nil, object, nil)

	return write.Execute()
}

// Update ...
func Update(
	auth *Auth,
	className string,
	objectID string,
	object types.M,
) (types.M, error) {

	enforceRoleSecurity("update", className, auth)

	var originalRestObject types.M

	var response types.M
	if TriggerExists(TypeBeforeSave, className) ||
		TriggerExists(TypeAfterSave, className) {
		// TODO 处理错误
		response, _ = Find(auth, className, types.M{"objectId": objectID}, types.M{})

		if utils.HasResults(response) == false {
			// TODO 未找到要更新的对象
		}

		result := utils.SliceInterface(response["results"])
		originalRestObject = utils.MapInterface(result[0])
		if originalRestObject == nil {
			// TODO 未找到要更新的对象
		}
	}

	write := NewWrite(auth, className, types.M{"objectId": objectID}, object, originalRestObject)

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
