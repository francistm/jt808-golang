package jt808

import (
	"testing"
	"time"

	"github.com/francistm/jt808-golang/message"
	"github.com/stretchr/testify/assert"
)

func TestMarshal0200(t *testing.T) {
	type args struct {
		mesgPack *message.MessagePack[message.MesgBody]
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "mesg 8100",
			args: args{
				mesgPack: &message.MessagePack[message.MesgBody]{
					PackHeader: message.PackHeader{
						MessageID:         0x8100,
						TerminalMobileNum: "123456789012",
						SerialNum:         126,
					},
					PackBody: &message.Body8100{
						SerialId: 0x1234,
						Result:   message.AckType_8100_Succeed,
						AuthCode: "123456",
					},
				},
			},
			want:    []byte{0x7e, 0x81, 0x0, 0x0, 0x9, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x0, 0x7d, 0x2, 0x12, 0x34, 0x0, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x5d, 0x7e},
			wantErr: false,
		},
		{
			name: "mesg 0200",
			args: args{
				mesgPack: &message.MessagePack[message.MesgBody]{
					PackHeader: message.PackHeader{
						MessageID:         0x0200,
						TerminalMobileNum: "123456789012",
						SerialNum:         126,
					},
					PackBody: func() *message.Body0200 {
						base := message.Body0200{
							Body0200Base: message.Body0200Base{
								WarnFlag:   1,
								StatusFlag: 2,
								Altitude:   40,
								Direction:  0,
							},
						}

						base.SetLatitude(12.222222)
						base.SetLongitude(132.444444)
						base.SetSpeed(6)
						base.SetTime(&[]time.Time{time.Unix(1539569410, 0)}[0]) // // 2018-10-15 10:10:10 UTC
						base.SetExtraMessage(map[uint8][]byte{
							0x01: {0x00, 0x00, 0x00, 0x64},
							0x02: {0x00, 0x7d},
						})

						return &base
					}(),
				},
			},
			want:    []byte{0x7e, 0x02, 0x00, 0x00, 0x26, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x00, 0x7d, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0xba, 0x7f, 0x0e, 0x07, 0xe4, 0xf1, 0x1c, 0x00, 0x28, 0x00, 0x3c, 0x00, 0x00, 0x18, 0x10, 0x15, 0x10, 0x10, 0x10, 0x01, 0x04, 0x00, 0x00, 0x00, 0x64, 0x02, 0x02, 0x00, 0x7d, 0x01, 0x13, 0x7e},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				got, err := Marshal(tt.args.mesgPack)

				if tt.wantErr {
					assert.Error(t, err)
				} else if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			})
		})
	}
}
