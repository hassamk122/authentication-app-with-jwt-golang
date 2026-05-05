package validation

import (
	"testing"

	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/dtos"
	"github.com/stretchr/testify/assert"
)

func TestValidate_IfUserRequest_Valid(t *testing.T) {
	user := dtos.CreateUserRequest{
		Username: "test",
		Email:    "test@mail.com",
		Password: "12345678",
	}

	err := Validate(user)

	assert.NoError(t, err)
}

func TestValidate_IfUserRequest_InValidUser(t *testing.T) {
	t.Run("If user has invalid name", func(t *testing.T) {
		not_valid_username := dtos.CreateUserRequest{
			Username: "te",
			Email:    "test@mail.com",
			Password: "12345678",
		}
		err := Validate(not_valid_username)
		assert.Error(t, err)
	})

	t.Run("If user has invalid mail length", func(t *testing.T) {
		not_valid_email := dtos.CreateUserRequest{
			Username: "test",
			Email:    "test@m",
			Password: "12345678",
		}
		err := Validate(not_valid_email)
		assert.Error(t, err)
	})

	t.Run("If user has invalid passoword length", func(t *testing.T) {
		not_valid_password := dtos.CreateUserRequest{
			Username: "test",
			Email:    "test@m",
			Password: "123456",
		}
		err := Validate(not_valid_password)
		assert.Error(t, err)
	})
}
