package medicalrecords

import (
	"github.com/citadel-corp/halosuster/internal/medicalpatients"
	"github.com/citadel-corp/halosuster/internal/user"
)

type ListMedicalRecordsResponse struct {
	IdentityDetail medicalpatients.MedicalPatientsResponse `json:"identityDetail"`
	Symptoms       string                                  `json:"symptoms"`
	Medications    string                                  `json:"medications"`
	CreatedAt      string                                  `json:"createdAt"`
	CreatedBy      user.UserResponse                       `json:"createdBy"`
}
