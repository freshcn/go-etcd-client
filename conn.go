package etcd

import (
	"context"
	"time"

	"github.com/freshcn/log"

	"github.com/coreos/etcd/clientv3"
)

var (
	// etcd 连接
	conn *clientv3.Client

	// RetryNum 重试次数
	RetryNum = 5

	// closeETCD 关闭etcd
	closeETCD bool
)

// getConn 获取etcd连接
// retry 是否为重试连接
func getConn(retry ...bool) *clientv3.Client {
	var (
		err error
	)

	if conn != nil {
		return conn
	}

	conn, err = clientv3.New(config)

	if err != nil {
		if len(retry) == 0 { // 当不为重试连接时
			for i := 0; i < RetryNum; i++ {
				conn = getConn(true)
				if conn != nil {
					break
				}
			}

			if conn == nil {
				log.Panic(err)
			}
		}
	}

	return conn
}

// getCtx 获取上下文
func getCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*5)
}

// stringSliceEqual 对比string的slice是否一样
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
