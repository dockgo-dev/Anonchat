package websocket

func (c *Client) Reading() {
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			return
		}
	}
}

func (c *Client) Writing(hub *Hub) {
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			hub.UnRegisterUser(c)
			return
		}

		hub.RegisterBC <- &Msg{
			ClientID:   c.ID,
			ClientName: c.Login,
			RoomID:     c.RoomID,
			Message:    msg,
		}
	}
}
