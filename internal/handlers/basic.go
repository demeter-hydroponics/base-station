package handlers

import (
    "net/http"
	"github.com/gorilla/websocket"
)

func check_origin(r *http.Request) bool {
    return true
}

var upgrader = websocket.Upgrader{CheckOrigin: check_origin} // use default options
