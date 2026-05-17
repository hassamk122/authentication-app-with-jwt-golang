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

func (h *Handler) RegisterHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var userReq dtos.CreateUserRequest
		if err := json.NewDecoder(req.Body).Decode(&userReq); err != nil {
			log.Println("Json decoding (handler layer)")
			utils.RespondWithError(res, http.StatusBadRequest, errs.ErrInvalidRequestPayload.Error())
			return
		}

		log.Println("Matched request format passing to validator (handler layer)")

		if err := validation.Validate(&userReq); err != nil {
			log.Println("validation failed (handler layer)")
			utils.RespondWithError(res, http.StatusBadRequest, err.Error())
			return
		}

		log.Println("Valid request passing to service (handler layer)")

		registerInfo, err := h.UserService.Register(ctx, userReq.Username, userReq.Email, userReq.Password)
		if errors.Is(err, errs.ErrEmailTaken) {
			log.Println("Email already taken (handler layer)")
			utils.RespondWithError(res, http.StatusConflict, "Email already taken")
			return
		}

		log.Println("Email not taken (handler layer)")

		if err != nil {
			log.Println("Something went wrong  (handler layer)", err)
			utils.RespondWithError(res, http.StatusInternalServerError, "Error creating user")
			return
		}

		log.Println("No errors found in registering (handler layer)")

		utils.SetAuthCookies(res, registerInfo.AccessToken, registerInfo.RefreshToken)

		log.Println("Cookies set (handler layer)")

		utils.RespondWithSuccess(res, http.StatusCreated, "User created successfully", registerInfo.User)

	}
}

func (h *Handler) LoginHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var userReq dtos.LoginUserRequest
		if err := json.NewDecoder(req.Body).Decode(&userReq); err != nil {
			log.Println("Json decoding (handler layer)")
			utils.RespondWithError(res, http.StatusBadRequest, errs.ErrInvalidRequestPayload.Error())
			return
		}

		log.Println("Matched request format passing to validator (handler layer)")

		if err := validation.Validate(&userReq); err != nil {
			log.Println("validation failed (handler layer)")
			utils.RespondWithError(res, http.StatusBadRequest, err.Error())
			return
		}

		log.Println("Valid request passing to service (handler layer)")

		loginInfo, err := h.UserService.Login(ctx, userReq.Email, userReq.Password)
		if errors.Is(err, errs.ErrInvalidCredentails) {
			log.Println("Invalid credentials (handler layer)")
			utils.RespondWithError(res, http.StatusBadRequest, "Invalid email or password")
			return
		}
		if err != nil {
			log.Println("Login failed (handler layer)")
			utils.RespondWithError(res, http.StatusForbidden, "Invalid email or password")
			return
		}

		log.Println("No errors found in login service (handler layer)")

		utils.SetAuthCookies(res, loginInfo.AccessToken, loginInfo.RefreshToken)

		log.Println("Cookies set (handler layer)")

		utils.RespondWithSuccess(res, http.StatusOK, "Login Successful", nil)
	}
}

func (h *Handler) LogoutHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		accessToken, err := req.Cookie("access_token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				utils.RespondWithError(res, http.StatusBadRequest, "Cookie not found")
			default:
				log.Println(err)
				utils.RespondWithError(res, http.StatusInternalServerError, "Server error")
			}
			return
		}

		err = h.UserService.Logout(ctx, accessToken)
		if err != nil {
			utils.RespondWithError(res, http.StatusBadRequest, "Unauthorized request")
		}

		utils.ClearAuthCookies(res)

		utils.RespondWithSuccess(res, http.StatusOK, "Logout Successful", nil)
	}
}
