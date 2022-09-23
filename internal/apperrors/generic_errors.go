package apperrors

import "fmt"

type (
	NotFoundError struct {
		msg string
	}

	ObjectInInvalidStateError struct {
		msg string
	}

	InternalError struct {
		msg string
	}

	UserLoginError struct {
		msg string
	}

	LoginAlreadyRegistered struct {
		msg string
	}
)

func NewNotFoundError(entity string, id uint64) error {
	return &NotFoundError{
		msg: fmt.Sprintf("%s %d not found", entity, id),
	}
}

func NewItemNotFoundInListError(itemID uint64, listID uint64) error {
	return &NotFoundError{
		msg: fmt.Sprintf("not found item %d in list %d", itemID, listID),
	}
}

func NewObjectInInvalidStateError(msg string) error {
	return &ObjectInInvalidStateError{msg}
}

func NewInternalError(msg string) error {
	return &InternalError{
		msg: fmt.Sprintf("%s. Contact the administrator.", msg),
	}
}

func NewUserLoginError() error {
	return &UserLoginError{msg: "Invalid user or password."}
}

func NewUserSSOLoginError() error {
	return &UserLoginError{msg: "Invalid SSO login or token."}
}

func NewLoginAlreadyRegisteredError(login string) error {
	return LoginAlreadyRegistered{msg: fmt.Sprintf("Login %s already in use.", login)}
}

func (e NotFoundError) Error() string {
	return e.msg
}

func (e ObjectInInvalidStateError) Error() string {
	return e.msg
}

func (e InternalError) Error() string {
	return e.msg
}

func (e UserLoginError) Error() string {
	return e.msg
}

func (e LoginAlreadyRegistered) Error() string {
	return e.msg
}
