package errors

import "fmt"

type ProductError struct {
	Msg string
}

func NewProductError(msg string) *ProductError {
	return &ProductError{msg}
}
func (a *ProductError) Error() string {
	return fmt.Sprintf("User Error : %s", a.Msg)
}

var (
	ErrDupCategory = NewProductError("category exists")
	ErrDupProduct  = NewProductError("product exists")
)
