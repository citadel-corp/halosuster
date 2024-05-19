package medicalpatients

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/citadel-corp/halosuster/internal/common/db"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	Create(ctx context.Context, medicalrecord *MedicalPatients) error
	GetByIdentityNumber(ctx context.Context, idNumber string) (*MedicalPatients, error)
	List(ctx context.Context, req ListPatientsPayload) ([]MedicalPatients, error)
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) Create(ctx context.Context, medicalpatient *MedicalPatients) error {
	q := `
        INSERT INTO medical_patients (id, identity_number, phone_number, name, birth_date, gender, identity_card_url)
        VALUES ($1, $2, $3, $4, $5, $6, $7);
    `
	_, err := d.db.DB().ExecContext(ctx, q, medicalpatient.ID, medicalpatient.IdentityNumber, medicalpatient.PhoneNumber,
		medicalpatient.Name, medicalpatient.Birthdate, medicalpatient.Gender, medicalpatient.IdentityCardUrl)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return ErrPatientNotFound
			default:
				return err
			}
		}
		return err
	}
	return nil
}

func (d *dbRepository) GetByIdentityNumber(ctx context.Context, idNumber string) (*MedicalPatients, error) {
	q := `
		SELECT id, identity_number, phone_number, name, birth_date, gender, identity_card_url, created_at
		FROM medical_patients
		WHERE identity_number = $1;
	`
	row := d.db.DB().QueryRowContext(ctx, q, idNumber)
	m := &MedicalPatients{}
	err := row.Scan(&m.ID, &m.IdentityNumber, &m.PhoneNumber, &m.Name, &m.Birthdate, &m.Gender, &m.IdentityCardUrl, &m.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPatientNotFound
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (d *dbRepository) List(ctx context.Context, req ListPatientsPayload) ([]MedicalPatients, error) {
	q := `
			SELECT id, identity_number, phone_number, name, birth_date, gender, identity_card_url, created_at
			FROM medical_patients
	`
	paramNo := 1
	params := make([]interface{}, 0)
	if req.IdentityNumber != "" {
		q += fmt.Sprintf("WHERE identity_number = $%d ", paramNo)
		paramNo += 1
		params = append(params, req.IdentityNumber)
	}
	if req.Name != "" {
		q += whereOrAnd(paramNo)
		q += fmt.Sprintf("LOWER(name) LIKE $%d ", paramNo)
		paramNo += 1
		params = append(params, "%"+req.Name+"%")
	}
	if req.PhoneNumber != "" {
		q += whereOrAnd(paramNo)
		q += fmt.Sprintf("phone_number LIKE $%d ", paramNo)
		paramNo += 1
		params = append(params, "%"+req.PhoneNumber+"%")
	}

	if req.CreatedAt == "asc" || req.CreatedAt == "desc" {
		q += `ORDER BY created_at ` + req.CreatedAt
	}

	q += fmt.Sprintf(" OFFSET $%d LIMIT $%d", paramNo, paramNo+1)
	params = append(params, req.Offset)
	params = append(params, req.Limit)

	rows, err := d.db.DB().QueryContext(ctx, q, params...)
	if err != nil {
		return nil, err
	}
	res := make([]MedicalPatients, 0)
	for rows.Next() {
		m := MedicalPatients{}
		err = rows.Scan(&m.ID, &m.IdentityNumber, &m.PhoneNumber, &m.Name, &m.Birthdate,
			&m.Gender, &m.IdentityCardUrl, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
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
