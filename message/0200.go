package message

import (
	"bytes"
	"io"
	"math"
	"sort"
	"time"
)

// Body0200 0x0200 消息体正文结构体
type Body0200 struct {
	Body0200Base
	RawExtraMessage []byte `jt808:",none"`

	parsedExtraMessage map[uint8][]byte // cache parsed rawExtraMessage
}

type Body0200Base struct {
	WarnFlag     uint32
	StatusFlag   uint32
	RawLatitude  uint32
	RawLongitude uint32
	Altitude     uint16
	RawSpeed     uint16
	Direction    uint16
	RawTime      string `jt808:"6,bcd"`
}

func (b *Body0200) Latitude() float64 {
	return float64(b.RawLatitude) / math.Pow(10, 6)
}

func (b *Body0200) SetLatitude(f float64) {
	b.RawLatitude = uint32(f * math.Pow(10, 6))
}

func (b *Body0200) Longitude() float64 {
	return float64(b.RawLongitude) / math.Pow(10, 6)
}

func (b *Body0200) SetLongitude(f float64) {
	b.RawLongitude = uint32(f * math.Pow(10, 6))
}

func (b *Body0200) Speed() float32 {
	return float32(b.RawSpeed) / 10
}

func (b *Body0200) SetSpeed(f float32) {
	b.RawSpeed = uint16(f * 10)
}

func (b *Body0200) Time() *time.Time {
	cstTime, err := time.ParseInLocation(timeLayout, b.RawTime, locationCST)

	if err != nil {
		// TODO: add logger here
		return nil
	}

	return &cstTime
}

func (b *Body0200) SetTime(t *time.Time) {
	b.RawTime = t.In(locationCST).Format(timeLayout)
}

func (b *Body0200) ExtraMessage() (map[uint8][]byte, error) {
	if b.parsedExtraMessage != nil {
		return b.parsedExtraMessage, nil
	}

	extraData := make(map[uint8][]byte)

	if len(b.RawExtraMessage) == 0 {
		return extraData, nil
	}

	reader := bytes.NewReader(b.RawExtraMessage)

	for {
		id, err := reader.ReadByte()

		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		dataLength, err := reader.ReadByte()

		if err != nil {
			return nil, err
		}

		extraData[id] = make([]byte, dataLength)

		if _, err := reader.Read(extraData[id]); err != nil {
			return nil, err
		}
	}

	b.parsedExtraMessage = extraData

	return extraData, nil
}

func (b *Body0200) SetExtraMessage(m map[uint8][]byte) error {
	writer := new(bytes.Buffer)

	dataIds := make([]uint8, 0, len(m))

	for dataId := range m {
		dataIds = append(dataIds, dataId)
	}

	sort.SliceStable(dataIds, func(i, j int) bool {
		return dataIds[i] < dataIds[j]
	})

	for _, dataId := range dataIds {
		data := m[dataId]
		dataSize := uint8(len(data))

		if err := writer.WriteByte(dataId); err != nil {
			return err
		}

		if err := writer.WriteByte(dataSize); err != nil {
			return err
		}

		if _, err := writer.Write(data); err != nil {
			return err
		}
	}

	// update cache
	b.parsedExtraMessage = m
	b.RawExtraMessage = writer.Bytes()

	return nil
}
