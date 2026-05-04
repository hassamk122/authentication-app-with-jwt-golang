package routes

import (
	"net/http"

	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/handlers"
)

func SetupHealthRoute(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("/health", handler.HealthHandler())
}
