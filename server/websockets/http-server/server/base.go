package server

import (
	"context"
	"log"
	"time"
)

// BaseServerProcessor implement sẵn Start/Stop/Restart
// để các server embed lại
type BaseServerProcessor struct {
	processor ServerProcessor
	cancel    context.CancelFunc
}

func (b *BaseServerProcessor) Init(p ServerProcessor) {
	b.processor = p
}

func (b *BaseServerProcessor) Start() error {
	log.Println("Starting server...")

	// chạy task trong goroutine riêng
	go func() {
		if err := b.processor.RunningTask(); err != nil {
			log.Printf("Server stopped with error: %v", err)
		}
	}()
	log.Println("Started server!!")
	return nil
}

func (b *BaseServerProcessor) Stop() error {
	log.Println("Stopping server...")
	// Ở đây base class không biết chi tiết stop,
	// có thể override trong HttpServer nếu cần shutdown http.Server
	return nil
}

func (b *BaseServerProcessor) Restart() error {
	log.Println("Restarting server...")
	if err := b.Stop(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return b.Start()
}
