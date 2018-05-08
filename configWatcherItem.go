package etcd

import (
	"strconv"
)

// ConfigWatcherItem 配置项目
type ConfigWatcherItem struct {
	Key  string
	data string
}

// SetData 设置值
func (c *ConfigWatcherItem) SetData(data string) {
	c.data = data
}

// String 返回字符
func (c ConfigWatcherItem) String() string {
	return c.data
}

// Int 返回为整数
func (c ConfigWatcherItem) Int() int {
	if rs, err := strconv.Atoi(c.data); err == nil {
		return rs
	}
	return 0
}

// Int64 返回int64
func (c ConfigWatcherItem) Int64() int64 {
	if rs, err := strconv.ParseInt(c.data, 10, 64); err == nil {
		return rs
	}
	return 0
}

// Float64 返回float64
func (c ConfigWatcherItem) Float64() float64 {
	if rs, err := strconv.ParseFloat(c.data, 64); err == nil {
		return rs
	}
	return 0
}
