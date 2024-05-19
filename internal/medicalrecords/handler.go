package medicalrecords

import (
	"errors"
	"net/http"

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
	var req PostMedicalRecord

	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

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
