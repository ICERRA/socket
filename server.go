package socket

import (
	"net"
	"runtime"
)

type Handler interface {
	Handle(conn net.Conn)
}

type server struct {
	listener net.Listener
}

func New(addr net.Addr) (*server, error) {
	listener, err := net.Listen(addr.Network(), addr.String())
	if err != nil {
		return nil, err
	}

	return &server{
		listener: listener,
	}, nil
}

func (s *server) Serve(handler Handler) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			runtime.Gosched()
			continue
		}

		go handler.Handle(conn)
	}
}

func (s *server) Close() error {
	return s.listener.Close()
}
