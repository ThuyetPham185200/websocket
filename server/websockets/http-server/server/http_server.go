package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type HttpServer struct {
	BaseServerProcessor
	httpServer *http.Server
}

func NewHttpServer(addr string, handler http.Handler) *HttpServer {
	s := &HttpServer{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
	s.Init(s) // 🔑 rất quan trọng: gắn HttpServer vào BaseServerProcessor
	return s
}

// RunningTask implement từ ServerProcessor
func (s *HttpServer) RunningTask() error {
	log.Printf("🌐 HTTP Server running at %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Override Stop để shutdown http.Server
func (s *HttpServer) Stop() error {
	log.Println("⏹️ Shutting down HTTP server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
