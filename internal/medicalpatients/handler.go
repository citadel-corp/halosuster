package medicalpatients

import (
	"errors"
	"net/http"

	"github.com/citadel-corp/halosuster/internal/common/request"
	"github.com/citadel-corp/halosuster/internal/common/response"
	"github.com/gorilla/schema"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateMedicalPatient(w http.ResponseWriter, r *http.Request) {
	var req PostMedicalPatients

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

	err = h.service.CreateMedicalPatients(r.Context(), req)
	if errors.Is(err, ErrPatientIdNumberAlreadyExists) {
		response.JSON(w, http.StatusConflict, response.ResponseBody{
			Message: "conflict",
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
		Message: "Patient registered successfully",
	})
}

func (h *Handler) ListMedicalPatient(w http.ResponseWriter, r *http.Request) {
	var req ListPatientsPayload

	newSchema := schema.NewDecoder()
	newSchema.IgnoreUnknownKeys(true)

	if err := newSchema.Decode(&req, r.URL.Query()); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{})
		return
	}

	patients, err := h.service.ListMedicalPatients(r.Context(), req)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, response.ResponseBody{
		Message: "Patients fetched successfully",
		Data:    patients,
	})
}
