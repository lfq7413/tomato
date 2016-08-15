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

func Test_SetEmailVerifyToken(t *testing.T) {
	var user types.M
	var expect types.M
	/*********************************************************/
	user = nil
	SetEmailVerifyToken(user)
	expect = nil
	if reflect.DeepEqual(expect, user) == false {
		t.Error("expect:", expect, "result:", user)
	}
	/*********************************************************/
	user = types.M{
		"username": "joe",
	}
	config.TConfig = &config.Config{
		VerifyUserEmails:                 false,
		EmailVerifyTokenValidityDuration: -1,
	}
	SetEmailVerifyToken(user)
	expect = types.M{
		"username": "joe",
	}
	if reflect.DeepEqual(expect, user) == false {
		t.Error("expect:", expect, "result:", user)
	}
	/*********************************************************/
	user = types.M{
		"username": "joe",
	}
	config.TConfig = &config.Config{
		VerifyUserEmails:                 true,
		EmailVerifyTokenValidityDuration: -1,
	}
	SetEmailVerifyToken(user)
	expect = types.M{
		"username":      "joe",
		"emailVerified": false,
	}
	if len(utils.S(user["_email_verify_token"])) != 32 {
		t.Error("expect:", 32, "result:", len(utils.S(user["_email_verify_token"])))
	}
	delete(user, "_email_verify_token")
	if reflect.DeepEqual(expect, user) == false {
		t.Error("expect:", expect, "result:", user)
	}
	/*********************************************************/
	user = types.M{
		"username": "joe",
	}
	config.TConfig = &config.Config{
		VerifyUserEmails:                 true,
		EmailVerifyTokenValidityDuration: 60,
	}
	SetEmailVerifyToken(user)
	expect = types.M{
		"username":      "joe",
		"emailVerified": false,
	}
	if len(utils.S(user["_email_verify_token"])) != 32 {
		t.Error("expect:", 32, "result:", len(utils.S(user["_email_verify_token"])))
	}
	delete(user, "_email_verify_token")
	if utils.S(user["_email_verify_token_expires_at"]) == "" {
		t.Error("expect:", "time", "result:", user["_email_verify_token_expires_at"])
	}
	delete(user, "_email_verify_token_expires_at")
	if reflect.DeepEqual(expect, user) == false {
		t.Error("expect:", expect, "result:", user)
	}
}

func Test_SendVerificationEmail(t *testing.T) {
	var user types.M
	config.TConfig = &config.Config{
		VerifyUserEmails: true,
		ServerURL:        "http://www.g.cn/",
	}
	adapter = mail.NewSMTPAdapter()
	user = types.M{
		"_email_verify_token": "abc",
		"username":            "joe",
		"mail":                "abc@g.cn",
	}
	SendVerificationEmail(user)
}

func Test_getUserIfNeeded(t *testing.T) {
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
	initEnv()
	schema = types.M{
		"fields": types.M{
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
	initEnv()
	schema = types.M{
		"fields": types.M{
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
	initEnv()
	schema = types.M{
		"fields": types.M{
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
	initEnv()
	schema = types.M{
		"fields": types.M{
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

func Test_defaultVerificationEmail(t *testing.T) {
	var options types.M
	var result types.M
	var expect types.M
	var text string
	/*********************************************************/
	options = nil
	result = defaultVerificationEmail(options)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*********************************************************/
	options = types.M{}
	result = defaultVerificationEmail(options)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*********************************************************/
	options = types.M{
		"user": types.M{
			"email": "123@g.com",
		},
		"appName": "tomato",
		"link":    "http://www.g.com",
	}
	result = defaultVerificationEmail(options)
	text = "Hi,\n\n"
	text += "You are being asked to confirm the e-mail address 123@g.com"
	text += " with tomato\n\n"
	text += "Click here to confirm it:\nhttp://www.g.com"
	expect = types.M{
		"text":    text,
		"to":      "123@g.com",
		"subject": "Please verify your e-mail for tomato",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_SendPasswordResetEmail(t *testing.T) {
	adapter = mail.NewSMTPAdapter()
	var schema types.M
	var object types.M
	var email string
	var result error
	var expect error
	/*********************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
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
	initEnv()
	schema = types.M{
		"fields": types.M{
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

func Test_setPasswordResetToken(t *testing.T) {
	var schema types.M
	var object types.M
	var email string
	var result types.M
	var expect types.M
	/*********************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
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
	initEnv()
	schema = types.M{
		"fields": types.M{
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
		"_id":      "1001",
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

func Test_defaultResetPasswordEmail(t *testing.T) {
	var options types.M
	var result types.M
	var expect types.M
	var text string
	/*********************************************************/
	options = nil
	result = defaultResetPasswordEmail(options)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*********************************************************/
	options = types.M{}
	result = defaultResetPasswordEmail(options)
	expect = nil
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	/*********************************************************/
	options = types.M{
		"user": types.M{
			"email": "123@g.com",
		},
		"appName": "tomato",
		"link":    "http://www.g.com",
	}
	result = defaultResetPasswordEmail(options)
	text = "Hi,\n\n"
	text += "You requested to reset your password for tomato\n\n"
	text += "Click here to reset it:\nhttp://www.g.com"
	expect = types.M{
		"text":    text,
		"to":      "123@g.com",
		"subject": "Password Reset for tomato",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
}

func Test_VerifyEmail(t *testing.T) {
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
	initEnv()
	schema = types.M{
		"fields": types.M{
			"username":      types.M{"type": "String"},
			"emailVerified": types.M{"type": "bool"},
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
	initEnv()
	schema = types.M{
		"fields": types.M{
			"username":      types.M{"type": "String"},
			"emailVerified": types.M{"type": "bool"},
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
	initEnv()
	schema = types.M{
		"fields": types.M{
			"username":      types.M{"type": "String"},
			"emailVerified": types.M{"type": "bool"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":                       "1001",
		"username":                       "joe",
		"_email_verify_token":            "abc1001",
		"emailVerified":                  false,
		"_email_verify_token_expires_at": utils.TimetoString(time.Now().UTC().Add(time.Second * 5)),
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

func Test_CheckResetTokenValidity(t *testing.T) {
	var schema, object types.M
	var username, token string
	var result types.M
	var expect types.M
	/*********************************************************/
	initEnv()
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
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
	initEnv()
	schema = types.M{
		"fields": types.M{
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
	initEnv()
	schema = types.M{
		"fields": types.M{
			"username": types.M{"type": "String"},
		},
	}
	orm.Adapter.CreateClass("_User", schema)
	object = types.M{
		"objectId":          "1001",
		"username":          "joe",
		"_perishable_token": "abc1001",
	}
	orm.Adapter.CreateObject("_User", schema, object)
	username = "joe"
	token = "abc1001"
	result = CheckResetTokenValidity(username, token)
	expect = types.M{
		"objectId": "1001",
		"username": "joe",
	}
	if reflect.DeepEqual(expect, result) == false {
		t.Error("expect:", expect, "result:", result)
	}
	orm.TomatoDBController.DeleteEverything()
}

func Test_UpdatePassword(t *testing.T) {
	// updateUserPassword
	// TODO
}

func Test_updateUserPassword(t *testing.T) {
	// TODO
}
