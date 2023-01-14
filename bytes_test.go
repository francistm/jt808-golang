package jt808

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_encodeBytes(t *testing.T) {
	unescaped := []byte{0x7d, 0x7e, 0x01, 0x02}
	escaped, err := encodeBytes(unescaped)

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x7d, 0x01, 0x7d, 0x02, 0x01, 0x02}, escaped)
}

func Test_decodeBytes(t *testing.T) {
	escaped := []byte{0x7e, 0x7d, 0x02, 0x7d, 0x01, 0x03, 0x04, 0x7e}
	unescaped, err := decodeBytes(escaped)

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x7e, 0x7e, 0x7d, 0x03, 0x04, 0x7e}, unescaped)
}

func Test_calculateChecksum(t *testing.T) {
	b := []byte{0x01, 0x02, 0x03}
	checksum, err := calculateChecksum(b)

	assert.NoError(t, err)
	assert.Equal(t, byte(0x00), checksum)
}
