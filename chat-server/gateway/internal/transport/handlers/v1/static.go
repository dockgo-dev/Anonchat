package v1

import "github.com/gin-gonic/gin"

type (
	DefaultResponse struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

type (
	RegisterRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	RegisterRoomRequest struct {
		Name    string `json:"name"`
		MaxUser int    `json:"max_user"`
	}
	RemoveRoomRequest struct {
		RoomID string `json:"room_id"`
	}

	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	TokenRequest struct {
		Token string `json:"token"`
	}
	UserData struct {
		Id    int64  `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
	}
)

func HandleStatus(ctx *gin.Context) {
	ctx.JSON(200, ResponseJSON(
		"success", "Mew, is work!", nil,
	))
}

func ResponseJSON(status string, message string, data interface{}) *DefaultResponse {
	return &DefaultResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
