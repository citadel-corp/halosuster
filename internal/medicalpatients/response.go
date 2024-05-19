package medicalpatients

type MedicalPatientsResponse struct {
	IdentityNumber      int64  `json:"identityNumber"`
	PhoneNumber         string `json:"phoneNumber"`
	Name                string `json:"name"`
	Birthdate           string `json:"birthDate"`
	Gender              string `json:"gender"`
	IdentityCardScanImg string `json:"identityCardScanImg"`
}
