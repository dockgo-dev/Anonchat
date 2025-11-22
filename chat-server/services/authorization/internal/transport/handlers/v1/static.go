package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gox7/notify/services/authorization/models"
)

func HandleStatus(ctx *gin.Context) {
	ctx.JSON(200, ResponseJSON(
		"success", "authorization is work",
		nil,
	))
}

func ResponseJSON(status string, message string, data any) models.DefaultResponse {
	return models.DefaultResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
