package medicalpatients

import "time"

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

var Genders []interface{} = []interface{}{Male, Female}

type MedicalPatients struct {
	ID string
	// UserID string
	IdentityNumber  string
	PhoneNumber     string
	Name            string
	Birthdate       time.Time
	Gender          Gender
	IdentityCardUrl string
	CreatedAt       time.Time
}
