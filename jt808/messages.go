package jt808

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"time"
)

var timeLayout = "060102150405"
var locationCST = time.FixedZone("Asia/Shanghai", 8*3600)

// MessagePack 通用的消息体结构
type MessagePack struct {
	PackBody      PackBody
	PackHeader    PackHeader
	Checksum      uint8
	ChecksumValid bool

	bodyBuf       []byte
	unmarshalFunc func([]byte) (PackBody, error)
}

// ConcatAndUnmarshal 拼接多个分段消息并解析
func (ptr *MessagePack) ConcatAndUnmarshal(packs ...*MessagePack) error {
	buf := bytes.NewBuffer(ptr.bodyBuf)

	if ptr.PackHeader.Package == nil {
		return errors.New("cannot concat packages without package property header")
	}

	concatPackageIndexList := make(map[uint16]bool)
	concatPackageIndexList[ptr.PackHeader.Package.CurrentIndex] = true

	for _, pack := range packs {
		if pack.PackHeader.Package == nil {
			return errors.New("package to be concat doesn't have package property header")
		}

		if !concatPackageIndexList[pack.PackHeader.Package.CurrentIndex] {
			buf.Write(pack.bodyBuf)
			concatPackageIndexList[pack.PackHeader.Package.CurrentIndex] = true
		}
	}

	for i := uint16(0); i < ptr.PackHeader.Package.TotalCount; i++ {
		if concatPackageIndexList[i+1] == false {
			return fmt.Errorf("missing package with index %d to concat and unmarshal message", i+1)
		}
	}

	packBody, err := ptr.unmarshalFunc(buf.Bytes())

	if err != nil {
		return err
	}

	ptr.PackBody = packBody

	// cleanup after unmarshal
	ptr.bodyBuf = nil
	ptr.unmarshalFunc = nil

	return nil
}

// PackBody 消息体包正文接口
type PackBody interface {
	marshalBody() ([]byte, error)
}

// Body0001 0x0001 消息体正文结构体
type Body0001 struct {
	AcknowledgeSerialID  uint16
	AcknowledgeMessageID uint16
	AcknowledgeType      uint8
}

func (b Body0001) marshalBody() ([]byte, error) {
	var buf bytes.Buffer

	if err := writeUint16(b.AcknowledgeSerialID, &buf); err != nil {
		return nil, err
	}

	if err := writeUint16(b.AcknowledgeMessageID, &buf); err != nil {
		return nil, err
	}

	if err := buf.WriteByte(b.AcknowledgeType); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func unmarshalBody0001(buf []byte) (PackBody, error) {
	var body Body0001
	reader := bytes.NewReader(buf)

	i, err := readUint16(reader)

	if err != nil {
		return nil, err
	}

	body.AcknowledgeMessageID = i

	i, err = readUint16(reader)

	if err != nil {
		return nil, err
	}

	body.AcknowledgeSerialID = i

	b, err := reader.ReadByte()

	if err != nil {
		return nil, err
	}

	body.AcknowledgeType = b

	return body, nil
}

// Body0200 0x0200 消息体正文结构体
type Body0200 struct {
	WarnFlag     uint32
	StatusFlag   uint32
	Latitude     float64
	Longitude    float64
	Altitude     uint16
	Speed        float32
	Direction    uint16
	Time         time.Time
	ExtraMessage map[uint8][]byte
}

func (b Body0200) marshalBody() ([]byte, error) {
	var buf bytes.Buffer

	if err := writeUint32(b.WarnFlag, &buf); err != nil {
		return nil, err
	}

	if err := writeUint32(b.StatusFlag, &buf); err != nil {
		return nil, err
	}

	timeStr := b.Time.In(locationCST).Format(timeLayout)
	latitudeInt := uint32(b.Latitude * math.Pow(10, 6))
	longitudeInt := uint32(b.Longitude * math.Pow(10, 6))
	speedInt := uint16(b.Speed * 10)

	if err := writeUint32(latitudeInt, &buf); err != nil {
		return nil, err
	}

	if err := writeUint32(longitudeInt, &buf); err != nil {
		return nil, err
	}

	if err := writeUint16(b.Altitude, &buf); err != nil {
		return nil, err
	}

	if err := writeUint16(speedInt, &buf); err != nil {
		return nil, err
	}

	if err := writeUint16(b.Direction, &buf); err != nil {
		return nil, err
	}

	if err := writeBCD(timeStr, &buf); err != nil {
		return nil, err
	}

	if b.ExtraMessage != nil {
		dataIDList := make([]uint8, 0, len(b.ExtraMessage))

		// TODO try to avoid twice loop
		for dataID := range b.ExtraMessage {
			dataIDList = append(dataIDList, dataID)
		}

		// sort the map with dataId
		// otherwise will get random extra data order in marshalled bytes
		sort.SliceStable(dataIDList, func(i, j int) bool {
			return dataIDList[i] < dataIDList[j]
		})

		for _, dataID := range dataIDList {
			dataContent := b.ExtraMessage[dataID]
			dataLength := uint8(len(dataContent))

			buf.WriteByte(dataID)
			buf.WriteByte(dataLength)
			buf.Write(dataContent)
		}
	}

	return buf.Bytes(), nil
}

func unmarshalBody0200(buf []byte) (PackBody, error) {
	var body Body0200
	reader := bytes.NewReader(buf)

	ui32, err := readUint32(reader)

	if err != nil {
		return nil, err
	}

	body.WarnFlag = ui32

	ui32, err = readUint32(reader)

	if err != nil {
		return nil, err
	}

	body.StatusFlag = ui32

	ui32, err = readUint32(reader)

	if err != nil {
		return nil, err
	}

	body.Latitude = float64(ui32) / math.Pow(10, 6)

	ui32, err = readUint32(reader)

	if err != nil {
		return nil, err
	}
	body.Longitude = float64(ui32) / math.Pow(10, 6)

	ui16, err := readUint16(reader)

	if err != nil {
		return nil, err
	}
	body.Altitude = ui16

	ui16, err = readUint16(reader)

	if err != nil {
		return nil, err
	}

	body.Speed = float32(ui16) / 10

	ui16, err = readUint16(reader)

	if err != nil {
		return nil, err
	}

	body.Direction = ui16

	s, err := readBCD(reader, 6)

	if err != nil {
		return nil, err
	}

	t, err := time.ParseInLocation(timeLayout, s, locationCST)

	if err != nil {
		return nil, err
	}

	body.Time = t

	// read extra messages in 0200 (if have)
	for {
		extraDataID, err := reader.ReadByte()

		if err == io.EOF {
			break
		}

		if body.ExtraMessage == nil {
			body.ExtraMessage = make(map[uint8][]byte)
		}

		extraDataLength, err := reader.ReadByte()

		if err != nil {
			return nil, err
		}

		body.ExtraMessage[extraDataID] = make([]byte, extraDataLength)

		if _, err := reader.Read(body.ExtraMessage[extraDataID]); err != nil {
			return nil, err
		}
	}

	return body, nil
}

// Body0801 0x0801 消息体正文结构体
type Body0801 struct {
	MediaID          uint32
	MediaType        uint8
	MediaContentType uint8
	EventID          uint8
	ChannelID        uint8
	PackBody0200     Body0200
	MediaContent     []byte
}

func (b Body0801) marshalBody() ([]byte, error) {
	panic("implement me")
}

func unmarshalBody0801(buf []byte) (PackBody, error) {
	var body Body0801
	reader := bytes.NewReader(buf)

	var mediaContentBuf bytes.Buffer
	var packBody0200Buf = make([]byte, 28)

	ui32, err := readUint32(reader)

	if err != nil {
		return nil, err
	}

	body.MediaID = ui32

	ui8, err := reader.ReadByte()

	if err != nil {
		return nil, err
	}

	body.MediaType = ui8

	ui8, err = reader.ReadByte()

	if err != nil {
		return nil, err
	}

	body.MediaContentType = ui8

	ui8, err = reader.ReadByte()

	if err != nil {
		return nil, err
	}

	body.EventID = ui8

	ui8, err = reader.ReadByte()

	if err != nil {
		return nil, err
	}

	body.ChannelID = ui8

	if _, err := reader.Read(packBody0200Buf); err != nil {
		return nil, err
	} else if body0200, err := unmarshalBody0200(packBody0200Buf); err != nil {
		return nil, err
	} else {
		body.PackBody0200 = body0200.(Body0200)
	}

	_, err = reader.WriteTo(&mediaContentBuf)

	if err != nil {
		return nil, err
	}

	body.MediaContent = mediaContentBuf.Bytes()

	return body, nil
}
