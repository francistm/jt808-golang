package jt808

import "errors"

var NotPackagedMessageError = errors.New("package to be concat doesn't have package property header")
var ConcatUnmarshalInvalidArgumentError = errors.New("require at least one partial PackMessage and output PackMessage")
