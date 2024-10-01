package models

type ChirpRequest struct {
	Body string `json:"body"`
}

type UserPostRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChirpPostRequest struct {
	Body   string `json:"body"`
	UserId string `json:"user_id"`
}

type LoginPostRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
