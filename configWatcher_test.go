package etcd

import (
	"fmt"
	"testing"
	"time"

	"github.com/freshcn/log"
)

func init() {
	log.SetDebug(true)
}

var (
	testConfig = Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}
)

func TestConfigWatcher(t *testing.T) {
	New(testConfig)
	test := NewConfigWatcher("config")
	defer test.Close()
	for i := 0; i < 5; i++ {
		Put(fmt.Sprintf("CONFIG/test%d", i), fmt.Sprintf("%d", time.Now().Unix()))
		time.Sleep(time.Second * 1)
		fmt.Println("config key", fmt.Sprintf("test%d", i), " data:", test.Get(fmt.Sprintf("test%d", i)).Int())
	}
}

func TestConfigWatcherHooks(t *testing.T) {
	New(testConfig)
	test := NewConfigWatcher("config")
	defer test.Close()
	// 添加前缀hook
	test.AddHook("test", func(new ConfigWatcherItem, old ConfigWatcherItem) {
		fmt.Println("```key:", new.Key, " New value:", new.Int(), " old value:", old.Int())
	}, true)

	var key string
	for i := 0; i < 5; i++ {
		key = fmt.Sprintf("CONFIG/test%d", 0)
		Put(key, fmt.Sprintf("%d", time.Now().Unix()))
		time.Sleep(time.Second * 1)
		if i == 3 {
			fmt.Println("=== del key ", key)
			Del(key)
		}
		fmt.Println("|||config key:", key, " data:", test.Get(fmt.Sprintf("test%d", 0)).Int())
	}
}
