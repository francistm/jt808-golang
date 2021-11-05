package jt808

// MessagePack 通用的消息体结构
type MessagePack struct {
	PackBody      *PackBody
	PackHeader    PackHeader
	Checksum      uint8
	ChecksumValid bool
}

//go:generate go run github.com/francistm/jt808-golang/tools/generator/message
type PackBody struct {
	body interface{}
}
