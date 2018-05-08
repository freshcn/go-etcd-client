package etcd

import (
	"github.com/freshcn/log"
)

type closer interface {
	// Close 关闭连接
	Close()
}

var closerArr = []closer{}

// Close 关闭ETCD
func Close() {
	closeETCD = true
	if len(closerArr) > 0 {
		for _, v := range closerArr {
			v.Close()
		}
		closerArr = []closer{}
		log.Warring("etcd closed.")
	}
}

// addCloser 添加关闭处理程序
func addCloser(handler closer) {
	closerArr = append(closerArr, handler)
}
