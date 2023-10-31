package message

import (
	"strings"

	"github.com/francistm/jt808-golang/internal/bytes"
)

// 终端鉴权
type Body0102_19 struct {
	AuthCodeSize  uint8
	AuthCode      string
	DeviceIMEI    string // 15 byte
	DeviceVersion string // 20 byte
}

func (body *Body0102_19) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer()

	buf.WriteUint8(uint8(len(body.AuthCode)))
	buf.WriteString(body.AuthCode)
	buf.WriteFixedString(body.DeviceIMEI, 15)
	buf.WriteFixedString(body.DeviceVersion, 20)

	return buf.Bytes(), nil
}

func (body *Body0102_19) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)
	authSize, err := reader.ReadUint8()

	if err != nil {
		return err
	}

	authCode, err := reader.ReadFixedString(int(authSize))

	if err != nil {
		return err
	}

	imei, err := reader.ReadFixedString(15)

	if err != nil {
		return err
	}

	version, err := reader.ReadFixedString(20)

	if err != nil {
		return err
	}

	body.AuthCodeSize = authSize
	body.AuthCode = authCode
	body.DeviceIMEI = strings.TrimRight(imei, "\x00")
	body.DeviceVersion = strings.TrimRight(version, "\x00")

	return nil
}
