package common

type LogicError struct {
	inner   error
	Message string
	Code    int
}

func (err LogicError) Error() string {
	return err.Message
}

func (err LogicError) Wrap(inner error) error {
	return LogicError{
		Code:    err.Code,
		Message: err.Message,
		inner:   inner,
	}
}

func (err LogicError) Unwrap() error {
	return err.inner
}

func (err LogicError) Is(target error) bool {
	return target.Error() == err.Message
}

var (
	ErrUnknown        = LogicError{Message: "Unknown error", Code: 500}
	ErrInput          = LogicError{Message: "Input error", Code: 400}
	ErrNotFound       = LogicError{Message: "Not found", Code: 404}
	ErrNotImplemented = LogicError{Message: "Not implemented", Code: 500}
	ErrUnauthorized   = LogicError{Message: "Unauthorized", Code: 401}
	ErrUnexpected     = LogicError{Message: "Unexpected error", Code: 500}
)
