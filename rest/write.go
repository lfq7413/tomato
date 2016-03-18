package rest

import (
	"strings"
	"time"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/utils"
)

// Write ...
type Write struct {
	auth         *Auth
	className    string
	query        map[string]interface{}
	data         map[string]interface{}
	originalData map[string]interface{}
	storage      map[string]interface{}
	runOptions   map[string]interface{}
	response     map[string]interface{}
	updatedAt    time.Time
}

// NewWrite 可用于 create 和 update ， create 时 query 为 nil
func NewWrite(
	auth *Auth,
	className string,
	query map[string]interface{},
	data map[string]interface{},
	originalData map[string]interface{},
) *Write {
	if query == nil && data["objectId"] != nil {
		// TODO objectId 无效
	}
	write := &Write{
		auth:         auth,
		className:    className,
		query:        query,
		data:         data,
		originalData: originalData,
		storage:      map[string]interface{}{},
		runOptions:   map[string]interface{}{},
		response:     nil,
		updatedAt:    time.Now().UTC(),
	}
	return write
}

// Execute ...
func (w *Write) Execute() map[string]interface{} {
	w.getUserAndRoleACL()
	w.validateClientClassCreation()
	w.validateSchema()
	w.handleInstallation()
	w.handleSession()
	w.validateAuthData()
	w.runBeforeTrigger()
	w.setRequiredFieldsIfNeeded()
	w.transformUser()
	w.expandFilesForExistingObjects()
	w.runDatabaseOperation()
	w.handleFollowup()
	w.runAfterTrigger()
	return w.response
}

func (w *Write) getUserAndRoleACL() error {
	if w.auth.IsMaster {
		return nil
	}
	w.runOptions["acl"] = []string{"*"}
	if w.auth.User != nil {
		roles := w.auth.GetUserRoles()
		roles = append(roles, w.auth.User.ID)
		if v, ok := w.runOptions["acl"].([]string); ok {
			v = utils.AppendString(v, roles)
		}
	}
	return nil
}

func (w *Write) validateClientClassCreation() error {
	sysClass := []string{"_User", "_Installation", "_Role", "_Session", "_Product"}
	if config.TConfig.AllowClientClassCreation {
		return nil
	}
	if w.auth.IsMaster {
		return nil
	}
	for _, v := range sysClass {
		if v == w.className {
			return nil
		}
	}
	if orm.CollectionExists(w.className) {
		return nil
	}
	// TODO 无法操作不存在的表
	return nil
}

func (w *Write) validateSchema() error {
	// TODO 判断是否可以进行操作
	orm.ValidateObject(w.className, w.data, w.query, w.runOptions)
	return nil
}

func (w *Write) handleInstallation() error {
	if w.response != nil || w.className != "_Installation" {
		return nil
	}

	if w.query == nil && w.data["deviceToken"] == nil && w.data["installationId"] == nil {
		// TODO 设备 id 不能为空
		return nil
	}

	if w.query == nil && w.data["deviceType"] == nil {
		// TODO 设备类型不能为空
		return nil
	}

	if w.data["deviceToken"] != nil && len(utils.String(w.data["deviceToken"])) == 64 {
		w.data["deviceToken"] = strings.ToLower(utils.String(w.data["deviceToken"]))
	}

	if w.data["installationId"] != nil {
		w.data["installationId"] = strings.ToLower(utils.String(w.data["installationId"]))
	}

	var idMatch map[string]interface{}

	if w.query != nil && w.query["objectId"] != nil {
		results := orm.Find("_Installation", map[string]interface{}{"objectId": w.query["objectId"]}, map[string]interface{}{})
		if results == nil || len(results) == 0 {
			// TODO 更新对象未找到
			return nil
		}
		idMatch = utils.MapInterface(results[0])
		if w.data["installationId"] != nil && idMatch["installationId"] != nil &&
			w.data["installationId"] != idMatch["installationId"] {
			//TODO installationId 不能修改
			return nil
		}
		if w.data["deviceToken"] != nil && idMatch["deviceToken"] != nil &&
			w.data["deviceToken"] != idMatch["deviceToken"] &&
			w.data["installationId"] == nil && idMatch["installationId"] == nil {
			//TODO deviceToken 不能修改
			return nil
		}
		if w.data["deviceType"] != nil && idMatch["deviceType"] != nil &&
			w.data["deviceType"] != idMatch["deviceType"] {
			//TODO deviceType 不能修改
			return nil
		}
		return nil
	}

	return nil
}

func (w *Write) handleSession() error {
	return nil
}

func (w *Write) validateAuthData() error {
	return nil
}

func (w *Write) runBeforeTrigger() error {
	return nil
}

func (w *Write) setRequiredFieldsIfNeeded() error {
	return nil
}

func (w *Write) transformUser() error {
	return nil
}

func (w *Write) expandFilesForExistingObjects() error {
	return nil
}

func (w *Write) runDatabaseOperation() error {
	return nil
}

func (w *Write) handleFollowup() error {
	return nil
}

func (w *Write) runAfterTrigger() error {
	return nil
}
