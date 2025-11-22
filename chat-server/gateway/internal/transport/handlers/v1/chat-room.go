package v1

import (
	"mew-gateway/internal/app"
	lib "mew-gateway/internal/libs"
	ws "mew-gateway/internal/websocket"
	"mew-gateway/pkg/tokens"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebsocketRoom(config *lib.Config, hub *ws.Hub) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomId := ctx.Query("room")
		access := ctx.Query("token")

		if roomId == "" {
			ctx.JSON(400, ResponseJSON(
				"error", "room parameter is required", nil,
			))
			return
		}

		connect, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			return
		}

		userData, err := app.Validate(config, access)
		if err != nil {
			connect.WriteMessage(websocket.TextMessage, []byte("error.token"))
			connect.Close()
			return
		}

		send := make(chan *ws.Msg)
		client := &ws.Client{
			ID:     userData.UserID,
			Login:  userData.Login,
			Conn:   connect,
			Send:   send,
			RoomID: roomId,
		}

		hub.RegisterUser(client)
		go client.Writing(hub)
		go client.Reading()
	}
}

func HandleCreateRoom(hub *ws.Hub) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RegisterRoomRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, ResponseJSON(
				"error", "invalid request body", nil,
			))
			return
		}

		user_id := ctx.GetInt64("user.id")
		if user_id == 0 {
			ctx.JSON(401, ResponseJSON(
				"error", "unauthorization", nil,
			))
			return
		}

		room_id := tokens.Generate(8)
		room_chan := make(chan *ws.Msg)
		room := &ws.Room{
			ID:        room_id,
			AdminID:   user_id,
			Name:      req.Name,
			MaxUsers:  req.MaxUser,
			Clients:   make(map[int64]*ws.Client),
			Broadcast: room_chan,
		}

		hub.RegisterRoom(room)
		ctx.JSON(201, ResponseJSON(
			"success", "create room", gin.H{"room_id": room_id},
		))
	}
}

func HandleRemoveRoom(hub *ws.Hub) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RemoveRoomRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, ResponseJSON(
				"error", "invalid request body", nil,
			))
			return
		}

		user_id := ctx.GetInt64("user.id")
		if user_id == 0 {
			ctx.JSON(401, ResponseJSON(
				"error", "unauthorization", nil,
			))
			return
		}

		room := hub.GetRoom(req.RoomID)
		if room == nil {
			ctx.JSON(404, ResponseJSON(
				"error", "room not found", nil,
			))
			return
		}

		if room.AdminID != user_id {
			ctx.JSON(403, ResponseJSON(
				"error", "none access for remove room", nil,
			))
			return
		}

		hub.UnRegisterRoom(room)
		ctx.JSON(200, ResponseJSON(
			"success", "remove room", nil,
		))
	}
}
