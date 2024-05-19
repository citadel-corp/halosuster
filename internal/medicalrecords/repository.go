package medicalrecords

import (
	"context"
	"fmt"

	"github.com/citadel-corp/halosuster/internal/common/db"
	"github.com/citadel-corp/halosuster/internal/medicalpatients"
	"github.com/citadel-corp/halosuster/internal/user"
)

type Repository interface {
	Create(ctx context.Context, medicalrecord *MedicalRecords) error
	List(ctx context.Context, req ListRecordsPayload) ([]ListMedicalRecordsResponse, error)
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) Create(ctx context.Context, medicalrecord *MedicalRecords) error {
	q := `
        INSERT INTO medical_records (id, user_id, patient_id, symptoms, medications)
        VALUES ($1, $2, $3, $4, $5);
    `
	_, err := d.db.DB().ExecContext(ctx, q, medicalrecord.ID, medicalrecord.UserID, medicalrecord.PatientId, medicalrecord.Symptoms, medicalrecord.Medications, medicalrecord.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (d *dbRepository) List(ctx context.Context, req ListRecordsPayload) ([]ListMedicalRecordsResponse, error) {
	q := `
			SELECT symptoms, medications, medical_records.created_at,
				users.id, users.nip, users.name,
				medical_patients.identity_number, medical_patients.phone_number,
                medical_patients.name, medical_patients.birth_date, medical_patients.gender,
                medical_patients.identity_card_url
			FROM medical_records
			LEFT JOIN users ON users.id = medical_records.user_id
			LEFT JOIN medical_patients ON medical_patients.id = patient_id
	`
	paramNo := 1
	params := make([]interface{}, 0)
	if req.IdentityNumber != "" {
		q += fmt.Sprintf("WHERE medical_patients.identity_number = $%d ", paramNo)
		paramNo += 1
		params = append(params, req.IdentityNumber)
	}
	// fmt.Println("req.UserId", req.UserId)
	if req.UserId != "" {
		q += whereOrAnd(paramNo)
		q += fmt.Sprintf("users.id = $%d ", paramNo)
		paramNo += 1
		params = append(params, req.UserId)
	}
	if req.NIP != "" {
		q += whereOrAnd(paramNo)
		q += fmt.Sprintf("LOWER(users.nip) = $%d ", paramNo)
		paramNo += 1
		params = append(params, req.NIP)
	}

	if req.CreatedAt == "asc" || req.CreatedAt == "desc" {
		q += `ORDER BY medical_records.created_at ` + req.CreatedAt
	}

	q += fmt.Sprintf(" OFFSET $%d LIMIT $%d", paramNo, paramNo+1)
	params = append(params, req.Offset)
	params = append(params, req.Limit)

	// fmt.Println("query:", q)

	rows, err := d.db.DB().QueryContext(ctx, q, params...)
	if err != nil {
		return nil, err
	}
	res := make([]ListMedicalRecordsResponse, 0)
	for rows.Next() {
		m := ListMedicalRecordsResponse{}
		p := medicalpatients.MedicalPatientsResponse{}
		u := user.UserResponse{}
		err = rows.Scan(&m.Symptoms, &m.Medications, &m.CreatedAt,
			&u.UserID, &u.NIP, &u.Name,
			&p.IdentityNumber, &p.PhoneNumber, &p.Name, &p.Birthdate,
			&p.Gender, &p.IdentityCardScanImg,
		)
		if err != nil {
			return nil, err
		}

		m.IdentityDetail = p
		m.CreatedBy = u
		res = append(res, m)
	}
	return res, nil
}

func whereOrAnd(paramNo int) string {
	if paramNo == 1 {
		return "WHERE "
	}
	return "OR "
}
