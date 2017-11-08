package errs

import (
	"strconv"

	"github.com/lfq7413/tomato/types"
)

// TomatoError ...
type TomatoError struct {
	Code    int
	Message string
}

func (e *TomatoError) Error() string {
	return `{"code": ` + strconv.Itoa(e.Code) + `,"error": "` + e.Message + `"}`
}

// E 组装 json 格式错误信息：
// {"code": 105,"error": "invalid field name: bl!ng"}
func E(code int, msg string) error {
	return &TomatoError{
		Code:    code,
		Message: msg,
	}
}

// ErrorToMap 把 error 转换为 types.M 格式，准备返回给客户端
func ErrorToMap(e error) types.M {
	if v, ok := e.(*TomatoError); ok {
		return types.M{
			"code":  v.Code,
			"error": v.Message,
		}
	}
	return types.M{
		"code":  OtherCause,
		"error": e.Error(),
	}
}

// ErrorMessageToMap 把错误转换为 types.M 格式，准备返回给客户端
func ErrorMessageToMap(code int, msg string) types.M {
	return types.M{"code": code, "error": msg}
}

// GetErrorCode 获取 error 中的 code
func GetErrorCode(e error) int {
	if v, ok := e.(*TomatoError); ok {
		return v.Code
	}
	return 0
}

// GetErrorMessage 获取 error 中的 Message
func GetErrorMessage(e error) string {
	if v, ok := e.(*TomatoError); ok {
		return v.Message
	}
	return e.Error()
}

// OtherCause ...
// Error code indicating some error other than those enumerated here.
const OtherCause = -1

// InternalServerError ...
// Internal server error.
const InternalServerError = 1

// ServiceUnavailable ...
// The service is currently unavailable.
const ServiceUnavailable = 2

// ClientDisconnected ...
// Connection failure.
const ClientDisconnected = 4

// ConnectionFailed ...
// Error code indicating the connection to the Parse servers failed.
const ConnectionFailed = 100

// UserInvalidLoginParams ...
// Invalid login parameters.
const UserInvalidLoginParams = 101

// ObjectNotFound ...
// Error code indicating the specified object doesn't exist.
const ObjectNotFound = 101

// InvalidQuery ...
// Error code indicating you tried to query with a datatype that doesn't
// support it, like exact matching an array or object.
const InvalidQuery = 102

// InvalidClassName ...
// Error code indicating a missing or invalid classname. Classnames are
// case-sensitive. They must start with a letter, and a-zA-Z0-9 are the
// only valid characters.
const InvalidClassName = 103

// MissingObjectID ...
// Error code indicating an unspecified object id.
const MissingObjectID = 104

// InvalidKeyName ...
// Error code indicating an invalid key name. Keys are case-sensitive. They
// must start with a letter, and a-zA-Z0-9 are the only valid characters.
const InvalidKeyName = 105

// InvalidPointer ...
// Error code indicating a malformed pointer. You should not see this unless
// you have been mucking about changing internal Parse code.
const InvalidPointer = 106

// InvalidJSON ...
// Error code indicating that badly formed JSON was received upstream. This
// either indicates you have done something unusual with modifying how
// things encode to JSON, or the network is failing badly.
const InvalidJSON = 107

// CommandUnavailable ...
// Error code indicating that the feature you tried to access is only
// available internally for testing purposes.
const CommandUnavailable = 108

// NotInitialized ...
// You must call Parse.initialize before using the Parse library.
const NotInitialized = 109

// IncorrectType ...
// Error code indicating that a field was set to an inconsistent type.
const IncorrectType = 111

// InvalidChannelName ...
// Error code indicating an invalid channel name. A channel name is either
// an empty string (the broadcast channel) or contains only a-zA-Z0-9
// characters and starts with a letter.
const InvalidChannelName = 112

// InvalidSubscriptionType ...
// Bad subscription type.
const InvalidSubscriptionType = 113

// InvalidDeviceToken ...
// The provided device token is invalid.
const InvalidDeviceToken = 114

// PushMisconfigured ...
// Error code indicating that push is misconfigured.
const PushMisconfigured = 115

// PushWhereAndChannels ...
// Can't set channels for a query-targeted push.
const PushWhereAndChannels = 115

// PushWhereAndType ...
// Can't set device type for a query-targeted push.
const PushWhereAndType = 115

// PushMissingData ...
// Push is missing a 'data' field.
const PushMissingData = 115

// PushMissingChannels ...
// Non-query push is missing a 'channels' field.
const PushMissingChannels = 115

// ClientPushDisabled ...
// Client-initiated push is not enabled.
const ClientPushDisabled = 115

// RestPushDisabled ...
// REST-initiated push is not enabled.
const RestPushDisabled = 115

// ClientPushWithURI ...
// Client-initiated push cannot use the "uri" option.
const ClientPushWithURI = 115

// PushQueryOrPayloadTooLarge ...
// Your push query or data payload is too large.
const PushQueryOrPayloadTooLarge = 115

// ObjectTooLarge ...
// Error code indicating that the object is too large.
const ObjectTooLarge = 116

// ExceededConfigParamsError ...
// You have reached the limit of 100 config parameters.
const ExceededConfigParamsError = 116

// InvalidLimitError ...
// An invalid value was set for the limit.
const InvalidLimitError = 117

// InvalidSkipError ...
// An invalid value was set for skip.
const InvalidSkipError = 118

// OperationForbidden ...
// Error code indicating that the operation isn't allowed for clients.
const OperationForbidden = 119

// CacheMiss ...
// Error code indicating the result was not found in the cache.
const CacheMiss = 120

// InvalidNestedKey ...
// Error code indicating that an invalid key was used in a nested
// JSONObject.
const InvalidNestedKey = 121

// InvalidFileName ...
// Error code indicating that an invalid filename was used for ParseFile.
// A valid file name contains only a-zA-Z0-9. characters and is between 1
// and 128 characters.
const InvalidFileName = 122

// InvalidACL ...
// Error code indicating an invalid ACL was provided.
const InvalidACL = 123

// Timeout ...
// Error code indicating that the request timed out on the server. Typically
// this indicates that the request is too expensive to run.
const Timeout = 124

// InvalidEmailAddress ...
// Error code indicating that the email address was invalid.
const InvalidEmailAddress = 125

// MissingContentType ...
// Error code indicating a missing content type.
const MissingContentType = 126

// MissingContentLength ...
// Error code indicating a missing content length.
const MissingContentLength = 127

// InvalidContentLength ...
// Error code indicating an invalid content length.
const InvalidContentLength = 128

// FileTooLarge ...
// Error code indicating a file that was too large.
const FileTooLarge = 129

// FileSaveError ...
// Error code indicating an error saving a file.
const FileSaveError = 130

// InvalidInstallationIDError ...
// Invalid installation id.
const InvalidInstallationIDError = 132

// InvalidDeviceTypeError ...
// Invalid device type.
const InvalidDeviceTypeError = 133

// InvalidChannelsArrayError ...
// Invalid channels array value.
const InvalidChannelsArrayError = 134

// MissingRequiredFieldError ...
// Required field is missing.
const MissingRequiredFieldError = 135

// ChangedImmutableFieldError ...
// An immutable field was changed.
const ChangedImmutableFieldError = 136

// DuplicateValue ...
// Error code indicating that a unique field was given a value that is
// already taken.
const DuplicateValue = 137

// InvalidExpirationError ...
// Invalid expiration value.
const InvalidExpirationError = 138

// InvalidRoleName ...
// Error code indicating that a role's name is invalid.
const InvalidRoleName = 139

// ReservedValue ...
// Field value is reserved.
const ReservedValue = 139

// ExceededQuota ...
// Error code indicating that an application quota was exceeded.  Upgrade to
// resolve.
const ExceededQuota = 140

// ScriptFailed ...
// Error code indicating that a Cloud Code script failed.
const ScriptFailed = 141

// FunctionNotFound ...
// Cloud function not found.
const FunctionNotFound = 141

// JobNotFound ...
// Background job not found.
const JobNotFound = 141

// SuccessErrorNotCalled ...
// success/error was not called.
const SuccessErrorNotCalled = 141

// MultupleSuccessErrorCalls ...
// Can't call success/error multiple times.
const MultupleSuccessErrorCalls = 141

// ValidationError ...
// Error code indicating that a Cloud Code validation failed.
const ValidationError = 142

// WebhookError ...
// Webhook error.
const WebhookError = 143

// ReceiptMissing ...
// Product purchase receipt is missing.
const ReceiptMissing = 143

// InvalidPurchaseReceipt ...
// Product purchase receipt is invalid.
const InvalidPurchaseReceipt = 144

// PaymentDisabled ...
// Payment is disabled on this device.
const PaymentDisabled = 145

// InvalidProductIdentifier ...
// The product identifier is invalid.
const InvalidProductIdentifier = 146

// ProductNotFoundInAppStore ...
// The product is not found in the App Store.
const ProductNotFoundInAppStore = 147

// InvalidServerResponse ...
// The Apple server response is not valid.
const InvalidServerResponse = 148

// ProductDownloadFilesystemError ...
// The product fails to download due to file system error.
const ProductDownloadFilesystemError = 149

// InvalidImageData ...
// Error code indicating that invalid image data was provided.
const InvalidImageData = 150

// UnsavedFileError ...
// Error code indicating an unsaved file.
const UnsavedFileError = 151

// InvalidPushTimeError ...
// Error code indicating an invalid push time.
const InvalidPushTimeError = 152

// FileDeleteError ...
// File could not be deleted.
const FileDeleteError = 153

// InefficientQueryError ...
// An inefficient query was rejected by the server.
const InefficientQueryError = 154

// RequestLimitExceeded ...
// Error code indicating that the application has exceeded its request
// limit.
const RequestLimitExceeded = 155

// MissingPushIDError ...
// A push id is missing. Deprecated.
const MissingPushIDError = 156

// MissingDeviceTypeError ...
// The device type field is missing. Deprecated.
const MissingDeviceTypeError = 157

// HostingError ...
// Hosting error.
const HostingError = 158

// TemporaryRejectionError ...
// An application's requests are temporary rejected by the server.
const TemporaryRejectionError = 159

// InvalidEventName ...
// Error code indicating an invalid event name.
const InvalidEventName = 160

// UsernameMissing ...
// Error code indicating that the username is missing or empty.
const UsernameMissing = 200

// PasswordMissing ...
// Error code indicating that the password is missing or empty.
const PasswordMissing = 201

// UsernameTaken ...
// Error code indicating that the username has already been taken.
const UsernameTaken = 202

// EmailTaken ...
// Error code indicating that the email has already been taken.
const EmailTaken = 203

// EmailMissing ...
// Error code indicating that the email is missing, but must be specified.
const EmailMissing = 204

// EmailNotFound ...
// Error code indicating that a user with the specified email was not found.
const EmailNotFound = 205

// SessionMissing ...
// Error code indicating that a user object without a valid session could
// not be altered.
const SessionMissing = 206

// MustCreateUserThroughSignup ...
// Error code indicating that a user can only be created through signup.
const MustCreateUserThroughSignup = 207

// AccountAlreadyLinked ...
// Error code indicating that an an account being linked is already linked
// to another user.
const AccountAlreadyLinked = 208

// InvalidSessionToken ...
// Error code indicating that the current session token is invalid.
const InvalidSessionToken = 209

// LinkedIDMissing ...
// Error code indicating that a user cannot be linked to an account because
// that account's id could not be found.
const LinkedIDMissing = 250

// InvalidLinkedSession ...
// Error code indicating that a user with a linked (e.g. Facebook) account
// has an invalid session.
const InvalidLinkedSession = 251

// InvalidGeneralAuthData ...
// Invalid auth data value used.
const InvalidGeneralAuthData = 251

// BadAnonymousID ...
// Anonymous id is not a valid lowercase UUID.
const BadAnonymousID = 251

// FacebookBadToken ...
// The supplied Facebook session token is expired or invalid.
const FacebookBadToken = 251

// FacebookBadID ...
// A user with a linked Facebook account has an invalid session.
const FacebookBadID = 251

// FacebookWrongAppID ...
// Unacceptable Facebook application id.
const FacebookWrongAppID = 251

// TwitterVerificationFailed ...
//Twitter credential verification failed.
const TwitterVerificationFailed = 251

// TwitterWrongID ...
// Submitted Twitter id does not match the id associated with the submitted access token.
const TwitterWrongID = 251

// TwitterWrongScreenName ...
// Submitted Twitter handle does not match the handle associated with the submitted access token.
const TwitterWrongScreenName = 251

// TwitterConnectFailure ...
// Twitter credentials could not be verified due to problems accessing the Twitter API.
const TwitterConnectFailure = 251

// UnsupportedService ...
// Error code indicating that a service being linked (e.g. Facebook or
// Twitter) is unsupported.
const UnsupportedService = 252

// UsernameSigninDisabled ...
// Authentication by username and password is not supported for this application.
const UsernameSigninDisabled = 252

// AnonymousSigninDisabled ...
// Anonymous users are not supported for this application.
const AnonymousSigninDisabled = 252

// FacebookSigninDisabled ...
// Authentication by Facebook is not supported for this application.
const FacebookSigninDisabled = 252

// TwitterSigninDisabled ...
// Authentication by Twitter is not supported for this application.
const TwitterSigninDisabled = 252

// InvalidAuthDataError ...
// An invalid authData value was passed.
const InvalidAuthDataError = 253

// ClassNotEmpty ...
// Class is not empty and cannot be dropped.
const ClassNotEmpty = 255

// AppNameInvalid ...
// App name is invalid.
const AppNameInvalid = 256

// AggregateError ...
// Error code indicating that there were multiple errors. Aggregate errors
// have an "errors" property, which is an array of error objects with more
// detail about each error that occurred.
const AggregateError = 600

// FileReadError ...
// Error code indicating the client was unable to read an input file.
const FileReadError = 601

// XDomainRequest ...
// Error code indicating a real error code is unavailable because
// we had to use an XDomainRequest object to allow CORS requests in
// Internet Explorer, which strips the body from HTTP responses that have
// a non-2XX status code.
const XDomainRequest = 602

// MissingAPIKeyError ...
// The request is missing an API key.
const MissingAPIKeyError = 902

// InvalidAPIKeyError ...
// The request is using an invalid API key.
const InvalidAPIKeyError = 903

// LinkingNotSupportedError ...
// Linking to an external account not supported yet with signup_or_login.
const LinkingNotSupportedError = 999
