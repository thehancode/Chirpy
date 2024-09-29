package models

import "time"

type ValidResponse struct {
	Valid bool `json:"valid"`
}

type ValidatedChirpResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type CreateUserResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type CreateChirpResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    string    `json:"user_id"`
}
