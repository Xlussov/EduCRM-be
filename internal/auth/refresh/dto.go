package refresh

type Request struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
