package dtos

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,min=7,max=30"`
	Password string `json:"password" validate:"required,min=8"`
}
