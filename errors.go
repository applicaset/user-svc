package usersvc

import "fmt"

type ErrUserNotFound struct {
	UUID string
}

func (err ErrUserNotFound) Error() string {
	return fmt.Sprintf("user with uuid '%s' not found", err.UUID)
}

type ErrInvalidAuthMethod struct {
	Name string
}

func (err ErrInvalidAuthMethod) Error() string {
	return fmt.Sprintf("auth method '%s' is not valid", err.Name)
}

type ErrInvalidAuthData struct {
}

func (err ErrInvalidAuthData) Error() string {
	return "auth data is not valid"
}
