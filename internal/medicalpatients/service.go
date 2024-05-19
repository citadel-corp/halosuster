package medicalpatients

import (
	"context"
	"strconv"

	"github.com/citadel-corp/halosuster/internal/common/id"
)

type Service interface {
	CreateMedicalPatients(ctx context.Context, req PostMedicalPatients) error
}

type medicalPatientsService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &medicalPatientsService{repository: repository}
}

func (s *medicalPatientsService) CreateMedicalPatients(ctx context.Context, req PostMedicalPatients) error {
	var err error
	idNumber := strconv.Itoa(int(req.IdentityNumber))

	medicalpatient := &MedicalPatients{
		ID:              id.GenerateStringID(16),
		IdentityNumber:  idNumber,
		PhoneNumber:     req.PhoneNumber,
		Name:            req.Name,
		Birthdate:       req.Birthdate,
		Gender:          req.Gender,
		IdentityCardUrl: req.IdentityCardScanImg,
	}
	err = s.repository.Create(ctx, medicalpatient)
	if err != nil {
		return err
	}

	return nil
}
