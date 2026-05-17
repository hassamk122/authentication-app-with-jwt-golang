package utils

import (
	"net/http"
	"os"
	"time"
)

func SetAuthCookies(res http.ResponseWriter, accessToken, refreshToken string) {

	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Path:     "/",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   ShouldBeSecure(),
	}

	http.SetCookie(res, &accessTokenCookie)

	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
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

func ClearAuthCookies(res http.ResponseWriter) {

	clearAccessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   ShouldBeSecure(),
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(res, clearAccessTokenCookie)

	clearRefreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth/refresh",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   ShouldBeSecure(),
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(res, clearRefreshTokenCookie)
}
