package errors

import "fmt"

type UserError struct {
	Msg string
}

func NewUserError(msg string) *UserError {
	return &UserError{msg}
}
func (a *UserError) Error() string {
	return fmt.Sprintf("User Error : %s", a.Msg)
}

var (
	ErrUserExists = NewUserError("User already exists")
)
