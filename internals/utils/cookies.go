package utils

import (
	"net/http"
	"os"

	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/dtos"
)

const (
	FIFTEEN_MINUTES_IN_SECONDS = 900
	THIRTY_DAYS_IN_SECONDS     = 2_592_000
)

func SetAuthCookies(res http.ResponseWriter, registerInfo dtos.RegisterInfo) {

	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    registerInfo.AccessToken,
		MaxAge:   FIFTEEN_MINUTES_IN_SECONDS,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   ShouldBeSecure(),
	}

	http.SetCookie(res, &accessTokenCookie)

	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    registerInfo.RefreshToken,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   THIRTY_DAYS_IN_SECONDS,
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
