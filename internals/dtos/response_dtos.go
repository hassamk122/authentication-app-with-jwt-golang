package dtos

import "github.com/hassamk122/authentication-app-with-jwt-golang/internals/store"

type RegisterInfo struct {
	User         *store.CreateUserRow
	AccessToken  string
	RefreshToken string
}

type LoginInfo struct {
	User         *store.GetUserByEmailRow
	AccessToken  string
	RefreshToken string
}
