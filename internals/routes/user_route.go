package routes

import (
	"net/http"

	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/handlers"
)

func SetupUserRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	userMux := http.NewServeMux()

	mux.Handle("/users/", http.StripPrefix("/users", userMux))

	userMux.Handle("POST /register", handler.RegisterHandler())
	userMux.Handle("POST /login", handler.LoginHandler())
}
