package httpx

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorBody{Code: code, Message: message},
	}
}

func WriteError(c *gin.Context, status int, code, message string) {
	c.JSON(status, NewErrorResponse(code, message))
}
