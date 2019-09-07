package socket

import (
	"sync"
)

var pool = &sync.Pool{
	New: func() interface{} {
		b := make([]byte, 2000)
		return b
	},
}

func Get() []byte {
	return pool.Get().([]byte)
}

func Put(buf []byte) {
	// reset length to 0
	pool.Put(buf[:0])
}
