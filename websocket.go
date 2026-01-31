package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // dev only
	},
}

// Store client connections (for example)
var clients = make(map[*websocket.Conn]bool)

// --- Function to send a message to a specific client ---
func sendMessage(ws *websocket.Conn, msg string) error {
	return ws.WriteMessage(websocket.TextMessage, []byte(msg))
}

// --- Handler ---
func handleWSClient(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer ws.Close()

	// Register client
	clients[ws] = true
	defer delete(clients, ws)

	// Send initial message
	if err := sendMessage(ws, "Hallo"); err != nil {
		log.Println("write error:", err)
		return
	}

	// --- Read loop keeps connection alive ---
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println("client disconnected:", err)
			return
		}
	}
}

func pingAllClients() {
	for ws := range clients {
		if err := sendMessage(ws, "ping"); err != nil {
			log.Println("error sending ping:", err)
			ws.Close()
			delete(clients, ws)
		}
	}
}
