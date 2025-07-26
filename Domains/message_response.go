package domains

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}
