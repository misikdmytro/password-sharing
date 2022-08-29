package model

type ErrorResponse struct {
	Message string `json:"message"`
}

type LinkResponse struct {
	Link string `json:"link"`
}

type PasswordResponse struct {
	Password string `json:"password"`
}
