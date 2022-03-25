package utils

import (
	"reflect"
	"strings"
	"unsafe"
)

func String(data []byte) string {
	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&data))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
	}))
}

func Bytes(data string) []byte {
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&data))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
}

func NotBlank(args ...string) bool {
	for _, s := range args {
		if len(s) == 0 {
			return false
		}

		if strings.TrimSpace(s) == "" {
			return false
		}
	}

	return true
}

func AnyBlank(args ...string) bool {
	return !NotBlank(args...)
}
