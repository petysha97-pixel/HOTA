package webcoket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader

func Chat(w http.ResponseWriter, r *http.Request) {

	//кон - постоянное подключение к чату
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Ошибка подключения к чату", http.StatusBadRequest)
		return
	}

	for {
		kind, message, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if err := conn.WriteMessage(kind, message); err != nil {
			return
		}
	}

}
