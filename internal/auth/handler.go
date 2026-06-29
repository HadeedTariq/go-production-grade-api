package auth

import (
	"net/http"

	authDto "github.com/HadeedTariq/go-production-grade-api/internal/auth/dto"
	"github.com/HadeedTariq/go-production-grade-api/internal/json"
	"github.com/HadeedTariq/go-production-grade-api/internal/validator"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req authDto.SignupRequest

	err := json.Read(r, &req)

	if err != nil {
		json.Write(w, http.StatusBadRequest, Response{Message: "Invalid JSON payload"})
		return
	}

	err = validator.Validate.Struct(req)

	if err != nil {
		json.Write(w, http.StatusBadRequest, Response{
			Message: validator.ParseValidationErrors(err),
		})
		return
	}
	msg, err := h.service.RegisterUser(r.Context(), req)

	if err != nil {
		json.Write(w, http.StatusBadRequest, Response{Message: err.Error()})
		return
	}

	json.Write(w, http.StatusCreated, Response{
		Message: msg,
	})
}

func (h *handler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		json.Write(w, http.StatusBadRequest, Response{Message: "Token is required"})
		return
	}

	msg, err := h.service.VerifyUser(r.Context(), token)

	if err != nil {
		json.Write(w, http.StatusBadRequest, Response{Message: err.Error()})
		return
	}

	json.Write(w, http.StatusCreated, Response{
		Message: msg,
	})
}

func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req authDto.SigninRequest

	err := json.Read(r, &req)

	if err != nil {
		json.Write(w, http.StatusBadRequest, Response{Message: "Invalid JSON payload"})
		return
	}
	err = validator.Validate.Struct(req)

	if err != nil {
		json.Write(w, http.StatusBadRequest, Response{
			Message: validator.ParseValidationErrors(err),
		})
		return
	}

	tokens, err := h.service.LoginUser(r.Context(), req)

	if err != nil {
		json.Write(w, http.StatusBadRequest, Response{Message: err.Error()})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    tokens.AccessToken,
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	json.Write(w, http.StatusCreated, Response{
		Message: "Login successful",
	})
}
