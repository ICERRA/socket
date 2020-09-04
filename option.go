package socket

import (
	"time"
)

type options struct {
	codec Codec

	onMessage func(message Message) error
	onError   func(error) bool // 出错的时候, 是否断开连接, 由调用方自由定制

	bufferSize    int // size of buffered channel
	maxReadLength int
	heartbeat     time.Duration
}

// Option sets server options.
type Option func(*options)

// CustomCodecOption returns a Option that will apply a custom Codec.
func CustomCodecOption(codec Codec) Option {
	return func(o *options) {
		o.codec = codec
	}
}

// HeartbeatOption returns a Option that is the size of buffered channel,
// for example an indicator of defaultBufferSize32 means a size of 256.
func BufferSizeOption(indicator int) Option {
	return func(o *options) {
		o.bufferSize = indicator
	}
}

// HeartbeatOption returns a Option that is the size of buffered channel,
// for example an indicator of defaultBufferSize32 means a size of 256.
func HeartbeatOption(heartbeat time.Duration) Option {
	return func(o *options) {
		o.heartbeat = heartbeat
	}
}

// MessageMaxSize returns a Option that will set buffer to receive message when new
// []byte arrived.
func MessageMaxSize(size int) Option {
	return func(o *options) {
		o.maxReadLength = size
	}
}

// OnErrorOption returns a Option that will set callback to call when error
// occurs.
func OnErrorOption(cb func(error) bool) Option {
	return func(o *options) {
		o.onError = cb
	}
}

// OnMessageOption returns a Option that will set callback to call when new
// message arrived.
func OnMessageOption(cb func(Message) error) Option {
	return func(o *options) {
		o.onMessage = cb
	}
}
