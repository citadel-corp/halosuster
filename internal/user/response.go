package user

import "time"

type UserAuthResponse struct {
	UserID      string  `json:"userId"`
	NIP         int     `json:"nip"`
	Name        string  `json:"name"`
	AccessToken *string `json:"accessToken"`
}

type UserResponse struct {
	UserID    string    `json:"userId"`
	NIP       int       `json:"nip"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}
