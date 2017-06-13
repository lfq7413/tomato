package rest

import (
	"reflect"
	"testing"
	"time"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/mail"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

func TestPostgres_getUserIfNeeded(t *testing.T) {
	var schema types.M
	var object types.M
	var user types.M
	var result types.M
	var expect types.M
	/*********************************************************/
	user = nil
	result = getUserIfNeeded(user)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*********************************************************/
	user = types.M{
		"username": "joe",
		"email":    "abc@g.cn",
	}
	result = getUserIfNeeded(user)
	expect = types.M{
		"username": "joe",
		"email":    "abc@g.cn",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*********************************************************/
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"username": types.M{"type": "String"},
			"email":    types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"email":    "abc@g.cn",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	user = types.M{
		"username": "jack",
	}
	result = getUserIfNeeded(user)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/*********************************************************/
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"username": types.M{"type": "String"},
			"email":    types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"email":    "abc@g.cn",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	user = types.M{
		"email": "aaa@g.cn",
	}
	result = getUserIfNeeded(user)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/*********************************************************/
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"username": types.M{"type": "String"},
			"email":    types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"email":    "abc@g.cn",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	object = types.M{
		"objectId": "1002",
		"username": "jack",
		"email":    "abc@g.cn",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	user = types.M{
		"email": "abc@g.cn",
	}
	result = getUserIfNeeded(user)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/*********************************************************/
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"username": types.M{"type": "String"},
			"email":    types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"email":    "abc@g.cn",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	user = types.M{
		"email": "abc@g.cn",
	}
	result = getUserIfNeeded(user)
	expect = types.M{
		"objectId": "1001",
		"username": "joe",
		"email":    "abc@g.cn",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func TestPostgres_SendPasswordResetEmail(t *testing.T) {
	adapter = mail.NewSMTPAdapter()
	var schema types.M
	var object types.M
	var email string
	var result error
	var expect error
	/*********************************************************/
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"username": types.M{"type": "String"},
			"email":    types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"email":    "abc@g.cn",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	email = "aa@g.cn"
	result = SendPasswordResetEmail(email)
	expect = errs.E(errs.EmailMissing, "you must provide an email")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/*********************************************************/
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"username": types.M{"type": "String"},
			"email":    types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"email":    "abc@g.cn",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	email = "abc@g.cn"
	result = SendPasswordResetEmail(email)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func TestPostgres_setPasswordResetToken(t *testing.T) {
	var schema types.M
	var object types.M
	var email string
	var result types.M
	var expect types.M
	/*********************************************************/
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"username": types.M{"type": "String"},
			"email":    types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"email":    "abc@g.cn",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	email = "aa@g.cn"
	result = setPasswordResetToken(email)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/*********************************************************/
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"username": types.M{"type": "String"},
			"email":    types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
		"email":    "abc@g.cn",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	email = "abc@g.cn"
	result = setPasswordResetToken(email)
	expect = types.M{
		"objectId": "1001",
		"username": "joe",
		"email":    "abc@g.cn",
	}
	if _, ok := result["_perishable_token"]; ok == false {
		t.Error("expect:", "_perishable_token", "result:", "nil")
	}
	delete(result, "_perishable_token")
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func TestPostgres_VerifyEmail(t *testing.T) {
	var schema, object types.M
	var username, token string
	var result bool
	var expect bool
	/*********************************************************/
	config.TConfig = &config.Config{
		VerifyUserEmails:                 false,
		EmailVerifyTokenValidityDuration: -1,
	}
	username = "joe"
	token = "abc"
	result = VerifyEmail(username, token)
	expect = false
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*********************************************************/
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId":            types.M{"type": "String"},
			"username":            types.M{"type": "String"},
			"_email_verify_token": types.M{"type": "String"},
			"emailVerified":       types.M{"type": "Boolean"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":            "1001",
		"username":            "joe",
		"_email_verify_token": "abc1001",
		"emailVerified":       false,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	config.TConfig = &config.Config{
		VerifyUserEmails:                 true,
		EmailVerifyTokenValidityDuration: -1,
	}
	username = "jack"
	token = "abc"
	result = VerifyEmail(username, token)
	expect = false
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/*********************************************************/
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId":            types.M{"type": "String"},
			"username":            types.M{"type": "String"},
			"_email_verify_token": types.M{"type": "String"},
			"emailVerified":       types.M{"type": "Boolean"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":            "1001",
		"username":            "joe",
		"_email_verify_token": "abc1001",
		"emailVerified":       false,
	}
	orm.Adapter.CreateObject("_User", schema, object)
	config.TConfig = &config.Config{
		VerifyUserEmails:                 true,
		EmailVerifyTokenValidityDuration: -1,
	}
	username = "joe"
	token = "abc1001"
	result = VerifyEmail(username, token)
	expect = true
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/*********************************************************/
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId":                       types.M{"type": "String"},
			"username":                       types.M{"type": "String"},
			"_email_verify_token":            types.M{"type": "String"},
			"_email_verify_token_expires_at": types.M{"type": "Date"},
			"emailVerified":                  types.M{"type": "Boolean"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":                       "1001",
		"username":                       "joe",
		"_email_verify_token":            "abc1001",
		"emailVerified":                  false,
		"_email_verify_token_expires_at": types.M{"__type": "Date", "iso": utils.TimetoString(time.Now().UTC().Add(time.Second * 5))},
	}
	orm.Adapter.CreateObject("_User", schema, object)
	config.TConfig = &config.Config{
		VerifyUserEmails:                 true,
		EmailVerifyTokenValidityDuration: 5,
	}
	username = "joe"
	token = "abc1001"
	result = VerifyEmail(username, token)
	expect = true
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func TestPostgres_CheckResetTokenValidity(t *testing.T) {
	var schema, object types.M
	var username, token string
	var result types.M
	var expect types.M
	/*********************************************************/
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId":          types.M{"type": "String"},
			"username":          types.M{"type": "String"},
			"_perishable_token": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":          "1001",
		"username":          "joe",
		"_perishable_token": "abc1001",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	username = "jack"
	token = "abc"
	result = CheckResetTokenValidity(username, token)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/*********************************************************/
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId": types.M{"type": "String"},
			"username": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId": "1001",
		"username": "joe",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	username = "joe"
	token = "abc"
	result = CheckResetTokenValidity(username, token)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
	/*********************************************************/
	tmpTimeStr := utils.TimetoString(time.Now().UTC().Add(1 * time.Hour))
	initPostgresEnv()
	schema = types.M{
		"fields": types.M{
			"objectId":                     types.M{"type": "String"},
			"username":                     types.M{"type": "String"},
			"_perishable_token":            types.M{"type": "String"},
			"_perishable_token_expires_at": types.M{"type": "Date"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":                     "1001",
		"username":                     "joe",
		"_perishable_token":            "abc1001",
		"_perishable_token_expires_at": types.M{"__type": "Date", "iso": tmpTimeStr},
	}
	orm.Adapter.CreateObject("_User", schema, object)
	username = "joe"
	token = "abc1001"
	result = CheckResetTokenValidity(username, token)
	expect = types.M{
		"objectId":                     "1001",
		"username":                     "joe",
		"_perishable_token":            "abc1001",
		"_perishable_token_expires_at": types.M{"__type": "Date", "iso": tmpTimeStr},
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}
