package etcd

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/freshcn/coroutine"
	"github.com/freshcn/log"
)

// ConfigWatcher 配置监听
// 自动加载对应的配置信息
type ConfigWatcher struct {
	// Prefix 配置的前缀
	Prefix string

	array map[string]ConfigWatcherItem

	close bool

	lock sync.Mutex
}

// NewConfigWatcher 创建一个新的配置监听对象
func NewConfigWatcher(prefix string) (rs ConfigWatcher) {
	rs = ConfigWatcher{
		Prefix: strings.ToUpper(prefix),
		array:  make(map[string]ConfigWatcherItem),
	}
	rs.startWatch()
	addCloser(&rs)
	return
}

// Close 关闭监听者
func (c *ConfigWatcher) Close() {
	log.Warringf("Config Watcher %s exit", c.Prefix)
	c.close = true
}

func (c *ConfigWatcher) startWatch() {
	// 开始监听
	coroutine.Run(func() {
		if watcher, err := NewWatcher(); err == nil {
			defer watcher.Close()
			rsp := watcher.Watch(context.Background(), c.Prefix, WithPrefix())

			var (
				key    string
				srcKey string
				// item   ConfigWatcherItem
				event  *clientv3.Event
				events clientv3.WatchResponse
			)

			for events = range rsp {
				if c.close {
					return
				}

				for _, event = range events.Events {
					if c.close {
						return
					}
					srcKey = (string)(event.Kv.Key)
					key = c.clearPrefix(srcKey)
					switch event.Type {
					case clientv3.EventTypeDelete: // 当key被删除时
						if c.Exists(key) {
							c.del(srcKey)
						}
					case clientv3.EventTypePut: // 当key更新时
						if v, exist := c.array[srcKey]; exist {
							v.SetData((string)(event.Kv.Value))
							c.save(srcKey, v)
						}
					}
				}
			}
		} else {
			log.Panic(err)
		}
	})

}

// clearPrefix 清除前缀
func (c *ConfigWatcher) clearPrefix(key string) string {
	return strings.TrimPrefix(key, fmt.Sprintf("%s/", c.Prefix))
}

// addPrefix 添加前缀
func (c *ConfigWatcher) addPrefix(key string) string {
	return fmt.Sprintf("%s/%s", c.Prefix, key)
}

// Get 获取配置值
func (c *ConfigWatcher) Get(key string) (item *ConfigWatcherItem) {
	key = c.addPrefix(key)

	if v, exists := c.array[key]; exists {
		item = &v
		return
	}

	if rs, err := Get(key); err == nil {
		item = &ConfigWatcherItem{
			data: rs,
			Key:  key,
		}
		c.save(key, *item)
	} else {
		log.Error(err)
	}
	return
}

func (c *ConfigWatcher) save(key string, item ConfigWatcherItem) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.array[key] = item
}

func (c *ConfigWatcher) del(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.array, key)
}

// Exists 判断配置是否存在
func (c *ConfigWatcher) Exists(key string) bool {
	if _, exists := c.array[c.addPrefix(key)]; exists {
		return true
	}
	return false
}
