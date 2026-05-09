package utils

import (
	"net/http"
	"os"
	"time"

	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/dtos"
)

func SetAuthCookies(res http.ResponseWriter, registerInfo dtos.RegisterInfo) {

	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    registerInfo.AccessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   ShouldBeSecure(),
	}

	http.SetCookie(res, &accessTokenCookie)

	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    registerInfo.RefreshToken,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().AddDate(0, 1, 0),
		HttpOnly: true,
		Secure:   ShouldBeSecure(),
		Path:     "/auth/refresh",
	}

	http.SetCookie(res, &refreshTokenCookie)
}

func ShouldBeSecure() bool {
	env := os.Getenv("ENVIRONMENT")
	return env != "DEV"
}
