package apperrors

import "fmt"

type (
	LastItemInListError struct {
		msg string
	}
)

func NewLastItemInListError(listID uint64, itemID uint64) error {
	return &LastItemInListError{
		msg: fmt.Sprintf("Item %d is the last in list %d, cannot remove.", itemID, listID),
	}
}

func (e LastItemInListError) Error() string {
	return e.msg
}
