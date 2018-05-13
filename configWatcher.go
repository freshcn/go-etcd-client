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

	dataMap map[string]ConfigWatcherItem

	close bool

	lock sync.RWMutex

	hooksChan chan configWatcherHook

	hooksLock   sync.RWMutex
	changeHooks []configWatcherHook
}

var (
	hooksCap = 10
)

// NewConfigWatcher 创建一个新的配置监听对象
func NewConfigWatcher(prefix string) (rs ConfigWatcher) {
	rs = ConfigWatcher{
		Prefix:      strings.ToUpper(prefix),
		dataMap:     make(map[string]ConfigWatcherItem),
		changeHooks: make([]configWatcherHook, 0, hooksCap),
		hooksChan:   make(chan configWatcherHook, hooksCap),
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

			for {
				if c.close {
					return
				}
				select {
				case events = <-rsp:
					for _, event = range events.Events {
						if c.close {
							return
						}
						srcKey = (string)(event.Kv.Key)
						key = c.clearPrefix(srcKey)
						switch event.Type {
						case clientv3.EventTypeDelete: // 当key被删除时
							if v := c.Get(key, false); v != nil {
								c.del(srcKey)
								c.doHook(srcKey, ConfigWatcherItem{Key: srcKey}, *v)
							}
						case clientv3.EventTypePut: // 当key更新时
							if v := c.Get(key, false); v != nil {
								old := *v
								v.SetData((string)(event.Kv.Value))
								c.save(srcKey, *v)
								c.doHook(srcKey, *v, old)
							} else {
								c.doHook(srcKey, ConfigWatcherItem{Key: srcKey, data: (string)(event.Kv.Value)}, ConfigWatcherItem{Key: srcKey})
							}
						}
					}
				case hook := <-c.hooksChan:
					// 将hook写入到hooks中
					c.addHook(hook)
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
// key 要获取的key，不需要添加前缀
// new 当不存在时，是否创建一个新对应, 默认会创建一个新的对象
func (c *ConfigWatcher) Get(key string, new ...bool) (item *ConfigWatcherItem) {

	key = c.addPrefix(key)

	var isNew = true
	if len(new) > 0 {
		isNew = new[0]
	}

	if v, exists := c.dataMap[key]; exists {
		item = &v
		return
	} else if !isNew { // 当不需要建新的对象时
		return
	}

	var unlock = true
	c.lock.RLock()
	defer func() {
		if unlock {
			c.lock.RUnlock()
		}
	}()

	if rs, err := Get(key); err == nil {
		item = &ConfigWatcherItem{
			data: rs,
			Key:  key,
		}
		unlock = false
		c.lock.RUnlock()
		c.save(key, *item)
	} else {
		log.Error(err)
	}
	return
}

func (c *ConfigWatcher) save(key string, item ConfigWatcherItem) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.dataMap[key] = item
}

func (c *ConfigWatcher) del(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.dataMap, key)
}

// Exists 判断配置是否存在
func (c *ConfigWatcher) Exists(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if _, exists := c.dataMap[c.addPrefix(key)]; exists {
		return true
	}
	return false
}

// AddHook 添加一个配置hook
// 当配置项目改变的时候，会调用这个hook
func (c *ConfigWatcher) AddHook(key string, hook ConfigWatcherHookFunc, withPrefix ...bool) bool {
	var prefix bool
	if len(withPrefix) > 0 {
		prefix = withPrefix[0]
	}

	var (
		tmpHook = configWatcherHook{
			Key:        c.addPrefix(key),
			Hook:       hook,
			WithPrefix: prefix,
		}
	)

	c.hooksChan <- tmpHook

	return true
}

// addHook 将数据写入到hooks中
func (c *ConfigWatcher) addHook(hook configWatcherHook) {
	c.hooksLock.Lock()
	defer c.hooksLock.Unlock()
	c.changeHooks = append(c.changeHooks, hook)
}

// doHook 执行Hook处理
func (c *ConfigWatcher) doHook(key string, new ConfigWatcherItem, old ConfigWatcherItem) {
	if len(c.changeHooks) > 0 {
		new.Key = c.clearPrefix(new.Key)
		old.Key = c.clearPrefix(old.Key)
		c.hooksLock.RLock()
		defer c.hooksLock.RUnlock()
		for _, v := range c.changeHooks {
			if v.WithPrefix { //　当为前缀匹配时
				if strings.Index(key, v.Key) == 0 {
					v.Hook(new, old)
				}
			} else if key == v.Key { // 完整匹配
				v.Hook(new, old)
			}
		}
	}
}
