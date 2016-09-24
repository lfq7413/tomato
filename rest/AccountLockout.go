package rest

import (
	"strconv"
	"time"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// AccountLockout 密码错误达到一定次数，锁定账户
type AccountLockout struct {
	username string
}

// NewAccountLockout ...
func NewAccountLockout(username string) *AccountLockout {
	return &AccountLockout{
		username: username,
	}
}

// HandleLoginAttempt 处理登录结果
func (a *AccountLockout) HandleLoginAttempt(loginSuccessful bool) error {
	if config.TConfig.EnableAccountLockout == false {
		return nil
	}
	err := a.notLocked()
	if err != nil {
		return err
	}
	if loginSuccessful {
		return a.setFailedLoginCount(0)
	}
	return a.handleFailedLoginAttempt()
}

// notLocked 检测账户是否已经被锁住
func (a *AccountLockout) notLocked() error {
	query := types.M{
		"username": a.username,
		"_account_lockout_expires_at": types.M{
			"$gt": types.M{
				"__type": "Date",
				"iso":    utils.TimetoString(time.Now().UTC()),
			},
		},
		"_failed_login_count": types.M{
			"$gte": config.TConfig.AccountLockoutThreshold,
		},
	}

	result, err := orm.TomatoDBController.Find("_User", query, types.M{})
	if err != nil {
		return err
	}
	if len(result) > 0 {
		msg := "Your account is locked due to multiple failed login attempts. Please try again after " +
			strconv.Itoa(config.TConfig.AccountLockoutDuration) + " minute(s)"
		return errs.E(errs.ObjectNotFound, msg)
	}
	return nil
}

// setFailedLoginCount 设置 _failed_login_count
func (a *AccountLockout) setFailedLoginCount(count int) error {
	query := types.M{
		"username": a.username,
	}
	updateFields := types.M{
		"_failed_login_count": count,
	}
	_, err := orm.TomatoDBController.Update("_User", query, updateFields, types.M{}, false)
	return err
}

// handleFailedLoginAttempt 处理失败的登录
func (a *AccountLockout) handleFailedLoginAttempt() error {
	err := a.initFailedLoginCount()
	if err != nil {
		return err
	}
	err = a.incrementFailedLoginCount()
	if err != nil {
		return err
	}
	err = a.setLockoutExpiration()
	return err
}

// initFailedLoginCount 如果 _failed_login_count 字段没有设置，则将其设置为 0
func (a *AccountLockout) initFailedLoginCount() error {
	failedLoginCountIsSet, err := a.isFailedLoginCountSet()
	if err != nil {
		return err
	}
	if failedLoginCountIsSet == false {
		return a.setFailedLoginCount(0)
	}
	return nil
}

// incrementFailedLoginCount _failed_login_count 字段增加 1
func (a *AccountLockout) incrementFailedLoginCount() error {
	query := types.M{
		"username": a.username,
	}
	updateFields := types.M{
		"_failed_login_count": types.M{
			"__op":   "Increment",
			"amount": 1,
		},
	}
	_, err := orm.TomatoDBController.Update("_User", query, updateFields, types.M{}, false)
	return err
}

// setLockoutExpiration 密码错误次数超限后，设置下次重试的时间
func (a *AccountLockout) setLockoutExpiration() error {
	query := types.M{
		"username":            a.username,
		"_failed_login_count": types.M{"$gte": config.TConfig.AccountLockoutThreshold},
	}
	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(config.TConfig.AccountLockoutDuration) * time.Minute)
	updateFields := types.M{
		"_account_lockout_expires_at": types.M{
			"__type": "Date",
			"iso":    utils.TimetoString(expiresAt),
		},
	}

	_, err := orm.TomatoDBController.Update("_User", query, updateFields, types.M{}, false)
	if err != nil {
		if errs.GetErrorCode(err) == errs.ObjectNotFound &&
			errs.GetErrorMessage(err) == "Object not found." {
			return nil
		}
		return err
	}

	return nil
}

// isFailedLoginCountSet 检测 _failed_login_count 字段是否存在
func (a *AccountLockout) isFailedLoginCountSet() (bool, error) {
	query := types.M{
		"username":            a.username,
		"_failed_login_count": types.M{"$exists": true},
	}
	result, err := orm.TomatoDBController.Find("_User", query, types.M{})
	if err != nil {
		return false, err
	}
	if len(result) > 0 {
		return true, nil
	}
	return false, nil
}
