package dto

import (
	"encoding/json"
	"time"
)

type RegisterForm struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginForm struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (u *UserInfo) ToPayloadMap() map[string]interface{} {
	// Convert to a json string

	data, err := json.Marshal(u)

	if err != nil {
		return nil
	}
	payload := map[string]interface{}{}
	// Convert to a map
	err = json.Unmarshal(data, &payload)

	if err != nil {
		return nil
	}

	return payload
}

type AuthResponse struct {
	User                  *UserInfo `json:"user,omitempty"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires"`
}

type UpdateProfileForm struct {
	PasswordOld string `json:"password_old"`
	PasswordNew string `json:"password_new"`
}

type LogoutForm struct {
	Token string `json:"refresh_token" validate:"required"`
}

type RefreshForm struct {
	Token string `json:"refresh_token" validate:"required"`
}
