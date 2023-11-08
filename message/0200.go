package message

import (
	"errors"
	"io"
	"math"
	"sort"
	"time"

	"github.com/francistm/jt808-golang/internal/bytes"
)

// 位置信息汇报
type Body0200 struct {
	Body0200Base
	RawExtraMessage []byte `jt808:"-1,raw"`

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

func (b *Body0200) Time() time.Time {
	t, _ := time.ParseInLocation(timeLayoutBCD, b.RawTime, timezoneCST)

	return t
}

func (b *Body0200) SetTime(tIn time.Time) {
	b.RawTime = tIn.In(timezoneCST).Format(timeLayoutBCD)
}

func (b *Body0200) ExtraMessage() (map[uint8][]byte, error) {
	if len(b.RawExtraMessage) == 0 {
		return nil, nil
	}

	if b.parsedExtraMessage == nil {
		reader := bytes.NewReader(b.RawExtraMessage)
		extraMesgs := make(map[uint8][]byte)

		for {
			id, err := reader.ReadByte()

			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				return nil, err
			}

			dataSize, err := reader.ReadByte()

			if err != nil {
				return nil, err
			}

			mesgData, err := reader.ReadFixedBytes(int(dataSize))

			if err != nil {
				return nil, err
			}

			extraMesgs[id] = mesgData
		}

		b.parsedExtraMessage = extraMesgs
	}

	return b.parsedExtraMessage, nil
}

func (b *Body0200) SetExtraMessage(m map[uint8][]byte) error {
	var (
		writer       = bytes.NewBuffer()
		extraMesgIds = make([]uint8, 0, len(m))
	)

	for extraMesgId := range m {
		extraMesgIds = append(extraMesgIds, extraMesgId)
	}

	sort.Slice(extraMesgIds, func(i, j int) bool {
		return extraMesgIds[i] < extraMesgIds[j]
	})

	for _, extraMesgId := range extraMesgIds {
		data := m[extraMesgId]
		dataSize := uint8(len(data))

		if err := writer.WriteByte(extraMesgId); err != nil {
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
