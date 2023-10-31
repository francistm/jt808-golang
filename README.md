[![Go](https://github.com/francistm/jt808-golang/actions/workflows/go.yml/badge.svg)](https://github.com/francistm/jt808-golang/actions/workflows/go.yml)

部分实现和测试用例参考 <https://github.com/SmallChi/JT808>

Project is WIP

# 消息定义
``` go
type Body0001 struct {
	AckMesgId   uint16
	AckSerialId uint16
	AckType     uint8
}
```
结构体字段顺序按照序列化的先后顺序。

## 字段类型映射

|Golang 类型|协议数据类型|tag定义 |
|-----------|------------|--------|
|uint8      |byte        |        |
|uint16     |word        |        |
|uint32     |dword       |        |
|[]byte     |byte[n]     |n,raw   |
|string     |bcd[n]      |n,bcd   |
|string     |string      |n,gbk   |

## 复杂类型解析
见 [./jt808/message/0200.go](https://github.com/francistm/jt808-golang/blob/c02868ec780de98aa3301ac24308a25532f2a7f6/jt808/message/0200.go) 。先将消息解析为 []byte 类型，然后在结构体中增加方法单独解析。

## 分包消息解析
如果消息头部中包含分包信息`MessagePack.Package`，则消息正文会被解析为`message.PartialPackBody`。
待消息全部接收完成后，使用 `jt808.ConcatUnmarshal(packs []*MessagePack[*message.PartialPackBody], target *MessagePack[T])` 方法一并进行解析。

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
go version go1.20.10 darwin/arm64

> sysctl -a | grep machdep.cpu.brand_string
machdep.cpu.brand_string: Apple M1 Pro

> go test -bench=. -benchmem  ./...
Benchmark_Unmarshal0001-8   	 1186509	      1006 ns/op	     776 B/op	      23 allocs/op
~~~
