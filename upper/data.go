package upper

type Service interface {
	Handle(srcport uint16, destport uint16, data *[]byte)
}

var (
	handlers = make(map[uint16]*Service)
)

func Recv(srcport uint16, destport uint16, data *[]byte) {
	if destport > 0 {
		h, ok := handlers[destport]
		if ok {
			(*h).Handle(srcport, destport, data)
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
