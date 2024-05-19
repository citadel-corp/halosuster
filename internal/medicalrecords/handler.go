package medicalrecords

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/citadel-corp/halosuster/internal/common/middleware"
	"github.com/citadel-corp/halosuster/internal/common/request"
	"github.com/citadel-corp/halosuster/internal/common/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateMedicalRecord(w http.ResponseWriter, r *http.Request) {
	var err error

	userId, err := getUserID(r)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, response.ResponseBody{
			Message: "unauthorized",
			Error:   err.Error(),
		})
		return
	}

	var req PostMedicalRecord

	err = request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	req.UserId = userId

	err = req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}

	err = h.service.CreateMedicalRecord(r.Context(), req)
	if errors.Is(err, ErrIdNumberDoesNotExist) {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "not found",
			Error:   err.Error(),
		})
		return
	}
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusCreated, response.ResponseBody{
		Message: "Record registered successfully",
	})
}

func getUserID(r *http.Request) (string, error) {
	if authValue, ok := r.Context().Value(middleware.ContextAuthKey{}).(string); ok {
		return authValue, nil
	} else {
		slog.Error("cannot parse auth value from context")
		return "", errors.New("cannot parse auth value from context")
	}
}
