package utils

import "sync"

var (
	userRegisterMap  = make(map[string]struct{})
	userRegisterLock sync.Mutex
)

func Lock(key string) bool {
	userRegisterLock.Lock()
	defer userRegisterLock.Unlock()

	if _, ok := userRegisterMap[key]; ok {
		return false
	}
	userRegisterMap[key] = struct{}{}
	return true
}

func Unlock(key string) {
	delete(userRegisterMap, key)
}
