package etcd

import (
	"context"
	"errors"

	"github.com/coreos/etcd/clientv3"
	"github.com/freshcn/log"
)

type (
	// Lease 租约
	Lease struct {
		lease   clientv3.Lease
		context context.Context
	}

	// LeaseGrantResponse 创建租约回执
	LeaseGrantResponse struct {
		ID    int64 `json:"id"`
		TTL   int64 `json:"ttl"`
		Error error `json:"error"`
	}

	// LeaseTimeToLiveResponse 获取租约的租约信息回执
	LeaseTimeToLiveResponse struct {
		LeaseGrantResponse
		GrantedTTL int64    `json:"granted-ttl"`
		Keys       [][]byte `json:"keys"`
	}
)

// NewLease 创建一个新的租约
func NewLease() (l Lease, err error) {
	conn := getConn()
	if conn == nil {
		log.Error("get connect to etcd has error.")
		err = errors.New("get connect to etcd has error")
		return
	}
	l.lease = clientv3.NewLease(conn)
	l.context = context.TODO()
	addCloser(l.lease)
	return
}

// SetContext 设置请求中的上下文
func (l *Lease) SetContext(ctx context.Context) *Lease {
	l.context = ctx
	return l
}

// Grant 生成一个新的租约
func (l *Lease) Grant(ttl int64) (rep LeaseGrantResponse, err error) {
	var (
		lrep *clientv3.LeaseGrantResponse
	)
	if lrep, err = l.lease.Grant(l.context, ttl); err == nil {
		rep.ID = int64(lrep.ID)
		rep.TTL = lrep.TTL
		return
	}
	rep.Error = err
	return
}

// Revoke 撤消一个租约
func (l *Lease) Revoke(id int64) (err error) {
	_, err = l.lease.Revoke(l.context, clientv3.LeaseID(id))
	return
}

// TimeToLive 返回指定租约的租约信息
func (l *Lease) TimeToLive(id int64, opts ...clientv3.LeaseOption) (rep LeaseTimeToLiveResponse, err error) {
	var (
		lrep *clientv3.LeaseTimeToLiveResponse
	)
	if lrep, err = l.lease.TimeToLive(l.context, clientv3.LeaseID(id), opts...); err == nil {
		rep.ID = int64(lrep.ID)
		rep.TTL = lrep.TTL
		rep.GrantedTTL = lrep.GrantedTTL
		rep.Keys = lrep.Keys
	}
	return
}

// KeepAlive 让租约永远有效
// 如果租约有效的信息一直没有发送成功，客户端将会每秒发送一次，直到成功为止。
// 详细说明请查看https://godoc.org/github.com/coreos/etcd/clientv3#Lease
func (l *Lease) KeepAlive(id int64) (rep LeaseGrantResponse, err error) {
	var (
		lrep <-chan *clientv3.LeaseKeepAliveResponse
	)

	if lrep, err = l.lease.KeepAlive(l.context, clientv3.LeaseID(id)); err == nil {
		ch := <-lrep
		rep.ID = int64(ch.ID)
		rep.TTL = ch.TTL
	}
	return
}

// KeepAliveOnce 为服务设置一次在线保持
func (l *Lease) KeepAliveOnce(id int64) (rep LeaseGrantResponse, err error) {
	var (
		lrep *clientv3.LeaseKeepAliveResponse
	)
	if lrep, err = l.lease.KeepAliveOnce(l.context, clientv3.LeaseID(id)); err == nil {
		rep.ID = int64(lrep.ID)
		rep.TTL = lrep.TTL
	}
	return
}

// GetLease 获取租约的操作对象
func (l *Lease) GetLease() clientv3.Lease {
	return l.lease
}

// Close 关闭租约
func (l *Lease) Close() error {
	log.Info("lease close")
	l.lease.Close()
	return nil
}
