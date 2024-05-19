package medicalpatients

import (
	"context"
	"strconv"
	"strings"

	"github.com/citadel-corp/halosuster/internal/common/id"
)

type Service interface {
	CreateMedicalPatients(ctx context.Context, req PostMedicalPatients) error
	ListMedicalPatients(ctx context.Context, req ListPatientsPayload) ([]MedicalPatients, error)
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

func (s *medicalPatientsService) ListMedicalPatients(ctx context.Context, req ListPatientsPayload) ([]MedicalPatients, error) {
	if req.Limit == 0 {
		req.Limit = 5
	}

	if req.CreatedAt == "" {
		req.CreatedAt = "desc"
	}

	req.PhoneNumber = strings.Replace(req.PhoneNumber, "+", "", 1)

	res, err := s.repository.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
