package models

type ChirpRequest struct {
	Body string `json:"body"`
}

type UserPostRequest struct {
	Email string `json:"email"`
}

type ChirpPostRequest struct {
	Body   string `json:"body"`
	UserId string `json:"user_id"`
}
