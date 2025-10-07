package websockets

import (
	"fmt"
	"net/http"
	"sync"
	"websocketserver/websockets/http-server/server"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// MessageHandler defines a custom callback for incoming messages
type MessageHandler func(conn *websocket.Conn, msgType int, data []byte)

// WebSocketServer manages multiple websocket clients
type WebSocketServer struct {
	httpserver *server.HttpServer
	upgrader   websocket.Upgrader

	connections sync.Map // map[string]*websocket.Conn
	handler     MessageHandler
}

// Constructor
func NewWebSocketServer(address string) *WebSocketServer {
	s := &WebSocketServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true }, // allow any origin
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/ws", s.handleUpgrade)
	s.httpserver = server.NewHttpServer(address, router)
	return s
}

// Start HTTP + WS server
func (s *WebSocketServer) Open() error {
	fmt.Println("Starting WebSocket server...")
	if err := s.httpserver.Start(); err != nil {
		return err
	}
	return nil
}

// Stop HTTP server and close all clients
func (s *WebSocketServer) Close() error {
	fmt.Println("Stopping WebSocket server...")

	s.connections.Range(func(key, value any) bool {
		conn := value.(*websocket.Conn)
		conn.Close()
		return true
	})

	if err := s.httpserver.Stop(); err != nil {
		return err
	}
	return nil
}

// RegisterHandler allows external modules to handle incoming messages
func (s *WebSocketServer) RegisterHandler(handler MessageHandler) {
	s.handler = handler
}

// Handle new WS connection
func (s *WebSocketServer) handleUpgrade(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	addr := ws.RemoteAddr().String()
	s.connections.Store(addr, ws)
	fmt.Printf("New client connected: %s\n", addr)

	go s.listenClient(addr, ws)
}

// Listen for messages from a specific client
func (s *WebSocketServer) listenClient(addr string, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		s.connections.Delete(addr)
		fmt.Printf("Client disconnected: %s\n", addr)
	}()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}

		if s.handler != nil {
			s.handler(conn, msgType, msg)
		} else {
			// default echo if no handler
			conn.WriteMessage(msgType, msg)
		}
	}
}
