package utils

import (
	"github.com/astaxie/beego/logs"
	"sync"
)

const (
	MAX_BYTES_LEN = 5120
)
var (
	bytesPool *sync.Pool
	make_cnt = 0
	pool_cnt = 0
)

func init() {
	bytesPool = &sync.Pool{
		New: func() interface{} {
			logs.Warn("Get ",MAX_BYTES_LEN, " bytes from New.")
			return make([]byte, MAX_BYTES_LEN)
		},
	}
}

func GetBytes(l int) []byte {
	bytes := bytesPool.Get().([]byte)
	//reset length before reuse
	bytes = bytes[:0]
	if cap(bytes) < l {
		//make_cnt++
		//logs.Info("Get ",l, " bytes from make. ",make_cnt)
		bytes = make([]byte, l)
		bytes = bytes[:0]
	}
	return bytes
}

func PutBytes(bytes *[]byte) {
	//logs.Info("Put ",len(*bytes), " bytes to pool.")
	bytesPool.Put(*bytes)
}

func ResetBytes(bytes *[]byte) {
	*bytes = (*bytes)[:0]
}