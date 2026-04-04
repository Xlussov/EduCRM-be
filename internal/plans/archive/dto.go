package archive

type Request struct {
	Status string `json:"status" validate:"required,oneof=ACTIVE ARCHIVED"`
}

type Response struct {
	Message string `json:"message"`
}
