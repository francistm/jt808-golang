package jt808

import (
	"github.com/francistm/jt808-golang/message"
)

func Marshal[T any](mesgPack *message.MessagePack[T]) ([]byte, error) {
	return mesgPack.MarshalBinary()
}
