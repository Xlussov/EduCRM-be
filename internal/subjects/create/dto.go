package create

type Request struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Response struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
