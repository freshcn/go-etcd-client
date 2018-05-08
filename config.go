package etcd

import (
	"time"

	"github.com/coreos/etcd/clientv3"
)

// Config 配置
type Config = clientv3.Config

var (
	config = Config{
		DialTimeout: time.Second * 5,
	}
)
