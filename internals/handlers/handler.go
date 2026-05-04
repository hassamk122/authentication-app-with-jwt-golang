package handlers

import "github.com/hassamk122/authentication-app-with-jwt-golang/internals/services"

type Handler struct {
	UserService services.UserService
}

func NewHandler(userService services.UserService) *Handler {
	return &Handler{
		UserService: userService,
	}
}
