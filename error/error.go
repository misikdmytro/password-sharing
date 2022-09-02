package error

import (
	"fmt"

	"github.com/misikdmitriy/password-sharing/model"
)

type ErrorCodes int

const (
	BadRequest          ErrorCodes = 40000
	PasswordNotFound               = 40401
	InternalServerError            = 50000
	InitDbError                    = 50001
	RandomizerError                = 50002
	DbQueryError                   = 50003
	DbCommandError                 = 50004
	EncodeError                    = 50005
	DecodeError                    = 50006
)

type PasswordSharingError struct {
	Code    ErrorCodes
	Message string
}

func (e *PasswordSharingError) Error() string {
	return fmt.Sprintf("password sharing error occured. code: %d. message: %s",
		e.Code,
		e.Message)
}

func (e *PasswordSharingError) ToResponse() (int, *model.ErrorResponse) {
	code := int(e.Code)
	return code / 100, &model.ErrorResponse{
		Code:    int(e.Code),
		Message: e.Message,
	}
}

func AsPasswordSharingError(err error) *PasswordSharingError {
	psError, ok := err.(*PasswordSharingError)
	if ok {
		return psError
	}

	return &PasswordSharingError{
		Code:    InternalServerError,
		Message: err.Error(),
	}
}

func BadRequestError() (int, *model.ErrorResponse) {
	return 400, &model.ErrorResponse{
		Code:    int(BadRequest),
		Message: "bad request",
	}
}
