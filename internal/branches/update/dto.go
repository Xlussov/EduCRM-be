package update

type Request struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	City    string `json:"city" validate:"required"`
}

type Response struct {
	Message string `json:"message"`
}
