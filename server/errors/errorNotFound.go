package customerrors

type NotFoundError struct {
	Msg string
}

func (e *NotFoundError) Error() string {
	return e.Msg
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{Msg: msg}
}
