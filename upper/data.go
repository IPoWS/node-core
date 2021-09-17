package upper

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type Service interface {
	Handle(srcport uint16, destport uint16, data *[]byte)
}

var (
	handlers = make(map[uint16]*Service)
	mu       sync.RWMutex
)

func Recv(srcport uint16, destport uint16, data *[]byte) {
	if destport > 0 {
		mu.RLock()
		h, ok := handlers[destport]
		mu.RUnlock()
		if ok {
			logrus.Infoln("[upper] handle data from", srcport, "to", destport, "with len", len(*data), "bytes.")
			(*h).Handle(srcport, destport, data)
		}
	}
}

func Register(port uint16, service *Service) bool {
	logrus.Debugln("[upper] reg start.")
	mu.RLock()
	_, ok := handlers[port]
	mu.RUnlock()
	ina := !ok
	logrus.Debugln("[upper] reg read ina:", ina)
	if ina {
		mu.Lock()
		handlers[port] = service
		mu.Unlock()
		logrus.Infoln("[upper] register service on port", port, "succeed.")
	}
	return ina
}

func Remove(port uint16) bool {
	mu.RLock()
	_, ok := handlers[port]
	mu.RUnlock()
	if ok {
		mu.Lock()
		delete(handlers, port)
		mu.Unlock()
		logrus.Infoln("[upper] del service on port", port, "succeed.")
	}
	return ok
}
