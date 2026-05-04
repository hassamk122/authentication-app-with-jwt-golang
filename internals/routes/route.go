package routes

import (
	"net/http"

	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/handlers"
)

func SetupRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	SetupHealthRoute(mux, handler)
}
