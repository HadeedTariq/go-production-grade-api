package auth

import (
	"net/http"
	"time"

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

func (h *handler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.AuthenticateUser(r.Context())

	if err != nil {
		json.Write(w, http.StatusBadRequest, Response{Message: err.Error()})
		return
	}

	json.Write(w, 200, user)
}

func (h *handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	msg, err := h.service.LogoutUser(r.Context())

	if err != nil {
		json.Write(w, http.StatusBadRequest, Response{Message: err.Error()})
		return
	}
	// Clear the accessToken cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    "",              // Clear the value
		Path:     "/",             // Must match the original path exactly
		MaxAge:   -1,              // -1 means delete cookie now
		Expires:  time.Unix(0, 0), // Set expiration date to the deep past (1970)
		HttpOnly: false,
		Secure:   true, // Match original security flags
		SameSite: http.SameSiteNoneMode,
	})

	// Clear the refreshToken cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	json.Write(w, 200, Response{
		Message: msg,
	})
}
