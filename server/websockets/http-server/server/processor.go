package server

// ServerProcessor định nghĩa interface chung
type ServerProcessor interface {
	Start() error
	Stop() error
	Restart() error
	RunningTask() error
}
