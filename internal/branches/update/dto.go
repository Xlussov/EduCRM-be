package update

type Request struct {
	Name    string `json:"name" validate:"required,min=2,max=50"`
	Address string `json:"address" validate:"required,min=2,max=255"`
	City    string `json:"city" validate:"required,min=2,max=50"`
}

type Response struct {
	Message string `json:"message"`
}
