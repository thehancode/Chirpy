package models

type ChirpRequest struct {
	Body string `json:"body"`
}

type UserPostRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPutRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChirpPostRequest struct {
	Body string `json:"body"`
}

type LoginPostRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds *int64 `json:"expires_in_seconds,omitempty"`
}

func (l *LoginPostRequest) SetDefaults() {
	if l.ExpiresInSeconds == nil {
		defaultExpiration := int64(3600)
		l.ExpiresInSeconds = &defaultExpiration
	}
}
