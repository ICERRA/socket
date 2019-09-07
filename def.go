package socket

import "github.com/pkg/errors"

var (
	ErrNetRead    = errors.New("net read error")
	ErrNetWrite   = errors.New("net write error")
	ErrHandlerMsg = errors.New("handle message error")
)
