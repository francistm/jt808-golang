package message

import (
	"time"
)

var (
	timeLayout  = "060102150405"
	locationCST = time.FixedZone("Asia/Shanghai", 8*3600)
)

// PartialPackBody
type PartialPackBody struct {
	RawBody []byte `jt808:",none"`
}
