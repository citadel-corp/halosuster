package medicalrecords

import "context"

type Service interface {
	CreateMedicalRecord(ctx context.Context, req PostMedicalRecords) (*MedicalRecordsResponse, error)
}

type userService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &userService{repository: repository}
}
