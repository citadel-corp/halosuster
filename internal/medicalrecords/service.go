package medicalrecords

import (
	"context"
	"strconv"

	"github.com/citadel-corp/halosuster/internal/common/id"
	"github.com/citadel-corp/halosuster/internal/medicalpatients"
)

type Service interface {
	CreateMedicalRecord(ctx context.Context, req PostMedicalRecord) error
}

type userService struct {
	repository        Repository
	patientRepository medicalpatients.Repository
}

func NewService(repository Repository, patientRepository medicalpatients.Repository) Service {
	return &userService{
		repository:        repository,
		patientRepository: patientRepository,
	}
}

func (s *userService) CreateMedicalRecord(ctx context.Context, req PostMedicalRecord) error {
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
