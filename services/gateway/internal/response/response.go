package response

type HealthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
