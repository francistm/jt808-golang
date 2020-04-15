package jt808

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"time"
)

var locationCST = time.FixedZone("Asia/Shanghai", 8*3600)

type MessagePack struct {
	PackBody      PackBody
	PackHeader    PackHeader
	Checksum      uint8
	ChecksumValid bool

	bodyBuf       []byte
	unmarshalFunc func([]byte) (PackBody, error)
}

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
			return errors.New(fmt.Sprintf("missing package with index %d to concat and unmarshal message", i+1))
		}
	}

	if packBody, err := ptr.unmarshalFunc(buf.Bytes()); err != nil {
		return err
	} else {
		ptr.PackBody = packBody
	}

	// cleanup after unmarshal
	ptr.bodyBuf = nil
	ptr.unmarshalFunc = nil

	return nil
}

type PackBody interface {
	marshalBody() ([]byte, error)
}

type Body0001 struct {
	AcknowledgeSerialId  uint16
	AcknowledgeMessageId uint16
	AcknowledgeType      uint8
}

func (b Body0001) marshalBody() ([]byte, error) {
	var buf bytes.Buffer

	if err := writeUint16(b.AcknowledgeSerialId, &buf); err != nil {
		return nil, err
	}

	if err := writeUint16(b.AcknowledgeMessageId, &buf); err != nil {
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

	if i, err := readUint16(reader); err != nil {
		return nil, err
	} else {
		body.AcknowledgeMessageId = i
	}

	if i, err := readUint16(reader); err != nil {
		return nil, err
	} else {
		body.AcknowledgeSerialId = i
	}

	if i, err := reader.ReadByte(); err != nil {
		return nil, err
	} else {
		body.AcknowledgeType = i
	}

	return body, nil
}

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
	panic("implement me")
}

func unmarshalBody0200(buf []byte) (PackBody, error) {
	var body Body0200
	reader := bytes.NewReader(buf)

	if i, err := readUint32(reader); err != nil {
		return nil, err
	} else {
		body.WarnFlag = i
	}

	if i, err := readUint32(reader); err != nil {
		return nil, err
	} else {
		body.StatusFlag = i
	}

	if i, err := readUint32(reader); err != nil {
		return nil, err
	} else {
		body.Latitude = float64(i) / math.Pow(10, 6)
	}

	if i, err := readUint32(reader); err != nil {
		return nil, err
	} else {
		body.Longitude = float64(i) / math.Pow(10, 6)
	}

	if i, err := readUint16(reader); err != nil {
		return nil, err
	} else {
		body.Altitude = i
	}

	if i, err := readUint16(reader); err != nil {
		return nil, err
	} else {
		body.Speed = float32(i) / 10
	}

	if i, err := readUint16(reader); err != nil {
		return nil, err
	} else {
		body.Direction = i
	}

	if s, err := readBCD(reader, 6); err != nil {
		return nil, err
	} else {
		t, err := time.ParseInLocation("060102150405", s, locationCST)

		if err != nil {
			return nil, err
		}

		body.Time = t
	}

	// read extra messages in 0200 (if have)
	for {
		extraDataId, err := reader.ReadByte()

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

		body.ExtraMessage[extraDataId] = make([]byte, extraDataLength)

		if _, err := reader.Read(body.ExtraMessage[extraDataId]); err != nil {
			return nil, err
		}
	}

	return body, nil
}

type Body0801 struct {
	MediaId          uint32
	MediaType        uint8
	MediaContentType uint8
	EventId          uint8
	ChannelId        uint8
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

	if i, err := readUint32(reader); err != nil {
		return nil, err
	} else {
		body.MediaId = i
	}

	if i, err := reader.ReadByte(); err != nil {
		return nil, err
	} else {
		body.MediaType = i
	}

	if i, err := reader.ReadByte(); err != nil {
		return nil, err
	} else {
		body.MediaContentType = i
	}

	if i, err := reader.ReadByte(); err != nil {
		return nil, err
	} else {
		body.EventId = i
	}

	if i, err := reader.ReadByte(); err != nil {
		return nil, err
	} else {
		body.ChannelId = i
	}

	if _, err := reader.Read(packBody0200Buf); err != nil {
		return nil, err
	} else if body0200, err := unmarshalBody0200(packBody0200Buf); err != nil {
		return nil, err
	} else {
		body.PackBody0200 = body0200.(Body0200)
	}

	if _, err := reader.WriteTo(&mediaContentBuf); err != nil {
		return nil, err
	} else {
		body.MediaContent = mediaContentBuf.Bytes()
	}

	return body, nil
}
