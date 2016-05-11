package rest

import (
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// Find 根据条件查找数据
// 返回格式如下：
// {
// 	"results":[
// 		{...},
// 	],
// 	"count":10
// }
func Find(
	auth *Auth,
	className string,
	where types.M,
	options types.M,
) (types.M, error) {

	err := enforceRoleSecurity("find", className, auth)
	if err != nil {
		return nil, err
	}
	query, err := NewQuery(auth, className, where, options)
	if err != nil {
		return nil, err
	}

	return query.Execute()
}

// Delete 删除指定对象
func Delete(
	auth *Auth,
	className string,
	objectID string,
) error {

	if className == "_User" && auth.CouldUpdateUserID(objectID) == false {
		return errs.E(errs.SessionMissing, "insufficient auth to delete user")
	}

	err := enforceRoleSecurity("delete", className, auth)
	if err != nil {
		return err
	}

	var inflatedObject types.M
	// 如果存在删前回调、或者删后回调、或者要删除的属于 _Session 类，则需要获取到要删除的对象数据
	if TriggerExists(TypeBeforeDelete, className) ||
		TriggerExists(TypeAfterDelete, className) ||
		(config.TConfig.LiveQuery != nil && config.TConfig.LiveQuery.HasLiveQuery(className)) ||
		className == "_Session" {
		response, err := Find(auth, className, types.M{"objectId": objectID}, types.M{})
		if err != nil || utils.HasResults(response) == false {
			return errs.E(errs.ObjectNotFound, "Object not found for delete.")
		}

		result := utils.SliceInterface(response["results"])
		inflatedObject = utils.MapInterface(result[0])
		if inflatedObject == nil {
			return errs.E(errs.ObjectNotFound, "Object not found for delete.")
		}
		inflatedObject["className"] = className
	}

	destroy := NewDestroy(auth, className, types.M{"objectId": objectID}, inflatedObject)

	return destroy.Execute()
}

// Create 创建对象
// 返回数据格式如下：
// {
// 	"status":201,
// 	"response":{...},
// 	"location":"http://..."
// }
func Create(
	auth *Auth,
	className string,
	object types.M,
) (types.M, error) {

	err := enforceRoleSecurity("create", className, auth)
	if err != nil {
		return nil, err
	}
	write, err := NewWrite(auth, className, nil, object, nil)
	if err != nil {
		return nil, err
	}

	return write.Execute()
}

// Update 更新对象
// 返回更新后的字段，一般只有 updatedAt
func Update(
	auth *Auth,
	className string,
	objectID string,
	object types.M,
) (types.M, error) {

	err := enforceRoleSecurity("update", className, auth)
	if err != nil {
		return nil, err
	}

	var originalRestObject types.M

	// 如果存在删前回调、或者删后回调，则需要获取到要删除的对象数据
	var response types.M
	if TriggerExists(TypeBeforeSave, className) ||
		TriggerExists(TypeAfterSave, className) ||
		(config.TConfig.LiveQuery != nil && config.TConfig.LiveQuery.HasLiveQuery(className)) {

		response, err = Find(auth, className, types.M{"objectId": objectID}, types.M{})
		if err != nil || utils.HasResults(response) == false {
			return nil, errs.E(errs.ObjectNotFound, "Object not found for update.")
		}

		result := utils.SliceInterface(response["results"])
		originalRestObject = utils.MapInterface(result[0])
		if originalRestObject == nil {
			return nil, errs.E(errs.ObjectNotFound, "Object not found for update.")
		}
	}

	write, err := NewWrite(auth, className, types.M{"objectId": objectID}, object, originalRestObject)
	if err != nil {
		return nil, err
	}

	return write.Execute()
}

// enforceRoleSecurity 对指定的类与操作进行安全校验
func enforceRoleSecurity(method string, className string, auth *Auth) error {
	// 非 Master 不得对 _Installation 进行删除操作
	if method == "delete" && className == "_Installation" && auth.IsMaster == false {
		return errs.E(errs.OperationForbidden, "Clients aren't allowed to perform the delete operation on the installation collection.")
	}
	return nil
}
