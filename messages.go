package jt808

// MessagePack 通用的消息体结构
type MessagePack[T any] struct {
	PackBody        T
	PackHeader      PackHeader
	Checksum        uint8
	IsChecksumValid bool
}
