package message

import "time"

// 查询服务器时间应答
type Body8004 struct {
	RawTime string `jt808:"6,bcd"`
}

func (b *Body8004) Time() time.Time {
	t, _ := time.ParseInLocation(timeLayoutBCD, b.RawTime, time.UTC)

	return t
}
