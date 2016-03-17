package rest

import (
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
	return nil
}

func (w *Write) handleInstallation() error {
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
