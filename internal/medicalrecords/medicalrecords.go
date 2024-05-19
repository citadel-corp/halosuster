package medicalrecords

import "time"

type MedicalRecords struct {
	ID string
	// UserID string
	PatientId   string
	Symptoms    string
	Medications string
	CreatedAt   time.Time
}
