package bytes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalcChecksum(t *testing.T) {
	b := []byte{0x01, 0x02, 0x03}
	checksum, err := CalcChecksum(b)

	if assert.NoError(t, err) {
		assert.Equal(t, byte(0x00), checksum)
	}
}
