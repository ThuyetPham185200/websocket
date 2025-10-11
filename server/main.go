package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"websocketserver/websockets"

	"github.com/gorilla/websocket"
)

var (
	mtx sync.RWMutex
)

func main() {
	// Create the WebSocket server
	wsServer := websockets.NewWebSocketServer("localhost:9000")

	clients := make([]*websockets.WSClientInfo, 0)

	// Register message handler
	wsServer.OnNewConnection(func(wsclient *websockets.WSClientInfo) {
		mtx.Lock()
		defer mtx.Unlock()
		clients = append(clients, wsclient)
		fmt.Printf("Client connected: %s\n", wsclient.Addr)
	})

	// Periodic broadcast loop (1Hz)
	go func() {
		for {
			time.Sleep(1 * time.Second)

			mtx.RLock()
			active := make([]*websockets.WSClientInfo, 0, len(clients))

			for _, wsc := range clients {
				if wsc.IsConnected() {
					data := "hello " + wsc.Addr + " at " + time.Now().Local().String()
					err := wsc.Send(websocket.TextMessage, []byte(data))
					if err != nil {
						fmt.Printf("send failed to %s: %v\n", wsc.Addr, err)
					}
					active = append(active, wsc)
				} else {
					fmt.Printf("removing disconnected client: %s\n", wsc.Addr)
					wsc.Close()
				}
			}

			mtx.RUnlock()

			// Update client list with active ones only
			mtx.Lock()
			clients = active
			mtx.Unlock()
		}
	}()

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
