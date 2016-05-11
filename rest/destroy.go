package rest

import (
	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
)

// Destroy 删除对象
type Destroy struct {
	auth         *Auth
	className    string
	query        types.M
	originalData types.M
}

// NewDestroy 组装 Destroy
func NewDestroy(
	auth *Auth,
	className string,
	query types.M,
	originalData types.M,
) *Destroy {
	destroy := &Destroy{
		auth:         auth,
		className:    className,
		query:        query,
		originalData: originalData,
	}
	return destroy
}

// Execute 执行删除请求
func (d *Destroy) Execute() error {
	d.handleSession()
	d.runBeforeTrigger()
	d.handleUserRoles()
	d.runDestroy()
	d.runAfterTrigger()
	return nil
}

// handleSession 处理 _Session 表的删除操作
func (d *Destroy) handleSession() error {
	if d.className != "_Session" {
		return nil
	}
	sessionToken := d.originalData["sessionToken"]
	if sessionToken != nil {
		usersCache.remove(d.originalData["sessionToken"].(string))
	}

	return nil
}

// runBeforeTrigger 执行删前回调
func (d *Destroy) runBeforeTrigger() error {

	config.TConfig.LiveQuery.OnAfterDelete(d.className, d.originalData, nil)

	RunTrigger(TypeBeforeDelete, d.className, d.auth, nil, d.originalData)

	return nil
}

// handleUserRoles 获取用户角色信息
func (d *Destroy) handleUserRoles() error {
	if d.auth.IsMaster == false {
		d.auth.GetUserRoles()
	}

	return nil
}

// runDestroy 添加 acl 字段，并执行删除对象操作
func (d *Destroy) runDestroy() error {
	options := types.M{}
	if d.auth.IsMaster == false {
		acl := []string{"*"}
		if d.auth.User != nil {
			acl = append(acl, d.auth.User["objectId"].(string))
			acl = append(acl, d.auth.UserRoles...)
		}
		options["acl"] = acl
	}
	return orm.Destroy(d.className, d.query, options)
}

// runAfterTrigger 执行删后回调
func (d *Destroy) runAfterTrigger() error {
	RunTrigger(TypeAfterDelete, d.className, d.auth, nil, d.originalData)

	return nil
}
