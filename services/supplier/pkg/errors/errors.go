package errors

import "fmt"

type SupplierError struct {
	Msg string
}

func NewSupplierError(msg string) *SupplierError {
	return &SupplierError{msg}
}
func (a *SupplierError) Error() string {
	return fmt.Sprintf("Supplier Error : %s", a.Msg)
}

var (
	ErrDupStore = NewSupplierError("duplicate store")
	ErrDupOffer = NewSupplierError("duplicate offers")
)
