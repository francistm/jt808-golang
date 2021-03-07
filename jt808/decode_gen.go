// Code generated by generator/decoder, DO NOT MODIFY MANUALLY

package jt808

import (
	"bytes"
	"fmt"
	message "github.com/francistm/jt808-golang/jt808/message"
)

func (messagePack *MessagePack) unmarshalBody(buf []byte) error {
	reader := bytes.NewReader(buf)

	if messagePack.PackHeader.Package != nil {
		messagePack.PackBody = new(message.PartialPackBody)
	} else {
		switch messagePack.PackHeader.MessageID {
		case uint16(0x1):
			messagePack.PackBody = new(message.Body0001)

		case uint16(0x200):
			messagePack.PackBody = new(message.Body0200)

		case uint16(0x801):
			messagePack.PackBody = new(message.Body0801)

		default:
			return fmt.Errorf("unsupported messageId: 0x%.4X", messagePack.PackHeader.MessageID)
		}
	}

	return unmarshalBody(reader, messagePack.PackBody)
}