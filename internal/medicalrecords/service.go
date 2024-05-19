package medicalrecords

import (
	"context"
	"strconv"

	"github.com/citadel-corp/halosuster/internal/common/id"
	"github.com/citadel-corp/halosuster/internal/medicalpatients"
)

type Service interface {
	CreateMedicalRecord(ctx context.Context, req PostMedicalRecord) error
	ListMedicalRecords(ctx context.Context, req ListRecordsPayload) ([]ListMedicalRecordsResponse, error)
}

type medicalRecordsService struct {
	repository        Repository
	patientRepository medicalpatients.Repository
}

func NewService(repository Repository, patientRepository medicalpatients.Repository) Service {
	return &medicalRecordsService{
		repository:        repository,
		patientRepository: patientRepository,
	}
}

func (s *medicalRecordsService) CreateMedicalRecord(ctx context.Context, req PostMedicalRecord) error {
	var err error
	idNumber := strconv.Itoa(int(req.IdentityNumber))

	// get patient by identity number
	patient, err := s.patientRepository.GetByIdentityNumber(ctx, idNumber)
	if err != nil {
		if err == medicalpatients.ErrPatientNotFound {
			return ErrIdNumberDoesNotExist
		}
		return err
	}

	medicalRecord := &MedicalRecords{
		ID:          id.GenerateStringID(16),
		UserID:      req.UserId,
		PatientId:   patient.ID,
		Symptoms:    req.Symptoms,
		Medications: req.Medications,
	}
	err = s.repository.Create(ctx, medicalRecord)
	if err != nil {
		return err
	}

	return nil
}

func (s *medicalRecordsService) ListMedicalRecords(ctx context.Context, req ListRecordsPayload) ([]ListMedicalRecordsResponse, error) {
	if req.Limit == 0 {
		req.Limit = 5
	}

	if req.CreatedAt == "" {
		req.CreatedAt = "desc"
	}

	res, err := s.repository.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
