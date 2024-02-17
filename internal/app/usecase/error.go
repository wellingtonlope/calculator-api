package usecase

const (
	ErrorTypeUnknown ErrorType = "unknown"
	ErrorTypeInvalid ErrorType = "invalid"
)

type (
	ErrorType string
	Error     struct {
		Message string
		Cause   error
		Type    ErrorType
	}
)

func (e Error) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Message
}

func NewError(message string, cause error, errorType ErrorType) Error {
	return Error{
		Message: message,
		Cause:   cause,
		Type:    errorType,
	}
}
