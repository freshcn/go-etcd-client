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

func TestConfigWatcher(t *testing.T) {
	New(Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	test := NewConfigWatcher("config")
	for i := 0; i < 5; i++ {
		Put(fmt.Sprintf("CONFIG/test%d", i), fmt.Sprintf("%d", time.Now().Unix()))
		time.Sleep(time.Second * 1)
		fmt.Println("config key", fmt.Sprintf("test%d", i), " data:", test.Get(fmt.Sprintf("test%d", i)).Int())
	}
}
