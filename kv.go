package etcd

import (
	"errors"
	"fmt"

	"github.com/coreos/etcd/clientv3"
)

const (
	// OPGet 获取
	OPGet = iota
	// OPPut 设置
	OPPut
	// OPDel 删除
	OPDel
)

// Put 设置数据
func Put(key, value string) (b bool, err error) {
	_, err = Do(clientv3.OpPut(key, value))
	if err != nil {
		return
	}
	b = true
	return
}

// Get 获取数据
func Get(key string) (data string, err error) {
	rs, err := Do(clientv3.OpGet(key))
	if err != nil {
		return
	}
	if rs.Get().Count > 0 {
		data = fmt.Sprintf("%s", rs.Get().Kvs[0].Value)
	}
	return
}

// Del 删除数据
func Del(key string) (b bool, err error) {
	_, err = Do(clientv3.OpDelete(key))
	if err != nil {
		return
	}
	b = true
	return
}

// Do 发出操作
func Do(op clientv3.Op) (res clientv3.OpResponse, err error) {
	conn := getConn()
	if conn == nil {
		err = errors.New("get etcd conn error")
		return
	}
	ctx, cancel := getCtx()
	res, err = conn.Do(ctx, op)
	cancel()
	return
}
