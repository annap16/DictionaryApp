package customerrors

type ForeignKeyError struct {
    Msg string
}

func (e *ForeignKeyError) Error() string {
    return e.Msg
}

func NewForeignKeyError(msg string) error {
    return &ForeignKeyError{Msg: msg}
}
