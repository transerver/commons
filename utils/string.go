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

func Blank(arg string) bool {
	if len(arg) == 0 {
		return true
	}

	if strings.TrimSpace(arg) == "" {
		return true
	}
	return false
}

func NotBlank(arg string) bool {
	return !Blank(arg)
}

// NotBlanks returns true if args has any not empty
func NotBlanks(args ...string) bool {
	for _, s := range args {
		if Blank(s) {
			return false
		}
	}

	return true
}

// Blanks returns true if args has any empty
func Blanks(args ...string) bool {
	return !NotBlanks(args...)
}

// BlanksAll returns true if args is all empty
func BlanksAll(args ...string) bool {
	for _, s := range args {
		blank := Blank(s)
		if !blank {
			return false
		}
	}

	return true
}
