package model

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type LinkResponse struct {
	Link string `json:"link"`
}

type PasswordResponse struct {
	Password string `json:"password"`
}

type HealthResponse struct {
	Healthy bool   `json:"healthy"`
	Reason  string `json:"reason"`
}
