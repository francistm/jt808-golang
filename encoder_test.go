package jt808

import (
	"testing"
	"time"

	"github.com/francistm/jt808-golang/message"
	"github.com/stretchr/testify/assert"
)

func TestMarshal0200(t *testing.T) {
	messagePack := new(message.MessagePack[*message.Body0200])

	messagePack.PackHeader = message.PackHeader{
		MessageID:        0x0200,
		TerminalMobileNo: "123456789012",
		SerialNo:         126,
	}

	body0200 := &message.Body0200{
		Body0200Base: message.Body0200Base{
			WarnFlag:   1,
			StatusFlag: 2,
			Altitude:   40,
			Direction:  0,
		},
	}

	utcTime := time.Unix(1539569410, 0) // 2018-10-15 10:10:10 UTC

	body0200.SetLatitude(12.222222)
	body0200.SetLongitude(132.444444)
	body0200.SetSpeed(6)
	body0200.SetTime(&utcTime)
	body0200.SetExtraMessage(map[uint8][]byte{
		0x01: {0x00, 0x00, 0x00, 0x64},
		0x02: {0x00, 0x7d},
	})

	messagePack.PackBody = body0200

	assert.NotPanics(t, func() {
		b, err := Marshal(messagePack)

		if assert.NoError(t, err) {
			assert.Equal(t, []byte{0x7e, 0x02, 0x00, 0x00, 0x26, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x00, 0x7d, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0xba, 0x7f, 0x0e, 0x07, 0xe4, 0xf1, 0x1c, 0x00, 0x28, 0x00, 0x3c, 0x00, 0x00, 0x18, 0x10, 0x15, 0x10, 0x10, 0x10, 0x01, 0x04, 0x00, 0x00, 0x00, 0x64, 0x02, 0x02, 0x00, 0x7d, 0x01, 0x13, 0x7e}, b)
		}
	})
}
