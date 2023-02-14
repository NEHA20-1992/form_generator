package validator

import (
	"errors"

	"github.com/NEHA20-1992/form_generator/api/model"
)

var ErrUserFirstNameRequired = errors.New("first name is required")
var ErrUserLastNameRequired = errors.New("last name is required")

func ValidateUser(data *model.User) (err error) {
	if data == nil {
		return
	}

	if data.Nickname == "" {
		err = ErrUserFirstNameRequired
		return
	}
	return
}
