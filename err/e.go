package err

import "errors"
import "strconv"

// E 组装 json 格式错误信息：
// {"code": 105,"error": "invalid field name: bl!ng"}
func E(code int, msg string) error {
	text := `{"code": ` + strconv.Itoa(code) + `,"error": "` + msg + `"}`
	return errors.New(text)
}

// OtherCause ...
// Error code indicating some error other than those enumerated here.
var OtherCause = -1

// InternalServerError ...
// Error code indicating that something has gone wrong with the server.
// If you get this error code, it is Parse's fault. Contact us at
// https://parse.com/help
var InternalServerError = 1

// ConnectionFailed ...
// Error code indicating the connection to the Parse servers failed.
var ConnectionFailed = 100

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

// PushMisconfigured ...
// Error code indicating that push is misconfigured.
var PushMisconfigured = 115

// ObjectTooLarge ...
// Error code indicating that the object is too large.
var ObjectTooLarge = 116

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

// DuplicateValue ...
// Error code indicating that a unique field was given a value that is
// already taken.
var DuplicateValue = 137

// InvalidRoleName ...
// Error code indicating that a role's name is invalid.
var InvalidRoleName = 139

// ExceededQuota ...
// Error code indicating that an application quota was exceeded.  Upgrade to
// resolve.
var ExceededQuota = 140

// ScriptFailed ...
// Error code indicating that a Cloud Code script failed.
var ScriptFailed = 141

// ValidationError ...
// Error code indicating that a Cloud Code validation failed.
var ValidationError = 142

// InvalidImageData ...
// Error code indicating that invalid image data was provided.
var InvalidImageData = 143

// UnsavedFileError ...
// Error code indicating an unsaved file.
var UnsavedFileError = 151

// InvalidPushTimeError ...
// Error code indicating an invalid push time.
var InvalidPushTimeError = 152

// FileDeleteError ...
// Error code indicating an error deleting a file.
var FileDeleteError = 153

// RequestLimitExceeded ...
// Error code indicating that the application has exceeded its request
// limit.
var RequestLimitExceeded = 155

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

// UnsupportedService ...
// Error code indicating that a service being linked (e.g. Facebook or
// Twitter) is unsupported.
var UnsupportedService = 252

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
