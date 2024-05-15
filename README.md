# mutexmap 基于 mutex 和 map 的 k-v-cache 的实现
跟syncmap不同，这个是一个rw-mutex和一个map的组合，目的是解决map的异步读写问题，这个比较鬼扯，查别人已有的代码也行，但不如自己顺手实现个完事

提供两个方案的实现：

[常规实现](/mutex_map.go)

[更高性能但更复杂的实现](/mutexmapcache/mutex_map_cache.go)

想用哪个就用哪个吧
