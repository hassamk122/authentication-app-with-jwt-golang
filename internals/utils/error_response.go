package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func RespondWithError(res http.ResponseWriter, statusCode int, msg string) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	json.NewEncoder(res).Encode(ErrorResponse{
		Message: msg,
	})
}

func RespondWithNotfound(res http.ResponseWriter) {
	RespondWithError(res, http.StatusNotFound, "Resource not found")
}
