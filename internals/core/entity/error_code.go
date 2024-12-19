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
	FileErrorNamespace           = errorx.NewNamespace("FileError")
)

var (
	AppInternalError   = ApplicationErrorNamespace.NewType("InternalError")
	ValidationError    = ApplicationErrorNamespace.NewType("ValidationError")
	BadRequest         = ApplicationErrorNamespace.NewType("BadRequest")
	TimeoutError       = ApplicationErrorNamespace.NewType("TimeoutError")
	AppConnectionError = ApplicationErrorNamespace.NewType("ConnectionError")
)

var (
	ConnectionError      = DatabaseErrorNamespace.NewType("ConnectionError")
	DuplicateEntry       = DatabaseErrorNamespace.NewType("DuplicateEntry")
	UnableToSave         = DatabaseErrorNamespace.NewType("UnableToSave")
	UnableToFindResource = DatabaseErrorNamespace.NewType("UnableToFindResource")
	UnableToRead         = DatabaseErrorNamespace.NewType("UableToRead")
)

var (
	InvalidCredentials = AuthenticationErrorNamespace.NewType("InvalidCredentials")
	Unauthorized       = AuthenticationErrorNamespace.NewType("Unauthorized")
	AuthInternalError  = AuthenticationErrorNamespace.NewType("InternalError")
)

var (
	FileTooLarge     = FileErrorNamespace.NewType("FileTooLarge")
	InvalidExtension = FileErrorNamespace.NewType("InvalidFileExtension")
	UnableToSaveFile = FileErrorNamespace.NewType("UnableToSave")
)

var (
	StatusCode = errorx.RegisterProperty("StatusCode")
)
