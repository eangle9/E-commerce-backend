package entity

import "github.com/joomcode/errorx"

const (
	NotFoundError = "NOT_FOUND_ERROR"
	// ValidationError = "VALIDATION_ERROR"
	InvalidRequest = "INVALID_REQUEST"
	// DatabaseError   = "DATABASE_ERROR"
	InternalError = "INTERNAL_SERVER_ERROR"
	// AuthenticationError = "AUTHENTICATION_ERROR"
	// Unauthorized = "UNAUTHORIZED"
	Success = "SUCCESS"
)

var (
	ApplicationErrorNamespace    = errorx.NewNamespace("ApplicationError")
	DatabaseErrorNamespace       = errorx.NewNamespace("DatabaseError")
	AuthenticationErrorNamespace = errorx.NewNamespace("AuthenticationError")
)

var (
	AppInternalError = ApplicationErrorNamespace.NewType("InternalError")
	ValidationError  = ApplicationErrorNamespace.NewType("ValidationError")
	BadRequest       = ApplicationErrorNamespace.NewType("BadRequest")
)

var (
	ConnectionError      = DatabaseErrorNamespace.NewType("ConnectionError")
	DuplicateEntry       = DatabaseErrorNamespace.NewType("DuplicateEntry")
	UnableToSave         = DatabaseErrorNamespace.NewType("UnableToSave")
	UnableToFindResource = DatabaseErrorNamespace.NewType("UnableToFindResource")
	UnableToRead         = DatabaseErrorNamespace.NewType("UableToRead")
	// DbInternalError      = DatabaseErrorNamespace.NewType("InternalError")
)

var (
	InvalidCredentials = AuthenticationErrorNamespace.NewType("InvalidCredentials")
	Unauthorized       = AuthenticationErrorNamespace.NewType("Unauthorized")
	AuthInternalError  = AuthenticationErrorNamespace.NewType("InternalError")
)

var (
	StatusCode = errorx.RegisterProperty("StatusCode")
)
