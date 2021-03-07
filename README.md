[![Build Status](https://www.travis-ci.com/francistm/jt808-golang.svg?branch=master)](https://www.travis-ci.com/francistm/jt808-golang)

部分实现和测试用例参考 <https://github.com/SmallChi/JT808>

Project is WIP

# 消息定义
``` go
type Body0001 struct {
	AcknowledgeMessageID uint16 `jt808:""`
	AcknowledgeSerialID  uint16 `jt808:""`
	AcknowledgeType      uint8  `jt808:""`
}
```
结构体字段顺序按照序列化的先后顺序。

## 字段类型映射

|Golang 类型|协议数据类型|tag定义 |
|-----------|------------|--------|
|uint8      |byte        |        |
|uint16     |word        |        |
|uint32     |dword       |        |
|[]byte     |byte[n]     |n,none  |
|string     |bcd[n]      |n,bcd   |
|string     |string      |n,gbk   |

## 复杂类型解析
见 [./jt808/message/0200.go](https://github.com/francistm/jt808-golang/blob/c02868ec780de98aa3301ac24308a25532f2a7f6/jt808/message/0200.go) 。先将消息解析为 []byte 类型，然后在结构体中增加方法单独解析。 

## 分包消息解析
如果消息头部中包含分包信息`MessagePack.Package`，则消息正文会被解析为`message.PartialPackBody`。
待消息全部接收完成后，使用 `jt808.ConcatUnmarshal(...packs *MessagePack)` 方法一并进行解析。

详情见 [./jt808/decoder_test.go:40](https://github.com/francistm/jt808-golang/blob/c02868ec780de98aa3301ac24308a25532f2a7f6/jt808/decoder_test.go#L40) 中的测试用例。

# 完成情况
|消息ID|完成情况|测试覆盖|中文说明|备注|
|------|--------|--------|--------|----|
|0x0001|[x]|[x]|终端通用应答||
|0x0200|[x]|[x]|位置信息汇报||
|0x0801|[x]|[x]|多媒体数据上传||

# Benchmark
~~~
> go version
go version go1.14.6 darwin/amd64

> sysctl -a | grep machdep.cpu.brand_string
machdep.cpu.brand_string: Intel(R) Core(TM) i5-6267U CPU @ 2.90GHz

> go test -bench=. -benchmem  ./...
BenchmarkUnmarshal0001-4   	  850024	      1380 ns/op	     392 B/op	      23 allocs/op
~~~

