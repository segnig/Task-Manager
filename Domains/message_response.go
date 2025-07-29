package domains

import "errors"

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}

var (
	// Task related errors
	ErrTaskNotFound     = errors.New("task not found")
	ErrTaskExists       = errors.New("task already exists")
	ErrInvalidTask      = errors.New("invalid task")
	ErrTaskValidation   = errors.New("task validation failed")
	ErrTaskUpdateFailed = errors.New("task update failed")
	ErrTaskDeleteFailed = errors.New("task delete failed")

	// Authorization errors
	ErrUnauthorized     = errors.New("unauthorized access")
	ErrPermissionDenied = errors.New("permission denied")

	// Repository errors
	ErrRepoFailure  = errors.New("repository operation failed")
	ErrRepoTimeout  = errors.New("repository timeout")
	ErrRepoConflict = errors.New("repository conflict")

	// Validation errors
	ErrValidation   = errors.New("validation error")
	ErrInvalidInput = errors.New("invalid input")
	ErrMissingField = errors.New("required field missing")
	ErrInvalidID    = errors.New("invalid ID format")

	// System errors
	ErrInternal        = errors.New("internal server error")
	ErrTimeout         = errors.New("operation timed out")
	ErrContextCanceled = errors.New("operation canceled")
)
