package message

// Body0801 0x0801 消息体正文结构体
type Body0801 struct {
	MediaID          uint32        `jt808:""`
	MediaType        uint8         `jt808:""`
	MediaContentType uint8         `jt808:""`
	EventID          uint8         `jt808:""`
	ChannelID        uint8         `jt808:""`
	PackBody0200     *Body0200Base `jt808:""`
	MediaContent     []byte        `jt808:",none"`
}
