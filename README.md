# Ginger
`Ginger` 是一个用Golang语言实现的TCP高性能网络库。
##特点
* 相对于C++,Go实现的网络库不必关心网络事件处理模型的具体实现，故而大大简化网络库开发。
* 相对于其他网络库，`Ginger`使用`Read+RecvBuffer`而非`ReadSlice`实现数据接收,从而提高网络库效率。
* 将协议解封功能解耦，便于业务协议开
