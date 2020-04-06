package jt808

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalHeader(t *testing.T) {
	var packHeader PackHeader
	var packHeaderBytes = []byte{0x01, 0x02, 0x00, 0x05, 0x01, 0x38, 0x12, 0x34, 0x56, 0x78, 0x00, 0x87}

	err := UnmarshalHeader(packHeaderBytes, &packHeader)

	assert.NoError(t, err)
	assert.Equal(t, uint16(0x0102), packHeader.MessageId)
	assert.Equal(t, uint16(0x0005), packHeader.Property.BodyByteLength)
	assert.Equal(t, false, packHeader.Property.IsEncrypted)
	assert.Equal(t, false, packHeader.Property.IsMultiplePackage)
	assert.Equal(t, "013812345678", packHeader.TerminalMobileNo)
	assert.Equal(t, uint16(0x0087), packHeader.SerialNo)
	assert.Nil(t, packHeader.Package)
}
