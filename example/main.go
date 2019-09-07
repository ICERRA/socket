package main

import (
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/Zereker/socket"
)

type codec struct {
}

func (c *codec) Decode(io.Reader) (socket.Message, error) {
	panic("implement me")
}

func (c *codec) Encode(socket.Message) ([]byte, error) {
	panic("implement me")
}

type handler struct {
	connID int64

	logger socket.Logger
}

func newHandler(connID int64) *handler {
	return &handler{connID: connID}
}

func (h *handler) Handle(conn net.Conn) {
	connID := atomic.AddInt64(&h.connID, 1)
	socket.NewConn(connID, conn,
		socket.CustomCodecOption(new(codec)),
		socket.OnConnectOption(func(conn net.Conn) error {
			return nil
		}),
		socket.OnMessageOption(func(message socket.Message) {
			h.logger.Info("recv message")
		}),
		socket.OnErrorOption(func(err error) {
			h.logger.Warn(err)
		}),
	)
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:12345")
	if err != nil {
		panic(err)
	}

	server, err := socket.New(addr)
	if err != nil {
		log.Println(err)
	}

	server.Serve(newHandler(time.Now().Unix()))
}
