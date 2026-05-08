package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/dtos"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/errs"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/utils"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/validation"
)

func (h *Handler) CreateUserHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var userReq dtos.CreateUserRequest
		if err := json.NewDecoder(req.Body).Decode(&userReq); err != nil {
			log.Println("Json decoding (handler layer)")
			utils.RespondWithError(res, http.StatusBadGateway, errs.ErrInvalidRequestPayload.Error())
			return
		}

		log.Println("Matched request format passing to validator (handler layer)")

		if err := validation.Validate(&userReq); err != nil {
			log.Println("validation failed (handler layer)")
			utils.RespondWithError(res, http.StatusBadRequest, err.Error())
			return
		}

		log.Println("Valid request passing to service (handler layer)")

		registerInfo, err := h.UserService.Register(ctx, userReq.Username, userReq.Email, userReq.Email)
		if errors.Is(err, errs.ErrEmailTaken) {
			log.Println("Email already taken (handler layer)")
			utils.RespondWithError(res, http.StatusConflict, "Email already taken")
			return
		}

		log.Println("Email not taken (handler layer)", registerInfo)

		if err != nil {
			log.Println("Something went wrong  (handler layer)")
			utils.RespondWithError(res, http.StatusInternalServerError, "Error creating user")
			return
		}

		log.Println("No errors found in registering (handler layer)")

		utils.RespondWithSuccess(res, http.StatusCreated, "User created successfully", nil)

	}
}
