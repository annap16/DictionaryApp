package customerrors

type DuplicateKeyError struct {
	Msg string
}

func (e *DuplicateKeyError) Error() string {
	return e.Msg
}

func NewDuplicateKeyError(msg string) error {
	return &DuplicateKeyError{Msg: msg}
}
