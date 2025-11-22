package websocket

func NewHub() *Hub {
	return &Hub{
		Rooms:       map[string]*Room{},
		RegisterU:   make(chan *Client),
		UnRegisterU: make(chan *Client),
		RegisterR:   make(chan *Room),
		UnRegisterR: make(chan *Room),
		RegisterBC:  make(chan *Msg),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.RegisterU:
			room := hub.Rooms[client.RoomID]
			if room == nil {
				room = &Room{
					ID:        client.RoomID,
					AdminID:   client.ID,
					Name:      "Room " + client.RoomID,
					MaxUsers:  100, // Дефолтный лимит
					Clients:   make(map[int64]*Client),
					Broadcast: make(chan *Msg),
				}
				hub.Rooms[client.RoomID] = room
			}

			if room.MaxUsers > 0 && len(room.Clients) >= room.MaxUsers {
				client.Conn.Close()
				continue
			}

			room.Clients[client.ID] = client
		case client := <-hub.UnRegisterU:
			room := hub.Rooms[client.RoomID]
			if room == nil {
				client.Conn.Close()
				continue
			}

			delete(room.Clients, client.ID)
			close(client.Send)
			client.Conn.Close()
		case room := <-hub.RegisterR:
			hub.Rooms[room.ID] = room
		case room := <-hub.UnRegisterR:
			delete(hub.Rooms, room.ID)
			close(room.Broadcast)
			for _, client := range room.Clients {
				client.Conn.Close()
			}
		case msg := <-hub.RegisterBC:
			room := hub.Rooms[msg.RoomID]
			if room == nil {
				continue
			}

			for _, client := range room.Clients {
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(room.Clients, client.ID)
				}
			}
		}
	}
}

func (hub *Hub) RegisterUser(client *Client) {
	hub.RegisterU <- client
}

func (hub *Hub) UnRegisterUser(client *Client) {
	hub.UnRegisterU <- client
}

func (hub *Hub) RegisterRoom(room *Room) {
	hub.RegisterR <- room
}

func (hub *Hub) UnRegisterRoom(room *Room) {
	hub.UnRegisterR <- room
}

func (hub *Hub) RegisterBreatcast(msg *Msg) {
	hub.RegisterBC <- msg
}

func (hub *Hub) GetRoom(room_id string) *Room {
	return hub.Rooms[room_id]
}
