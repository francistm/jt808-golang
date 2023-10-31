package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Body0102_UnmarshalBinary(t *testing.T) {
	data := []byte{0x05, 0x34, 0x35, 0x36, 0x31, 0x32, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x61, 0x62, 0x63, 0x64, 0x65, 0x76, 0x32, 0x2E, 0x30, 0x2E, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	body0102 := new(Body0102_19)

	err := body0102.UnmarshalBinary(data)

	if assert.NoError(t, err) {
		assert.Equal(t, "45612", body0102.AuthCode)
		assert.Equal(t, "1234567890abcde", body0102.DeviceIMEI)
		assert.Equal(t, "v2.0.0", body0102.DeviceVersion)
	}
}

func Test_Body0102_MarshalBinary(t *testing.T) {
	body0102 := &Body0102_19{
		AuthCode:      "45612",
		DeviceIMEI:    "1234567890abcde",
		DeviceVersion: "v2.0.0",
	}

	got, err := body0102.MarshalBinary()

	if assert.NoError(t, err) {
		assert.Equal(t, []byte{0x05, 0x34, 0x35, 0x36, 0x31, 0x32, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x61, 0x62, 0x63, 0x64, 0x65, 0x76, 0x32, 0x2E, 0x30, 0x2E, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, got)
	}
}