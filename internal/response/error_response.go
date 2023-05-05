package response

type errorResponse struct {
	Message string `json:"error"`
}

func ErrorResponse(message string) errorResponse {
	return errorResponse{Message: message}
}
