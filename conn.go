package socket

import (
	"bufio"
	"bytes"
	"context"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
)

// conn represents a client connection to a TCP server.
type conn struct {
	connID int64
	addr   string

	rawConn net.Conn
	reader  *bufio.Reader

	opts options

	recvMsg chan []byte
	sendMsg chan Message

	ctx    context.Context
	cancel context.CancelFunc

	err   error
	group *errgroup.Group
}

const BufferSize32 = 32

// NewConn returns a new client connection which has not started to
// serve requests yet.
func NewConn(connID int64, conn net.Conn, opt ...Option) *conn {
	var opts options
	for _, o := range opt {
		o(&opts)
	}

	if opts.bufferSize <= 0 {
		opts.bufferSize = BufferSize32
	}

	return newClientConnWithOptions(connID, conn, opts)
}

func newClientConnWithOptions(connID int64, c net.Conn, opts options) *conn {
	parentCtx, cancel := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(parentCtx)

	cc := &conn{
		connID: connID,
		addr:   c.RemoteAddr().String(),

		rawConn: c,
		reader:  bufio.NewReader(c),

		opts: opts,

		recvMsg: make(chan []byte, opts.bufferSize),
		sendMsg: make(chan Message, opts.bufferSize),

		ctx:    ctx,
		cancel: cancel,
		group:  group,
	}

	return cc
}

// ConnID returns the net ID of client connection.
func (c *conn) ConnID() int64 {
	return c.connID
}

// Start starts the client connection, creating go-routines for reading,
// writing and handlng.
func (c *conn) Start() {
	if c.opts.onConnect != nil {
		if err := c.opts.onConnect(c.rawConn); err != nil {
			return
		}
	}

	c.group.Go(c.readLoop)
	c.group.Go(c.writeLoop)
	c.group.Go(c.handleLoop)

	if err := c.group.Wait(); err != nil {
		c.err = err
		c.cancel()
		c.rawConn.Close()
	}
}

// Write writes a message to the client.
func (c *conn) Write(message Message) error {
	select {
	case <-c.ctx.Done():
		return errors.New("tcp closed")
	case c.sendMsg <- message:
		return nil
	default:
		return errors.New("chan blocked")
	}
}

// RemoteAddr returns the remote address of server connection.
func (c *conn) RemoteAddr() net.Addr {
	return c.rawConn.RemoteAddr()
}

// LocalAddr returns the local address of server connection.
func (c *conn) LocalAddr() net.Addr {
	return c.rawConn.LocalAddr()
}

/* readLoop() blocking read from connection, deserialize bytes into message,
then find corresponding handler, put it into channel */
func (c *conn) readLoop() error {
	data := make([]byte, c.opts.bufferSize)

	for {
		select {
		case <-c.ctx.Done():
			return nil
		default:
			n, err := c.reader.Read(data)
			if err != nil {
				if c.opts.onError != nil {
					c.opts.onError(errors.Wrap(ErrNetRead, err.Error()))
				}
				return err
			}

			select {
			case c.recvMsg <- data[0:n]:
				data = data[:0]
			default:
				log.Println("recv chan blocked")
			}
		}
	}
}

func (c *conn) handleLoop() error {
	for {
		select {
		case <-c.ctx.Done():
			return nil
		case data := <-c.recvMsg:
			message, err := c.opts.codec.Decode(bytes.NewBuffer(data))
			if err != nil {
				if c.opts.onError != nil {
					c.opts.onError(errors.Wrap(ErrHandlerMsg, err.Error()))
				}
				return err
			}

			if c.opts.onMessage != nil {
				c.opts.onMessage(message)
			}
		}

	}
}

/* writeLoop() receive message from channel, serialize it into bytes,
then blocking write into connection */
func (c *conn) writeLoop() error {
	for {
		select {
		case <-c.ctx.Done():
			return nil
		case msg := <-c.sendMsg:
			data, err := c.opts.codec.Encode(msg)
			if err != nil {
				return err
			}

			if _, err = c.rawConn.Write(data); err != nil {
				if c.opts.onError != nil {
					c.opts.onError(errors.Wrap(ErrNetWrite, err.Error()))
				}

				return err
			}
		}
	}
}
