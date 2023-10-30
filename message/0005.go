package message

// 终端补传分包请求
type Body0005 struct {
	SerialId   uint16
	TotalCount uint16
	RawMesgIds []byte `jt808:"-1,raw"`
}
