package entity

type ErrorKind int

const (
	ErrValidation ErrorKind = iota
	ErrNotFound
	ErrUnauthorized
	ErrRateLimited
	ErrInternal
)

type DomainError struct {
	Kind    ErrorKind
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewValidationError(msg string) *DomainError {
	return &DomainError{Kind: ErrValidation, Message: msg}
}

func NewNotFoundError(msg string) *DomainError {
	return &DomainError{Kind: ErrNotFound, Message: msg}
}
