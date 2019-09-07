package socket

import "net"

type options struct {
	codec Codec

	onConnect func(net.Conn) error
	onMessage func(Message)
	onError   func(error) // 只要网络出错， 或者encode decode 出错， 都会收到该通知

	bufferSize int // size of buffered channel
}

// Option sets server options.
type Option func(*options)

// CustomCodecOption returns a Option that will apply a custom Codec.
func CustomCodecOption(codec Codec) Option {
	return func(o *options) {
		o.codec = codec
	}
}

// BufferSizeOption returns a Option that is the size of buffered channel,
// for example an indicator of BufferSize32 means a size of 256.
func BufferSizeOption(indicator int) Option {
	return func(o *options) {
		o.bufferSize = indicator
	}
}

// OnConnectOption returns a Option that will set callback to call when new
// client connected.
func OnConnectOption(cb func(net.Conn) error) Option {
	return func(o *options) {
		o.onConnect = cb
	}
}

// OnMessageOption returns a Option that will set callback to call when new
// message arrived.
func OnMessageOption(cb func(Message)) Option {
	return func(o *options) {
		o.onMessage = cb
	}
}

// OnErrorOption returns a Option that will set callback to call when error
// occurs.
func OnErrorOption(cb func(error)) Option {
	return func(o *options) {
		o.onError = cb
	}
}
