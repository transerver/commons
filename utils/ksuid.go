package utils

import (
	"github.com/segmentio/ksuid"
	"sync"
)

type KsUid struct {
	mutex *sync.Mutex
	uid   ksuid.KSUID
}

func (id *KsUid) Next() ksuid.KSUID {
	id.mutex.Lock()
	defer id.mutex.Unlock()

	next := id.uid.Next()
	id.uid = next
	return next
}

func NewKsUid() *KsUid {
	return &KsUid{mutex: new(sync.Mutex), uid: ksuid.New()}
}
