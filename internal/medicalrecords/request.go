package medicalrecords

import (
	"fmt"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type PostMedicalRecord struct {
	IdentityNumber int64  `json:"identityNumber"`
	Symptoms       string `json:"symptoms"`
	Medications    string `json:"medications"`
}

func (p PostMedicalRecord) Validate() error {
	idNumber := strconv.Itoa(int(p.IdentityNumber))
	if len(idNumber) != 16 {
		return fmt.Errorf("%s: %s", "identityNumber", "must be 16 characters")
	}

	return validation.ValidateStruct(&p,
		validation.Field(&p.IdentityNumber, validation.Required),
		validation.Field(&p.Symptoms, validation.Required, validation.Length(1, 2000)),
		validation.Field(&p.Medications, validation.Required, validation.Length(1, 2000)),
	)
}
