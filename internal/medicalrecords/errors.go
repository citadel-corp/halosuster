package medicalrecords

import "errors"

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrIdNumberDoesNotExist = errors.New("identity number does not exist")
)
