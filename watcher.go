package etcd

import (
	"errors"

	"github.com/coreos/etcd/clientv3"
)

// NewWatcher 查看数据的变化
func NewWatcher() (watcher clientv3.Watcher, err error) {
	conn := getConn()
	if conn == nil {
		err = errors.New("get etcd conn error")
		return
	}
	watcher = clientv3.NewWatcher(conn)
	return
}

// WithRange key范围获取
func WithRange(key string) clientv3.OpOption {
	return clientv3.WithRange(key)
}

// WithPrefix 前缀观察
func WithPrefix() clientv3.OpOption {
	return clientv3.WithPrefix()
}

// WithProgressNotify 通知监听过
func WithProgressNotify() clientv3.OpOption {
	return clientv3.WithProgressNotify()
}
