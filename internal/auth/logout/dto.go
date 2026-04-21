package logout

type Request struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type Response struct {
	Message string `json:"message"`
}
