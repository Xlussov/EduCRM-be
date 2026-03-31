package update

type Request struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Response struct {
	Message string `json:"message"`
}
