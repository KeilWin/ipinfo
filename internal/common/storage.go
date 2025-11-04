package common

type Storage interface {
	StartUp() error
	ShutDown() error
}
