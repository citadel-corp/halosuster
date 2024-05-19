package medicalpatients

import "errors"

var (
	ErrPatientNotFound              = errors.New("patient not found")
	ErrPatientIdNumberAlreadyExists = errors.New("identity number already exists")
)
