package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mbyd916/rtchat/pkg/websocket"
)

func serveWs(p *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("websocket endpoint hit")

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
		return
	}

	c := &websocket.Client{
		Conn: conn,
		Pool: p,
	}

	p.Register <- c
	c.Read()
}

func setupRoutes() {
	p := websocket.NewPool()
	go p.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(p, w, r)
	})
}

func main() {
        log.Println("start ChatApp_v0.0.1");
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
