package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader upgrades a http connection to Websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

// reader reads from websocket connection and writes the same message back in loop
func reader(conn *websocket.Conn) {
	for {

		msgType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(fmt.Errorf("read err: %w", err))
			return
		}

		log.Printf("msg received: %s\n", p)

		err = conn.WriteMessage(msgType, p)
		if err != nil {
			log.Println(fmt.Errorf("write error: %w", err))
			return
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	log.Printf("request from host: %v", r.Host)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(fmt.Errorf("upgrade err: %w", err))
		return
	}

	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "simple server")
	})

	http.HandleFunc("/ws", serveWs)
}

func main() {
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
