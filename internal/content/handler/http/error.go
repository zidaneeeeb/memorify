package http

import (
	"errors"
	"hbdtoyou/internal/content"
)

// Followings are the known errors from Content HTTP handlers.
var (
	// errBadRequest is returned when the given request is
	// bad/invalid.
	errBadRequest = errors.New("BAD_REQUEST")

	// errExpiredToken is returned when the given token is
	// expired.
	errExpiredToken = errors.New("EXPIRED_TOKEN")

	// errInternalServer is returned when there is an
	// unexpected error encountered when processing a request.
	errInternalServer = errors.New("INTERNAL_SERVER_ERROR")

	// errUserAlreadyExist is returned when the given
	// user already exist based on the predefined
	// unique constraints.
	errUserAlreadyExist = errors.New("USER_ALREADY_EXIST")

	// errDataNotFound is returned when the desired data is
	// not found.
	errDataNotFound = errors.New("DATA_NOT_FOUND")

	// errInvalidPassword is returned when the given password
	// is invalid.
	errInvalidPassword = errors.New("INVALID_PASSWORD")

	// errInvalidToken is returned when the given token is
	// invalid.
	errInvalidToken = errors.New("INVALID_TOKEN")

	// errInvalidUserID is returned when the given user ID is
	// invalid.
	errInvalidUserID = errors.New("INVALID_USER_ID")

	// errInvalidContentID is returned when the given content ID is
	// invalid.
	errInvalidContentID = errors.New("INVALID_CONTENT_ID")

	// errInvalidContentName is returned when the given content name is
	// invalid.
	errInvalidContentName = errors.New("INVALID_CONTENT_NAME")

	// errInvalidDetailContentJSONText is returned when the given detail content JSON text is
	// invalid.
	errInvalidDetailContentJSONText = errors.New("INVALID_DETAIL_CONTENT_JSON_TEXT")

	// errInvalidTemplateID is returned when the given template id is
	// invalid.
	errInvalidTemplateID = errors.New("INVALID_TEMPLATE_ID")

	// errInvalidContentType is returned when the given content type is
	// invalid.
	errInvalidContentType = errors.New("INVALID_CONTENT_TYPE")

	// errInvalidContentStatus is returned when the given content status is
	// invalid.
	errInvalidContentStatus = errors.New("INVALID_CONTENT_STATUS")

	// errInvalidContentAccess is returned when the given content access is
	// invalid.
	errInvalidContentAccess = errors.New("INVALID_CONTENT_ACCESS")

	// errInvalidUsername is returned when the given username
	// is invalid.
	errInvalidUsername = errors.New("INVALID_USERNAME")

	// errMethodNotAllowed is returned when accessing not
	// allowed HTTP method.
	errMethodNotAllowed = errors.New("METHOD_NOT_ALLOWED")

	// errRequestTimeout is returned when processing time has
	// reached the timeout limit.
	errRequestTimeout = errors.New("REQUEST_TIMEOUT")

	// errSourceNotProvided is returned when there is no
	// source provided in the request.
	errSourceNotProvided = errors.New("SOURCE_NOT_PROVIDED")

	// errTooManyRequest is returned when the given request is
	// exceeding the maximum allowed.
	errTooManyRequest = errors.New("TOO_MANY_REQUEST")

	// errUnauthorizedAccess is returned when the request
	// is unaothorized.
	errUnauthorizedAccess = errors.New("UNAUTHORIZED_ACCESS")
)

var (
	// mapHTTPError maps service error into HTTP error that
	// categorize as bad request error.
	//
	// Internal server error-related should not be mapped here,
	// and the handler should just return `errInternal` as the
	// error instead
	mapHTTPError = map[error]error{
		content.ErrInvalidContentID:             errInvalidContentID,
		content.ErrInvalidTemplateID:            errInvalidTemplateID,
		content.ErrInvalidContentType:           errInvalidContentType,
		content.ErrInvalidContentStatus:         errInvalidContentStatus,
		content.ErrInvalidContentName:           errInvalidContentName,
		content.ErrInvalidDetailContentJSONText: errInvalidDetailContentJSONText,
		content.ErrInvalidUserID:                errInvalidUserID,
		content.ErrDataNotFound:                 errDataNotFound,
		content.ErrInvalidContentAccess:         errInvalidContentAccess,
	}
)
