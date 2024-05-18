package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/citadel-corp/halosuster/internal/common/db"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByNIP(ctx context.Context, nip int) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context, req ListUserPayload) ([]*User, error)
	Update(ctx context.Context, user *User) error
	DeleteByID(ctx context.Context, id string) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

// Create implements Repository.
func (d *dbRepository) Create(ctx context.Context, user *User) error {
	createUserQuery := `
		INSERT INTO users (
			id, name, nip, user_type, hashed_password, identity_card_url
		) VALUES (
			$1, $2, $3, $4, $5, $6
		);
	`
	_, err := d.db.DB().ExecContext(ctx, createUserQuery, user.ID, user.Name, user.NIP, user.UserType, user.HashedPassword, user.IdentityCardURL)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return ErrNIPAlreadyExists
			default:
				return err
			}
		}
		return err
	}
	return nil
}

// GetByNIP implements Repository.
func (d *dbRepository) GetByNIP(ctx context.Context, nip int) (*User, error) {
	getUserQuery := `
		SELECT id, name, nip, user_type, hashed_password, identity_card_url, created_at
		FROM users
		WHERE nip = $1;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, nip)
	u := &User{}
	err := row.Scan(&u.ID, &u.Name, &u.NIP, &u.UserType, &u.HashedPassword, &u.IdentityCardURL, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (d *dbRepository) GetByID(ctx context.Context, id string) (*User, error) {
	getUserQuery := `
		SELECT id, name, nip, user_type, hashed_password, identity_card_url, created_at
		FROM users
		WHERE id = $1;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, id)
	u := &User{}
	err := row.Scan(&u.ID, &u.Name, &u.NIP, &u.UserType, &u.HashedPassword, &u.IdentityCardURL, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// List implements Repository.
func (d *dbRepository) List(ctx context.Context, req ListUserPayload) ([]*User, error) {
	paramNo := 1
	listQuery := "SELECT id, name, nip, user_type, hashed_password, identity_card_url, created_at FROM users WHERE "
	params := make([]interface{}, 0)
	if req.UserID != "" {
		listQuery += fmt.Sprintf("id = $%d AND ", paramNo)
		paramNo += 1
		params = append(params, req.UserID)
	}
	if req.Name != "" {
		listQuery += fmt.Sprintf("name LIKE '%%%s%%' AND ", req.Name)
		paramNo += 1
		params = append(params, req.Name)
	}
	if req.nipStr != "" {
		listQuery += fmt.Sprintf("nip LIKE '%%%s%%' AND ", req.nipStr)
		paramNo += 1
		params = append(params, req.nipStr)
	}
	switch req.RoleType {
	case ITType:
		listQuery += fmt.Sprintf("user_type = $%d AND ", paramNo)
		paramNo += 1
		params = append(params, "IT")
	case NurseType:
		listQuery += fmt.Sprintf("user_type = $%d AND ", paramNo)
		paramNo += 1
		params = append(params, "Nurse")
	}
	switch req.CreatedAtType {
	case Ascending:
		listQuery += fmt.Sprintf(" ORDER BY created_at ASC ")
	case Descending:
		listQuery += fmt.Sprintf(" ORDER BY created_at DESC ")
	}
	if strings.HasSuffix(listQuery, "AND ") {
		listQuery, _ = strings.CutSuffix(listQuery, "AND ")
	}
	listQuery += fmt.Sprintf(" LIMIT %d OFFSET %d;", req.Limit, req.Offset)
	if strings.Contains(listQuery, "WHERE  ORDER") {
		listQuery = strings.Replace(listQuery, "WHERE  ORDER", "ORDER", 1)
	}
	rows, err := d.db.DB().QueryContext(ctx, listQuery, params...)
	if err != nil {
		return nil, err
	}
	res := make([]*User, 0)
	for rows.Next() {
		u := &User{}
		err = rows.Scan(&u.ID, &u.Name, &u.NIP, &u.UserType, &u.HashedPassword, &u.IdentityCardURL, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

// Update implements Repository.
func (d *dbRepository) Update(ctx context.Context, user *User) error {
	q := `
        UPDATE users
        SET name = $1, nip = $2, user_type = $3, hashed_password = $4, identity_card_url = $5
        WHERE id = $6;
    `
	row, err := d.db.DB().ExecContext(ctx, q, user.Name, user.NIP, user.UserType, user.HashedPassword, user.IdentityCardURL, user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// DeleteByID implements Repository.
func (d *dbRepository) DeleteByID(ctx context.Context, id string) error {
	q := `
        DELETE FROM users
        WHERE id = $1;
    `
	row, err := d.db.DB().ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
