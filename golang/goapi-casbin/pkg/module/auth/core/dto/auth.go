package dto

import "time"

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

type TokenInfo struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
}

type AuthResponse struct {
	User  UserInfo `json:"user"`
	Token string   `json:"token"`
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
