package update

type Request struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"omitempty,min=2,max=500"`
}

type Response struct {
	Message string `json:"message"`
}
