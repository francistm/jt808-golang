package jt808

// MessagePack 通用的消息体结构
type MessagePack struct {
	PackBody      interface{}
	PackHeader    PackHeader
	Checksum      uint8
	ChecksumValid bool
}
