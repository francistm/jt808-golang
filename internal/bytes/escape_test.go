package bytes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Escape(t *testing.T) {
	unescaped := []byte{0x7d, 0x7e, 0x01, 0x02}
	escaped, err := Escape(unescaped)

	if assert.NoError(t, err) {
		assert.Equal(t, []byte{0x7d, 0x01, 0x7d, 0x02, 0x01, 0x02}, escaped)
	}
}

func Test_Unescape(t *testing.T) {
	escaped := []byte{0x7e, 0x7d, 0x02, 0x7d, 0x01, 0x03, 0x04, 0x7e}
	unescaped, err := Unescape(escaped)

	if assert.NoError(t, err) {
		assert.Equal(t, []byte{0x7e, 0x7e, 0x7d, 0x03, 0x04, 0x7e}, unescaped)
	}
}
