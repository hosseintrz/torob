package errors

import "fmt"

type AuthError struct {
	Msg string
}

func NewAuthError(msg string) *AuthError {
	return &AuthError{msg}
}
func (a *AuthError) Error() string {
	return fmt.Sprintf("User Error : %s", a.Msg)
}

var (
	ErrInvalidToken = NewAuthError("token is invalid")
)
