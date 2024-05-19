package medicalpatients

import (
	"context"
	"errors"

	"github.com/citadel-corp/halosuster/internal/common/db"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	Create(ctx context.Context, medicalrecord *MedicalPatients) error
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
				return ErrPatientIdNumberAlreadyExists
			default:
				return err
			}
		}
		return err
	}
	return nil
}
