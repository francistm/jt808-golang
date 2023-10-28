package message

// Body0801 0x0801 消息体正文结构体
type Body0801 struct {
	MediaID          uint32
	MediaType        uint8
	MediaContentType uint8
	EventID          uint8
	ChannelID        uint8
	PackBody0200     *Body0200Base
	MediaContent     []byte `jt808:",none"`
}
