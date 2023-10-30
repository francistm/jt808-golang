package message

import (
	"time"
)

var (
	timeLayoutBCD = "060102150405"
	timezoneCST   = time.FixedZone("Asia/Shanghai", 8*3600)
)

type PartialPackBody struct {
	RawBody []byte `jt808:"-1,raw"`
}
