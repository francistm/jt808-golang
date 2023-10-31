package bytes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Buffer_WriteBCD(t *testing.T) {
	type args struct {
		s    string
		size int
	}

	tests := []struct {
		name     string
		args     args
		wantData []byte
		wantErr  bool
	}{
		{
			name: "string less than size",
			args: args{
				s:    "123",
				size: 4,
			},
			wantData: []byte{0x00, 0x00, 0x01, 0x23},
			wantErr:  false,
		},
		{
			name: "string equal size",
			args: args{
				s:    "12345678",
				size: 4,
			},
			wantData: []byte{0x12, 0x34, 0x56, 0x78},
		},
		{
			name: "string greater than size",
			args: args{
				s:    "123456789",
				size: 4,
			},
			wantData: []byte{0x12, 0x34, 0x56, 0x78},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := NewBuffer()

			err := buffer.WriteBCD(tt.args.s, tt.args.size)

			if tt.wantErr {
				assert.Error(t, err)
			} else if assert.NoError(t, err) {
				assert.Equal(t, tt.wantData, buffer.Bytes())
			}
		})
	}
}
