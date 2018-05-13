package etcd

// ConfigWatcherHookFunc 配置消息变更hook函数
// 可注册一个hook,当对应的配置变化的时候，就会被触化
type ConfigWatcherHookFunc func(new ConfigWatcherItem, old ConfigWatcherItem)

// configWatcherHook 配置变更hook的数据对象
type configWatcherHook struct {
	// Key 要处理的key
	Key string
	// WithPreFix 是否为前缀匹配
	WithPrefix bool
	// Hook 要运行的函数
	Hook ConfigWatcherHookFunc
}
