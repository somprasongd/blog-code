package dto

type NewUserForm struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"omitempty,oneof='admin' 'user'"`
}

type UserId struct {
	ID string `validate:"required,uuid4"`
}

type UpdateUserPasswordForm struct {
	Password string `json:"password"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserResponses []*UserResponse
