# ginger
`ginger` 是一个用Golang语言实现的高性能网络库。
相对于C++,Go实现的网络库不必关心网络事件处理模型的具体实现，故而大大简化网络库开发。
相对于其他网络库，`ginger`使用Read+RecvBuffer而非ReadSlice实现数据接收,从而提高网络库效率。
