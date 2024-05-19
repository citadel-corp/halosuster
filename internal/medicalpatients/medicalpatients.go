package medicalpatients

import "time"

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

var Genders []interface{} = []interface{}{Male, Female}

type MedicalPatients struct {
	ID string `json:"-"`
	// UserID string
	IdentityNumber  string    `json:"identityNumber"`
	PhoneNumber     string    `json:"phoneNumber"`
	Name            string    `json:"name"`
	Birthdate       time.Time `json:"birthDate"`
	Gender          Gender    `json:"gender"`
	IdentityCardUrl string    `json:"-"`
	CreatedAt       time.Time `json:"createdAt"`
}
