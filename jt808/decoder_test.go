package jt808

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	var messagePack MessagePack
	b := []byte{0x7e, 0x00, 0x01, 0x00, 0x05, 0x01, 0x86, 0x57, 0x40, 0x59, 0x79, 0x00, 0x8f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x99, 0x7e}
	err := Unmarshal(b, &messagePack)

	assert.NoError(t, err)
	assert.Equal(t, uint16(0x0001), messagePack.PackHeader.MessageId)
	assert.Equal(t, uint16(0x0005), messagePack.PackHeader.Property.BodyByteLength)
	assert.Equal(t, uint16(0x1011), messagePack.PackBody.(Body0001).AcknowledgeMessageId)
	assert.Equal(t, uint16(0x1213), messagePack.PackBody.(Body0001).AcknowledgeSerialId)
	assert.Equal(t, uint8(0x14), messagePack.PackBody.(Body0001).AcknowledgeType)
	assert.Equal(t, uint8(0x99), messagePack.Checksum)
}
