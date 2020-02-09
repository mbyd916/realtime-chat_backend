package websocket

import "log"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (p *Pool) Start() {
	for {
		select {
		case c := <-p.Register:
			p.Clients[c] = true
			log.Println("size of connection pool: ", len(p.Clients))
			for c := range p.Clients {
				log.Printf("[notify register] client: %+v", c)
				c.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			}
		case c := <-p.Unregister:
			delete(p.Clients, c)
			log.Println("size of connection pool: ", len(p.Clients))
			for c := range p.Clients {
				log.Printf("[notify unregister] client: %+v", c)
				c.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			}
		case msg := <-p.Broadcast:
			log.Println("sending message to all clients in pool")
			for c := range p.Clients {
				if err := c.Conn.WriteJSON(msg); err != nil {
					log.Printf("err: %v", err)
					return
				}
			}

		}
	}

}
