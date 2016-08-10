package rest

import (
	"reflect"
	"testing"

	"github.com/lfq7413/tomato/types"
)

func Test_SetEmailVerifyToken(t *testing.T) {
	// TODO
}

func Test_SendVerificationEmail(t *testing.T) {
	// getUserIfNeeded
	// TODO
}

func Test_getUserIfNeeded(t *testing.T) {
	// TODO
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
	// setPasswordResetToken
	// TODO
}

func Test_setPasswordResetToken(t *testing.T) {
	// TODO
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
	// TODO
}

func Test_CheckResetTokenValidity(t *testing.T) {
	// TODO
}

func Test_UpdatePassword(t *testing.T) {
	// CheckResetTokenValidity
	// updateUserPassword
	// TODO
}

func Test_updateUserPassword(t *testing.T) {
	// TODO
}
