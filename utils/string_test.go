package utils

import (
	"github.com/transerver/commons/logger"
	"testing"
)

func TestToString(t *testing.T) {
	data := []byte(RandomString(32))
	logger.Debugf("Data: %s", data)
	str := String(data)
	logger.Debug(str)
	logger.Debug(len(str))

	buf := Bytes(str)
	logger.Debugf("%s", buf)
	logger.Debug(len(buf))
	logger.Debug(cap(buf))
}
