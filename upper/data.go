package upper

type Service interface {
	Handle(data *[]byte)
}

var (
	handlers = make(map[uint16]*Service)
)

func Recv(port uint16, data *[]byte) {
	if port > 0 {
		h, ok := handlers[port]
		if ok {
			(*h).Handle(data)
		}
	}
}

func Register(port uint16, service *Service) bool {
	_, ok := handlers[port]
	ina := !ok
	if ina {
		handlers[port] = service
	}
	return ina
}

func Remove(port uint16) bool {
	_, ok := handlers[port]
	if ok {
		delete(handlers, port)
	}
	return ok
}
