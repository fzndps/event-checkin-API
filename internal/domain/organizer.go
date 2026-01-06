package domain

import "time"

type Organizer struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Name         string `json:"name" binding:"required,min=3"`
	PasswordHash string `json:"password" binding:"required,min=3"`
}

type LoginRequest struct {
	Email        string `json:"email" binding:"required,email"`
	PasswordHash string `json:"password" binding:"required,min=3"`
}

type OrganizerLoginResponse struct {
	Token     string     `json:"token"`
	Organizer *Organizer `json:"organizer"`
}

type AuthResponse struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

func (o *Organizer) ToResponse() *AuthResponse {
	return &AuthResponse{
		ID:           o.ID,
		Email:        o.Email,
		Name:         o.Name,
		PasswordHash: o.PasswordHash,
		CreatedAt:    o.CreatedAt,
	}
}
