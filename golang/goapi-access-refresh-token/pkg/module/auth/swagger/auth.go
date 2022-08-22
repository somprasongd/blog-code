package swagger

type RegisterForm struct {
	// Required: true
	Email string `json:"email" example:"user@mail.com"`
	// Required: true
	Password string `json:"password" example:"password1234"`
}

type LoginForm struct {
	// Required: true
	Email string `json:"email" example:"user@mail.com"`
	// Required: true
	Password string `json:"password" example:"password1234"`
}

type UserInfo struct {
	ID    string `json:"id" example:"bfbc2a69-9825-4a0e-a8d6-ffb985dc719c"`
	Email string `json:"email" example:"user@mail.com"`
	Role  string `json:"role" example:"user"`
}

type AuthResponse struct {
	User  UserInfo `json:"user"`
	Token string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAbWFpbC5jb20iLCJpYXQiOjE2NTk0MzI5NTYsInN1YiI6Ijk2YWUzNWM0LTE0Y2ItNDAzMy1iYTMwLTVkYTBmNjA2NjFiNCJ9.spR28QwRVbmOjJPu6iwRhA19jOpxYtgpRRsiaNWGTYk"`
}

type UserSampleData struct {
	Data UserInfo `json:"user"`
}

type AuthSampleData struct {
	Data AuthResponse `json:"auth"`
}

type ErrorDetailRegister struct {
	Target  string `json:"target" example:"email"`
	Message string `json:"message" example:"email field is required"`
}

type ErrRegisterSampleData struct {
	Code    string                `json:"code" example:"422"`
	Message string                `json:"message" example:"invalid data see details"`
	Details []ErrorDetailRegister `json:"details"`
}

type ErrorDetailLogin struct {
	Target  string `json:"target" example:"password"`
	Message string `json:"message" example:"password field is required"`
}

type ErrLoginSampleData struct {
	Code    string             `json:"code" example:"422"`
	Message string             `json:"message" example:"invalid data see details"`
	Details []ErrorDetailLogin `json:"details"`
}

type ErrorDetailUpdate struct {
	Target  string `json:"target" example:"password_old"`
	Message string `json:"message" example:"password_old field is required"`
}

type ErrUpdateSampleData struct {
	Code    string              `json:"code" example:"422"`
	Message string              `json:"message" example:"invalid data see details"`
	Details []ErrorDetailUpdate `json:"details"`
}
