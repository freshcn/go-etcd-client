本项目为etcd的客户端一个简单的封装

通过调用`etcd.New()`来开始，结束的时候需要调用`etcd.Close()`来关闭etcd的连接

ConfigWatcher可以提供对某些配置项目的监听，做到动态的调整