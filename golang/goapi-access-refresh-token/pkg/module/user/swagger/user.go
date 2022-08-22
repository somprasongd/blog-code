package swagger

type UserRepsonse struct {
	ID    string `json:"id" example:"bfbc2a69-9825-4a0e-a8d6-ffb985dc719c"`
	Email string `json:"email" example:"user@mail.com"`
	Role  string `json:"role" example:"user"`
}

type ListUserRepsonse []UserRepsonse

type UserSampleData struct {
	Data UserRepsonse `json:"user"`
}

type UserSampleListData struct {
	Data ListUserRepsonse `json:"users"`
}

type CreateUserFrom struct {
	// Required: true
	Email string `json:"email" example:"user@mail.com"`
	// Required: true
	Password string `json:"password" example:"password1234"`
}

type UpdateUserPasswordForm struct {
	// Required: true
	PasswordOld string `json:"password_old" example:"password1234"`
	// Required: true
	PasswordNew string `json:"password_new" example:"1234password"`
}

type ErrorDetailCreate struct {
	Target  string `json:"target" example:"email"`
	Message string `json:"message" example:"email field is required"`
}

type ErrCreateSampleData struct {
	Code    string              `json:"code" example:"422"`
	Message string              `json:"message" example:"invalid data see details"`
	Details []ErrorDetailCreate `json:"details"`
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
