package response

type successResponse struct {
	Message string `json:"message"`
}

func SuccessResponse(message string) successResponse {
	return successResponse{Message: message}
}
