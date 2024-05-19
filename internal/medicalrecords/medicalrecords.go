package medicalrecords

import "time"

type MedicalRecords struct {
	ID string
	// UserID string
	IdentityNumber string
	Symptoms       string
	Medications    string
	CreatedAt      time.Time
}
