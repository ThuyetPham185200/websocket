package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all connections
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	for {
		// Read message from browser
		content := " Cai dit me may!"
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			break
		}
		msg = append(msg, []byte(content)...)
		fmt.Printf("Thang nao goi bo: %s\n", msg)

		// Write message back to browser
		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Println("write error:", err)
			break
		}
		content = " Dit me may phat nua!"
		msg = append(msg, []byte(content)...)
		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Println("write error:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
