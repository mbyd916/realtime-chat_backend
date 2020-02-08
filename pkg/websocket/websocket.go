package websocket

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader upgrades a http connection to Websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, nil)
}

// Reader reads from websocket connection and writes the same message back in loop
func Reader(conn *websocket.Conn) {
	for {

		msgType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(fmt.Errorf("read err: %w", err))
			return
		}

		log.Printf("msg received: %s", p)

		err = conn.WriteMessage(msgType, p)
		if err != nil {
			log.Println(fmt.Errorf("write error: %w", err))
			return
		}
	}
}

func Writer(conn *websocket.Conn) {
	for {

		log.Println("Sending...")
		msgType, r, err := conn.NextReader()
		if err != nil {

			log.Println(fmt.Errorf("next reader err: %w", err))
			return
		}

		w, err := conn.NextWriter(msgType)
		if err != nil {

			log.Println(fmt.Errorf("next writer err: %w", err))
			return
		}

		if _, err := io.Copy(w, r); err != nil {
			log.Println(fmt.Errorf("copy err: %w", err))
			return
		}

		log.Println("msg copied")
		if err := w.Close(); err != nil {
			log.Println(fmt.Errorf("close err: %w", err))
		}
	}

}
