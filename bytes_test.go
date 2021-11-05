package jt808

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEscapeChars(t *testing.T) {
	unescaped := []byte{0x7d, 0x7e, 0x01, 0x02}
	escaped, err := escapeChars(unescaped)

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x7d, 0x01, 0x7d, 0x02, 0x01, 0x02}, escaped)
}

func TestUnescapeChars(t *testing.T) {
	escaped := []byte{0x7e, 0x7d, 0x02, 0x7d, 0x01, 0x03, 0x04, 0x7e}
	unescaped, err := unescapeChars(escaped)

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x7e, 0x7e, 0x7d, 0x03, 0x04, 0x7e}, unescaped)
}

func TestComputeChecksum(t *testing.T) {
	b := []byte{0x01, 0x02, 0x03}
	checksum, err := computeChecksum(b)

	assert.NoError(t, err)
	assert.Equal(t, byte(0x00), checksum)
}
