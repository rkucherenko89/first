package customErrors

import "fmt"

type ErrorProductNotFound struct {
	Msg string
	ID  int
}

func (e *ErrorProductNotFound) Error() string {
	return fmt.Sprintf("Product not found with ID: %d", e.ID)
}

func NewErrorProductNotFound(msg string, id int) *ErrorProductNotFound {
	return &ErrorProductNotFound{msg, id}
}
