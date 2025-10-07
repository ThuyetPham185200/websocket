package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"websocketserver/websockets"

	"github.com/gorilla/websocket"
)

func main() {
	// Create the WebSocket server
	wsServer := websockets.NewWebSocketServer("localhost:9000")

	// Register message handler
	wsServer.RegisterHandler(func(conn *websocket.Conn, msgType int, data []byte) {
		fmt.Println("Received:", string(data))
		conn.WriteMessage(msgType, []byte("OK: "+string(data)))
	})

	// Channel to listen for Ctrl+C or kill signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run server in separate goroutine
	go func() {
		if err := wsServer.Open(); err != nil {
			fmt.Println("WebSocket server error:", err)
		}
	}()

	fmt.Println("WebSocket server started on ws://localhost:9000/ws")

	select {
	case <-stop:
		fmt.Println("\nShutting down WebSocket server...")
		wsServer.Close()
	}

	// Optional: wait a bit for cleanup
	time.Sleep(2 * time.Second)
	fmt.Println("Shutdown complete.")
}
