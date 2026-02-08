package dto

type RegisterUserPayload struct {
	Name     string `validate:"required"`
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type LoginUserPayload struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type UserWithTokenResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
