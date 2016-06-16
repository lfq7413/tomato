package errs

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/lfq7413/tomato/types"
)

// E 组装 json 格式错误信息：
// {"code": 105,"error": "invalid field name: bl!ng"}
func E(code int, msg string) error {
	text := `{"code": ` + strconv.Itoa(code) + `,"error": "` + msg + `"}`
	return errors.New(text)
}

// ErrorToMap 把 error 转换为 types.M 格式，准备返回给客户端
func ErrorToMap(e error) types.M {
	var result types.M
	errMsg := e.Error()
	err := json.Unmarshal([]byte(errMsg), &result)
	if err != nil {
		result = types.M{"code": OtherCause, "error": e.Error()}
	}
	return result
}

// ErrorMessageToMap 把错误转换为 types.M 格式，准备返回给客户端
func ErrorMessageToMap(code int, msg string) types.M {
	return types.M{"code": code, "error": msg}
}

// GetErrorCode 获取 error 中的 code
func GetErrorCode(e error) int {
	var result types.M
	errMsg := e.Error()
	err := json.Unmarshal([]byte(errMsg), &result)
	if err != nil {
		return 0
	}
	if c, ok := result["code"].(float64); ok {
		return int(c)
	} else if c, ok := result["code"].(int); ok {
		return c
	}
	return 0
}

// GetErrorMessage 获取 error 中的 Message
func GetErrorMessage(e error) string {
	var result types.M
	errMsg := e.Error()
	err := json.Unmarshal([]byte(errMsg), &result)
	if err != nil {
		return errMsg
	}
	if msg, ok := result["error"].(string); ok {
		return msg
	}
	return errMsg
}

// OtherCause ...
// Error code indicating some error other than those enumerated here.
var OtherCause = -1

// InternalServerError ...
// Internal server error.
var InternalServerError = 1

// ServiceUnavailable ...
// The service is currently unavailable.
var ServiceUnavailable = 2

// ClientDisconnected ...
// Connection failure.
var ClientDisconnected = 4

// ConnectionFailed ...
// Error code indicating the connection to the Parse servers failed.
var ConnectionFailed = 100

// UserInvalidLoginParams ...
// Invalid login parameters.
var UserInvalidLoginParams = 101

// ObjectNotFound ...
// Error code indicating the specified object doesn't exist.
var ObjectNotFound = 101

// InvalidQuery ...
// Error code indicating you tried to query with a datatype that doesn't
// support it, like exact matching an array or object.
var InvalidQuery = 102

// InvalidClassName ...
// Error code indicating a missing or invalid classname. Classnames are
// case-sensitive. They must start with a letter, and a-zA-Z0-9 are the
// only valid characters.
var InvalidClassName = 103

// MissingObjectID ...
// Error code indicating an unspecified object id.
var MissingObjectID = 104

// InvalidKeyName ...
// Error code indicating an invalid key name. Keys are case-sensitive. They
// must start with a letter, and a-zA-Z0-9 are the only valid characters.
var InvalidKeyName = 105

// InvalidPointer ...
// Error code indicating a malformed pointer. You should not see this unless
// you have been mucking about changing internal Parse code.
var InvalidPointer = 106

// InvalidJSON ...
// Error code indicating that badly formed JSON was received upstream. This
// either indicates you have done something unusual with modifying how
// things encode to JSON, or the network is failing badly.
var InvalidJSON = 107

// CommandUnavailable ...
// Error code indicating that the feature you tried to access is only
// available internally for testing purposes.
var CommandUnavailable = 108

// NotInitialized ...
// You must call Parse.initialize before using the Parse library.
var NotInitialized = 109

// IncorrectType ...
// Error code indicating that a field was set to an inconsistent type.
var IncorrectType = 111

// InvalidChannelName ...
// Error code indicating an invalid channel name. A channel name is either
// an empty string (the broadcast channel) or contains only a-zA-Z0-9
// characters and starts with a letter.
var InvalidChannelName = 112

// InvalidSubscriptionType ...
// Bad subscription type.
var InvalidSubscriptionType = 113

// InvalidDeviceToken ...
// The provided device token is invalid.
var InvalidDeviceToken = 114

// PushMisconfigured ...
// Error code indicating that push is misconfigured.
var PushMisconfigured = 115

// PushWhereAndChannels ...
// Can't set channels for a query-targeted push.
var PushWhereAndChannels = 115

// PushWhereAndType ...
// Can't set device type for a query-targeted push.
var PushWhereAndType = 115

// PushMissingData ...
// Push is missing a 'data' field.
var PushMissingData = 115

// PushMissingChannels ...
// Non-query push is missing a 'channels' field.
var PushMissingChannels = 115

// ClientPushDisabled ...
// Client-initiated push is not enabled.
var ClientPushDisabled = 115

// RestPushDisabled ...
// REST-initiated push is not enabled.
var RestPushDisabled = 115

// ClientPushWithURI ...
// Client-initiated push cannot use the "uri" option.
var ClientPushWithURI = 115

// PushQueryOrPayloadTooLarge ...
// Your push query or data payload is too large.
var PushQueryOrPayloadTooLarge = 115

// ObjectTooLarge ...
// Error code indicating that the object is too large.
var ObjectTooLarge = 116

// ExceededConfigParamsError ...
// You have reached the limit of 100 config parameters.
var ExceededConfigParamsError = 116

// InvalidLimitError ...
// An invalid value was set for the limit.
var InvalidLimitError = 117

// InvalidSkipError ...
// An invalid value was set for skip.
var InvalidSkipError = 118

// OperationForbidden ...
// Error code indicating that the operation isn't allowed for clients.
var OperationForbidden = 119

// CacheMiss ...
// Error code indicating the result was not found in the cache.
var CacheMiss = 120

// InvalidNestedKey ...
// Error code indicating that an invalid key was used in a nested
// JSONObject.
var InvalidNestedKey = 121

// InvalidFileName ...
// Error code indicating that an invalid filename was used for ParseFile.
// A valid file name contains only a-zA-Z0-9. characters and is between 1
// and 128 characters.
var InvalidFileName = 122

// InvalidAcl ...
// Error code indicating an invalid ACL was provided.
var InvalidAcl = 123

// Timeout ...
// Error code indicating that the request timed out on the server. Typically
// this indicates that the request is too expensive to run.
var Timeout = 124

// InvalidEmailAddress ...
// Error code indicating that the email address was invalid.
var InvalidEmailAddress = 125

// MissingContentType ...
// Error code indicating a missing content type.
var MissingContentType = 126

// MissingContentLength ...
// Error code indicating a missing content length.
var MissingContentLength = 127

// InvalidContentLength ...
// Error code indicating an invalid content length.
var InvalidContentLength = 128

// FileTooLarge ...
// Error code indicating a file that was too large.
var FileTooLarge = 129

// FileSaveError ...
// Error code indicating an error saving a file.
var FileSaveError = 130

// FileDeleteError ...
// File could not be deleted.
var FileDeleteError = 131

// InvalidInstallationIDError ...
// Invalid installation id.
var InvalidInstallationIDError = 132

// InvalidDeviceTypeError ...
// Invalid device type.
var InvalidDeviceTypeError = 133

// InvalidChannelsArrayError ...
// Invalid channels array value.
var InvalidChannelsArrayError = 134

// MissingRequiredFieldError ...
// Required field is missing.
var MissingRequiredFieldError = 135

// ChangedImmutableFieldError ...
// An immutable field was changed.
var ChangedImmutableFieldError = 136

// DuplicateValue ...
// Error code indicating that a unique field was given a value that is
// already taken.
var DuplicateValue = 137

// InvalidExpirationError ...
// Invalid expiration value.
var InvalidExpirationError = 138

// InvalidRoleName ...
// Error code indicating that a role's name is invalid.
var InvalidRoleName = 139

// ReservedValue ...
// Field value is reserved.
var ReservedValue = 139

// ExceededQuota ...
// Error code indicating that an application quota was exceeded.  Upgrade to
// resolve.
var ExceededQuota = 140

// ScriptFailed ...
// Error code indicating that a Cloud Code script failed.
var ScriptFailed = 141

// FunctionNotFound ...
// Cloud function not found.
var FunctionNotFound = 141

// JobNotFound ...
// Background job not found.
var JobNotFound = 141

// SuccessErrorNotCalled ...
// success/error was not called.
var SuccessErrorNotCalled = 141

// MultupleSuccessErrorCalls ...
// Can't call success/error multiple times.
var MultupleSuccessErrorCalls = 141

// ValidationError ...
// Error code indicating that a Cloud Code validation failed.
var ValidationError = 142

// WebhookError ...
// Webhook error.
var WebhookError = 143

// ReceiptMissing ...
// Product purchase receipt is missing.
var ReceiptMissing = 143

// InvalidPurchaseReceipt ...
// Product purchase receipt is invalid.
var InvalidPurchaseReceipt = 144

// PaymentDisabled ...
// Payment is disabled on this device.
var PaymentDisabled = 145

// InvalidProductIdentifier ...
// The product identifier is invalid.
var InvalidProductIdentifier = 146

// ProductNotFoundInAppStore ...
// The product is not found in the App Store.
var ProductNotFoundInAppStore = 147

// InvalidServerResponse ...
// The Apple server response is not valid.
var InvalidServerResponse = 148

// ProductDownloadFilesystemError ...
// The product fails to download due to file system error.
var ProductDownloadFilesystemError = 149

// InvalidImageData ...
// Error code indicating that invalid image data was provided.
var InvalidImageData = 150

// UnsavedFileError ...
// Error code indicating an unsaved file.
var UnsavedFileError = 151

// InvalidPushTimeError ...
// Error code indicating an invalid push time.
var InvalidPushTimeError = 152

// InefficientQueryError ...
// An inefficient query was rejected by the server.
var InefficientQueryError = 154

// RequestLimitExceeded ...
// Error code indicating that the application has exceeded its request
// limit.
var RequestLimitExceeded = 155

// MissingPushIDError ...
// A push id is missing. Deprecated.
var MissingPushIDError = 156

// MissingDeviceTypeError ...
// The device type field is missing. Deprecated.
var MissingDeviceTypeError = 157

// HostingError ...
// Hosting error.
var HostingError = 158

// TemporaryRejectionError ...
// An application's requests are temporary rejected by the server.
var TemporaryRejectionError = 159

// InvalidEventName ...
// Error code indicating an invalid event name.
var InvalidEventName = 160

// UsernameMissing ...
// Error code indicating that the username is missing or empty.
var UsernameMissing = 200

// PasswordMissing ...
// Error code indicating that the password is missing or empty.
var PasswordMissing = 201

// UsernameTaken ...
// Error code indicating that the username has already been taken.
var UsernameTaken = 202

// EmailTaken ...
// Error code indicating that the email has already been taken.
var EmailTaken = 203

// EmailMissing ...
// Error code indicating that the email is missing, but must be specified.
var EmailMissing = 204

// EmailNotFound ...
// Error code indicating that a user with the specified email was not found.
var EmailNotFound = 205

// SessionMissing ...
// Error code indicating that a user object without a valid session could
// not be altered.
var SessionMissing = 206

// MustCreateUserThroughSignup ...
// Error code indicating that a user can only be created through signup.
var MustCreateUserThroughSignup = 207

// AccountAlreadyLinked ...
// Error code indicating that an an account being linked is already linked
// to another user.
var AccountAlreadyLinked = 208

// InvalidSessionToken ...
// Error code indicating that the current session token is invalid.
var InvalidSessionToken = 209

// LinkedIDMissing ...
// Error code indicating that a user cannot be linked to an account because
// that account's id could not be found.
var LinkedIDMissing = 250

// InvalidLinkedSession ...
// Error code indicating that a user with a linked (e.g. Facebook) account
// has an invalid session.
var InvalidLinkedSession = 251

// InvalidGeneralAuthData ...
// Invalid auth data value used.
var InvalidGeneralAuthData = 251

// BadAnonymousID ...
// Anonymous id is not a valid lowercase UUID.
var BadAnonymousID = 251

// FacebookBadToken ...
// The supplied Facebook session token is expired or invalid.
var FacebookBadToken = 251

// FacebookBadID ...
// A user with a linked Facebook account has an invalid session.
var FacebookBadID = 251

// FacebookWrongAppID ...
// Unacceptable Facebook application id.
var FacebookWrongAppID = 251

// TwitterVerificationFailed ...
//Twitter credential verification failed.
var TwitterVerificationFailed = 251

// TwitterWrongID ...
// Submitted Twitter id does not match the id associated with the submitted access token.
var TwitterWrongID = 251

// TwitterWrongScreenName ...
// Submitted Twitter handle does not match the handle associated with the submitted access token.
var TwitterWrongScreenName = 251

// TwitterConnectFailure ...
// Twitter credentials could not be verified due to problems accessing the Twitter API.
var TwitterConnectFailure = 251

// UnsupportedService ...
// Error code indicating that a service being linked (e.g. Facebook or
// Twitter) is unsupported.
var UnsupportedService = 252

// UsernameSigninDisabled ...
// Authentication by username and password is not supported for this application.
var UsernameSigninDisabled = 252

// AnonymousSigninDisabled ...
// Anonymous users are not supported for this application.
var AnonymousSigninDisabled = 252

// FacebookSigninDisabled ...
// Authentication by Facebook is not supported for this application.
var FacebookSigninDisabled = 252

// TwitterSigninDisabled ...
// Authentication by Twitter is not supported for this application.
var TwitterSigninDisabled = 252

// InvalidAuthDataError ...
// An invalid authData value was passed.
var InvalidAuthDataError = 253

// ClassNotEmpty ...
// Class is not empty and cannot be dropped.
var ClassNotEmpty = 255

// AppNameInvalid ...
// App name is invalid.
var AppNameInvalid = 256

// AggregateError ...
// Error code indicating that there were multiple errors. Aggregate errors
// have an "errors" property, which is an array of error objects with more
// detail about each error that occurred.
var AggregateError = 600

// FileReadError ...
// Error code indicating the client was unable to read an input file.
var FileReadError = 601

// XDomainRequest ...
// Error code indicating a real error code is unavailable because
// we had to use an XDomainRequest object to allow CORS requests in
// Internet Explorer, which strips the body from HTTP responses that have
// a non-2XX status code.
var XDomainRequest = 602

// MissingAPIKeyError ...
// The request is missing an API key.
var MissingAPIKeyError = 902

// InvalidAPIKeyError ...
// The request is using an invalid API key.
var InvalidAPIKeyError = 903

// LinkingNotSupportedError ...
// Linking to an external account not supported yet with signup_or_login.
var LinkingNotSupportedError = 999
