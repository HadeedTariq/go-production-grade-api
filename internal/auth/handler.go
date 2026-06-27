package auth

import (
	"net/http"

	authDto "github.com/HadeedTariq/go-production-grade-api/internal/auth/dto"
	"github.com/HadeedTariq/go-production-grade-api/internal/json"
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
	// ~ so according to me the validation related stuff is handled over there
	var req authDto.SignupRequest

	err := json.Read(r, &req)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// ~ so over there not any error occur over there
	msg, err := h.service.RegisterUser(r.Context(), req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.Write(w, http.StatusCreated, Response{
		Message: msg,
	})
}
