package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PackHeader_MarshalBinary(t *testing.T) {
	packHeader := PackHeader{
		MessageID: 0x0102,
		Property: PackHeaderProperty{
			BodyByteLength: 5,
		},
		TerminalMobileNo: "013812345678",
		SerialNo:         0x0087,
	}

	got, err := packHeader.MarshalBinary()

	if assert.NoError(t, err) {
		assert.Equal(t, []byte{0x01, 0x02, 0x00, 0x05, 0x01, 0x38, 0x12, 0x34, 0x56, 0x78, 0x00, 0x87}, got)
	}
}

func Test_PackHeader_UnmarshalBinary(t *testing.T) {
	var packHeader PackHeader
	var packHeaderData = []byte{0x01, 0x02, 0x00, 0x05, 0x01, 0x38, 0x12, 0x34, 0x56, 0x78, 0x00, 0x87}

	err := packHeader.UnmarshalBinary(packHeaderData)

	if assert.NoError(t, err) {
		assert.Equal(t, uint16(0x0102), packHeader.MessageID)
		assert.Equal(t, uint16(0x0005), packHeader.Property.BodyByteLength)
		assert.Equal(t, false, packHeader.Property.IsEncrypted)
		assert.Equal(t, false, packHeader.Property.IsMultiplePackage)
		assert.Equal(t, "013812345678", packHeader.TerminalMobileNo)
		assert.Equal(t, uint16(0x0087), packHeader.SerialNo)
		assert.Nil(t, packHeader.Package)
	}
}

func Test_PackHeader_UnmarshalBinary_MultiPkg(t *testing.T) {
	var packHeader PackHeader
	var packHeaderData = []byte{0x08, 0x01, 0x23, 0x24, 0x01, 0x38, 0x12, 0x34, 0x56, 0x78, 0x20, 0x31, 0x00, 0x0a, 0x00, 0x01}

	err := packHeader.UnmarshalBinary(packHeaderData)

	if assert.NoError(t, err) {
		assert.Equal(t, uint16(0x0801), packHeader.MessageID)
		if assert.NotNil(t, packHeader.Package) {
			assert.Equal(t, uint16(1), packHeader.Package.Index)
			assert.Equal(t, uint16(10), packHeader.Package.TotalCount)
		}
	}
}
