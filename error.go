package jt808

import "errors"

var ErrNotPackagedMessage = errors.New("package to be concat doesn't have package property header")
var ErrConcatUnmarshalInvalidArgument = errors.New("require at least one partial PackMessage and output PackMessage")
