package utils

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func RespondWithSuccess(res http.ResponseWriter, statusCode int, msg string, data any) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	json.NewEncoder(res).Encode(
		SuccessResponse{
			Message: msg,
			Data:    data,
		},
	)
}
