package main

import (
	"fmt"
	"net/http"

	"github.com/mbyd916/rtchat/pkg/websocket"
)

func serveWs(w http.ResponseWriter, r *http.Request) {

	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
		return
	}

	go websocket.Writer(ws)

	websocket.Reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/ws", serveWs)
}

func main() {
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
