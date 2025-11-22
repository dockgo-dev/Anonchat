package websocket

import "github.com/gorilla/websocket"

type (
	Hub struct {
		Rooms map[string]*Room

		RegisterU   chan *Client
		UnRegisterU chan *Client
		RegisterR   chan *Room
		UnRegisterR chan *Room
		RegisterBC  chan *Msg
	}

	Room struct {
		ID        string
		AdminID   int64
		Name      string
		MaxUsers  int
		Clients   map[int64]*Client
		Broadcast chan *Msg
	}

	Client struct {
		ID     int64
		Login  string
		Conn   *websocket.Conn
		Send   chan *Msg
		RoomID string
	}

	Msg struct {
		ClientID   int64  `json:"client_id"`
		ClientName string `json:"client_name"`
		RoomID     string `json:"room_id"`
		Message    []byte `json:"message"`
	}
)
