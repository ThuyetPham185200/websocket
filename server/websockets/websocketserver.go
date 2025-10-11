package websockets

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"websocketserver/websockets/http-server/server"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Client infor
type WSClientInfo struct {
	Addr string
	conn *websocket.Conn
}

func NewWSClientInfo(addr string, con *websocket.Conn) *WSClientInfo {
	return &WSClientInfo{
		Addr: addr,
		conn: con,
	}
}

func (ws *WSClientInfo) Send(messageType int, data []byte) error {
	return ws.conn.WriteMessage(messageType, data)
}

func (ws *WSClientInfo) Recv() (messageType int, p []byte, err error) {
	return ws.conn.ReadMessage()
}

func (ws *WSClientInfo) Close() error {
	return ws.conn.Close()
}

func (ws *WSClientInfo) IsConnected() bool {
	deadline := time.Now().Add(time.Second)
	err := ws.conn.WriteControl(websocket.PingMessage, []byte{}, deadline)
	return err == nil || !strings.Contains(err.Error(), "use of closed network connection")
}

// MessageHandler defines a custom callback for incoming messages
type MessageHandler func(wsclient *WSClientInfo)

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
		conn := value.(*WSClientInfo)
		conn.Close()
		return true
	})

	if err := s.httpserver.Stop(); err != nil {
		return err
	}
	return nil
}

// RegisterHandler allows external modules to handle incoming messages
func (s *WebSocketServer) OnNewConnection(handler MessageHandler) {
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
	wsc := NewWSClientInfo(addr, ws)
	s.connections.Store(addr, wsc)
	fmt.Printf("New client connected: %s\n", addr)

	go s.newConnection(wsc)
}

// Listen for messages from a specific client
func (s *WebSocketServer) newConnection(wsclient *WSClientInfo) {
	if s.handler != nil {
		s.handler(wsclient)
	}
}
