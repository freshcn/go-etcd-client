package etcd

import (
	"fmt"
	"net/url"
	"time"

	"github.com/freshcn/log"

	"github.com/freshcn/coroutine"
)

var (
	// endpoints etcd的endpoints信息
	endpoints []string
)

// New 开启etcd节点更新
func New(conf Config) {
	config = conf
	if config.DialTimeout == 0 {
		config.DialTimeout = time.Second * 5
	}

	// 初始化情况
	var initChan = make(chan bool)

	// 开始定时更新etcd的节点
	coroutine.Run(func() {
		conn = getConn()
		initChan <- true
		if conn != nil {
			defer conn.Close()
			defer Close()
		}

		for {
			if closeETCD {
				return
			}
			getEndPoints()
			time.Sleep(time.Second * 5)
		}
	})

	<-initChan

}

// getEndPoints 定时更新etcd的节点信息
func getEndPoints() {
	ctx, cancel := getCtx()
	resp, err := conn.MemberList(ctx)
	cancel()
	if err == nil {
		var (
			tmpURL *url.URL

			tmpEndpoints = make([]string, len(resp.Members))
		)
		for i, v := range resp.Members {
			if len(v.GetClientURLs()) == 0 {
				continue
			}
			if tmpURL, err = url.Parse((v.GetClientURLs())[0]); err == nil {
				tmpEndpoints[i] = fmt.Sprintf("%s:%s", tmpURL.Hostname(), tmpURL.Port())
			}
		}
		if !stringSliceEqual(endpoints, tmpEndpoints) {
			conn.SetEndpoints(tmpEndpoints...)
			endpoints = tmpEndpoints
			log.Info("etcd节点变为:", endpoints)
		}
	} else {
		log.Error(err)
	}

}
