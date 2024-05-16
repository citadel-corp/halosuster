package user

import "time"

type User struct {
	ID              string
	NIP             int
	Name            string
	UserType        UserType
	IdentityCardURL *string
	HashedPassword  *string
	CreatedAt       time.Time
}

type UserType string

const (
	IT    UserType = "IT"
	Nurse UserType = "Nurse"
)
