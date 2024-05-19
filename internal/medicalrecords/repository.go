package medicalrecords

import (
	"context"

	"github.com/citadel-corp/halosuster/internal/common/db"
)

type Repository interface {
	Create(ctx context.Context, medicalrecord *MedicalRecords) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) Create(ctx context.Context, medicalrecord *MedicalRecords) error {
	q := `
        INSERT INTO medical_records (id, identity_number, symptoms, medications, created_at)
        VALUES ($1, $2, $3, $4, $5);
    `
	_, err := d.db.DB().ExecContext(ctx, q, medicalrecord.ID, medicalrecord.IdentityNumber, medicalrecord.Symptoms, medicalrecord.Medications, medicalrecord.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
