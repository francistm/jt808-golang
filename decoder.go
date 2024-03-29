package jt808

import (
	"fmt"
	"sort"

	"github.com/francistm/jt808-golang/internal/bytes"
	"github.com/francistm/jt808-golang/internal/decode"
	"github.com/francistm/jt808-golang/message"
)

//go:generate go run github.com/francistm/jt808-golang/tools/generator/decoder

// Unmarshal 由二进制解析一个完整的消息包
func Unmarshal[T message.MesgBody](data []byte, mesgPack *message.MessagePack[T]) error {
	return mesgPack.UnmarshalBinary(data)
}

// ConcatUnmarshal 拼接多个分段消息并解析
func ConcatUnmarshal[T message.MesgBody](packs []*message.MessagePack[*message.PartialPackBody], pack *message.MessagePack[T]) error {
	if len(packs) < 2 {
		return ErrIncompletedPkgMesg
	}

	if packs[0].PackHeader.Package == nil {
		return ErrNotPkgMesg
	}

	var (
		mesgBodyBuf    = bytes.NewBuffer()
		mesgBodyReader *bytes.Reader
		mesgId         = packs[0].PackHeader.MessageID
	)

	if len(packs) != int(packs[0].PackHeader.Package.Total) {
		return ErrIncompletedPkgMesg
	}

	sort.Slice(packs, func(i, j int) bool {
		var (
			packsLeft  = packs[i]
			packsRight = packs[j]
		)

		if packsLeft.PackHeader.Package == nil {
			return false
		}

		return packsLeft.PackHeader.Package.Index < packsRight.PackHeader.Package.Index
	})

	for i, pack := range packs {
		if pack.PackHeader.Package == nil {
			return ErrNotPkgMesg
		}

		if pack.PackHeader.MessageID != mesgId {
			return fmt.Errorf("message at %d is not type of %.4X", i+1, mesgId)
		}

		if pack.PackHeader.Package.Index != uint16(i+1) {
			return fmt.Errorf("message at %d is not the %dth message", i+1, i+1)
		}

		mesgId = pack.PackHeader.MessageID
		mesgBodyBuf.Write(pack.PackBody.RawBody)
	}

	pack.PackHeader = packs[0].PackHeader
	pack.PackHeader.Package = nil
	pack.PackHeader.Property.IsMultiplePackage = false
	pack.PackHeader.Property.BodyByteLength = uint16(mesgBodyBuf.Len())

	packBody, err := pack.NewPackBodyFromMesgId()

	if err != nil {
		return err
	}

	typedPackBody, ok := packBody.(T)

	if !ok {
		return fmt.Errorf("can't convert body from %T to %T", packBody, pack.PackBody)
	}

	pack.Checksum = 0
	pack.PackBody = typedPackBody
	mesgBodyReader = bytes.NewReader(mesgBodyBuf.Bytes())

	if err := decode.UnmarshalStruct(mesgBodyReader, packBody); err != nil {
		return err
	}

	return nil
}
