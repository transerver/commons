package utils

import (
	"encoding/base64"
	"github.com/segmentio/ksuid"
	"testing"
	"time"
)

func TestGenerateKsUid(t *testing.T) {
	k, _ := ksuid.NewRandom()
	nowTime := time.Now().Round(1 * time.Minute)
	xTime := k.Time().Round(1 * time.Minute)
	if xTime != nowTime {
		t.Fatal(xTime, "!=", nowTime)
	}
	t.Log(k.String())
	t.Log(k.String())
	t.Log(k.String())
	t.Log(k.String())
	t.Log(k.String())
	t.Log(k.String())
	t.Log(k.String())
	t.Log(k.String())
	t.Log(k.String())
	t.Log(k.Next())
	t.Log(k.Next())
}

func TestKsUid(t *testing.T) {
	uid := NewKsUid()
	t.Log(uid.Next())
	t.Log(uid.Next().Time().UnixMilli())
	t.Log(uid.Next().Next().Time().UnixMilli())
	t.Log(uid.Next())
	t.Log(uid.Next())
	t.Log(uid.Next())
	t.Log(uid.Next())
	t.Log(uid.Next())
	t.Log(uid.Next())
	t.Log(uid.Next())
	t.Log(uid.Next())
	t.Log(uid.Next())
	t.Log(uid.Next())
	t.Log(uid.Next())
	t.Log(uid.Next())
	toString := base64.StdEncoding.EncodeToString(uid.Next().Payload())
	t.Log(toString)
}
